// cache_mailconfig_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//加载邮件配置信息
func (cdb *CacheDB) LoadMailConfigs(items *[]*battery.DBMailInfoConfig) (err error) {
	c := cdb.dbCommon.OpenTable(xybusiness.DB_TABLE_MAILINFO_CONFIG, mgo.Strong)
	defer c.Close()
	query := c.Find(bson.M{})
	err = query.All(items)
	err = xyerror.DBError(err)

	return
}
