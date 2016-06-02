// xydb_mission
package batterydb

import (
	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//查询玩家某类任务信息
// uid string 玩家id
// missionType battery.MissionType 任务类型
// userMissionStat battery.UserMissionState 任务状态
//return:
// userMissions []*battery.UserMission 玩家任务列表
func (db *BatteryDB) QueryUserMission(uid string, missionType battery.MissionType, userMissionStat battery.UserMissionState, userMissions []*battery.UserMission) (err error) {
	condition := bson.M{"uid": uid, "type": missionType}
	if battery.UserMissionState_UserMissionState_Unkown != userMissionStat {
		condition["state"] = userMissionStat
	}
	selector := bson.M{"_id":0,"xxx_unrecognized" :0}
	err = db.GetAllData(xybusiness.DB_TABLE_USER_MISSION, condition, selector, 0, &userMissions, mgo.Strong)
	return
}

//查询玩家某类任务信息
// uid string 玩家id
// missionId uint64 任务id
//return:
// userMission *battery.DBUserMission 任务详细信息
func (db *BatteryDB) QueryUserMissionDetail(uid string, missionId uint64, userMission *battery.DBUserMission) (err error) {
	condition := bson.M{"uid": uid, "mid": missionId}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_USER_MISSION, condition, selector, userMission, mgo.Strong)
	return
}

//根据类型、状态、时间戳查询玩家任务信息
// uid string 玩家id
// missionTypes []battery.MissionType 任务类型
// missionStates []battery.UserMissionState 任务状态
// now int64 当前时间戳（只有日常任务查询的时候需要根据时间查询）
// userMissions *[]*battery.DBUserMission 玩家任务信息列表
func (db *BatteryDB) QueryUserMissionByTypesAndStatesAndTimestamp(uid string, missionTypes []battery.MissionType, missionStates []battery.UserMissionState, now int64, userMissions *[]*battery.DBUserMission) (err error) {
	condition := bson.M{"uid": uid}

	if len(missionStates) > 1 {
		condition["state"] = bson.M{"$in": missionStates}
	} else if len(missionStates) == 1 {
		condition["state"] = missionStates[0]
	}

	if len(missionTypes) > 1 {
		condition["type"] = bson.M{"$in": missionTypes}
	} else if len(missionTypes) == 1 {
		condition["type"] = missionTypes[0]
	}

	//如果设置了查找时间，则需要添加过滤条件
	if now != 0 {
		beginTimestamp, endTimestamp := xyutil.TodayTimeRange(now)
		condition["timestamps.0.timestamp"] = bson.M{"$gte": beginTimestamp, "$lte": endTimestamp}
	}

	selector := bson.M{"_id":0,"xxx_unrecognized" :0}

	err = db.GetAllData(xybusiness.DB_TABLE_USER_MISSION, condition, selector, 0, userMissions, mgo.Strong)

	return
}

//刷新玩家任务状态
// uid string 玩家id
// missions []uint64 任务列表
// state battery.UserMissionState 任务状态信息
// now int64 当前时间戳
func (db *BatteryDB) UpdateUserMissionState(uid string, missions []uint64, state battery.UserMissionState, now int64) (err error) {
	condition := bson.M{}
	condition["uid"] = uid
	condition["mid"] = bson.M{"$in": missions}
	fields := bson.M{}
	pushs := bson.M{}
	fields["state"] = state
	timestamp := &battery.UserMissionTimestamp{
		State:     &state,
		Timestamp: proto.Int64(now),
	}
	pushs["timestamps"] = timestamp
	err = db.UpdateMultipleFieldsWithPush(xybusiness.DB_TABLE_USER_MISSION, condition, fields, pushs, true)
	return
}

//增加玩家任务信息
// userMission *battery.DBUserMission 任务信息指针
func (db *BatteryDB) AddUserMission(userMission *battery.DBUserMission) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_USER_MISSION, userMission)
	return
}

//刷新玩家任务信息
// userMission *battery.DBUserMission 任务信息指针
func (db *BatteryDB) UpdateUserMission(userMission *battery.DBUserMission) (err error) {
	uid := userMission.GetUid()
	mid := userMission.GetMid()

	if uid == "" || mid == 0 {
		err = xyerror.ErrBadInputData
		return
	}

	condition := bson.M{"uid": uid, "mid": mid}
	fields := bson.M{}

	if nil != userMission.QuotaStats && len(userMission.QuotaStats) > 0 {
		fields["quotastats"] = userMission.QuotaStats
	}

	if nil != userMission.State {
		fields["state"] = userMission.State
	}

	if nil != userMission.Timestamps && len(userMission.Timestamps) > 0 {
		fields["timestamps"] = userMission.Timestamps
	}

	if 0 != userMission.GetTimestamp() {
		fields["timestamp"] = userMission.GetTimestamp()
	}

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_USER_MISSION, condition, fields, false)

	return
}

//设置玩家任务状态为过期
// uid string 玩家标识
// missionType battery.MissionType 任务类型
func (db *BatteryDB) DeactivateUserMissions(uid string, missionType battery.MissionType) (err error) {
	condition := bson.M{"uid": uid, "type": missionType, "state": battery.UserMissionState_UserMissionState_Active}
	return db.UpdateOneField(xybusiness.DB_TABLE_USER_MISSION, condition, "state", battery.UserMissionState_UserMissionState_Expired, true)
}

//查询玩家doneCollected任务状态
// uid string 玩家id
// missionType battery.MissionType 任务类型
func (db *BatteryDB) QueryUserDoneCollectedMission(uid string, missionType battery.MissionType) (userDoneCollectedMission *battery.DBUserDoneCollectedMission, err error) {
	userDoneCollectedMission = &battery.DBUserDoneCollectedMission{}
	condition := bson.M{"uid": uid, "type": missionType}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_USER_DONECOLLECTED_MISSION, condition, selector, userDoneCollectedMission, mgo.Strong)
	return
}

//修改玩家doneCollected任务信息
// uid string 玩家id
// missionType battery.MissionType 任务类型
func (db *BatteryDB) UpsertUserDoneCollectedMission(userDoneCollectedMission *battery.DBUserDoneCollectedMission) error {
	condition := bson.M{"uid": userDoneCollectedMission.GetUid(), "type": userDoneCollectedMission.GetType()}
	return db.UpsertData(xybusiness.DB_TABLE_USER_DONECOLLECTED_MISSION, condition, userDoneCollectedMission)
}
