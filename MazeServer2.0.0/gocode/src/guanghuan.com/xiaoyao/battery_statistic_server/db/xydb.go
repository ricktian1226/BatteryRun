package batterydb

import (
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
