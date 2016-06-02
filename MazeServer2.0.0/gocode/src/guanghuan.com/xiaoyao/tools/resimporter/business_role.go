// business_role
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strings"
)

type RoleInfos struct{}

func (p *RoleInfos) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	const (
		RoleID = iota
		MaxLevel
		JigsawID
		IsDefaultOwn
	)

	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	config := new(battery.DBRoleInfoConfig)

	for i, v := range subs {
		switch i {
		case RoleID:
			convertAtoui64(v, &(config.Id))
		case MaxLevel:
			convertAtoi(v, &(config.MaxLevel))
		case JigsawID:
			convertAtoui64(v, &(config.JigsawId))
		case IsDefaultOwn:
			convertParseBool(v, &(config.IsDefaultOwn))
		}
	}

	config.IsValid = proto.Bool(true)

	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.RoleInfoConfig = config
	*ops = append(*ops, op)

}
