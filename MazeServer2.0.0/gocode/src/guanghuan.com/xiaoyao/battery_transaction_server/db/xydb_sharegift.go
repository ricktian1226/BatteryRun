package batterydb

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) QueryAllShareInfo(uid string, userSharedInfo *[]*battery.DBUserShareInfo) (err error) {

    condition := bson.M{"uid": uid}
    selector := bson.M{"_id": 0}
    err = db.GetAllData(xybusiness.DB_TABLE_USERSHARED_RECORD, condition, selector, 0, userSharedInfo, mgo.Strong)
    return
}

func (db *BatteryDB) UpsertUserInfo(uid string, id uint32, userShareInfo *battery.DBUserShareInfo) (err error) {
    condition := bson.M{"uid": uid, "id": id}
    err = db.UpsertData(xybusiness.DB_TABLE_USERSHARED_RECORD, condition, userShareInfo)
    return
}

func (db *BatteryDB) QueryShareInfo(uid string, id uint32, userShareInfo *battery.DBUserShareInfo) (err error) {
    condition := bson.M{"uid": uid, "id": id}
    selector := bson.M{"_id": 0}
    err = db.GetOneData(xybusiness.DB_TABLE_USERSHARED_RECORD, condition, selector, userShareInfo, mgo.Strong)
    return
}
func (db *BatteryDB) AddUserInfo(userShareInfo *battery.DBUserShareInfo) (err error) {

    return db.AddData(xybusiness.DB_TABLE_USERSHARED_RECORD, userShareInfo)

}
