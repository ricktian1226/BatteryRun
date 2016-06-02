// parser
package server

import (
	//"fmt"
	"encoding/binary"
	xylog "guanghuan.com/xiaoyao/common/log"
	"strconv"
	"unsafe"
)

//分析器状态枚举值
const (
	OP_START = iota
	OP_MSG_CMD
	OP_MSG_CMD_1
	OP_MSG_INDEX
	OP_MSG_INDEX_1
	OP_MSG_EXT
	OP_MSG_EXT_1
	OP_MSG_EXT_2
	OP_MSG_EXT_3
	OP_MSG_LEN
	OP_MSG_LEN_1
	OP_MSG_LEN_2
	OP_MSG_LEN_3
	OP_MSG_CONTENT
	OP_MSG_END
)

//协议头各个字节的偏移，参见header的结构
const (
	BYTE_MSG_CMD = iota
	BYTE_MSG_CMD_1
	BYTE_MSG_INDEX
	BYTE_MSG_INDEX_1
	BYTE_MSG_EXT
	BYTE_MSG_EXT_1
	BYTE_MSG_EXT_2
	BYTE_MSG_EXT_3
	BYTE_MSG_LEN
	BYTE_MSG_LEN_1
	BYTE_MSG_LEN_2
	BYTE_MSG_LEN_3
)

//协议头
type Header struct {
	Cmd    uint16 //命令码
	Index  uint16 //避免重放攻击的校验
	Ext    uint32 //扩展字段
	Length uint32 //报文长度
}

//报文分析器
type Parser struct {
	C         ClientInterface //该分析器所属的会话
	state     int
	lengthGot uint32 //已获取的报文长度
	MsgBuf    []byte //报文临时缓存
	Header           //报文头信息
}

func (p *Parser) setByte( /*pos int,*/ state int, b byte) {
	p.MsgBuf = append(p.MsgBuf, b)
	p.state = state
}

func (p *Parser) getHeader() error {
	//获取cmd
	p.Cmd = binary.BigEndian.Uint16(p.MsgBuf[BYTE_MSG_CMD : BYTE_MSG_CMD+(int)(unsafe.Sizeof(p.Cmd))])
	//获取index
	p.Index = binary.BigEndian.Uint16(p.MsgBuf[BYTE_MSG_INDEX : BYTE_MSG_INDEX+(int)(unsafe.Sizeof(p.Index))])
	//获取ext
	p.Ext = binary.BigEndian.Uint32(p.MsgBuf[BYTE_MSG_EXT : BYTE_MSG_EXT+(int)(unsafe.Sizeof(p.Ext))])
	//获取length
	p.Length = binary.BigEndian.Uint32(p.MsgBuf[BYTE_MSG_LEN : BYTE_MSG_LEN+(int)(unsafe.Sizeof(p.Length))])

	//报文长度如果超过最大长度，则视为非法
	if p.Length > PROTO_MAX_CONTENT_SIZE {
		xylog.Error("Get proto content too long (" + strconv.Itoa((int)(p.Length)) + " > " + strconv.Itoa(PROTO_MAX_CONTENT_SIZE) + " ).")
		return protoErr(p.state, BYTE_MSG_LEN+(int)(unsafe.Sizeof(p.Length)), p.MsgBuf)
	}

	//把缓存清了
	p.MsgBuf = nil

	return nil
}

func (p *Parser) Do(buf []byte, c ClientInterface) error {
	var i int

	len := len(buf)

	for i = 0; i < len; i++ {
		switch p.state {
		case OP_START:
			p.setByte(OP_MSG_CMD, buf[i])
		case OP_MSG_CMD:
			p.setByte(OP_MSG_CMD_1, buf[i])
		case OP_MSG_CMD_1:
			p.setByte(OP_MSG_INDEX, buf[i])
		case OP_MSG_INDEX:
			p.setByte(OP_MSG_INDEX_1, buf[i])
		case OP_MSG_INDEX_1:
			p.setByte(OP_MSG_EXT, buf[i])
		case OP_MSG_EXT:
			p.setByte(OP_MSG_EXT_1, buf[i])
		case OP_MSG_EXT_1:
			p.setByte(OP_MSG_EXT_2, buf[i])
		case OP_MSG_EXT_2:
			p.setByte(OP_MSG_EXT_3, buf[i])
		case OP_MSG_EXT_3:
			p.setByte(OP_MSG_LEN, buf[i])
		case OP_MSG_LEN:
			p.setByte(OP_MSG_LEN_1, buf[i])
		case OP_MSG_LEN_1:
			p.setByte(OP_MSG_LEN_2, buf[i])
		case OP_MSG_LEN_2:
			p.setByte(OP_MSG_LEN_3, buf[i])
			//取到长度，解析长度
			if err := p.getHeader(); err != nil {
				return err
			}

			//如果消息体为空，取到消息头就开始处理
			if p.Length == 0 {
				p.C.ProcessMsg() //处理消息
				p.MsgBuf, p.Length, p.lengthGot, p.state = nil, 0, 0, OP_START
			}

		case OP_MSG_LEN_3:
			p.MsgBuf = append(p.MsgBuf, buf[i])
			p.lengthGot += 1

			switch {
			case p.Length <= p.lengthGot: //报文正文全部取到，开始处理
				p.C.ProcessMsg() //处理消息
				p.MsgBuf, p.Length, p.lengthGot, p.state = nil, 0, 0, OP_START

			default:
				continue
			}

		}
	}

	return nil
}
