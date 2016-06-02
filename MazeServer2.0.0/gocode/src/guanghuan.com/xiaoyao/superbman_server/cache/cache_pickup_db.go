// cache_pickup_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (ldb *CacheDB) LoadPickUps(pickUps *[]*battery.DBPickUpItem) (err error) {

	//xylog.Debug("(ldb *CacheDB) LoadPickUpItems db : %v", ldb.db)

	c := ldb.dbCommon.OpenTable(xybusiness.DB_TABLE_PICKUP_WEIGHT, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"valid": true}
	query := c.Find(queryStr)
	err = query.All(pickUps)
	err = xyerror.DBError(err)

	return
}
