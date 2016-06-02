// client
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	xymmopb "guanghuan.com/xiaoyao/battery_mmo_server/pb"
	server "guanghuan.com/xiaoyao/battery_mmo_server/server"
	xylogs "guanghuan.com/xiaoyao/common/log"
	"net"
	"sync/atomic"
	"time"
	"unsafe"
)

type Client struct {
	cid        uint32
	nc         net.Conn //客户端对应的连接
	parser     server.Parser
	inqueue    chan server.MsgEntry //待处理的消息队列
	outqueue   chan server.MsgEntry //待处理的消息队列
	referCount uint32               //引用计数，销毁会话的前提的是引用计数为0，否则就必须等待
	done       chan bool            //会话主协程需要等相应的几个协程都退出才能退出
	gameMap    uint32               //地图编码
	startTime  uint64               //开始时间
}

// Lock should be held
func (c *Client) initClient() {

	c.parser.C = c

	c.inqueue = make(chan server.MsgEntry, 100)
	c.outqueue = make(chan server.MsgEntry, 100)

	//处理消息队列
	go c.processInbound()
	go c.processOutbound()
	//loop处理tcp消息
	go c.readLoop()
}

//处理消息队列
func (c *Client) processInbound() {
	atomic.AddUint32(&(c.referCount), 1)
	xylogs.Debug("go#1 Ready to get inbound msg")

	for msg := range c.inqueue {
		xylogs.Debug("go#1 Get inbound msg : %v", msg)
		switch (xymmopb.CMD_CODE)(msg.Cmd) {
		case xymmopb.CMD_CODE_CMD_S_GenerateUid:
			c.processGenerateUid(&msg)
		case xymmopb.CMD_CODE_CMD_S_GameStart:
			c.processGameStart(&msg)
		case xymmopb.CMD_CODE_CMD_S_Action:
			c.processActions(&msg)
		case xymmopb.CMD_CODE_CMD_S_GameOver:
			c.processGameOver()
		default:
			xylogs.Error("Unkown CMD %d ", msg.Cmd)
		}
	}

	c.done <- true
}

//处理消息队列
func (c *Client) processOutbound() {
	atomic.AddUint32(&(c.referCount), 1)
	buf := make([]byte, 12)
	xylogs.Debug("go#2 Ready to get outbound msg")
	for msg := range c.outqueue {
		xylogs.Debug("go#2 Get outbound msg %v ", msg)
		binary.BigEndian.PutUint16(buf, msg.Cmd)
		binary.BigEndian.PutUint16(buf[server.BYTE_MSG_INDEX:server.BYTE_MSG_INDEX+(int)(unsafe.Sizeof(c.parser.Header.Index))], 0)
		binary.BigEndian.PutUint32(buf[server.BYTE_MSG_EXT:server.BYTE_MSG_EXT+(int)(unsafe.Sizeof(c.parser.Header.Ext))], 0)
		binary.BigEndian.PutUint32(buf[server.BYTE_MSG_LEN:server.BYTE_MSG_LEN+(int)(unsafe.Sizeof(c.parser.Header.Length))], (uint32)(len(msg.MsgBuf)))
		bufTmp := append(buf[0:12], (msg.MsgBuf)...)
		//write是阻塞的？
		xylogs.Debug("go#2 Write outbound msg %v ", bufTmp)
		if _, err := c.nc.Write(bufTmp); err != nil {
			xylogs.Error("%s", err.Error())
			continue
		}
	}

	c.done <- true
}

func (c *Client) readLoop() {

	//
	time.Sleep(time.Second * 2)

	xylogs.Debug("go#3 ReadLoop begin")
	defer xylogs.Debug("go#3 ReadLoop end")

	b := make([]byte, 10240)

	atomic.AddUint32(&(c.referCount), 1)

	//主循环，读取字节流
	for {
		n, err := c.nc.Read(b)
		if err != nil {
			c.Close()
			xylogs.Error("go#3 Read error, will be closed.")
			return
		}

		xylogs.Debug("Read msg : %v", b[:n])

		//解析字节流并且进行消息处理，如果字节流不符合消息格式，认为连接非法，直接关闭这个连接
		if err := c.parser.Do(b[:n], c); err != nil {
			//todo record  msg err
			c.Close()
			xylogs.Error("parse error, will be closed.")
			return
		}
	}

	c.done <- true
}

func (c *Client) ProcessMsg() {

	msg := server.MsgEntry{
		Cmd:    c.parser.Cmd,
		MsgBuf: c.parser.MsgBuf,
	}
	xylogs.Debug("Gonna push inbound msg : %v", msg)

	c.inqueue <- msg

	xylogs.Debug("push inbound msg : %v", msg)

	return
}

func (c *Client) processGenerateUid(msg *server.MsgEntry) {

	req := &xymmopb.Msg_S_Uid{}
	err := proto.Unmarshal(msg.MsgBuf, req)

	xylogs.Debug("msg.MsgBuf : %v\n", msg.MsgBuf)

	if err != nil {
		xylogs.Error("err : %s", err.Error())
	}

	player.Uid = *(req.Uid)
	c.cid = *(req.Uid)
	xylogs.Debug("Get uid : %d", player.Uid)

	c.SendClientStart()
}

func (c *Client) SendClientStart() {
	req := &xymmopb.Msg_C_Start{
		Player: &xymmopb.Player{
			Uid: proto.Uint32(player.Uid),
		},
	}
	var (
		buf []byte
		err error
	)
	if buf, err = proto.Marshal(req); err != nil {
		xylogs.Debug("Err : %s", err.Error())
	}

	msg := server.MsgEntry{
		Cmd:    (uint16)(xymmopb.CMD_CODE_CMD_C_Start),
		MsgBuf: buf,
	}

	c.outqueue <- msg
}

func (c *Client) processGameStart(msg *server.MsgEntry) {

	req := &xymmopb.Msg_S_GameStart{}
	err := proto.Unmarshal(msg.MsgBuf, req)

	xylogs.Debug("msg.MsgBuf : %v\n", msg.MsgBuf)

	if err != nil {
		xylogs.Error("err : %s", err.Error())
	}

	c.gameMap = *(req.Map.Id)
	xylogs.Debug("GameMap : %d", c.gameMap)
	c.startTime = *(req.Time)
	xylogs.Debug("StartTime : %d", c.startTime)
	xylogs.Debug("There are %d players", len(req.Players))
	for i, p := range req.Players {
		xylogs.Debug("Player[%d] : %d", i, *(p.Uid))
	}

	//上报一个action
	time.Sleep(time.Second * 2)

	c.sendAction()
}

func (c *Client) sendAction() {

	player := &xymmopb.Player{
		Uid: proto.Uint32(player.Uid),
	}
	timestamp := proto.Int64(10000000)
	direction := proto.Int32((int32)(xymmopb.Direction_UP))
	coordinate := &xymmopb.Coordinate{
		X: proto.Int32(10),
		Y: proto.Int32(11),
	}

	action := &xymmopb.Action{
		Player:     player,
		Timestamp:  timestamp,
		Direction:  (*xymmopb.Direction)(direction),
		Coordinate: coordinate,
	}

	var actions []*xymmopb.Action

	req := &xymmopb.Msg_CS_Action{
		Actions: append(actions, action),
	}

	var (
		buf []byte
		err error
	)
	if buf, err = proto.Marshal(req); err != nil {
		xylogs.Debug("Err : %s", err.Error())
	}

	msg := server.MsgEntry{
		Cmd: (uint16)(xymmopb.CMD_CODE_CMD_C_Action),
		//		Cid:    c.cid,
		MsgBuf: buf,
	}

	c.outqueue <- msg
}

func (c *Client) processActions(msg *server.MsgEntry) {

	req := &xymmopb.Msg_CS_Action{}
	err := proto.Unmarshal(msg.MsgBuf, req)

	if err != nil {
		xylogs.Error("err : %s", err.Error())
	}

	for i, a := range req.Actions {
		xylogs.Debug("Action[%d] : %s", i, a.String())
	}

}

func (c *Client) processGameOver() {

	n := 10

	xylogs.Debug("Get GameOver msg, I will sendAction again in %d secs.", n)

	time.Sleep(time.Second * time.Duration(n))

	c.sendAction()
}

func (c *Client) Close() {

}
