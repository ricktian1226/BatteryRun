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

import "time"

// 推送类型定义
const (
	PushType_Power_Full = iota // 体力满通知
)

// 推送状态
const (
	PushState_WaitForPush = 0 // 等待推送
	PushState_Pushed      = 1 // 已经推送
)

// 向推送数据表中插入新的推送请求
func (db *BatteryDB) InsertNewPushNotification(device_id string, push_time int64, msg string, push_type int32, name string) error {
	// 删除相同类型的
	db.removeSameTypeNotification(device_id, push_type)

	// 插入新的通知
	var notification battery.NotificationForPush
	notification.DeviceId = proto.String(device_id)
	notification.Msg = proto.String(msg)
	notification.Type = proto.Int32(push_type)
	notification.PushTime = proto.Int64(push_time)
	notification.Timestamp = proto.Int64(time.Now().Unix())
	notification.Name = proto.String(name)
	notification.State = proto.Int32(PushState_WaitForPush)
	notification.NotificationId = proto.String(xyutil.NewId())

	//dbsession := db.session.Copy()
	//defer dbsession.Close()

	//c := dbsession.DB(DefaultDB).C(DB_TABLE_PUSH_Natice)
	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_Natice, mgo.Strong)
	defer c.Close()

	err := c.Insert(notification)

	err = xyerror.DBError(err)
	return err
}

// 移除同类型的
func (db *BatteryDB) removeSameTypeNotification(device_id string, push_type int32) error {
	//dbsession := db.session.Copy()
	//defer dbsession.Close()

	//c := dbsession.DB(DefaultDB).C(DB_TABLE_PUSH_Natice)
	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_Natice, mgo.Strong)
	defer c.Close()
	queryStr := bson.M{"type": push_type, "deviceid": device_id}
	err := c.Remove(queryStr)

	err = xyerror.DBError(err)
	return err
}

// 更新推送状态
func (db *BatteryDB) UpdateNotificationState(notice_id string, state int32) error {
	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_Natice, mgo.Strong)
	defer c.Close()
	queryStr := bson.M{"notificationid": notice_id}
	var query = c.Find(queryStr)

	count, err := query.Count()
	if err != nil {
		return err
	}

	if count > 0 {
		var notice battery.NotificationForPush
		err = query.One(&notice)
		if err != nil {
			return err
		}

		notice.State = proto.Int32(state)
		err = c.Update(queryStr, notice)
	}

	return err
}

// 根据设备ID和插入的时间删除对应的推送信息
func (db *BatteryDB) RemoveNotificationByDeviceIdAndTime(device_id string, insert_time int64) error {
	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_Natice, mgo.Strong)
	defer c.Close()
	queryStr := bson.M{"deviceid": device_id, "inserttime": insert_time}
	err := c.Remove(queryStr)

	err = xyerror.DBError(err)
	return err
}

// 插入推送通知操作记录
func (db *BatteryDB) InsertNewPushRecord(device_id string, push_time int64, insert_time int64, msg string, push_type int32, success bool, log string) error {
	var pushRecord battery.PushRecord
	pushRecord.DeviceId = proto.String(device_id)
	pushRecord.Msg = proto.String(msg)
	pushRecord.PushTime = proto.Int64(push_time)
	pushRecord.Timestamp = proto.Int64(insert_time)
	pushRecord.Type = proto.Int32(push_type)
	pushRecord.SendTime = proto.Int64(time.Now().Unix())
	pushRecord.Success = proto.Bool(success)
	pushRecord.Log = proto.String(log)

	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_RECORD, mgo.Strong)
	defer c.Close()
	err := c.Insert(pushRecord)

	err = xyerror.DBError(err)
	return err
}

// 获取到时间推送的待推送信息
func (db *BatteryDB) GetNotificationByTime(push_time int64) (Notices []battery.NotificationForPush, err error) {
	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_Natice, mgo.Monotonic)
	defer c.Close()
	queryStr := bson.M{"pushtime": bson.M{"$lt": push_time}, "state": PushState_WaitForPush}
	var query = c.Find(queryStr)
	err = query.All(&Notices)

	err = xyerror.DBError(err)
	return
}

func (db *BatteryDB) RecoverNotificationById(notice_id string) error {
	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_Natice, mgo.Strong)
	defer c.Close()
	queryStr := bson.M{"notificationid": notice_id}
	var query = c.Find(queryStr)

	count, err := query.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		var notice battery.NotificationForPush
		err = query.One(&notice)
		if err != nil {
			return err
		}

		// 如果已经过期，删除并插入推送记录
		if notice.GetPushTime() < time.Now().Unix() {
			db.InsertNewPushRecord(notice.GetDeviceId(), notice.GetPushTime(), notice.GetTimestamp(), notice.GetMsg(), notice.GetType(), false, "can not recover for timeout")
			// 暂时不删除
			c.Remove(queryStr)
		} else {
			notice.State = proto.Int32(PushState_WaitForPush)
			err = c.Update(queryStr, notice)
		}
	}
	return err
}

func (db *BatteryDB) RemoveNotificationById(notice_id string, desp string, success bool) error {
	c := db.OpenTable(xybusiness.DB_TABLE_PUSH_Natice, mgo.Strong)
	defer c.Close()
	queryStr := bson.M{"notificationid": notice_id}
	var query = c.Find(queryStr)

	count, err := query.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		var notice battery.NotificationForPush
		err = query.One(&notice)
		if err != nil {
			return err
		}

		db.InsertNewPushRecord(notice.GetDeviceId(), notice.GetPushTime(), notice.GetTimestamp(), notice.GetMsg(), notice.GetType(), success, desp)
		// 暂时不删除
		//c.Remove(queryStr)
	}
	return nil
}
