package batterydb

import (
	//"fmt"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) AddAnnouncementConfig(a *battery.DBAnnouncementConfig) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_ANNOUNCEMENT_CONFIG, a)
	return
}

//清空公告配置表
func (db *BatteryDB) RemoveAllAnnouncementConfig() (err error) {
	err = db.RemoveAllData(xybusiness.DB_TABLE_ANNOUNCEMENT_CONFIG, bson.M{})
	return
}
