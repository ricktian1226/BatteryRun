// business_role
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strings"
)

const (
	RoleLevelID = iota
	RoleLevelGoldBonus
	RoleLevelScoreBonus
	Price
)

type RoleLevelMallItems struct{}

func (p *RoleLevelMallItems) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.Debug("line : %s", line)

	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	mallItem := new(battery.DBMallItem)
	roleBonus := new(battery.DBRoleLevelBonusItem)

	for i, v := range subs {
		switch i {
		case RoleLevelID:
			convertAtoui64(v, &(mallItem.Id))
			roleBonus.Id = mallItem.Id
		case RoleLevelGoldBonus:
			convertAtoi(v, &(roleBonus.GoldBonus))
		case RoleLevelScoreBonus:
			convertAtoi(v, &(roleBonus.ScoreBonus))
		case Price:
			getMoneys(lineNum, v, &mallItem.Price)
		}
	}

	// mallType
	var mallType battery.MallType = battery.MallType_Mall_Exchange
	mallItem.MallType = &mallType

	//Items
	mallItem.Items = make([]*battery.PropItem, 0)
	propItem := &battery.PropItem{}
	propItem.Id = proto.Uint64(mallItem.GetId())
	var itemType battery.PropType = battery.PropType_PROP_ROLE
	propItem.Type = &itemType
	propItem.Amount = proto.Uint32(1)
	mallItem.Items = append(mallItem.Items, propItem)

	//Discount
	mallItem.Discount = proto.Uint32(100)

	//BestDeal
	mallItem.Bestdeal = proto.Bool(false)

	//TeSell
	mallItem.Tesell = proto.Bool(false)

	//IsValid
	mallItem.Valid = proto.Bool(true)

	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.Mallitem = mallItem
	op.RoleLevelBonus = roleBonus

	*ops = append(*ops, op)

}
