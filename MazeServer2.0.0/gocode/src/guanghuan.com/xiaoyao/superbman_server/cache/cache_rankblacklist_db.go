package xybusinesscache

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *CacheDB) LoadBlacklist(list *[]BlackList) (err error) {
    c := db.dbCommon.OpenTable(xybusiness.DB_TABLE_BLACKLIST, mgo.Strong)
    defer c.Close()
    queryStr := bson.M{}
    query := c.Find(queryStr)
    err = query.All(list)
    err = xyerror.DBError(err)
    return
}
