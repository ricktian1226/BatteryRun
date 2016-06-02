// business_announcement
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//"strconv"
	"strings"
	"time"
)

const (
	AnnouncementConfigId = iota
	AnnouncementConfigTitle
	AnnouncementConfigMessage
	AnnouncementConfigDescription
	AnnouncementConfigBeginTime
	AnnouncementConfigEndTime
)

type Announcement struct{}

func (a *Announcement) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.DebugNoId("Announcement line %d : %s", lineNum, line)
	//道具
	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	announcement := new(battery.DBAnnouncementConfig)

	for i, v := range subs {
		switch i {
		case AnnouncementConfigId:
			convertAtoui64(v, &announcement.Id)
		case AnnouncementConfigTitle:
			announcement.Title = proto.String(v)
		case AnnouncementConfigMessage:
			announcement.Message = proto.String(v)
		case AnnouncementConfigDescription:
			announcement.Description = proto.String(v)
		case AnnouncementConfigBeginTime:
			timestamp, err := convertAtoTimestamp(v)
			if err != nil {
				xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
				break
			}
			if 0 == timestamp {
				timestamp = DEF_ENDTIME
			}
			announcement.BeginTime = proto.Int64(timestamp)
			tmp := time.Unix(timestamp, 0)
			announcement.BeginTimeStr = proto.String(fmt.Sprintf("%04d%02d%02d%02d%02d%02d", tmp.Year(), tmp.Month(), tmp.Day(),
				tmp.Hour(), tmp.Minute(), tmp.Second()))
			announcement.CreateDate = proto.String(xyutil.CurTimeStr())
		case AnnouncementConfigEndTime:
			timestamp, err := convertAtoTimestamp(v)
			if err != nil {
				xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
				break
			}
			if 0 == timestamp {
				timestamp = DEF_ENDTIME
			}
			announcement.EndTime = proto.Int64(timestamp)
			tmp := time.Unix(timestamp, 0)
			announcement.EndTimeStr = proto.String(fmt.Sprintf("%04d%02d%02d%02d%02d%02d", tmp.Year(), tmp.Month(), tmp.Day(),
				tmp.Hour(), tmp.Minute(), tmp.Second()))
		}
	}

	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.AnnouncementItem = announcement
	*ops = append(*ops, op)

}
