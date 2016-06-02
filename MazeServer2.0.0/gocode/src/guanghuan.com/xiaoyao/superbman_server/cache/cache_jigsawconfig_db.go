// cache_jigsawconfig_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (cdb *CacheDB) LoadJigsawConfigs(items *[]*battery.JigsawConfig) (err error) {
	c := cdb.dbCommon.OpenTable(xybusiness.DB_TABLE_JIGSAW_CONFIG, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"jigsawid": bson.M{"$gte": 0}}
	query := c.Find(queryStr)
	err = query.All(items)
	err = xyerror.DBError(err)

	return
}
