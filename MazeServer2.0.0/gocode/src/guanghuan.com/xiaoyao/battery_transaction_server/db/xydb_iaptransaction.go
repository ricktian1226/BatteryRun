package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) AddIapTransaction(t *battery.IapTransaction) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_IAPTRANSACTION, *t)
	return
}

func (db *BatteryDB) IsTransactionExsit(tid string) (bool, error) {
	return db.IsRecordExistingWithField(xybusiness.DB_TABLE_IAPTRANSACTION, "transactionid", tid, mgo.Strong)
}

//查询匹配tids的所有事务信息
// tids []string 事务id列表
// transactions *[]*battery.IapTransaction 事务列表
func (db *BatteryDB) QueryIapTransactions(tids []string, transactions *[]*battery.IapTransaction) (err error) {
	condition := bson.M{"transactionid": bson.M{"$in": tids}}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetAllData(xybusiness.DB_TABLE_IAPTRANSACTION, condition, selector, 0, transactions, mgo.Strong)
	return
}

//查询玩家iap商品购买次数
// uid string 玩家标识
// iapGoodId string iap商品标识
func (db *BatteryDB) GetIapShoppingCount(uid, iapGoodId string) (count int, err error) {
	condition := bson.M{"uid": uid, "itemid": iapGoodId}
	count, err = db.GetRecordCount(xybusiness.DB_TABLE_IAPTRANSACTION, condition, mgo.Strong)
	return
}
