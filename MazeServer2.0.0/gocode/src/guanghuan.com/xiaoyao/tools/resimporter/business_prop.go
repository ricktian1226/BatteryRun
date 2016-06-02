// business_prop
package main

import (
    proto "code.google.com/p/goprotobuf/proto"
    //xylog "guanghuan.com/xiaoyao/common/log"
    xyutil "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "strings"
)

const (
    PropID = iota
    PropType
    PropItems
    PropResolveValue
    PropLottoValue
    PropValid
)

type Props struct{}

func (p *Props) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

    //xylog.Debug("line : %s", line)

    //道具
    subs := strings.Split(line, SEP)

    if 0 >= len(subs) {
        return
    }

    prop := new(battery.Prop)

    for i, v := range subs {
        switch i {
        case PropID:
            //xylog.Debug("prop[%d] id : %s", i, v)
            id := new(int32)
            convertAtoi(v, &id)
            prop.Id = proto.Uint64(uint64(*id))
        case PropType:
            //xylog.Debug("prop[%d] type : %s", i, v)
            t := new(int32)
            convertAtoi(v, &t)
            pt := battery.PropType(*t)
            prop.Type = &pt
        case PropResolveValue:
            //xylog.Debug("prop[%d] resolvevalue : %s", i, v)
            getMoneys(lineNum, v, &prop.Resolvevalue)
        case PropItems:
            //xylog.Debug("prop[%d] items : %s", i, v)
            getPropItems(lineNum, subs[i], &prop.Items)

        case PropLottoValue:
            convertAtoui(subs[i], &prop.Lottovalue)
        }
    }

    //校验一下道具id和道具类型是否匹配
    propItemTmp := &battery.PropItem{
        Id:   prop.Id,
        Type: prop.Type,
    }
    checkProp(propItemTmp)

    op := &battery.ResOpItem{
        Optype: optype,
    }

    prop.Valid = proto.Bool(true)

    prop.Createdate = proto.String(xyutil.CurTimeStr())
    op.Prop = prop
    *ops = append(*ops, op)

}

const (
    NewAccountSource = iota
    NewAccountPropItems
    NewAccountPropDispenseType
    NewAccountPropMailId
)

type NewAccountProp struct{}

func (p *NewAccountProp) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
    //道具
    subs := strings.Split(line, SEP)

    if 0 >= len(subs) {
        return
    }

    newAccountProp := new(battery.DBNewAccountProp)

    for i, v := range subs {
        switch i {
        case NewAccountSource:
            source := new(int32)
            convertAtoi(v, &source)
            newAccountProp.Source = battery.ID_SOURCE(*source).Enum()
        case NewAccountPropItems:
            getPropItems(lineNum, subs[i], &newAccountProp.Items)
        case NewAccountPropDispenseType:
            dispenseType := new(int32)
            convertAtoi(v, &dispenseType)
            newAccountProp.DispenseType = battery.DISPENSE_TYPE(*dispenseType).Enum()
        case NewAccountPropMailId:
            convertAtoi(v, &newAccountProp.MailId)
        }
    }

    op := &battery.ResOpItem{
        Optype: optype,
    }

    op.NewAccountProp = newAccountProp
    *ops = append(*ops, op)

}
