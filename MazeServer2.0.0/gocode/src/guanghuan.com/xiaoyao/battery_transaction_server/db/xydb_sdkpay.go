package batterydb

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    //    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) AddOrder(i *battery.SDKOrderInfo) (err error) {
    err = db.AddData(xybusiness.DB_TABLE_SDKORDER, i)
    return
}

func (db *BatteryDB) GetSDKOrder(uid, orderid string, order *battery.SDKOrderInfo) (err error) {
    condition := bson.M{"uid": uid, "orderid": orderid}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
    err = db.GetOneData(xybusiness.DB_TABLE_SDKORDER, condition, selector, order, mgo.Strong)
    return
}

func (db *BatteryDB) GetAllUnfinishedSDKOrder(uid string, orders *[]*battery.SDKOrderInfo) (err error) {
    condition := bson.M{"uid": uid, "state": 1}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
    err = db.GetAllData(xybusiness.DB_TABLE_SDKORDER, condition, selector, 0, orders, mgo.Strong)
    return
}

func (db *BatteryDB) UpdateSDKOrder(uid, orderid string, order *battery.SDKOrderInfo) (err error) {
    condition := bson.M{"uid": uid, "orderid": orderid}
    err = db.UpsertData(xybusiness.DB_TABLE_SDKORDER, condition, order)
    return
}

func (db *BatteryDB) IsOrderIdExist(orderId string) (bool, error) {
    return db.IsRecordExistingWithField(xybusiness.DB_TABLE_SDKORDER, "orderid", orderId, mgo.Strong)
}
