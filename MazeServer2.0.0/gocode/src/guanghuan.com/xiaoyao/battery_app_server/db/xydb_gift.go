package batterydb

import (
//	"fmt"
//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//xyerror "guanghuan.com/xiaoyao/superbman_server/error"
////	"log"
//xylog "guanghuan.com/xiaoyao/common/log"
// xyutil "guanghuan.com/xiaoyao/common/util"
)

//func (db *BatteryDB) AddGift(gift battery.Gift) (err error) {

//	err = db.AddData(DB_TABLE_GIFT, gift)
//	if err != nil {
//		xylog.Error("[DB] add new gift failed: %v", err)
//	}
//	return
//}

//func (db *BatteryDB) GetGift(gift_gid string) (gift battery.Gift, err error) {
//	var qstr interface{}
//	qstr = bson.M{"giftid": gift_gid}

//	err = db.GetOneData(DB_TABLE_GIFT, qstr, &gift)
//	return
//}
//func (db *BatteryDB) UpdateGift(gift battery.Gift) (err error) {
//	gift_gid := gift.GetGiftId()

//	var qstr interface{}
//	qstr = bson.M{"giftid": gift_gid}
//	err = db.UpdateData(DB_TABLE_GIFT, qstr, gift)
//	return
//}

//func (db *BatteryDB) ConfirmGiftFinished(gift battery.Gift) (err error) {
//	if gift.GiftId == nil {
//		return xyerror.DBErrFailedDueToClientError
//	}
//	gift_id := gift.GetGiftId()

//	c1 := db.OpenTable(DB_TABLE_GIFT_DONE)
//	defer c1.Close()

//	c2 := c1.OpenATable(DB_TABLE_GIFT)
//	defer c2.Close()

//	err = c1.Insert(gift)
//	if err == xyerror.DBErrOK {
//		var qstr interface{}
//		qstr = bson.M{"giftid": gift_id}
//		err = c2.Remove(qstr)
//		if err != xyerror.DBErrOK {
//			xylog.Warning("WARNING: gift(%s) was confirmed not removed from original collection", gift_id)
//		}
//	}
//	return
//}

//func (db *BatteryDB) GetAllGifts(from_gid string, to_gid string, op_type battery.Gift_OpType, max_count int) (gifts []*battery.Gift, err error) {

//	if from_gid == "" && to_gid == "" {
//		xylog.Warning("WARNING: either from_id or to_id should be set")
//		err = xyerror.DBErrFailedDueToClientError
//		return
//	}

//	c := db.OpenTable(DB_TABLE_GIFT)
//	defer c.Close()

//	var qstr interface{}
//	if from_gid != "" && to_gid != "" {
//		qstr = bson.M{"optype": op_type, "fromid": from_gid, "toid": to_gid}
//	} else if from_gid != "" {
//		qstr = bson.M{"optype": op_type, "fromid": from_gid}
//	} else if to_gid != "" {
//		qstr = bson.M{"optype": op_type, "toid": to_gid}
//	}
//	xylog.Debug("query: %v", qstr)
//	var query *mgo.Query
//	if max_count > 0 {
//		query = c.Find(qstr).Limit(max_count)
//	} else {
//		query = c.Find(qstr)
//	}
//	err = query.All(&gifts)

//	return
//}

//func (db *BatteryDB) GetAllGiftsByUid(uid string, old_gift_first bool, max_count int) (gifts []*battery.Gift, err error) {

//	c := db.OpenTable(DB_TABLE_GIFT)
//	defer c.Close()

//	var qstr interface{} = db.queryAllGiftsCondition(uid)
//	xylog.Debug("query: %v", qstr)

//	query := c.Find(qstr)
//	if max_count > 0 {
//		query = query.Limit(max_count)
//	}
//	if old_gift_first {
//		// 旧的在前
//		query = query.Sort("createdate")
//	} else {
//		// 新的在前
//		query = query.Sort("-createdate")
//	}
//	err = query.All(&gifts)

//	err = xyerror.DBError(err)
//	return
//}

//func (db *BatteryDB) GetGiftsCountByUid(uid string) (count int, err error) {
//	var condition interface{} = db.queryAllGiftsCondition(uid)
//	xylog.Debug("query: %v", condition)

//	count, err = db.GetRecordCount(DB_TABLE_GIFT, condition)

//	err = xyerror.DBError(err)
//	return
//}

//// 构造与某id相关的请求的查询条件
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
//	err = db.GetOneData(DB_TABLE_GIFT, condition, &gift)
//	if err != nil {
//		// 尝试从已经确认的表里查询
//		err = db.GetOneData(DB_TABLE_GIFT_DONE, condition, &gift)
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
//	err = db.GetFirstData(DB_TABLE_GIFT_LOG, condition, sort, &gift_log)
//	return
//}
