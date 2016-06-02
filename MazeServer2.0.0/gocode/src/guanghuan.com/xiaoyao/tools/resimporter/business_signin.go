// business_signin
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
	SignInActivity_Id = iota
	SignInActivity_Type
	SignInActivity_GoalValue
	SignInActivity_BeginTime
	SignInActivity_EndTime
	SignInActivity_AutoCollect
	SignInActivity_Valid
)

type SignInActivity struct{}

func (s *SignInActivity) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	activity := new(battery.DBSignInActivity)

	for i, v := range subs {
		switch i {
		case SignInActivity_Id:
			convertAtoui64(v, &(activity.Id))
		case SignInActivity_Type:
			t := new(int32)
			convertAtoi(v, &t)
			at := battery.SignInActivityType(*t)
			activity.Type = at.Enum()
		case SignInActivity_GoalValue:
			convertAtoui(v, &(activity.GoalValue))
		case SignInActivity_BeginTime:
			timeTmp, err := convertAtoTimestamp(subs[i])
			if err != nil {
				xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
				break
			}
			activity.BeginTime = proto.Int64(timeTmp)

		case SignInActivity_EndTime:
			timeTmp, err := convertAtoTimestamp(subs[i])
			if err != nil {
				xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
				break
			}
			activity.EndTime = proto.Int64(timeTmp)
			if 0 == activity.GetEndTime() {
				activity.EndTime = proto.Int64(DEF_ENDTIME)
			}
		case SignInActivity_AutoCollect:
			b, _ := strconv.ParseBool(v)
			activity.AutoCollect = proto.Bool(b)
		case SignInActivity_Valid:
			b, _ := strconv.ParseBool(v)
			activity.Valid = proto.Bool(b)
		}
	}

	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.Activity = activity
	*ops = append(*ops, op)
}

const (
	SignInItem_Id = iota
	SignInItem_Value
	SignInItem_Award
	SignInItem_Valid
)

type SignInItem struct{}

func (s *SignInItem) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	signInItem := new(battery.DBSignInItem)

	for i, v := range subs {
		switch i {
		case SignInItem_Id:
			convertAtoui64(v, &(signInItem.Id))
		case SignInItem_Value:
			convertAtoui(v, &(signInItem.Value))
		case SignInItem_Award:
			signInItem.Award = make([]*battery.PropItem, 0)
			getPropItems(lineNum, v, &signInItem.Award)
			//case SignInItem_Valid:
			//	b, _ := strconv.ParseBool(v)
			//	signInItem.Valid = proto.Bool(b)
		}
	}

	signInItem.Valid = proto.Bool(true)

	op := &battery.ResOpItem{
		Optype: optype,
	}

	signInItem.Createdate = proto.String(xyutil.CurTimeStr())
	op.SigninItem = signInItem
	*ops = append(*ops, op)
}
