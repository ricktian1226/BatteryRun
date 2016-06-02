// cache_tip_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (cdb *CacheDB) LoadTips(tips *[]*battery.DBTip) (err error) {
	condition := bson.M{}
	selector := bson.M{"id": 1, "language": 1, "title": 1, "content": 1}
	return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_TIP_CONFIG, condition, selector, 0, tips, mgo.Strong)
}
