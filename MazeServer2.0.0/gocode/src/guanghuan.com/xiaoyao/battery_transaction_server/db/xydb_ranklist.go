package batterydb

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) GetRankList(uid string, ranklist *battery.RankLisk) (err error) {
    condition := bson.M{"uid": uid}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
    err = db.GetOneData(xybusiness.DB_TABLE_RANKLIST, condition, selector, ranklist, mgo.Monotonic)
    return
}

func (db *BatteryDB) UpsertRankList(uid string, ranklist *battery.RankLisk) (err error) {
    condition := bson.M{"uid": uid}
    err = db.UpsertData(xybusiness.DB_TABLE_RANKLIST, condition, ranklist)
    return
}
