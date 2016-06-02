package batterydb

import (
	"gopkg.in/mgo.v2/bson"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//增加提示配置信息
func (db *BatteryDB) AddTipConfig(a *battery.DBTip) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_TIP_CONFIG, a)
	return
}

//清空广告配置信息
func (db *BatteryDB) RemoveAllTipConfig() (err error) {
	err = db.RemoveAllData(xybusiness.DB_TABLE_TIP_CONFIG, bson.M{})
	return
}
