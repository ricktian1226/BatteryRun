// xydb_prop
package batterydb

import (
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"gopkg.in/mgo.v2/bson"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// 增加道具信息
func (db *BatteryDB) AddProp(p *battery.Prop) error {
	return db.AddData(xybusiness.DB_TABLE_PROP, *p)
}

//增加收集物配置信息
func (db *BatteryDB) AddPickUp(p *battery.DBPickUpItem) error {
	return db.AddData(xybusiness.DB_TABLE_PICKUP_WEIGHT, *p)
}

//upsert道具信息
func (db *BatteryDB) UpsertProp(p *battery.Prop) error {
	condition := bson.M{"id": p.GetId()}
	return db.UpsertData(xybusiness.DB_TABLE_PROP, condition, *p)
}

//删除道具信息
func (db *BatteryDB) DelProp(p *battery.Prop) (err error) {
	var condition interface{}
	condition = bson.M{"id": p.GetId()}
	fields := bson.M{}
	fields["valid"] = false

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_PROP, condition, fields, true)
	if err != nil {
		xylog.Error(xylog.DefaultLogId, "[DB] fail to delete prop: %v ", err)
	}
	return
}

//清空道具配置表
func (db *BatteryDB) RemoveAllProps() (err error) {
	err = db.RemoveAllData(xybusiness.DB_TABLE_PROP, bson.M{})
	return
}

//清空收集物权重配置表
func (db *BatteryDB) RemoveAllPickUpWeightItems() (err error) {
	err = db.RemoveAllData(xybusiness.DB_TABLE_PICKUP_WEIGHT, bson.M{})
	return
}

//修改道具信息
func (db *BatteryDB) ModProp(p *battery.Prop) (err error) {
	fields := bson.M{}

	if nil != p.Type {
		fields["type"] = p.GetType()
	}

	if nil != p.Items {
		fields["items"] = p.GetItems()
	}

	if nil != p.Lottovalue {
		fields["value"] = p.GetLottovalue()
	}

	if nil != p.Valid {
		fields["valid"] = p.GetValid()
	}

	condition := bson.M{"id": p.GetId()}

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_PROP, condition, fields, true)
	if err != nil {
		xylog.Error(xylog.DefaultLogId, "[DB] fail to update prop: %v ", err)
	}
	return
}

//增加运营操作日志
func (db *BatteryDB) AddMaintenanceLog(l *battery.MaintenanceLog) (err error) {
	return db.AddData(xybusiness.DB_TABLE_MAINTENANCE_LOG, l)
}

//增加提示配置信息
func (db *BatteryDB) AddNewAccountPropConfig(a *battery.DBNewAccountProp) (err error) {
	xylog.DebugNoId("AddNewAccountPropConfig %v", a)
	err = db.AddData(xybusiness.DB_TABLE_NEWACCOUNTPROP_CONFIG, a)
	return
}

//清空广告配置信息
func (db *BatteryDB) RemoveAllNewAccountPropConfig() (err error) {
	err = db.RemoveAllData(xybusiness.DB_TABLE_NEWACCOUNTPROP_CONFIG, bson.M{})
	return
}
