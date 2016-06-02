// client
package server

import (
	proto "code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	xymmopb "guanghuan.com/xiaoyao/battery_mmo_server/pb"
	xylogs "guanghuan.com/xiaoyao/common/log"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBufSize = 32768
	msgScratchSize = 512
)

const (
	STATE_NOK = iota //会话状态：不可用
	STATE_OK         //会话状态：可用
)

const (
	NO_IN_ROOM = iota //会话状态：不可用
	IN_ROOM           //会话状态：可用
)

type stats struct {
	inMsgs   int64
	outMsgs  int64
	inBytes  int64
	outBytes int64
}

type MsgEntry struct {
	Cmd    uint16 //命令码
	MsgBuf []byte //报文临时缓存
}

type ClientInterface interface {
	ProcessMsg()
}

type ClientID uint32

type Client struct {
	mu         sync.RWMutex
	cid        ClientID
	nc         net.Conn //客户端对应的连接
	srv        *Server
	roomstate  uint32 //玩家的房间状态，0 表示游离态，1 表示在房间中
	room       *Room  //客户端所属的room
	parser     Parser
	inqueue    chan MsgEntry //待处理的消息队列
	outqueue   chan MsgEntry //待处理的消息队列
	referCount uint32        //引用计数，销毁会话的前提的是引用计数为0，否则就必须等待
	//done       chan bool     //会话主协程需要等相应的几个协程都退出才能退出
	//state      uint32        //会话状态，为了使用atomic，定义为uint32，而不用bool
	pingtimer PingTimer //心跳监控
	stats
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (c *Client) logDebug(info string) {
	xylogs.Debug("Client #%d %s "+info, c.cid, clientConnStr(c.nc))
}

func (c *Client) logError(info string) {
	xylogs.Error("Client #%d %s "+info, c.cid, clientConnStr(c.nc))
}

func (c *Client) logInfo(info string) {
	xylogs.Info("Client #%d %s "+info, c.cid, clientConnStr(c.nc))
}

func (c *Client) logWarn(info string) {
	xylogs.Warning("Client #%d %s "+info, c.cid, clientConnStr(c.nc))
}

func (c *Client) logTrace(info string) {
	xylogs.Trace("Client #%d %s "+info, c.cid, clientConnStr(c.nc))
}

func (c *Client) logFatal(info string) {
	xylogs.Fatal("Client #%d %s "+info, c.cid, clientConnStr(c.nc))
}

func clientConnStr(conn net.Conn) interface{} {
	if ip, ok := conn.(*net.TCPConn); ok {
		addr := ip.RemoteAddr().(*net.TCPAddr)
		return []string{fmt.Sprintf("%v:%d", addr.IP, addr.Port)}
	}
	return "N/A"
}

// Lock should be held
func (c *Client) initClient() {
	s := c.srv

	c.cid = (ClientID)(atomic.AddUint32((*uint32)(&s.gcid), 1))

	c.parser.C = c

	c.roomstate = NO_IN_ROOM //初始未在房间中

	//处理消息队列
	go c.processInbound()
	go c.processOutbound()

	//loop处理tcp消息
	go c.readLoop()

	//会话id下发给客户端
	c.sendCid()

	//协程1 专门用来做心跳定时处理
	c.pingtimer.SetDier(c)
	c.pingtimer.Set()

}

//处理消息队列
func (c *Client) processInbound() {
	c.logDebug("Ready to get inbound msg")

	atomic.AddUint32(&(c.referCount), 1)

	for msg := range c.inqueue {
		c.logDebug("Get inbound msg : " + xymmopb.CMD_CODE_name[(int32)(msg.Cmd)])
		switch (xymmopb.CMD_CODE)(msg.Cmd) {
		case xymmopb.CMD_CODE_CMD_C_Start:
			c.processClientStartMsg()
		case xymmopb.CMD_CODE_CMD_C_Action:
			c.processActionMsg(&msg)
		case xymmopb.CMD_CODE_CMD_C_PING:
			c.processPingMsg(&msg)
		default:
			xylogs.Error("Unkown CMD %d ", msg.Cmd)
		}
	}
}

//处理消息队列
func (c *Client) processOutbound() {
	buf := make([]byte, 12)
	c.logDebug("Ready to get outbound msg")

	atomic.AddUint32(&(c.referCount), 1)

	for msg := range c.outqueue {
		c.logDebug("Get outbound msg : " + xymmopb.CMD_CODE_name[(int32)(msg.Cmd)])
		binary.BigEndian.PutUint16(buf, msg.Cmd)
		binary.BigEndian.PutUint16(buf[2:4], 0)
		binary.BigEndian.PutUint32(buf[4:8], 0)
		binary.BigEndian.PutUint32(buf[8:12], (uint32)(len(msg.MsgBuf)))
		buf = append(buf[0:12], (msg.MsgBuf)...)
		//如果写失败，则连接跳出处理循环
		if _, err := c.nc.Write(buf); err != nil {
			c.logError("Write failed, " + err.Error())
			break
		}
	}
}

func (c *Client) readLoop() {

	c.logDebug("readLoop begin")
	defer c.logDebug("readLoop end")

	b := make([]byte, defaultBufSize)

	atomic.AddUint32(&(c.referCount), 1)

	//主循环，读取字节流
	c.logDebug("Read begin")
	for {
		//如果会话状态为不可用。则退出
		/*if STATE_OK != atomic.LoadUint32(&(c.state)) {
			c.logDebug("Client STATE_NOK Stop readLoop.")
			break
		}*/

		n, err := c.nc.Read(b)
		if err != nil {
			xylogs.Debug("Read Error : " + err.Error() + ", Stop readLoop.")
			break
		}

		//解析字节流并且进行消息处理，如果字节流不符合消息格式，认为连接非法，直接关闭这个连接
		if err := c.parser.Do(b[:n], c); err != nil {
			c.logError("Parse Error, Stop readLoop.")
			break
		}
	}

	//告诉主协程，readLoop协程退出
	//c.done <- true
	c.Close()
}

func (c *Client) ProcessMsg() {

	msg := MsgEntry{
		Cmd:    c.parser.Cmd,
		MsgBuf: c.parser.MsgBuf,
	}

	c.inqueue <- msg

	//重置一下心跳时间
	//c.pingtimer.RefleshTime()

	return
}

func (c *Client) Close() {

	//先退出房间(如果已加入房间，删除房间下的引用信息)
	if c.room != nil {
		c.room.RemoveClient(c)
		c.LeaveRoom()
	}

	/*if STATE_OK == atomic.LoadUint32(&(c.state)) {
		atomic.StoreUint32(&(c.state), STATE_NOK)
	}*/

	//销毁通道
	close(c.inqueue)
	close(c.outqueue)
	//close(c.done)

	//关闭tcp连接
	c.clearConnection()

	//把服务器信息中该客户端的相关信息删除
	if c.srv != nil {
		c.srv.removeClient(c)
	}
}

func (c *Client) clearConnection() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nc == nil {
		return
	}
	c.nc.Close()
	c.nc = nil
}

//生成uid下发给客户端
func (c *Client) sendCid() {
	req := &xymmopb.Msg_S_Uid{
		Uid: (*uint32)(&c.cid),
	}

	buf, err := proto.Marshal(req)
	if err != nil {
		c.logDebug("proto.Marshal failed.")
		return
	}

	var msg MsgEntry
	msg.Cmd = (uint16)(xymmopb.CMD_CODE_CMD_S_GenerateUid)
	msg.MsgBuf = buf

	c.outqueue <- msg
}

//处理客户端开始请求
func (c *Client) processClientStartMsg() {

	//先将玩家加入房间
	gRoom.AddClient(c)

	//尝试启动房间游戏
	if ok, _ := gRoom.Start(); ok {
		xylogs.Debug("Room#%d Start.", gRoom.rid)
	} else {
		xylogs.Debug("Room#%d NoStart.", gRoom.rid)
	}
}

//处理动作请求
func (c *Client) processActionMsg(msg *MsgEntry) {
	var msgOut MsgEntry
	msgOut.Cmd = (uint16)(xymmopb.CMD_CODE_CMD_S_Action)
	msgOut.MsgBuf = msg.MsgBuf

	//向所在房间的所有玩家广播动作
	c.mu.RLock()
	defer c.mu.RUnlock()

	//加个保护，玩家在联机游戏中时动作消息才需要
	if IN_ROOM == c.roomstate {
		c.room.Broadcast(&msgOut)
	}
}

//处理心跳请求
func (c *Client) processPingMsg(msg *MsgEntry) {
	var msgOut MsgEntry
	msgOut.Cmd = (uint16)(xymmopb.CMD_CODE_CMD_S_PING)
	msgOut.MsgBuf = msg.MsgBuf

	c.outqueue <- msgOut
}

//心跳过期的处理函数
func (c *Client) Die() {
	//只要关闭socket就会触犯一系列的会话销毁
	c.clearConnection()
}

//加入房间
func (c *Client) JoinRoom(r *Room) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.room = r
	c.roomstate = IN_ROOM
}

//离开房间
func (c *Client) LeaveRoom() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.room = nil
	c.roomstate = NO_IN_ROOM
}
