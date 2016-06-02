// xydb_checkpoint
package batterydb

import (
    //xylog "guanghuan.com/xiaoyao/common/log"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//查询玩家区间记忆点信息
// uid string 玩家id
// beginCheckPointId uint32 起始记忆点
// endCheckPointId uint32 结束记忆点
// dbUserCheckPoints *[]*battery.DBUserCheckPoint 玩家记忆点信息的指针地址
func (db *BatteryDB) QueryUserCheckPoints(uid string, beginCheckPointId, endCheckPointId uint32, dbUserCheckPoints *[]*battery.DBUserCheckPoint) (err error) {
    condition := bson.M{"uid": uid, "checkpointid": bson.M{"$gte": beginCheckPointId, "$lte": endCheckPointId}}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
    err = db.GetAllData(xybusiness.DB_TABLE_USER_CHECK_POINT, condition, selector, 0, dbUserCheckPoints, mgo.Strong)
    err = xyerror.DBError(err)
    return
}

//查询玩家对应记忆点的详细信息
// uid string 玩家id
// checkPointId uint32 记忆点id
// dbUserCheckPoint *battery.DBUserCheckPoint 玩家记忆点信息的指针地址
func (db *BatteryDB) QueryUserCheckPointDetail(uid string, checkPointId uint32, dbUserCheckPoint *battery.DBUserCheckPoint) (err error) {
    condition := bson.M{"uid": uid, "checkpointid": checkPointId}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
    err = db.GetOneData(xybusiness.DB_TABLE_USER_CHECK_POINT, condition, selector, dbUserCheckPoint, mgo.Strong)
    err = xyerror.DBError(err)
    return
}

//获取玩家好友uid列表
// sids []string 好友sid列表
// uids *[]string 好友uid列表
//func (db *BatteryDB) QueryFriendsUid(sids []string, idSource battery.ID_SOURCE, uids *[]string) (err error) {
//	condition := bson.M{"source": idSource, "sid": bson.M{"$in": sids}}
//	tpids := make([]*battery.IDMap, 0)
//	err = db.GetAllData(DB_TABLE_TPID_MAP, condition, 0, &tpids)
//	if err == xyerror.ErrOK {
//		for _, tpid := range tpids {
//			*uids = append(*uids, tpid.GetGid())
//		}
//	}
//	return
//}

//查询好友记忆点排行榜
// checkPointId uint32 记忆点id
// uids []string 好友uid列表
// dbUserCheckPoint *[]*battery.DBUserCheckPoint 保存checkpoint信息列表
func (db *BatteryDB) QueryCheckPointFriendsRank(checkPointId uint32, uids []string, dbUserCheckPoints *[]*battery.DBUserCheckPoint) (err error) {
    collection := db.OpenTable(xybusiness.DB_TABLE_USER_CHECK_POINT, mgo.Monotonic)
    defer collection.Close()
    condition := bson.M{"checkpointid": checkPointId, "uid": bson.M{"$in": uids}}
    query := collection.Find(condition)
    sorter := "-score"
    query = query.Sort(sorter)
    err = query.All(dbUserCheckPoints)
    err = xyerror.DBError(err)
    return
}

//刷新玩家记忆点信息
// dbDetail *battery.DBUserCheckPoint 新的玩家记忆点信息
func (db *BatteryDB) UpsertUserCheckPoint(dbDetail *battery.DBUserCheckPoint) (err error) {
    condition := bson.M{"uid": dbDetail.GetUid(), "checkpointid": dbDetail.GetCheckPointId()}
    err = db.UpsertData(xybusiness.DB_TABLE_USER_CHECK_POINT, condition, *dbDetail)
    return
}

//增加玩家记忆点提交日志
func (db *BatteryDB) AddCheckPointLog(log *battery.CheckPointLog) (err error) {
    err = db.AddData(xybusiness.DB_TABLE_CHECKPOINT_LOG, *log)
    return
}

//查询玩家所有记忆点（不包含零号记忆点）的各种参数和
// uid string 玩家id
// dbUserCheckPoint *battery.DBUserCheckPoint 返回的记忆点
func (db *BatteryDB) QueryUserCheckPointsSum(uid string, dbUserCheckPoint *battery.DBUserCheckPoint) (err error) {

    collection := db.OpenTable(xybusiness.DB_TABLE_USER_CHECK_POINT, mgo.Strong)
    defer collection.Close()
    pipeLine := []bson.M{
        bson.M{"$match": bson.M{"uid": uid, "checkpointid": bson.M{"$gt": 0}}},
        bson.M{"$group": bson.M{
            "_id":              bson.M{"uid": "$uid"},
            "score":            bson.M{"$sum": "$score"},
            "charge":           bson.M{"$sum": "$charge"},
            "coin":             bson.M{"$sum": "$coin"},
            "collectionscount": bson.M{"$sum": "$collectionscount"},
        },
        },
        bson.M{"$project": bson.M{"uid": "$_id.uid", "score": "$score", "charge": "$charge", "coin": "$coin", "collectionscount": "$collectionscount"}},
    }
    err = collection.Pipe(pipeLine).One(dbUserCheckPoint)
    err = xyerror.DBError(err)
    return
}

//查询玩家0号记忆点的数据，0号记忆点保存的是玩家所有记忆点数据的和(零号就是关键，保存的是关键数据——纯属调侃)
// uid string 玩家id
// dbUserCheckPoint *battery.DBUserCheckPoint
func (db *BatteryDB) QueryUserNOZeroCheckPoints(uid string, dbUserCheckPoint *battery.DBUserCheckPoint) (err error) {
    return db.QueryUserCheckPointDetail(uid, uint32(0), dbUserCheckPoint)
}
