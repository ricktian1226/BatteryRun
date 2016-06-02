// business_goods
package main

import (
    proto "code.google.com/p/goprotobuf/proto"
    "fmt"
    xylog "guanghuan.com/xiaoyao/common/log"
    xyutil "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "strconv"
    "strings"
    "time"
)

const (
    MallItemId = iota
    MallItemMallType
    MallItemSubType
    MallItemPosIndex
    MallItemDiscount
    MallItemPrice
    MallItemIapId
    MallItemItems
    MallItemAmountPerUser
    MallItemAmountPerRound
    MallItemAmountPerDay
    MallItemBestDeal
    MallItemTeSell
    MallItemExpireDate
    MallItemValid
    MallItemIcon
    MallItemName
    MallItemDescription
    MallItemLabel
    MallItemDisplayIds
    MallItemGiftIds
    MallItemMultiple
)

type Goods struct{}

type MAPId2Flag map[uint64]*battery.MallItemFlag

func (p *Goods) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {

    //道具
    subs := strings.Split(line, SEP)

    if 0 >= len(subs) {
        return
    }

    mapId2Flag := make(MAPId2Flag, 0)

    mallitem := new(battery.DBMallItem)

    for i, v := range subs {
        switch i {
        case MallItemId:
            id := new(int32)
            convertAtoi(v, &id)
            mallitem.Id = proto.Uint64(uint64(*id))
        case MallItemMallType:
            t := new(int32)
            convertAtoi(v, &t)
            pt := battery.MallType(*t)
            mallitem.MallType = &pt
        case MallItemSubType:
            t := new(int32)
            convertAtoi(v, &t)
            pt := battery.MallSubType(*t)
            mallitem.MallSubType = &pt
        case MallItemPosIndex:
            t := new(uint32)
            convertAtoui(v, &t)
            mallitem.PosIndex = proto.Uint32(*t)
        case MallItemPrice:
            getMoneys(lineNum, v, &mallitem.Price)
        case MallItemIapId:
            mallitem.Iapid = proto.String(v)
        case MallItemItems:
            getPropItems(lineNum, subs[i], &(mallitem.Items))
            for _, item := range mallitem.Items { //为每个商品子项保存一个标志
                mapId2Flag[item.GetId()] = &battery.MallItemFlag{
                    Id:   proto.Uint64(item.GetId()),
                    Flag: proto.Int64(0),
                }
            }

        case MallItemDiscount:
            convertAtoui(subs[i], &mallitem.Discount)
        case MallItemAmountPerUser:
            convertAtoui(subs[i], &mallitem.Amountperuser)
        case MallItemAmountPerRound:
            convertAtoui(subs[i], &mallitem.Amountpergame)
        case MallItemAmountPerDay:
            convertAtoui(subs[i], &mallitem.Amountperday)
        case MallItemBestDeal:
            b, _ := strconv.ParseBool(subs[i])
            mallitem.Bestdeal = proto.Bool(b)
        case MallItemTeSell:
            b, _ := strconv.ParseBool(subs[i])
            mallitem.Tesell = proto.Bool(b)
        case MallItemValid:
            b, _ := strconv.ParseBool(subs[i])
            mallitem.Valid = proto.Bool(b)
        case MallItemExpireDate:
            timestamp, err := convertAtoTimestampDate(subs[i])
            if err != nil {
                xylog.ErrorNoId("line %d, i %d : %s", lineNum, i, err.Error())
                break
            }
            if 0 == timestamp {
                timestamp = DEF_ENDTIME
            }
            mallitem.Expiretimestamp = proto.Int64(timestamp)
            tmp := time.Unix(timestamp, 0)
            mallitem.Expiredate = proto.String(fmt.Sprintf("%04d%02d%02d%02d%02d%02d", tmp.Year(), tmp.Month(), tmp.Day(),
                tmp.Hour(), tmp.Minute(), tmp.Second()))
        case MallItemIcon:
            convertAtoi(subs[i], &mallitem.Icon)
        case MallItemName:
            mallitem.Name = proto.String(v)
        case MallItemDescription:
            mallitem.Description = proto.String(v)
        case MallItemLabel:
            convertAtoi(subs[i], &mallitem.Label)
        case MallItemDisplayIds:
            ids := getIds(v)
            setMallItemFlag(mapId2Flag, ids, battery.MALL_ITEM_FLAG_MALL_ITEM_FLAG_DISPLAY)
        case MallItemGiftIds:
            ids := getIds(v)
            setMallItemFlag(mapId2Flag, ids, battery.MALL_ITEM_FLAG_MALL_ITEM_FLAG_GIFT)
        case MallItemMultiple:
            convertAtoi(v, &mallitem.Multiple)
        }

        mallitem.Createdate = proto.String(xyutil.CurTimeStr())
    }

    for _, flag := range mapId2Flag {
        mallitem.ItemFlags = append(mallitem.ItemFlags, flag)
    }

    op := &battery.ResOpItem{
        Optype: optype,
    }

    //xylog.DebugNoId("mallitem : %v", mallitem)

    op.Mallitem = mallitem
    *ops = append(*ops, op)

}

const (
    CheckPointId = iota
    GoodsId
)

type UnlockGoodsConfig struct{}

func (u *UnlockGoodsConfig) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
    subs := strings.Split(line, SEP)

    if 0 >= len(subs) {
        return
    }
    unlockConfig := new(battery.DBCheckPointUnlockGoodsConfig)
    for i, v := range subs {
        switch i {
        case CheckPointId:
            convertAtoui(v, &unlockConfig.CheckPointId)
        case GoodsId:
            convertAtoui64(v, &unlockConfig.GoodsId)
        }
    }
    op := &battery.ResOpItem{
        Optype: optype,
    }
    op.UnlockGoodsConfig = unlockConfig
    *ops = append(*ops, op)
    xylog.DebugNoId("%v", unlockConfig)
}
