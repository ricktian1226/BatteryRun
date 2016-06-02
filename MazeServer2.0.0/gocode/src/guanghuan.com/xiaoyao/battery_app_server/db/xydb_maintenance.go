package batterydb

import (
	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//upsert封号玩家信息
// 如果玩家封号信息已存在，会覆盖原来的解封时间戳
// uid string 玩家标识
// timestamp int64 解封时间戳
func (db *BatteryDB) UpsertBannedUser(uid string, timestamp int64) error {
	banned := &battery.DBUserBlackListItem{
		Uid:       proto.String(uid),
		Timestamp: proto.Int64(timestamp),
	}
	xylog.DebugNoId("UpsertBannedUser(%v)", banned)
	return db.UpsertData(xybusiness.DB_TABLE_BANNED_USER, bson.M{"uid": uid}, banned)
}

//删除封号玩家信息
// uid string 玩家标识
func (db *BatteryDB) RemoveBannedUser(uid string) error {
	return db.RemoveAllData(xybusiness.DB_TABLE_BANNED_USER, bson.M{"uid": uid})
}

//查询封号玩家信息
// now int64 当前时间戳
func (db *BatteryDB) QueryBannedUsers(now int64) (uids []string, err error) {
	uids = make([]string, 0)
	items := make([]*battery.DBUserBlackListItem, 0)
	condition := bson.M{"$or": []interface{}{bson.M{"timestamp": bson.M{"$lt": 0}}, bson.M{"timestamp": bson.M{"$gte": now}}}}
	selector := bson.M{"uid": 1}
	err = db.GetAllData(xybusiness.DB_TABLE_BANNED_USER, condition, selector, 0, &items, mgo.Strong)
	if err != xyerror.ErrOK {
		return
	}

	for _, item := range items {
		uids = append(uids, item.GetUid())
	}

	return
}
