package batterydb

import (
	//proto "code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//玩家阅读系统邮件
// uid string 玩家id
// mailID int32 邮件id
func (db *BatteryDB) ReadSystemMail(uid string, mailID int32) (err error) {
	var tempSystemMailList *battery.DBSystemMailListItem
	tempSystemMailList, err = db.GetSystemMailList(uid)
	if err != xyerror.ErrOK {
		//xylog.Error("[%s] GetSystemMailList failed : %v", uid, err)
		return
	} else {
		for _, mailItem := range tempSystemMailList.MailInfoList {
			id := mailItem.GetMailId()
			if id == mailID {
				//获取邮件的配置信息
				config := xybusinesscache.DefMailConfigCacheManager.MailConfig(id)
				if nil == config {
					xylog.Debug(uid, "failed to find MailConfig of %d,pls check.", id)
					err = xyerror.ErrBadInputData
					return
				}

				if config.GetPropID() <= xybusinesscache.INVALID_PROPID { //非礼包邮件
					if mailItem.GetIsRead() {
						xylog.Error(uid, "mail(%d) already read, pls check", mailID)
						err = xyerror.ErrBadInputData
						return
					} else { //标记为已经阅读 在下次查询邮件的时候如果不是配置邮件会被删除
						*mailItem.IsRead = true
						err = db.UpsertMailInfo(tempSystemMailList)
						return
					}
				} else { //带奖励的邮件只能领取，不能read
					xylog.Error(uid, "read mail(%d) with prop, are you sure?", mailID)
					err = xyerror.ErrBadInputData
					return
				}
			}
		}
	}

	return
}

//更新玩家系统邮件信息
func (db *BatteryDB) UpsertMailInfo(systemMailListItem *battery.DBSystemMailListItem) (err error) {
	condition := bson.M{"uid": systemMailListItem.GetUid()}
	err = db.UpsertData(xybusiness.DB_TABLE_SYSTEMMAIL, condition, systemMailListItem)
	return
}

//获取玩家系统邮件系统，马上更新
// uid string 玩家id
//return:
// systemMailList *battery.DBSystemMailListItem 系统邮件信息
// err error 操作结果
func (db *BatteryDB) GetSystemMailList(uid string) (systemMailList *battery.DBSystemMailListItem, err error) {
	systemMailList = &battery.DBSystemMailListItem{}
	condition := bson.M{"uid": uid}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_SYSTEMMAIL, condition, selector, &systemMailList, mgo.Strong)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			//第一次创建邮件列表 需要保存一个 ID 用于运营渠道添加的邮件
			systemMailList.Uid = &uid
			//systemMailList.ChangeID = proto.Int32(c_StartChangeID + 1)
			err = db.UpsertMailInfo(systemMailList)
		}
	}
	return
}

//获取玩家系统邮件系统，但是不马上更新，与调用者后续的修改合并。以此来减少对数据的upsert，优化性能
// uid string 玩家id
//return:
// systemMailList *battery.DBSystemMailListItem 系统邮件信息
// chage bool 是否有变更
// err error 操作结果
func (db *BatteryDB) GetSystemMailListWithoutChange(uid string) (systemMailList *battery.DBSystemMailListItem, change bool, err error) {
	systemMailList, change, err = &battery.DBSystemMailListItem{}, false, xyerror.ErrOK
	condition := bson.M{"uid": uid}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_SYSTEMMAIL, condition, selector, &systemMailList, mgo.Strong)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			//第一次创建邮件列表 需要保存一个 ID 用于运营渠道添加的邮件
			systemMailList.Uid = &uid
			//systemMailList.ChangeID = proto.Int32(c_StartChangeID + 1)
			change = true
			err = xyerror.ErrOK
		}
	}
	return
}
