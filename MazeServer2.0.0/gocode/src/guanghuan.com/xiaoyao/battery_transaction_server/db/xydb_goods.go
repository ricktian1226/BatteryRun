package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) AddShoppingLog(gl *battery.ShoppingLog) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_SHOPPING_LOG, *gl)
	return
}

func (db *BatteryDB) AddShoppingTransaction(st *battery.ShoppingTransaction) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_SHOPPING_TRANSACTION, *st)
	return
}

func (db *BatteryDB) AddReceipt(receipt battery.Receipt) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_RECEIPT, receipt)
	return
}

//func (db *BatteryDB) GetShoppingCount(uid string, goodsId uint64, gameId string) (count int, err error) {
//	condition := bson.M{"uid": uid, "gameid": gameId, "goodsid": goodsId}
//	count, err = db.GetRecordCount(xybusiness.DB_TABLE_SHOPPING_TRANSACTION, condition, mgo.Strong)
//	return
//}

func (db *BatteryDB) GetShoppingCount(condition bson.M) (count int, err error) {
	//condition := bson.M{"uid": uid, "gameid": gameId, "goodsid": goodsId}
	count, err = db.GetRecordCount(xybusiness.DB_TABLE_SHOPPING_TRANSACTION, condition, mgo.Strong)
	return
}
