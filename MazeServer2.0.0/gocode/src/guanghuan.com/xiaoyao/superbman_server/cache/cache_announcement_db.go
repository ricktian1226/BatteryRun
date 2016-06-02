package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"time"
)

//加载当天的公告配置信息
func (cdb *CacheDB) LoadAnnouncements(dbAnnouncements *[]*battery.DBAnnouncementConfig) (err error) {
	now := time.Now().Unix()
	condition := bson.M{"begintime": bson.M{"$lte": now}, "endtime": bson.M{"$gte": now}}
	selector := bson.M{"id": 1, "title": 1, "message": 1, "description": 1}
	return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_ANNOUNCEMENT_CONFIG, condition, selector, 0, dbAnnouncements, mgo.Strong)
}
