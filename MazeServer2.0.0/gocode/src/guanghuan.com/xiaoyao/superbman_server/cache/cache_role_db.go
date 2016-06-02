// cache_role_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//获取角色加成配置信息
func (cdb *CacheDB) LoadRoleConfig(items *[]*battery.DBRoleInfoConfig) (err error) {
	selector := bson.M{"_id": 0, "isvalid": 0, "xxx_unrecognized": 0}
	err = cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_ROLEINFO_CONFIG, bson.M{}, selector, 0, items, mgo.Strong)
	return
}
