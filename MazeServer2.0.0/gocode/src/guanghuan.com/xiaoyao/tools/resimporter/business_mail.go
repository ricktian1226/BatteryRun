// business_beforegameprop
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	configparse "guanghuan.com/xiaoyao/tools/configparse"
	"strconv"
	"strings"
	"time"
)

type MailConfigs struct{}

func (p *MailConfigs) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	const (
		MailID = iota
		Title
		Message
		Description
		Type
		PropID
		StartTime
		EndTime
	)

	//xylog.Debug("line : %s", line)

	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	var tempTime time.Time

	var err error
	config := new(battery.SystemMailInfoConfig)

	for i, v := range subs {
		switch i {
		case MailID:
			//xylog.Debug("MailConfig [%d] MailID : %s", i, v)
			convertAtoi(v, &(config.MailID))
		case Title:
			//xylog.Debug("MailConfig [%d] Title : %s", i, v)
			config.Title = proto.String(v)
		case Message:
			//xylog.Debug("MailConfig [%d] Message : %s", i, v)
			config.Message = proto.String(v)
		case Description:
			//xylog.Debug("MailConfig [%d] Message : %s", i, v)
			if len(v) > 0 {
				config.Description = proto.String(v)
			}
		case Type:
			//xylog.Debug("MailConfig [%d] Type : %s", i, v)
			var mailtype battery.SystemMailType = battery.SystemMailType_SystemMailType_Gift
			var typeValue int64 = 0
			if len(v) > 0 {
				typeValue, err = strconv.ParseInt(v, 10, 32)
				if typeValue <= 0 {
					xylog.ErrorNoId("sub typeValue <= 0 : [%s]", v)
					continue
				}
				mailtype = battery.SystemMailType(typeValue)
			} else {
				xylog.ErrorNoId("sub type <= 0 : [%s]", v)
				continue
			}
			config.Mailtype = &mailtype
		case PropID:
			//xylog.Debug("MailConfig [%d] PropID : %s", i, v)
			convertAtoui64(v, &(config.PropID))
		case StartTime:
			//xylog.Debug("MailConfig [%d] StartTime : %s", i, v)
			var startTime int64 = 0
			if len(v) > 0 {
				tempTime, err = time.Parse(configparse.TimeFormat, v)
				startTime = tempTime.Unix()
				if err != nil {
					xylog.ErrorNoId("time.Parse StartTime error strValue = [%s]", v)
					continue
				} else {
					//xylog.Debug("StartTime = %v", tempTime)
					//xylog.Debug("StartTime = %d", startTime)
				}
			} else {
				xylog.ErrorNoId("time.Parse StartTime error strValue = [%s]", v)
				continue
			}
			config.StartTime = proto.Int64(startTime)
			//xylog.Debug("StartTime = %d", *config.StartTime)
		case EndTime:
			//xylog.Debug("MailConfig [%d] EndTime : %s", i, v)
			var tempTime2 time.Time
			var endTime int64 = 0
			if len(v) > 0 {
				tempTime2, err = time.Parse(configparse.TimeFormat, v)
				endTime = tempTime2.Unix()
				if err != nil {
					xylog.ErrorNoId("time.Parse EndTime error strValue = [%s]", v)
					continue
				} else {
					//xylog.Debug("EndTime = %v", tempTime2)
					//xylog.Debug("EndTime = %d", endTime)
				}
			} else {
				xylog.ErrorNoId("time.Parse EndTime error strValue = [%s]", v)
				continue
			}
			config.EndTime = proto.Int64(endTime)
			//xylog.Debug("EndTime = %d", *config.EndTime)
		}
	}

	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.MailConfig = config
	*ops = append(*ops, op)

}
