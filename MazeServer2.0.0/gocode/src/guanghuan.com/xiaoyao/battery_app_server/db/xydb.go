package batterydb

import (
	//xydb "guanghuan.com/xiaoyao/common/db"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

type BatteryDB struct {
	//xydb.XYDB
	xybusiness.XYBusinessDB
}

func NewBatteryDB(dburl, dbname string) (db *BatteryDB) {
	db = &BatteryDB{
		XYBusinessDB: *xybusiness.NewXYBusinessDB(dburl, dbname),
	}
	return
}
