// xydb_signin
package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//查询玩家签到活动信息
// uid string 玩家id
//return:
// userActivitys *[]*battery.DBUserSignInActivity 玩家活动信息列表
func (db *BatteryDB) QueryUserSignInActivitys(uid string, userActivitys *[]*battery.DBUserSignInActivity) (err error) {
	condition := bson.M{"uid": uid}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetAllData(xybusiness.DB_TABLE_SIGN_IN_RECORD, condition, selector, 0, userActivitys, mgo.Strong)
	return
}

//查询玩家签到活动详细信息
// uid string 玩家id
// id uint64 签到活动id
//return:
// userActivity *battery.DBUserSignInActivity 玩家签到活动信息明细
func (db *BatteryDB) QueryUserSignInActivityDetail(uid string, id uint64, userActivity *battery.DBUserSignInActivity) (err error) {
	condition := bson.M{"uid": uid, "id": id}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_SIGN_IN_RECORD, condition, selector, userActivity, mgo.Strong)
	return
}

//增加玩家签到活动信息
// userActivity *battery.DBUserSignInActivity 签到活动信息明细
func (db *BatteryDB) AddUserSignInActivity(userActivity *battery.DBUserSignInActivity) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_SIGN_IN_RECORD, userActivity)
	return
}

//更新玩家签到活动信息
// uid string 玩家id
// userActivity *battery.DBUserSignInActivity 签到活动信息明细
func (db *BatteryDB) UpsertUserSignInActivity(uid string, userActivity *battery.DBUserSignInActivity) (err error) {
	condition := bson.M{"uid": uid}
	err = db.UpsertData(xybusiness.DB_TABLE_SIGN_IN_RECORD, condition, userActivity)
	return
}
