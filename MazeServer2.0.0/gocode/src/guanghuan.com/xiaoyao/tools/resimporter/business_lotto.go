// business_lotto
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strconv"
	"strings"
)

//抽奖内部价值对应的权重列表
const (
	WBeginValue = iota //权重区间起始值
	WEndValue          //权重区间结束值
	WList              //权重列表
	WValid             //是否可用
)

//特殊抽奖配置信息
const (
	SpecificLottoSerialNum = iota //特殊抽奖序列号
	SpecificLottoPropList         //礼包列表
	SpecificLottoSelected         //选中的格子
	SpecificLottoValid            //是否可用
)

//游戏阶段映射表
const (
	//SScore      = iota //游戏分数
	SQuotaId    = iota //指标
	SQuotaValue        //指标值
	SStage             //阶段
	SValid             //是否可用
)

//抽奖奖池信息
const (
	SlotItemSlotId = iota //抽奖格子编码
	SlotItemDaType        //抽奖类型
	SlotItemPropID        //礼包id
	SlotItemWeight        //权重
	SlotItemStage         //阶段编号（游戏后抽奖用）
	SlotItemValid         //是否可用
)

const (
	LottoSlotCount = 8
)

type LottoSlotItems struct{}

func (si *LottoSlotItems) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, SEP)
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
		case SlotItemValid:
			valid, _ := strconv.ParseBool(v)
			slotItem.Valid = proto.Bool(valid)

		}
	}
	op := &battery.ResOpItem{
		Optype: optype,
	}

	slotItem.Createdate = proto.String(xyutil.CurTimeStr())
	op.Slotitem = slotItem
	*ops = append(*ops, op)

}

type LottoWeight struct{}

func (si *LottoWeight) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, SEP)

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
			ws := strings.Split(v, SEP1)
			weight.Weightlist = make([]uint32, 0)
			for _, w := range ws {
				//xylog.Debug("i %d w [%s] len(w) %d", i, w, len(w))
				if len(w) > 0 {
					tmp := new(int32)
					convertAtoi(w, &tmp)
					weight.Weightlist = append(weight.Weightlist, uint32(*tmp))
				}
			}

			if len(weight.Weightlist) != LottoSlotCount {
				xylog.ErrorNoId("weight list %s is invalid", v)
				continue
			}

		case WValid:
			valid, _ := strconv.ParseBool(v)
			weight.Valid = proto.Bool(valid)

		}

	}

	//xylog.Debug("weight: %v", weight)

	op := &battery.ResOpItem{
		Optype: optype,
	}

	weight.Createdate = proto.String(xyutil.CurTimeStr())
	op.Weight = weight
	*ops = append(*ops, op)
}

type LottoSerialNumSlot struct{}

func (si *LottoSerialNumSlot) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, SEP)

	if 0 >= len(subs) {
		return
	}

	s := new(battery.LottoSerialNumSlot)

	for i, v := range subs {
		switch i {
		case SpecificLottoSerialNum:
			convertAtoi(v, &s.SerialNum)
		case SpecificLottoPropList:
			ws := strings.Split(v, SEP1)
			s.PropList = make([]uint64, 0)
			for _, w := range ws {
				//xylog.Debug("i %d w [%s] len(w) %d", i, w, len(w))
				if len(w) > 0 {
					tmp := new(uint64)
					convertAtoui64(w, &tmp)
					s.PropList = append(s.PropList, uint64(*tmp))
				}
			}

			if len(s.PropList) != LottoSlotCount {
				xylog.ErrorNoId("prop list %s is invalid", v)
				continue
			}
		case SpecificLottoSelected:
			convertAtoui(v, &s.Selected)
		case SpecificLottoValid:
			valid, _ := strconv.ParseBool(v)
			s.Valid = proto.Bool(valid)

		}

	}

	//xylog.Debug("weight: %v", weight)

	op := &battery.ResOpItem{
		Optype: optype,
	}

	s.Createdate = proto.String(xyutil.CurTimeStr())
	op.LottoSerialNumSlot = s
	*ops = append(*ops, op)
}

type LottoStage struct{}

func (si *LottoStage) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, SEP)
	if len(subs) <= 0 {
		return
	}

	stage := new(battery.LottoStageItem)

	for i, v := range subs {
		xylog.DebugNoId("subs[%d] : %v", i, v)

		switch i {
		case SQuotaId:
			quotaId := new(int32)
			convertAtoi(v, &quotaId)
			stage.QuotaId = battery.QuotaEnum(*quotaId).Enum()
		case SQuotaValue:
			convertAtoui64(v, &(stage.QuotaValue))
		case SStage:
			convertAtoui(v, &(stage.Stage))
		}
	}

	stage.Valid = proto.Bool(true)

	op := &battery.ResOpItem{
		Optype: optype,
	}

	stage.Createdate = proto.String(xyutil.CurTimeStr())
	op.Stage = stage

	*ops = append(*ops, op)
}
