package batterydb

//import (
//	"gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	xylog "guanghuan.com/xiaoyao/common/log"
//	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
//	"guanghuan.com/xiaoyao/superbman_server/server"
//	"time"
//)

//获取公告最大的条数
//const MaxAnnouncementCount = 10

////公告合法性校验
//func (db *BatteryDB) CheckAnnouncement(a *battery.DBAnnouncement) bool {
//	if a != nil && 0 >= a.GetId() {
//		return false
//	}
//	return true
//}

//增加公告信息
//func (db *BatteryDB) AddAnnouncement(a *battery.DBAnnouncement) (err error) {

//	// 只检查Uid是否为空
//	if db.CheckAnnouncement(a) {
//		err = db.AddData(xybusiness.DB_TABLE_ANNOUNCEMENT, a)
//		if err != xyerror.ErrOK {
//			xylog.Error("[DB] fail to add new announcement: %v ", err)
//		}
//	} else {
//		xylog.Error("[DB] fail to add new announcement")
//		err = xyerror.ErrBadInputData
//	}

//	err = xyerror.DBError(err)
//	return err
//}

//增加公告信息（多条）
//func (db *BatteryDB) AddAnnouncements(as []interface{}) (err error) {

//	tbl := db.OpenTable(xybusiness.DB_TABLE_ANNOUNCEMENT, mgo.Strong)
//	defer tbl.Close()

//	err = tbl.Insert(as...)

//	err = xyerror.DBError(err)

//	return err
//}

//删除公告信息
//func (db *BatteryDB) DelAnnouncement(a *battery.DBAnnouncement) (err error) {
//	if db.CheckAnnouncement(a) {
//		var condition interface{}
//		condition = bson.M{"id": a.GetId()}
//		err = db.RemoveData(xybusiness.DB_TABLE_ANNOUNCEMENT, condition)
//		if err != xyerror.ErrOK {
//			xylog.Error("[DB] fail to delete announcement: %v ", err)
//		}
//	} else {
//		xylog.Error("[DB] fail to delete announcement")
//		err = xyerror.ErrBadInputData
//	}
//	return err
//}

//刷新公告信息
//func (db *BatteryDB) UpdateAnnouncement(a *battery.DBAnnouncement) (err error) {
//	if db.CheckAnnouncement(a) {

//		fields := bson.M{}

//		if nil != a.SubmitTime {
//			fields["submittime"] = a.GetSubmitTime()
//		}

//		if nil != a.BeginTime {
//			fields["begintime"] = a.GetBeginTime()
//		}

//		if nil != a.EndTime {
//			fields["endtime"] = a.GetEndTime()
//		}

//		if nil != a.Title {
//			fields["title"] = a.GetTitle()
//		}

//		if nil != a.Content {
//			fields["content"] = a.GetContent()
//		}

//		if nil != a.Uid {
//			fields["uid"] = a.GetUid()
//		}

//		condition := bson.M{"id": a.GetId()}

//		//err = db.UpdateDataWithField(DB_TABLE_ANNOUNCEMENT, "id", a.GetId(), a)
//		err = db.UpdateMultipleFields(DB_TABLE_ANNOUNCEMENT, condition, fields, true)
//		if err != nil {
//			xylog.Error("[DB] fail to update announcement: %v ", err)
//		}
//	} else {
//		xylog.Error("[DB] fail to update announcement")
//		err = xyerror.ErrBadInputData
//	}
//	return err
//}

//按照id查找公告信息
// a *battery.Announcement 查询条件，取id
//func (db *BatteryDB) QueryAnnouncementById(a *battery.DBAnnouncement, result *battery.DBAnnouncement) (err error) {
//	if db.CheckAnnouncement(a) {
//		var condition interface{}
//		condition = bson.M{"id": *a.Id}
//		err = db.GetOneData(xybusiness.DB_TABLE_ANNOUNCEMENT, condition, result, mgo.Strong)
//		xylog.Debug("QueryAnnouncementByTime err(%v), condition(%v), results(%v)", err, condition, result)
//		xylog.Debug("result : %v", result)
//		if err != nil {
//			xylog.Error("[DB] fail to QueryAnnouncementById: %v ", err)
//		}
//	} else {
//		xylog.Error("[DB] fail to QueryAnnouncementById, input id 0.")
//		err = xyerror.ErrBadInputData
//	}
//	return
//}

//根据当前的时间戳获取公告信息
//return:
// results []battery.Announcement 公告信息列表
// err error 操作错误
//func (db *BatteryDB) QueryAnnouncementByTime() (results []*battery.DBAnnouncement, err error) {

//	results = make([]*battery.DBAnnouncement, 0)

//	//获取当前时间戳
//	time := time.Now().Unix()

//	var condition interface{}
//	condition = bson.M{"begintime": bson.M{"$lte": time}, "endtime": bson.M{"$gte": time, "$gt": 111}, "state": battery.ANNOUNCEMENT_STATE_ANNOUNCEMENT_STATE_VALID}
//	err = db.GetAllData(xybusiness.DB_TABLE_ANNOUNCEMENT, condition, MaxAnnouncementCount, &results, mgo.Monotonic)
//	//if err != xyerror.ErrOK {
//	//	xylog.Error("[DB] fail to QueryAnnouncementByTime: %v ", err)
//	//}

//	//xylog.Debug("QueryAnnouncementByTime err(%v), condition(%v), results(%v)", err, condition, results)

//	return
//}

//根据多种类型的公告
// uid string 玩家id
// types []int64
//func (db *BatteryDB) QueryAnnouncementByTypes(uid string, types []int64) (announcements []*battery.DBAnnouncement, err error) {
//	condition := bson.M{"endtime": bson.M{"$in": types}, "uid": uid}
//	//var announcements []*battery.Announcement
//	err = db.GetAllData(xybusiness.DB_TABLE_ANNOUNCEMENT, condition, 0, &announcements, mgo.Strong)
//	return
//}

//根据类型删除玩家的公告信息
// uid string 玩家id
//return:
//  err error 返回错误
//func (db *BatteryDB) DeleteAnnouncementByTypes(uid string, types []int64) (err error) {
//	condition := bson.M{"endtime": bson.M{"$in": types}, "uid": uid}
//	err = db.RemoveAllData(xybusiness.DB_TABLE_ANNOUNCEMENT, condition)
//	return
//}

//按照时间段查询公告信息
//func (db *BatteryDB) QueryAnnouncementByTimeRange(a *battery.DBAnnouncement) (results []battery.DBAnnouncement, err error) {
//	begin, end := a.GetBeginTime(), a.GetEndTime()
//	if (0 == begin && 0 == end) || begin > end {
//		xylog.Error("[DB] fail to QueryAnnouncementByTime, input begin 0 end 0 or begin > end.")
//		err = xyerror.ErrBadInputData
//	} else {
//		condition := make(bson.M)
//		if 0 != begin {
//			condition["begintime"] = bson.M{"$gte": begin}
//		}
//		if 0 != end {
//			condition["endtime"] = bson.M{"$lte": end, "$gt": 111}
//		}
//		results = make([]battery.DBAnnouncement, 0)
//		//最多取10条
//		err = db.GetAllData(xybusiness.DB_TABLE_ANNOUNCEMENT, condition, MaxAnnouncementCount, &results, mgo.Strong)
//		if err != nil {
//			xylog.Error("[DB] fail to QueryAnnouncementByTime: %v ", err)
//		}
//		xylog.Debug("QueryAnnouncementByTime err(%v), condition(%v), results(%v)", err, condition, results)
//	}

//	return
//}

//生成公告id
//func (db *BatteryDB) AnnouncementID(timeNow *time.Time) (id uint64) {
//	id = uint64(timeNow.Year()*1e13) +
//		uint64(timeNow.Month()*1e11) +
//		uint64(timeNow.Day()*1e9) +
//		uint64(timeNow.Hour()*1e7) +
//		uint64(timeNow.Minute()*1e5) +
//		uint64(timeNow.Second()*1e3) +
//		uint64(timeNow.Nanosecond())/uint64(time.Millisecond)

//	return
//}
