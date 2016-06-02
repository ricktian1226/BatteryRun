// cache_bgpropconfig_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (cdb *CacheDB) LoadBeforeGameRandomWeights(items *[]*battery.DBBeforeGameRandomGoodWeight) (err error) {
	condition := bson.M{"valid": true}
	selector := bson.M{"goodid": 1, "weight": 1}
	err = DefCacheDB.dbCommon.GetAllData(xybusiness.DB_TABLE_BEFOREGAME_RANDOM_WEIGHT, condition, selector, 0, items, mgo.Strong)
	return
}
