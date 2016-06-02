// cache_rolebonus_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//获取角色加成配置信息
func (cdb *CacheDB) LoadRoleLevelBonus(items *[]*battery.DBRoleLevelBonusItem) (err error) {
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_ROLE_LEVEL_BONUS, bson.M{}, selector, 0, items, mgo.Strong)

	return
}
