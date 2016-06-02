package batterydb

import (
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	//"guanghuan.com/xiaoyao/common/apn"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

type BatteryDB struct {
	xybusiness.XYBusinessDB
}

func NewBatteryDB(dburl, dbname string) (db *BatteryDB) {
	db = &BatteryDB{
		XYBusinessDB: *xybusiness.NewXYBusinessDB(dburl, dbname),
	}
	return
}
