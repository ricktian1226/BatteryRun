package db

import (
	//xydb "guanghuan.com/xiaoyao/common/db"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

type BatteryDB struct {
	xybusiness.XYBusinessDB
}

const (
	BATTERYAPPID = "10004" // 游戏后台id
)

func NewBatteryDB(dburl, dbname string) (db *BatteryDB) {
	db = &BatteryDB{
		XYBusinessDB: *xybusiness.NewXYBusinessDB(dburl, dbname),
	}
	return
}
