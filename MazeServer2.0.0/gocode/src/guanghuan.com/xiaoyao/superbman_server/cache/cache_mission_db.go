// cache_mission_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (ldb *CacheDB) LoadMissionItems(missionItems *[]*battery.MissionItem) (err error) {

	c := ldb.dbCommon.OpenTable(xybusiness.DB_TABLE_MISSION, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"valid": true}
	query := c.Find(queryStr)
	err = query.All(missionItems)
	err = xyerror.DBError(err)

	return
}
