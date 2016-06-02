// business
package main

import (
	"bufio"
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FILE_NAME_PROPS     = "props.conf"
	FILE_NAME_SLOTITEMS = "slotitems.conf"
	FILE_NAME_WEIGHTS   = "weights.conf"
	FILE_NAME_STAGES    = "stages.conf"
	FILE_NAME_GOODS     = "goods.conf"
)

func GetProperty(ops *[]*battery.LottoResOpItem) (err error) {

	optype := battery.RES_OP_TYPE_OP_ADD
	{
		m := &Props{}
		GetOps(FILE_NAME_PROPS, ops, &optype, m)
	}

	{
		m := &LottoSlotItems{}
		GetOps(FILE_NAME_SLOTITEMS, ops, &optype, m)
	}

	{
		m := &LottoWeight{}
		GetOps(FILE_NAME_WEIGHTS, ops, &optype, m)
	}

	{
		m := &LottoStage{}
		GetOps(FILE_NAME_STAGES, ops, &optype, m)
	}

	{
		m := &Goods{}
		GetOps(FILE_NAME_GOODS, ops, &optype, m)
	}

	xylog.Debug("WOps : %v ", ops)

	return
}

const (
	PropID = iota
	PropType
	PropResolveValue
	PropItems
	PropLottoValue
	PropValid
)

const (
	SlotItemSlotId = iota
	SlotItemDaType
	SlotItemPropID
	SlotItemWeight
	SlotItemStage
	SlotItemValid
)

const (
	WBeginValue = iota
	WEndValue
	WList
	WValid
)

const (
	SScore = iota
	SStage
	SValid
)

const (
	MallItemId = iota
	MallItemMallType
	MallItemDiscount
	MallItemPrice
	MallItemItems
	MallItemAmountPerUser
	MallItemAmountPerGame
	MallItemBestDeal
	MallItemExpireTimeStamp
	MallItemAmountRemain
	MallItemValid
)

type OpsInterface interface {
	Analyse(line string, ops *[]*battery.LottoResOpItem, optype *battery.RES_OP_TYPE)
}

type Goods struct{}

func (p *Goods) Analyse(line string, ops *[]*battery.LottoResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.Debug("line : %s", line)

	//道具
	subs := strings.Split(line, ",")

	if 0 >= len(subs) {
		return
	}

	mallitem := new(battery.MallItem)

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
			mallitem.Malltype = &pt
		case MallItemPrice:
			if nil == mallitem.Price {
				mallitem.Price = make([]*battery.MoneyItem, 0)
			}

			priceStr := strings.Split(subs[i], ":")
		LOOP2:
			for _, subPriceStr := range priceStr {
				mi := &battery.MoneyItem{}
				subPrice := strings.Split(subPriceStr, "_")
				if len(subPrice) < 2 {
					break LOOP2
				}
				mt := new(int32)
				amount := new(uint32)
				convertAtoi(subPrice[0], &mt)
				convertAtoui(subPrice[1], &amount)
				moneyType := battery.MoneyType(*mt)
				mi.Type = &moneyType
				mi.Amount = proto.Uint32(*amount)

				mallitem.Price = append(mallitem.Price, mi)

			}
		case MallItemItems:
			if nil == mallitem.Items {
				mallitem.Items = make([]*battery.PropItem, 0)
			}

			itemsStr := strings.Split(subs[i], ":")
		LOOP3:
			for _, subItemsStr := range itemsStr {
				item := &battery.PropItem{}
				subItem := strings.Split(subItemsStr, "_")
				if len(subItem) < 3 {
					break LOOP3
				}
				id := new(uint64)
				amount := new(uint32)
				propType := new(int32)
				convertAtoui64(subItem[0], &id)
				convertAtoui(subItem[1], &amount)
				convertAtoi(subItem[2], &propType)
				pType := battery.PropType(*propType)
				item.Type = &pType
				item.Amount = proto.Uint32(*amount)
				item.Id = id

				mallitem.Items = append(mallitem.Items, item)

			}

		case MallItemDiscount:
			convertAtoui(subs[i], &mallitem.Discount)
		case MallItemAmountPerUser:
			convertAtoui(subs[i], &mallitem.Amountperuser)
		case MallItemAmountPerGame:
			convertAtoui(subs[i], &mallitem.Amountpergame)
		case MallItemAmountRemain:
			convertAtoui(subs[i], &mallitem.Amountremain)
		case MallItemBestDeal:
			b, _ := strconv.ParseBool(subs[i])
			mallitem.Bestdeal = proto.Bool(b)
		case MallItemValid:
			b, _ := strconv.ParseBool(subs[i])
			mallitem.Valid = proto.Bool(b)
		case MallItemExpireTimeStamp:
			convertAtoi64(subs[i], &mallitem.Expiretimestamp)
			tmp := time.Unix(mallitem.GetExpiretimestamp(), 0)
			mallitem.Expiredate = proto.String(fmt.Sprintf("%04d%02d%02d%02d%02d%02d", tmp.Year(), tmp.Month(), tmp.Day(),
				tmp.Hour(), tmp.Minute(), tmp.Second()))
		}
		mallitem.Createdate = proto.String(xyutil.CurTimeStr())
	}

	op := &battery.LottoResOpItem{
		Optype: optype,
	}

	op.Mallitem = mallitem
	*ops = append(*ops, op)

}

type Props struct{}

func (p *Props) Analyse(line string, ops *[]*battery.LottoResOpItem, optype *battery.RES_OP_TYPE) {

	//xylog.Debug("line : %s", line)

	//道具
	subs := strings.Split(line, ",")

	if 0 >= len(subs) {
		return
	}

	prop := new(battery.Prop)

	for i, v := range subs {
		switch i {
		case PropID:
			xylog.Debug("prop[%d] id : %s", i, v)
			id := new(int32)
			convertAtoi(v, &id)
			prop.Id = proto.Uint64(uint64(*id))
		case PropType:
			xylog.Debug("prop[%d] type : %s", i, v)
			t := new(int32)
			convertAtoi(v, &t)
			pt := battery.PropType(*t)
			prop.Type = &pt
		case PropResolveValue:
			xylog.Debug("prop[%d] resolvevalue : %s", i, v)
			if nil == prop.Resolvevalue {
				prop.Resolvevalue = make([]*battery.MoneyItem, 0)
			}

			resolveStr := strings.Split(subs[i], ":")
		LOOP1:
			for _, subResolveStr := range resolveStr {
				mi := &battery.MoneyItem{}
				subPrice := strings.Split(subResolveStr, "_")
				if len(subPrice) < 2 {
					break LOOP1
				}
				mt := new(int32)
				amount := new(uint32)
				convertAtoi(subPrice[0], &mt)
				convertAtoui(subPrice[1], &amount)
				moneyType := battery.MoneyType(*mt)
				mi.Type = &moneyType
				mi.Amount = proto.Uint32(*amount)

				prop.Resolvevalue = append(prop.Resolvevalue, mi)

			}
		case PropItems:
			xylog.Debug("prop[%d] items : %s", i, v)
			if nil == prop.Items {
				prop.Items = make([]*battery.PropItem, 0)
			}

			resolveStr := strings.Split(subs[i], ":")
		LOOP3:
			for _, subResolveStr := range resolveStr {
				mi := &battery.PropItem{}
				subPrice := strings.Split(subResolveStr, "_")
				if len(subPrice) < 2 {
					break LOOP3
				}
				mt := new(int32)
				amount := new(uint32)
				convertAtoi(subPrice[0], &mt)
				convertAtoui(subPrice[1], &amount)
				id := uint64(*mt)
				mi.Id = &id
				mi.Amount = proto.Uint32(*amount)

				prop.Items = append(prop.Items, mi)

			}
		case PropLottoValue:
			convertAtoui(subs[i], &prop.Lottovalue)
		}
	}

	op := &battery.LottoResOpItem{
		Optype: optype,
	}

	prop.Valid = proto.Bool(true)

	op.Prop = prop
	*ops = append(*ops, op)

}

type LottoSlotItems struct{}

func (si *LottoSlotItems) Analyse(line string, ops *[]*battery.LottoResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, ",")
	if 0 >= len(subs) {
		return
	}
	slotItem := new(battery.LottoSlotItem)

	for i, v := range subs {
		switch i {
		case SlotItemSlotId:
			convertAtoui(v, &slotItem.Slotid)
		case SlotItemDaType:
			datype := new(int32)
			convertAtoi(v, &datype)
			datypeTmp := battery.DrawAwardType(*datype)
			slotItem.Datype = &datypeTmp
		case SlotItemPropID:
			convertAtoui64(v, &slotItem.Propid)
		case SlotItemWeight:
			convertAtoui(v, &slotItem.Weight)
		case SlotItemStage:
			convertAtoui(v, &slotItem.Stage)
		}

		op := &battery.LottoResOpItem{
			Optype: optype,
		}
		slotItem.Valid = proto.Bool(true)

		op.Slotitem = slotItem
		*ops = append(*ops, op)

	}
}

type LottoWeight struct{}

func (si *LottoWeight) Analyse(line string, ops *[]*battery.LottoResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, ",")
	if 0 >= len(subs) {
		return
	}
	weight := new(battery.LottoWeight)

	for i, v := range subs {
		switch i {
		case WBeginValue:
			convertAtoi(v, &weight.Beginvalue)
		case WEndValue:
			convertAtoi(v, &weight.Endvalue)
		case WList:
			ws := strings.Split(v, ":")
			if len(ws) != 8 {
				continue
			} else {
				weight.Weightlist = make([]uint32, 8)
				for i, w := range ws {
					tmp := new(int32)
					convertAtoi(w, &tmp)
					weight.Weightlist[i] = uint32(*tmp)
				}
			}
		case WValid:
			valid, _ := strconv.ParseBool(v)
			weight.Valid = proto.Bool(valid)

		}

		op := &battery.LottoResOpItem{
			Optype: optype,
		}

		op.Weight = weight
		*ops = append(*ops, op)

	}
}

type LottoStage struct{}

func (si *LottoStage) Analyse(line string, ops *[]*battery.LottoResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, ",")
	if len(subs) <= 0 {
		return
	}

	stage := new(battery.LottoStageItem)

	for i, v := range subs {
		xylog.Debug("subs[%d] : %v", i, v)

		switch i {
		case SScore:
			convertAtoui(v, &(stage.Score))
		case SStage:
			convertAtoui(v, &(stage.Stage))
		}
	}

	stage.Valid = proto.Bool(true)

	op := &battery.LottoResOpItem{
		Optype: optype,
	}
	op.Stage = stage

	*ops = append(*ops, op)
}

func GetOps(filename string, ops *[]*battery.LottoResOpItem, optype *battery.RES_OP_TYPE, opinterface OpsInterface) {
	file, err := os.Open(filename)
	if err != nil {
		xylog.Error("Open file %s failed", FILE_NAME_STAGES)
		return
	}
	defer file.Close()
	rb := bufio.NewReader(file)
	for {
		line, err := rb.ReadString('\n')
		if err != nil {
			return
		}
		line = line[:len(line)-2]

		xylog.Debug("line : %v", line)
		if line[0] == '#' || len(line) <= 0 { //跳过注释和空行
			continue
		}

		opinterface.Analyse(line, ops, optype)
	}
	return
}

func convertAtoi(a string, i **int32) {
	if len(a) > 0 {
		value, err := strconv.Atoi(a)
		//xylog.Debug("value : %d, err : %v", value, err)
		if err == nil {
			*i = proto.Int32(int32(value))
		}
	}
}

func convertAtoui(a string, i **uint32) {
	if len(a) > 0 {
		value, err := strconv.Atoi(a)
		//xylog.Debug("value : %d, err : %v", value, err)
		if err == nil {
			*i = proto.Uint32(uint32(value))
		}
	}
}

func convertAtoui64(a string, i **uint64) {
	if len(a) > 0 {
		value, err := strconv.Atoi(a)
		//xylog.Debug("value : %d, err : %v", value, err)
		if err == nil {
			*i = proto.Uint64(uint64(value))
		}
	}
}

func convertAtoi64(a string, i **int64) {
	if len(a) > 0 {
		value, err := strconv.Atoi(a)
		//xylog.Debug("value : %d, err : %v", value, err)
		if err == nil {
			*i = proto.Int64(int64(value))
		}
	}
}
