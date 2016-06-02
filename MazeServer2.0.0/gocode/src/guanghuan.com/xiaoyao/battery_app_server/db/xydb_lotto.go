// xydb_lotto
package batterydb

import (
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//"time"
)

//增加抽奖槽位信息
func (db *BatteryDB) AddSlotItem(si *battery.LottoSlotItem) error {
	return db.AddData(xybusiness.DB_TABLE_LOTTO_SLOTITEMS, *si)
}

//upsert抽奖槽位信息
func (db *BatteryDB) UpsertSlotItem(si *battery.LottoSlotItem) (err error) {
	condition := bson.M{"slotid": si.GetSlotid(), "propid": si.GetPropid(), "datype": si.GetDatype()}
	return db.UpsertData(xybusiness.DB_TABLE_LOTTO_SLOTITEMS, condition, *si)
}

//删除抽奖槽位信息
func (db *BatteryDB) DelSlotItem(si *battery.LottoSlotItem) (err error) {

	//使用的伪删除
	var condition interface{}
	condition = bson.M{"slotid": si.GetSlotid(), "propid": si.GetPropid()}

	fields := bson.M{}
	fields["valid"] = false

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_LOTTO_SLOTITEMS, condition, fields, true)
	if err != nil {
		xylog.Error(xylog.DefaultLogId, "[DB] fail to update slotitem: %v ", err)
	}
	return
}

//清空抽奖槽位信息
func (db *BatteryDB) RemoveAllSlotItems() error {
	return db.RemoveAllData(xybusiness.DB_TABLE_LOTTO_SLOTITEMS, bson.M{})
}

//修改槽位信息
func (db *BatteryDB) ModSlotItem(si *battery.LottoSlotItem) (err error) {
	fields := bson.M{}

	if nil != si.Weight {
		fields["weight"] = si.GetWeight()
	}

	if nil != si.Valid {
		fields["valid"] = si.GetValid()
	}

	condition := bson.M{"slotid": si.GetSlotid(), "propid": si.GetPropid()}

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_PROP, condition, fields, true)
	if err != nil {
		xylog.Error(xylog.DefaultLogId, "[DB] fail to update prop: %v ", err)
	}
	return
}

//增加系统抽奖权重信息
func (db *BatteryDB) AddWeight(w *battery.LottoWeight) error {
	return db.AddData(xybusiness.DB_TABLE_LOTTO_WEIGHT, w)
}

//upsert系统抽奖权重信息
func (db *BatteryDB) UpsertWeight(w *battery.LottoWeight) error {
	condition := bson.M{}
	if nil != w.Beginvalue {
		condition["beginvalue"] = w.GetBeginvalue()
	}
	if nil != w.Endvalue {
		condition["endvalue"] = w.GetEndvalue()
	}
	return db.UpsertData(xybusiness.DB_TABLE_LOTTO_WEIGHT, condition, w)
}

//删除系统抽奖信息
func (db *BatteryDB) DelWeight(w *battery.LottoWeight) (err error) {

	//使用的伪删除
	var condition interface{}
	condition = bson.M{"beginvalue": w.GetBeginvalue(), "endvalue": w.GetBeginvalue()}

	fields := bson.M{}
	fields["valid"] = false

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_LOTTO_WEIGHT, condition, fields, true)
	if err != nil {
		xylog.Error(xylog.DefaultLogId, "[DB] fail to update slotitem: %v ", err)
	}
	return
}

//清空系统抽奖权重信息
func (db *BatteryDB) RemoveAllWeights() error {
	return db.RemoveAllData(xybusiness.DB_TABLE_LOTTO_WEIGHT, bson.M{})
}

//修改系统抽奖信息
func (db *BatteryDB) ModWeight(w *battery.LottoWeight) (err error) {
	fields := bson.M{}

	if nil != w.Weightlist {
		fields["weightlist"] = w.GetWeightlist()
	}

	if nil != w.Valid {
		fields["valid"] = w.GetValid()
	}

	condition := bson.M{"beginvalue": w.GetBeginvalue(), "endvalue": w.GetEndvalue()}

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_LOTTO_WEIGHT, condition, fields, true)
	if err != nil {
		xylog.Error(xylog.DefaultLogId, "[DB] fail to update prop: %v ", err)
	}

	return
}

//增加系统抽奖权重信息
func (db *BatteryDB) AddSerialNumSlot(l *battery.LottoSerialNumSlot) error {
	xylog.Debug(xylog.DefaultLogId, "AddSerialNumSlot : %v", l)
	return db.AddData(xybusiness.DB_TABLE_LOTTO_SERIALNUM_SLOT, l)
}

//upsert系统抽奖权重信息
//func (db *BatteryDB) UpsertWeight(w *battery.LottoWeight) error {
//	condition := bson.M{}
//	if nil != w.Beginvalue {
//		condition["beginvalue"] = w.GetBeginvalue()
//	}
//	if nil != w.Endvalue {
//		condition["endvalue"] = w.GetEndvalue()
//	}
//	return db.UpsertData(xybusiness.DB_TABLE_LOTTO_WEIGHT, condition, w)
//}

//删除系统抽奖信息
//func (db *BatteryDB) DelWeight(w *battery.LottoWeight) (err error) {

//	//使用的伪删除
//	var condition interface{}
//	condition = bson.M{"beginvalue": w.GetBeginvalue(), "endvalue": w.GetBeginvalue()}

//	fields := bson.M{}
//	fields["valid"] = false

//	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_LOTTO_WEIGHT, condition, fields, true)
//	if err != nil {
//		xylog.Error("[DB] fail to update slotitem: %v ", err)
//	}
//	return
//}

//清空系统抽奖权重信息
func (db *BatteryDB) RemoveAllSerialNumSlots() error {
	return db.RemoveAllData(xybusiness.DB_TABLE_LOTTO_SERIALNUM_SLOT, bson.M{})
}

//增加游戏后抽奖阶段信息
//func (db *BatteryDB) AddStage(s *battery.LottoStageItem) (err error) {
//	return db.AddData(xybusiness.DB_TABLE_LOTTO_STAGE, s)
//}

//upsert游戏后抽奖阶段信息
//func (db *BatteryDB) UpsertStage(s *battery.LottoStageItem) (err error) {
//	condition := bson.M{"quotaid": s.GetQuotaId(), "quotavalue": s.GetQuotaValue()}
//	return db.UpsertData(xybusiness.DB_TABLE_LOTTO_STAGE, condition, s)
//}

//清空游戏后抽奖阶段信息
//func (db *BatteryDB) RemoveAllStages() error {
//	return db.RemoveAllData(xybusiness.DB_TABLE_LOTTO_STAGE, bson.M{})
//}
