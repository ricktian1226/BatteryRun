// business_mission

package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strconv"
	"strings"
)

const (
	MissionConfig_id = iota
	MissionConfig_type
	MissionConfig_relatedMissions
	MissionConfig_relatedProps
	MissionConfig_quotas
	MissionConfig_rewards
	MissionConfig_begintime
	MissionConfig_endtime
	MissionConfig_autocollect
	MissionConfig_valid
	MissionConfig_tipid
	MissionConfig_tipdesc
	MissionConfig_expiredrestart
	MissionConfig_priority
)

const (
	MissionConfig_Quota_Id = iota
	MissionConfig_Quota_CycleType
	MissionConfig_Quota_CycleValue
	MissionConfig_Quota_GoalValue
	MissionConfig_Quota_ConstraintEqual
	MissionConfig_Quota_Max
)

type MissionConfig struct{}

func (m *MissionConfig) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	missionItem := &battery.MissionItem{}

	for i, v := range subs {
		switch i {
		case MissionConfig_id:
			id := new(int32)
			convertAtoi(v, &id)
			missionItem.Id = proto.Uint64(uint64(*id))

		case MissionConfig_type:
			missionType := new(int32)
			convertAtoi(v, &missionType)
			missionItem.Type = battery.MissionType(*missionType).Enum()

		case MissionConfig_relatedMissions:
			missionItem.RelatedMissions = make([]uint64, 0)
			missionStrs := strings.Split(subs[i], SEP1)
			for _, ms := range missionStrs {

				if len(ms) <= 0 {
					continue
				}

				mid := new(uint64)
				convertAtoui64(ms, &mid)
				missionItem.RelatedMissions = append(missionItem.RelatedMissions, *mid)
			}
		case MissionConfig_relatedProps:
			missionItem.RelatedProps = make([]*battery.PropItem, 0)
			getPropItems(lineNum, subs[i], &missionItem.RelatedProps)

		case MissionConfig_quotas:
			missionItem.Quotas = make([]*battery.MissionQuotaItem, 0)
			quotaStrs := strings.Split(subs[i], SEP1)
			for _, quotaStr := range quotaStrs {
				quotaItems := strings.Split(quotaStr, SEP2)
				if len(quotaItems) != MissionConfig_Quota_Max {
					continue
				}

				missionQuotaItem := &battery.MissionQuotaItem{}
				for j, quotaItem := range quotaItems {
					switch j {
					case MissionConfig_Quota_Id:
						id := new(int32)
						convertAtoi(quotaItem, &id)
						missionQuotaItem.Id = battery.QuotaEnum(*id).Enum()
					case MissionConfig_Quota_CycleType:
						cycleType := new(int32)
						convertAtoi(quotaItem, &cycleType)
						missionQuotaItem.CycleType = battery.QuotaCycleType(*cycleType).Enum()
					case MissionConfig_Quota_CycleValue:
						convertAtoi64(quotaItem, &missionQuotaItem.CycleValue)
					case MissionConfig_Quota_GoalValue:
						convertAtoui64(quotaItem, &missionQuotaItem.GoalValue)
					case MissionConfig_Quota_ConstraintEqual:
						convertParseBool(quotaItem, &missionQuotaItem.ConstraintEqual)
					}
				}
				missionItem.Quotas = append(missionItem.Quotas, missionQuotaItem)
			}
		case MissionConfig_rewards:
			missionItem.Rewards = make([]*battery.PropItem, 0)
			getPropItems(lineNum, subs[i], &missionItem.Rewards)
		case MissionConfig_begintime:
			timeTmp, err := convertAtoTimestamp(subs[i])
			if err != nil {
				xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
				continue
			}

			if 0 == timeTmp {
				timeTmp = DEF_BEGINTIME
			}
			missionItem.Begintime = proto.Int64(timeTmp)
		case MissionConfig_endtime:
			timeTmp, err := convertAtoTimestamp(subs[i])
			if err != nil {
				xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
				continue
			}

			if 0 == timeTmp {
				timeTmp = DEF_ENDTIME
			}
			missionItem.Endtime = proto.Int64(timeTmp)
		case MissionConfig_autocollect:
			b, _ := strconv.ParseBool(subs[i])
			missionItem.AutoCollect = proto.Bool(b)
		case MissionConfig_valid:
			b, _ := strconv.ParseBool(subs[i])
			missionItem.Valid = proto.Bool(b)
		case MissionConfig_tipid:
			convertAtoui(v, &(missionItem.Tipid))
		case MissionConfig_tipdesc:
			missionItem.TipDesc = proto.String(v)
		case MissionConfig_priority:
			if len(v) > 0 {
				convertAtoui(v, &(missionItem.Priority))
			} else {
				missionItem.Priority = proto.Uint32(battery.Default_MissionItem_Priority)
			}

		case MissionConfig_expiredrestart:
			b, err := strconv.ParseBool(subs[i])
			if err == nil {
				missionItem.ExpiredRestart = proto.Bool(b)
			} else {
				missionItem.ExpiredRestart = proto.Bool(false)
			}
		}
	}

	op := &battery.ResOpItem{
		Optype: optype,
	}

	missionItem.Createdate = proto.String(xyutil.CurTimeStr())
	//xylog.Debug("mission : %v", missionItem)
	op.MissionItem = missionItem
	*ops = append(*ops, op)
}
