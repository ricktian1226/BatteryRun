// business_share
package main

import (
    proto "code.google.com/p/goprotobuf/proto"

    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "strconv"
    "strings"
)

const (
    SHAREAWARDID = iota
    SHAREAWARDCOUNTER
    SHAREAWARDREWARDS
)

type ShareAwrads struct{}

func (s *ShareAwrads) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
    subs := strings.Split(line, SEP)
    if 0 >= len(subs) {
        return
    }
    shareAwards := new(battery.DBShareAward)
    for i, v := range subs {
        switch i {
        case SHAREAWARDID:
            convertAtoui(v, &shareAwards.Id)
        case SHAREAWARDCOUNTER:
            convertAtoi(v, &shareAwards.Counter)
        case SHAREAWARDREWARDS:
            getPropItems(lineNum, subs[i], &shareAwards.Items)
        }
    }
    op := &battery.ResOpItem{
        Optype: optype,
    }
    op.ShareAwards = shareAwards
    *ops = append(*ops, op)
}

const (
    SHAREACTIVITYID = iota
    SHAREAVTICITYPE
    SHAREACTIVITYGOALVALUE
    SHAREACTIVITYDAILYLIMIT
    SHAREACTIVITYDISPENSETYPE
    SHAREACTIVITYSTARTTIME
    SHAREACTIVITYENDTIME
    SHAREACTIVITYRESTART
    SHAREACTIVITYVALID
    SHAREACTIVITYMAILID
    SHAREAWARDSBENGINTIME
    SHAREAWRADSENDTIME
)

type ShareActivity struct{}

func (s *ShareActivity) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
    subs := strings.Split(line, SEP)
    if 0 >= len(subs) {
        return
    }
    shareActicity := new(battery.DBShareActivity)
    for i, v := range subs {
        switch i {
        case SHAREACTIVITYID:
            convertAtoui(v, &shareActicity.Id)
        case SHAREAVTICITYPE:
            sharetype := new(int32)
            convertAtoi(v, &sharetype)
            shareActicity.ShareType = battery.SHARE_TYPE(*sharetype).Enum()
        case SHAREACTIVITYGOALVALUE:
            convertAtoi(v, &shareActicity.GoalValue)
        case SHAREACTIVITYDAILYLIMIT:
            convertAtoi(v, &shareActicity.DailyLimit)
        case SHAREACTIVITYDISPENSETYPE:
            dispenseType := new(int32)
            convertAtoi(v, &dispenseType)
            shareActicity.DispenseType = battery.DISPENSE_TYPE(*dispenseType).Enum()
        case SHAREACTIVITYSTARTTIME:
            timeTmp, err := convertAtoTimestamp(subs[i])
            if err != nil {
                xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
                break
            }
            if timeTmp == 0 {
                timeTmp = DEF_BEGINTIME
            }
            shareActicity.ActivityStartTime = proto.Int64(timeTmp)
        case SHAREACTIVITYENDTIME:
            timeTmp, err := convertAtoTimestamp(subs[i])
            if err != nil {
                xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
                break
            }
            if timeTmp == 0 {
                timeTmp = DEF_ENDTIME
            }
            shareActicity.ActivityEndTime = proto.Int64(timeTmp)
        case SHAREACTIVITYRESTART:
            b, _ := strconv.ParseBool(v)
            shareActicity.Restart = proto.Bool(b)
        case SHAREACTIVITYVALID:
            b, _ := strconv.ParseBool(v)
            shareActicity.Valid = proto.Bool(b)
        case SHAREACTIVITYMAILID:
            convertAtoi(v, &shareActicity.MailID)
        case SHAREAWARDSBENGINTIME:
            timeTmp, err := convertAtoTimestamp(subs[i])
            if err != nil {
                xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
                break
            }
            if timeTmp == 0 {
                timeTmp = DEF_BEGINTIME
            }
            shareActicity.AwardsBeginTime = proto.Int64(timeTmp)
        case SHAREAWRADSENDTIME:
            timeTmp, err := convertAtoTimestamp(subs[i])
            if err != nil {
                xylog.ErrorNoId("line %d : %s", lineNum, err.Error())
                break
            }
            if timeTmp == 0 {
                timeTmp = DEF_ENDTIME
            }
            shareActicity.AwardsEndTime = proto.Int64(timeTmp)
        }

    }
    op := &battery.ResOpItem{
        Optype: optype,
    }
    xylog.DebugNoId("sharetivity %v ", shareActicity)
    op.ShareActivity = shareActicity
    *ops = append(*ops, op)
}
