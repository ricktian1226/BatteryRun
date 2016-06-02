// cache_checkpoint_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//加载前几名玩家记忆点数据
// top uint32 前几名
// checkPointId uint32 记忆点id
// checkPoints *[]*battery.DBUserCheckPoint 玩家记忆点数据列表
func (cdb *CacheDB) LoadCheckPoint(top int, checkPointId uint32, checkPoints *[]*battery.DBUserCheckPoint, platform battery.PLATFORM_TYPE) (err error) {
	collection := cdb.dbCommon.OpenTable(xybusiness.DB_TABLE_USER_CHECK_POINT, mgo.Monotonic)
	defer collection.Close()
	condition := bson.M{"checkpointid": checkPointId, "platformtype": platform}
	query := collection.Find(condition)
	sorter := "-score"
	query = query.Sort(sorter).Limit(top)
	err = query.All(checkPoints)
	err = xyerror.DBError(err)
	return
}

func (cdb *CacheDB) LoadRankList(top int, ranklist *[]*battery.RankLisk) (err error) {
	collection := cdb.dbCommon.OpenTable(xybusiness.DB_TABLE_RANKLIST, mgo.Monotonic)
	defer collection.Close()
	condition := bson.M{"platformtype": battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID}
	query := collection.Find(condition)
	sorter := "-score"
	query = query.Sort(sorter).Limit(top)
	err = query.All(ranklist)
	err = xyerror.DBError(err)
	return
}
