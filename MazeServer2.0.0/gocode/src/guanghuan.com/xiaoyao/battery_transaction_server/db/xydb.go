package batterydb

import (
	//xydb "guanghuan.com/xiaoyao/common/db"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//--------- Attention 测试代码，不要使用 ---------
var DefBatteryDB *BatteryDB

type empty struct{}
type Set map[interface{}]empty

func NewTestBatteryDB(dburl string, dbname string) {
	DefBatteryDB = NewBatteryDB(dburl, dbname)
	DefBatteryDB.OpenDB()
}

//---------------------------------------------

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
