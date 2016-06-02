package batterydb

import (
    "gopkg.in/mgo.v2/bson"

    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) AddShareActivityConfig(a *battery.DBShareActivity) (err error) {
    xylog.DebugNoId("Add ShareActivityConfig %v", a)
    err = db.AddData(xybusiness.DB_TABLE_SHARE_ACTIVITY, a)
    return
}

func (db *BatteryDB) AddShareAwardsConfig(a *battery.DBShareAward) (err error) {
    xylog.DebugNoId("Add ShareAwardsConfig %v", a)
    err = db.AddData(xybusiness.DB_TABLE_SHARE_AWARDS, a)
    return
}
func (db *BatteryDB) RemoveAllShareActivityConfig() (err error) {
    err = db.RemoveAllData(xybusiness.DB_TABLE_SHARE_ACTIVITY, bson.M{})
    return
}

func (db *BatteryDB) RemoveAllShareAwardsConfig() (err error) {
    err = db.RemoveAllData(xybusiness.DB_TABLE_SHARE_AWARDS, bson.M{})
    return
}
