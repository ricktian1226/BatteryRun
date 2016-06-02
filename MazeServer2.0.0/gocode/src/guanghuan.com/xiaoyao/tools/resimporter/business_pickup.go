// business_pickup
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strings"
)

const (
	PickUp_CheckPointId = iota
	PickUp_PropType
	PickUp_PropId
	PickUp_Weight
)

type PickUps struct{}

func (p *PickUps) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.Debug("line : %s", line)

	//道具
	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	pickUpItem := new(battery.DBPickUpItem)

	for i, v := range subs {
		switch i {
		case PickUp_CheckPointId:
			convertAtoui(v, &pickUpItem.CheckPointId)
		case PickUp_PropType:
			t := new(int32)
			convertAtoi(v, &t)
			pt := battery.PropType(*t)
			pickUpItem.PropType = &pt
		case PickUp_PropId:
			convertAtoui64(v, &pickUpItem.PropId)
		case PickUp_Weight:
			convertAtoui(v, &pickUpItem.Weight)
		}
	}

	//校验一下道具id和道具类型是否匹配
	propItemTmp := &battery.PropItem{
		Id:   pickUpItem.PropId,
		Type: pickUpItem.PropType,
	}
	checkProp(propItemTmp)

	op := &battery.ResOpItem{
		Optype: optype,
	}
	pickUpItem.Valid = proto.Bool(true)
	pickUpItem.Createdate = proto.String(xyutil.CurTimeStr())
	op.PickUp = pickUpItem
	*ops = append(*ops, op)

}
