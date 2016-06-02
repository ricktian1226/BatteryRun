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

//增加广告配置信息
func (db *BatteryDB) AddAdvertisementConfig(a *battery.Advertisement) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_ADVERTISEMENT_CONFIG, a)
	return
}

//清空广告配置信息
func (db *BatteryDB) RemoveAllAdvertisementConfig() (err error) {
	err = db.RemoveAllData(xybusiness.DB_TABLE_ADVERTISEMENT_CONFIG, bson.M{})
	return
}

//增加广告位配置信息
func (db *BatteryDB) AddAdvertisementSpaceConfig(a *battery.AdvertisementSpace) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_ADVERTISEMENTSPACE_CONFIG, a)
	return
}

//清空广告位配置信息
func (db *BatteryDB) RemoveAllAdvertisementSpaceConfig() (err error) {
	err = db.RemoveAllData(xybusiness.DB_TABLE_ADVERTISEMENTSPACE_CONFIG, bson.M{})
	return
}
