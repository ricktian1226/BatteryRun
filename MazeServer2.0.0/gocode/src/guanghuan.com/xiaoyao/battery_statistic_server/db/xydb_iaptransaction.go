package batterydb

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//查询未上报成功的内购上报信息
func (db *BatteryDB) QueryUndoneIapStatistics() (statistics []*battery.DBIapStatistic, err error) {
    statistics = make([]*battery.DBIapStatistic, 0)
    condition := bson.M{"done": false}
    selector := bson.M{"transactionid": 1, "detail": 1, "retrycount": 1, "retrytimestamp": 1}
    err = db.GetAllData(xybusiness.DB_TABLE_IAPSTATISTIC, condition, selector, 0, &statistics, mgo.Strong)
    return
}

//刷新内购上报信息
func (db *BatteryDB) UpsertIapStatistic(statistic *battery.DBIapStatistic) (err error) {
    return db.UpsertData(xybusiness.DB_TABLE_IAPSTATISTIC, bson.M{"transactionid": statistic.GetTransactionId()}, statistic)
}
