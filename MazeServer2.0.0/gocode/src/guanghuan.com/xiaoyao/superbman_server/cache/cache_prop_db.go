// cache_prop_db
package xybusinesscache

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    //"guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (ldb *CacheDB) LoadProps(props *[]battery.Prop) (err error) {

    //xylog.Debug("(ldb *CacheDB) LoadProps db : %v", ldb.db)

    c := ldb.dbCommon.OpenTable(xybusiness.DB_TABLE_PROP, mgo.Strong)
    defer c.Close()

    queryStr := bson.M{"valid": true}
    query := c.Find(queryStr)
    err = query.All(props)
    err = xyerror.DBError(err)

    return
}

// 加载登录礼包配置信息
func (cdb *CacheDB) LoadNewAccountProps(dbNewAccountProps *[]*battery.DBNewAccountProp) (err error) {
    condition := bson.M{}
    selector := bson.M{}
    return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_NEWACCOUNTPROP_CONFIG, condition, selector, 0, dbNewAccountProps, mgo.Strong)
}

// 加载分享礼包信息
func (db *CacheDB) LoadSharedActivity(dbSharedActivity *[]*battery.DBShareActivity) (err error) {
    condition := bson.M{}
    selector := bson.M{}
    return db.dbCommon.GetAllData(xybusiness.DB_TABLE_SHARE_ACTIVITY, condition, selector, 0, dbSharedActivity, mgo.Strong)
}

func (db *CacheDB) LoadSharedAwards(dbSharedAwards *[]*battery.DBShareAward) (err error) {
    condition := bson.M{}
    selector := bson.M{}
    return db.dbCommon.GetAllData(xybusiness.DB_TABLE_SHARE_AWARDS, condition, selector, 0, dbSharedAwards, mgo.Strong)
}
