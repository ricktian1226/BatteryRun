// business_tip
package main

import (
	//"strconv"
	"strings"

	proto "code.google.com/p/goprotobuf/proto"

	//xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

//提示信息
const (
	TipId       = iota //提示信息标识
	TipLanguage        //提示语言类型
	TipTitle           //提示标题
	TipContent         //提示内容
)

type Tip struct{}

func (t *Tip) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, SEP)
	if 0 >= len(subs) {
		return
	}
	tip := new(battery.DBTip)

	for i, v := range subs {
		switch i {
		case TipId:
			id := new(int32)
			convertAtoi(v, &id)
			tip.Id = battery.TIP_IDENTITY(*id).Enum()
		case TipLanguage:
			language := new(int32)
			convertAtoi(v, &language)
			tip.Language = battery.LANGUAGE_TYPE(*language).Enum()
		case TipTitle:
			tip.Title = proto.String(v)
		case TipContent:
			tip.Content = proto.String(v)
		}
	}
	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.Tip = tip
	*ops = append(*ops, op)

}
