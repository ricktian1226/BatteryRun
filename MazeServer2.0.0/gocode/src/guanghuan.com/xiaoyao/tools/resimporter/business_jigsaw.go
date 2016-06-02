// business_beforegameprop
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strconv"
	"strings"
)

const (
	JigsawID = iota
	Items
	UnlockProps
)

type JigsawConfigs struct{}

func (p *JigsawConfigs) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.Debug("line : %s", line)

	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	var err error
	config := new(battery.JigsawConfig)

	for i, v := range subs {
		switch i {
		case JigsawID:
			//xylog.Debug("JigsawConfigs [%d] ID : %s", i, v)
			convertAtoui64(v, &(config.Jigsawid))
		case Items:
			//xylog.Debug("JigsawConfigs [%d] Items : %s", i, v)
			itemStrList := strings.Split(v, SEP1)
			for _, itemStr := range itemStrList {
				var itemid uint64

				itemid, err = strconv.ParseUint(itemStr, 10, 64)
				if err != nil {
					continue
				}

				//xylog.Debug("itemid : [%d]", itemid)
				config.Jigsawidlist = append(config.Jigsawidlist, itemid)
			}
		case UnlockProps:
			//xylog.Debug("JigsawConfigs [%d] UnlockProps : %s", i, v)
			PropsStrList := strings.Split(v, SEP1)
			for _, PropStr := range PropsStrList {
				propItem := &battery.PropItem{}
				propPara := strings.Split(PropStr, SEP2)
				if len(propPara) < 2 {
					break
				}
				var propid uint64
				var amount uint64

				propid, err = strconv.ParseUint(propPara[0], 10, 64)
				if err != nil {
					continue
				}
				amount, err = strconv.ParseUint(propPara[1], 10, 32)
				if err != nil {
					continue
				}
				//xylog.Debug("propid : [%d]", propid)
				propItem.Id = &propid
				propItem.Amount = proto.Uint32(uint32(amount))
				config.Unlockprops = append(config.Unlockprops, propItem)
			}
		}
	}

	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.JigsawConfig = config
	*ops = append(*ops, op)

}
