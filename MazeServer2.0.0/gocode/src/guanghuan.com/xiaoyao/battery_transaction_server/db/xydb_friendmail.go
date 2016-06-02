package batterydb

import (
	proto "code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"time"
)

//获取好友邮件
func (db *BatteryDB) GetFriendMailList(uid string, mailValidSecs int64, mailMaxCount int) (friendMailList []*battery.FriendMailInfo, err error) {
	//删除过期邮件
	curTime := time.Now().Unix()
	curTime -= mailValidSecs
	tbl := db.OpenTable(xybusiness.DB_TABLE_FRIENDMAIL, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"uid": uid, "createtime": bson.M{"$lt": curTime}}
	_, err = tbl.RemoveAll(condition)

	//获取当前邮件
	if err == xyerror.ErrOK {
		condition := bson.M{"uid": uid}
		query := tbl.Find(condition).Sort("-createtime").Limit(mailMaxCount) //最多30封邮件
		err = query.All(&friendMailList)
	} else {
		xylog.Error(uid, "GetFriendMailList RemoveData failed : %v", err)
		return
	}

	return
}

//获取玩家好友邮件的数目
func (db *BatteryDB) GetFriendMailCount(uid string) (count int, err error) {
	condition := bson.M{"uid": uid}
	//xylog.Debug(uid, "[%s] GetFriendMailCount condition : %v", uid, condition)
	count, err = db.GetRecordCount(xybusiness.DB_TABLE_FRIENDMAIL, condition, mgo.Strong)
	return
}

//增加好友邮件
// uid string 玩家id
// friendUid string 好友uid
// mailType battery.FriendMailType 邮件类型
func (db *BatteryDB) AddFriendMail(uid string, friendUid string, mailType battery.FriendMailType) (err error) {
	friendMail := battery.FriendMailInfo{
		Uid:        proto.String(uid),
		FriendId:   proto.String(friendUid),
		Mailtype:   mailType.Enum(),
		CreateTime: proto.Int64(time.Now().Unix()),
	}

	err = db.AddData(xybusiness.DB_TABLE_FRIENDMAIL, friendMail)
	return
}

//好友邮件是否存在
// uid string 玩家id
// friendUid string 好友uid
// mailType battery.FriendMailType 邮件类型
// createTime int64 邮件创建时间
func (db *BatteryDB) IsFriendMailExisting(uid string, friendUid string, mailType battery.FriendMailType, createTime int64) (isExist bool) {
	condition := bson.M{"uid": uid, "friendid": friendUid, "mailtype": mailType, "createtime": createTime}
	isExist, _ = db.IsRecordExisting(xybusiness.DB_TABLE_FRIENDMAIL, condition, mgo.Strong)
	return
}

//删除好友邮件
// uid string 玩家id
// friendUid string 好友uid
// mailType battery.FriendMailType 邮件类型
// createTime int64 创建时间
func (db *BatteryDB) RemoveFriendMail(uid string, friendUid string, mailType battery.FriendMailType, createTime int64) (err error) {
	condition := bson.M{"uid": uid, "friendid": friendUid, "mailtype": mailType, "createtime": createTime}
	err = db.RemoveData(xybusiness.DB_TABLE_FRIENDMAIL, condition)
	return
}

//获取体力赠送邮件数
// uid string 玩家id
//return:
// count int 邮件数目
func (db *BatteryDB) GetStaminaGiveMailCount(uid string) (count int, err error) {
	condition := bson.M{"uid": uid, "mailtype": battery.FriendMailType_FriendMailType_StaminaGive}
	count, err = db.GetRecordCount(xybusiness.DB_TABLE_FRIENDMAIL, condition, mgo.Strong)
	return
}

//删除所有体力赠送邮件
// uid string 玩家id
func (db *BatteryDB) RemoveAllStaminaGiveMail(uid string) (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_FRIENDMAIL, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"uid": uid, "mailtype": battery.FriendMailType_FriendMailType_StaminaGive}
	_, err = tbl.RemoveAll(condition)
	return
}

//删除所有体力请求邮件
// uid string 玩家id
func (db *BatteryDB) RemoveAllStaminaApplyMail(uid string) (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_FRIENDMAIL, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"uid": uid, "mailtype": battery.FriendMailType_FriendMailType_StaminaApply}
	_, err = tbl.RemoveAll(condition)
	return
}

//应答所有好友体力赠送请求
func (db *BatteryDB) GiveStaminaToAllFriendApplyMailList(uid string) (friendUids []string, err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_FRIENDMAIL, mgo.Strong)
	defer tbl.Close()

	//查询所有的好友邮件
	friendMailList := make([]*battery.FriendMailInfo, 0)
	condition := bson.M{"uid": uid, "mailtype": battery.FriendMailType_FriendMailType_StaminaApply}
	query := tbl.Find(condition).Sort("-createtime")
	err = query.All(&friendMailList)
	if err != xyerror.ErrOK {
		//xylog.Error("[%s] query all friend apply mail failed : %v", uid, err)
		return
	}
	//为所有好友体力申请请求增加对应邮件，批量插入邮件信息
	mails := make([]interface{}, 0)
	friendUids = make([]string, 0)
	for _, mailinfo := range friendMailList {
		mail := &battery.FriendMailInfo{
			Uid:        proto.String(mailinfo.GetFriendId()),
			FriendId:   &uid,
			Mailtype:   battery.FriendMailType_FriendMailType_StaminaGive.Enum(),
			CreateTime: proto.Int64(time.Now().Unix()),
		}
		mails = append(mails, mail)
		friendUids = append(friendUids, mailinfo.GetFriendId())
	}

	if len(mails) > 0 { //如果
		bulk := tbl.Bulk()
		bulk.Unordered()
		bulk.Insert(mails...)
		_, err = bulk.Run()
		if err != xyerror.ErrOK {
			xylog.Error(uid, "bulk insert friend mails failed : %v", err)
			return
		}

		//删除所有的好友体力请求邮件
		_, err = tbl.RemoveAll(condition)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "remove friend apply mail failed : %v", err)
			return
		}
	}

	return
}

//查询玩家的好友邮件时间戳信息
// uid string 玩家id
// friendUid string 玩家好友id
func (db *BatteryDB) GetStaminaGiveLogItem(uid, friendUid string) (staminaGiveLogItem battery.StaminaGiveLogItem, err error) {
	//只是查询赠送体力和请求体力的时间戳就可以了
	condition := bson.M{"uid": uid, "frienduid": friendUid}
	selector := bson.M{ /*"uid": 1, "frienduid": 1, */ "staminagivelasttime": 1, "staminaapplylasttime": 1}
	err = db.GetOneData(xybusiness.DB_TABLE_STAMINA_GIVEAPPLY_LOG, condition, selector, &staminaGiveLogItem, mgo.Strong)
	return
}

// UpsertStaminaGiveLogItem 查询玩家的好友邮件时间戳信息
// uid string 玩家id
// friendUid string 玩家好友id
func (db *BatteryDB) UpsertStaminaGiveLogItem(staminaGiveLogItem *battery.StaminaGiveLogItem) (err error) {
	return db.UpsertData(xybusiness.DB_TABLE_STAMINA_GIVEAPPLY_LOG, bson.M{"uid": staminaGiveLogItem.GetUid(), "frienduid": staminaGiveLogItem.GetFriendUid()}, staminaGiveLogItem)
}

// UpsertFriendMailCount 修改玩家好友邮件数目
// uid string 玩家id
// count int32 玩家好友邮件数
func (db *BatteryDB) UpsertFriendMailCount(uid string, count int32) (err error) {
	friendMailCount := &battery.DBFriendMailCount{
		Uid:   &uid,
		Count: &count,
	}
	return db.UpsertData(xybusiness.DB_TABLE_FRIENDMAILCOUNT, bson.M{"uid": uid}, friendMailCount)
}

// RemoveFriendMailCount 删除玩家好友邮件数目记录
// uid string 玩家id
func (db *BatteryDB) RemoveFriendMailCount(uid string) (err error) {
	return db.RemoveData(xybusiness.DB_TABLE_FRIENDMAILCOUNT, bson.M{"uid": uid})
}
