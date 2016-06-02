// business_beforegameprop
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strings"
)

const (
	ItemID = iota
	Weight
)

type BeforeGameProps struct{}

func (p *BeforeGameProps) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.Debug("line : %s", line)

	//赛前道具权重
	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	config := new(battery.DBBeforeGameRandomGoodWeight)

	for i, v := range subs {
		switch i {
		case ItemID:
			//xylog.Debug("RandomGood [%d] ID : %s", i, v)
			convertAtoui64(v, &(config.GoodId))
		case Weight:
			//xylog.Debug("RandomGood [%d] Weight : %s", i, v)
			convertAtoui(v, &(config.Weight))
		}
	}

	config.Valid = proto.Bool(true)

	op := &battery.ResOpItem{
		Optype: optype,
	}

	config.Createdate = proto.String(xyutil.CurTimeStr())
	op.BeforeGameRandomGoodWeight = config
	*ops = append(*ops, op)

}
