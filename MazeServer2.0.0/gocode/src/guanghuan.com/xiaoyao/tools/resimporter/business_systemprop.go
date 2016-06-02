// business_prop
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strings"
)

const (
	ID = iota
	Value
	//ExpiredLimitation
)

type SystemProps struct{}

func (p *SystemProps) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.Debug("line : %s", line)

	//道具
	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	runeconfig := new(battery.RuneConfig)

	for i, v := range subs {
		switch i {
		case ID:
			propid := new(int32)
			convertAtoi(v, &propid)
			runeconfig.Propid = proto.Uint64(uint64(*propid))
		case Value:
			value := new(int32)
			convertAtoi(v, &value)
			runeconfig.Value = proto.Int32(*value)
			//case ExpiredLimitation:
			//	convertAtoi64(v, &(runeconfig.ExpiredLimitation))
		}
	}

	op := &battery.ResOpItem{
		Optype: optype,
	}

	runeconfig.Createdate = proto.String(xyutil.CurTimeStr())

	op.RuneConfig = runeconfig
	*ops = append(*ops, op)

}
