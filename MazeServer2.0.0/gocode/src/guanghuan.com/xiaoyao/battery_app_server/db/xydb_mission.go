// xydb_mission
package batterydb

import (
	//xylog "guanghuan.com/xiaoyao/common/log"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//增加任务配置信息
func (db *BatteryDB) AddMissionItem(item *battery.MissionItem) (err error) {
	return db.AddData(xybusiness.DB_TABLE_MISSION, item)
}

//Upsert任务配置信息
func (db *BatteryDB) UpsertMissionItem(item *battery.MissionItem) (err error) {
	condition := bson.M{"id": item.GetId()}
	return db.UpsertData(xybusiness.DB_TABLE_MISSION, condition, item)
}

//删除所有任务配置信息
func (db *BatteryDB) RemoveAllMissionItems() error {
	return db.RemoveAllData(xybusiness.DB_TABLE_MISSION, bson.M{})
}
