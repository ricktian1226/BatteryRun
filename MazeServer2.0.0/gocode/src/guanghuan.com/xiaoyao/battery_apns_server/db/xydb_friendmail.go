// 玩家好友邮件数目相关推送
package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// QueryFriendMailCount 查询玩家好友邮件数目信息
func (db *BatteryDB) QueryFriendMailCount(limit int32) (friendMailCounts []*battery.DBFriendMailCount, err error) {
	friendMailCounts = make([]*battery.DBFriendMailCount, 0)
	selector := bson.M{"uid": 1}
	condition := bson.M{"count": bson.M{"$gte": limit}}
	err = db.GetAllData(xybusiness.DB_TABLE_FRIENDMAILCOUNT, condition, selector, 0, &friendMailCounts, mgo.Strong)
	return
}
