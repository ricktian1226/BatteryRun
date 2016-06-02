package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//加载广告配置信息
func (cdb *CacheDB) LoadAdvertisements(advertisements *[]*battery.Advertisement) (err error) {
	condition := bson.M{}
	selector := bson.M{"id": 1, "viewurl": 1, "materialurl": 1, "clickurl": 1}
	return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_ADVERTISEMENT_CONFIG, condition, selector, 0, advertisements, mgo.Strong)
}

//加载广告位配置信息
func (cdb *CacheDB) LoadAdvertisementSpaces(advertisementSpaces *[]*battery.AdvertisementSpace) (err error) {
	condition := bson.M{}
	selector := bson.M{"id": 1, "items": 1, "enable": 1, "flags": 1}
	return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_ADVERTISEMENTSPACE_CONFIG, condition, selector, 0, advertisementSpaces, mgo.Strong)
}
