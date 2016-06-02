// client_test
package main

import (
	//proto "code.google.com/p/goprotobuf/proto"
	//"encoding/binary"
	"flag"
	//	"fmt"
	//xymmopb "guanghuan.com/xiaoyao/battery_mmo_server/pb"
	xylog "guanghuan.com/xiaoyao/common/log"
	"net"
	"os"
	"time"
)

type Player struct {
	Uid uint32
}

var player Player
var c Client

func main() {
	xylog.DefConfig.ProcessCmd()
	flag.Parse()
	xylog.ApplyConfig(xylog.DefConfig)
	xylog.Info("Config: %s", xylog.DefConfig.String())

	srvStr := "192.168.1.205:12345"
	conn, err := net.Dial("tcp", srvStr)
	if err != nil {
		xylog.Error("Connect to %s failed.", srvStr)
		os.Exit(0)
	}
	defer conn.Close()

	c.nc = conn

	c.initClient()

	for {
		//开始

		time.Sleep(time.Second * 3)
	}

}

/*func test(conn net.Conn) {

	//先休眠5秒钟，保证uid初始化
	time.Sleep(time.Second * 5)

	var cmd xymmopb.CMD_CODE

	//game start
	{
		cmd = xymmopb.CMD_CODE_CMD_C_Start
		xylog.Debug("===========test %s begin===========", xymmopb.CMD_CODE_name[int32(cmd)])
		defer xylog.Debug("===========test %s end===========", xymmopb.CMD_CODE_name[int32(cmd)])
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
			xylog.Debug("Err : %s", err.Error())
		}

		sendMsg(buf, cmd, conn)
	}

}

func sendMsg(buf []byte, cmd xymmopb.CMD_CODE, nc net.Conn) {
	bufTmp := make([]byte, 12)
	binary.BigEndian.PutUint16(bufTmp, (uint16)(cmd))
	binary.BigEndian.PutUint16(bufTmp[2:4], 0)
	binary.BigEndian.PutUint32(bufTmp[4:8], 0)
	binary.BigEndian.PutUint32(bufTmp[8:12], (uint32)(len(buf)))
	bufTmp = append(bufTmp, buf...)

	if _, err := nc.Write(bufTmp); err != nil {
		xylog.Error("%s Write %s failed", xymmopb.CMD_CODE_name[int32(cmd)], bufTmp)
	}
}

func readLoop(conn net.Conn) {
	buf := make([]byte, 10240)
	var bufTmp []byte
	for {
		xylog.Debug("%s read loop", conn.LocalAddr().String())

		n, err := conn.Read(buf)
		if err != nil {
			xylog.Error("Read Error : %s", err.Error())
			break
		}

		if n < 12 {
			xylog.Warning("len(msg) < 12, skip")
			continue
		}

		cmd := binary.BigEndian.Uint16(buf[:2])
		index := binary.BigEndian.Uint16(buf[2:4])
		ext := binary.BigEndian.Uint32(buf[4:8])
		length := binary.BigEndian.Uint32(buf[8:12])
		bufTmp = buf[12:(12 + length)]
		xylog.Debug("Get msg (cmd %d, index %d, ext %d, length %d, bufTmp : %v)\n", cmd, index, ext, length, bufTmp[12:12+length])
		switch (xymmopb.CMD_CODE)(cmd) {
		case xymmopb.CMD_CODE_CMD_S_GenerateUid:
			processGenerateUid(bufTmp)
			buf = nil
		case xymmopb.CMD_CODE_CMD_GameStart:
			processGameStart(bufTmp)
			buf = nil
		default:
			xylog.Error("Unkown CMD %d", cmd)
		}
		bufTmp = nil
	}

}

func processGenerateUid(buf []byte) {

	req := &xymmopb.MsgUid{}
	err := proto.Unmarshal(buf, req)

	if err != nil {
		xylog.Error("err : %s", err.Error())
	}

	player.Uid = *(req.Uid)
	xylog.Debug("Get uid : %d", player.Uid)
}

func processGameStart(buf []byte) {
	req := &xymmopb.MsgGameStart{}
	err := proto.Unmarshal(buf, req)

	if err != nil {
		xylog.Error("err : %s", err.Error())
	}

	xylog.Debug("Get MsgGameStart msg : Map : %d, players : %v", *req.Map, req.Players)
}*/
