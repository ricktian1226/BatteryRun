package batterydb

//import (
//	//	"fmt"
//	"gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
//	"guanghuan.com/xiaoyao/superbman_server/server"
//	//	"log"
//	xylog "guanghuan.com/xiaoyao/common/log"
//	// xyutil "guanghuan.com/xiaoyao/common/util"
//)

//func (db *BatteryDB) GetGiftsCountByUid(uid string) (count int, err error) {
//	var condition interface{} = db.queryAllGiftsCondition(uid)
//	xylog.Debug(uid, "query: %v", condition)

//	count, err = db.GetRecordCount(xybusiness.DB_TABLE_GIFT, condition, mgo.Strong)

//	err = xyerror.DBError(err)
//	return
//}

////// 构造与某id相关的请求的查询条件
//func (db *BatteryDB) queryAllGiftsCondition(uid string) (condition interface{}) {
//	ask_str := bson.M{"optype": battery.Gift_Op_Ask, "fromid": uid}
//	give_str := bson.M{"optype": battery.Gift_Op_Give, "toid": uid}

//	condition = bson.M{"$or": []interface{}{ask_str, give_str}}

//	return
//}

//// 查询最近一次 A 向 B 赠送的体力记录
//func (db *BatteryDB) GetLastGiftSentAfterDate(from_id string, to_id string, after_date int64) (gift battery.Gift, err error) {
//	var condition interface{}
//	if after_date > 0 {
//		condition = bson.M{"fromid": from_id, "toid": to_id,
//			"createdate": bson.M{"$gt": after_date}}
//	} else {
//		// 如果date不是一个有效时间，
//		condition = bson.M{"fromid": from_id, "toid": to_id}
//	}
//	xylog.Debug("query: %v", condition)
//	err = db.GetOneData(DB_TABLE_GIFT, condition, &gift, mgo.Strong)
//	if err != nil {
//		// 尝试从已经确认的表里查询
//		err = db.GetOneData(DB_TABLE_GIFT_DONE, condition, &gift, mgo.Strong)
//	}

//	return
//}

//func (db *BatteryDB) AddGiftLog(gl battery.GiftLog) (err error) {
//	err = db.AddData(DB_TABLE_GIFT_LOG, gl)
//	return err
//}
//func (db *BatteryDB) GetLastGiftOpLog(uid string, friend string, op_type battery.Gift_OpType, fail_reason int32) (gift_log battery.GiftLog, err error) {
//	condition := bson.M{"uid": uid, "friendid": friend, "optype": op_type, "failreason": fail_reason}
//	sort := "-opdate"
//	err = db.GetFirstData(DB_TABLE_GIFT_LOG, condition, sort, &gift_log, mgo.Strong)
//	return
//}
