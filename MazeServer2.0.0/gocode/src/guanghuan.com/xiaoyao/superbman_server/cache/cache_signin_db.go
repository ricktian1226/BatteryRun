// cache_signin_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (cdb *CacheDB) LoadSignInActivitys(activitys *[]*battery.DBSignInActivity) (err error) {
	c := cdb.dbCommon.OpenTable(xybusiness.DB_TABLE_SIGN_IN_ACTIVITY, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"valid": true}
	query := c.Find(queryStr)
	err = query.All(activitys)
	err = xyerror.DBError(err)

	return
}

func (cdb *CacheDB) LoadSignInItems(items *[]*battery.DBSignInItem) (err error) {
	c := cdb.dbCommon.OpenTable(xybusiness.DB_TABLE_SIGN_IN_ITEM, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"valid": true}
	query := c.Find(queryStr)
	err = query.All(items)
	err = xyerror.DBError(err)

	return
}
