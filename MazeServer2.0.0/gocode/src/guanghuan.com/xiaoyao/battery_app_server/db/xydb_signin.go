// xydb_signin
package batterydb

import (
	//xylog "guanghuan.com/xiaoyao/common/log"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//增加签到活动信息
func (db *BatteryDB) AddSignInActivity(activity *battery.DBSignInActivity) (err error) {
	return db.AddData(xybusiness.DB_TABLE_SIGN_IN_ACTIVITY, activity)
}

//upsert 签到活动信息
func (db *BatteryDB) UpsertSignInActivity(activity *battery.DBSignInActivity) (err error) {
	condition := bson.M{"id": activity.GetId()}
	return db.UpsertData(xybusiness.DB_TABLE_SIGN_IN_ACTIVITY, condition, activity)
}

//删除所有的活动信息
func (db *BatteryDB) RemoveAllSignInActivitys() error {
	return db.RemoveAllData(xybusiness.DB_TABLE_SIGN_IN_ACTIVITY, bson.M{})
}

//增加签到奖励配置信息
func (db *BatteryDB) AddSignInItem(item *battery.DBSignInItem) (err error) {
	return db.AddData(xybusiness.DB_TABLE_SIGN_IN_ITEM, item)
}

//upsert 签到奖励配置信息
func (db *BatteryDB) UpsertSignInItem(item *battery.DBSignInItem) (err error) {
	condition := bson.M{"id": item.GetId(), "value": item.GetValue()}
	return db.UpsertData(xybusiness.DB_TABLE_SIGN_IN_ITEM, condition, item)
}

//删除所有的活动奖励信息
func (db *BatteryDB) RemoveAllSignInItems() error {
	return db.RemoveAllData(xybusiness.DB_TABLE_SIGN_IN_ITEM, bson.M{})
}
