package batteryapi

// xyapi_rolelist

//购买角色(返回购买确认信息（成功则扣除相关消耗：碎片、金币、宝石），失败则提示错误信息)

//返回游戏结算结果(在游戏结算中应该有了加进去就行)

import (
	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"time"
)

const c_FailoverTime int64 = 10

const (
	NOTIFICATION_CONTENT_FRIEND_CALL = "小伙伴在召唤你~"
)

//好友邮件的消息接口
func (api *XYAPI) OperationFriendMailInfoList(req *battery.FriendMailListRequest, resp *battery.FriendMailListResponse) (err error) {

	var (
		uid       = req.GetUid()
		errStruct = xyerror.DefaultError()
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	errStruct.Code = battery.ErrorCode_NoError.Enum()

	//初始化返回消息
	resp.Uid = req.Uid
	resp.Cmd = req.Cmd
	resp.Source = req.Source
	resp.FriendSid = req.FriendSid
	resp.CreateTime = req.CreateTime

	switch req.GetCmd() {

	case battery.FriendMailCmd_FriendMailCmd_MailList: //查询玩家好友邮件
		resp.MailInfoItemList, errStruct = api.getUserFriendMails(uid, req.GetSource())

	case battery.FriendMailCmd_FriendMailCmd_StaminaGive: //赠送好友体力（from 好友列表）
		var waitTime int64
		waitTime, errStruct = api.giveStamina(uid, req.GetFriendSid(), req.GetSource())
		resp.WaitTime = &waitTime
		if waitTime < -1 {
			resp.CreateTime = proto.Int64(time.Now().Unix())
		}

	case battery.FriendMailCmd_FriendMailCmd_StaminaGetApply: //请求好友赠送体力（from 好友列表）
		var waitTime int64
		waitTime, errStruct = api.applyStamina(uid, req.GetFriendSid(), req.GetSource())
		resp.WaitTime = &waitTime
		if waitTime < -1 {
			resp.CreateTime = proto.Int64(time.Now().Unix())
		}

	case battery.FriendMailCmd_FriendMailCmd_MailStaminaGive: //赠送好友体力（from 邮件）
		errStruct = api.giveStaminaFromMail(uid, req.GetFriendSid(), req.GetSource(), req.GetCreateTime())

	case battery.FriendMailCmd_FriendMailCmd_StaminatGet: //领取体力（from 邮件）
		var staminaNum int32
		staminaNum, errStruct = api.collectStaminaFromMail(uid, req.GetFriendSid(), req.GetCreateTime(), req.GetSource())
		if errStruct.GetCode() == battery.ErrorCode_NoError {
			resp.StaminaNum = proto.Int32(staminaNum)
		}

	case battery.FriendMailCmd_FriendMailCmd_StaminatGetAll: //领取全部（一键确认 from 邮件）
		var staminaNum int32
		staminaNum, errStruct = api.confirmAllMail(uid)
		if errStruct.GetCode() == battery.ErrorCode_NoError {
			resp.StaminaNum = proto.Int32(staminaNum)
		}

	}

	errStruct.Desc = nil
	resp.Error = errStruct

	return
}

//查询好友uid
// friendSid string 好友的sid
//return:
// friendUid string 好友uid
func (api *XYAPI) getFriendUid(uid, friendSid string, source battery.ID_SOURCE) (friendUid string, errStruct *battery.Error) {
	errStruct = new(battery.Error)
	errStruct.Code = battery.ErrorCode_NoError.Enum()

	if friendSid == "" {
		errStruct.Code = battery.ErrorCode_BadInputData.Enum()
		return
	} else {
		var err error
		friendUid, err = api.SidToGid(friendSid, source)
		if friendUid == "" || err != xyerror.ErrOK {
			xylog.Error(uid, "invalid sid : %s", friendSid)
			errStruct.Code = battery.ErrorCode_GetUidError.Enum()
			return
		}
	}
	return
}

//获取玩家的好友邮件
// uid string 玩家id
// source battery.ID_SOURCE 第三方账户来源
//return:
// mailInfoItems []*battery.FriendMailInfo 邮件列表
func (api *XYAPI) getUserFriendMails(uid string, source battery.ID_SOURCE) (mailInfoItems []*battery.FriendMailInfo, errStruct *battery.Error) {

	errStruct = new(battery.Error)
	errStruct.Code = battery.ErrorCode_NoError.Enum()

	mailInfoItemsTmp, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).GetFriendMailList(uid, DefConfigCache.Configs().DefaultGiftValidTimeSec, int(DefConfigCache.Configs().DefaultGiftAcceptMaxCount))
	if err == xyerror.ErrOK {
		xylog.Debug(uid, "user friendmails from db : %v", mailInfoItemsTmp)
		uids := make([]string, 0)
		for _, mailInfoItem := range mailInfoItemsTmp {
			uids = append(uids, mailInfoItem.GetFriendId())
		}

		xylog.Debug(uid, "friends' uid : %v", uids)

		tpids := make([]*battery.IDMap, 0)
		selector := bson.M{}
		err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).XYBusinessDB.QueryTpidsByUids(selector, uids, &tpids)
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_QueryTpidError.Enum()
			return
		}
		mapUid2Sid := make(map[string]string, 0)
		for _, tpid := range tpids {
			mapUid2Sid[tpid.GetGid()] = tpid.GetSid()
		}

		xylog.Debug(uid, "mapUid2Sid : %v", mapUid2Sid)

		for _, mailInfo := range mailInfoItemsTmp {
			if sid, ok := mapUid2Sid[mailInfo.GetFriendId()]; ok {
				mailInfo.FriendId = proto.String(sid)
				item := mailInfo
				mailInfoItems = append(mailInfoItems, item)

			}
		}

		xylog.Debug(uid, "user mailInfoItems : %v", mailInfoItems)

	} else {
		errStruct.Code = battery.ErrorCode_QueryUserFriendMailError.Enum()
		xylog.Error(uid, "QueryUserFriendMailError : %v", err)
		return
	}
	return
}

//赠送好友体力（from 好友列表）
// uid string 玩家id
// friendSid string 好友sid
// source battery.ID_SOURCE 好友来源
//return:
// waitTime int64 下次赠送等待时间
func (api *XYAPI) giveStamina(uid, friendSid string, source battery.ID_SOURCE) (waitTime int64, errStruct *battery.Error) {
	errStruct = new(battery.Error)
	errStruct.Code = battery.ErrorCode_NoError.Enum()

	var (
		friendUid string
		err       error
	)

	friendUid, errStruct = api.getFriendUid(uid, friendSid, source)
	if friendUid != "" && errStruct.GetCode() == battery.ErrorCode_NoError {
		waitTime, err = api.subGiveStamina(uid, friendUid)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "subGiveStamina failed : %v, waitTime : %d", err, waitTime)
			errStruct.Code = battery.ErrorCode_GiveStaminaToFriendError.Enum()
			return
		}

		//修改玩家的好友邮件数目
		api.changeFriendMailCount(friendUid)

	}

	return
}

//向好友申请体力（from 好友列表）
// uid string 玩家id
// friendSid string 好友sid
// source battery.ID_SOURCE 好友来源
//return:
// waitTime int64 下次赠送等待时间
func (api *XYAPI) applyStamina(uid, friendSid string, source battery.ID_SOURCE) (waitTime int64, errStruct *battery.Error) {
	errStruct = new(battery.Error)
	errStruct.Code = battery.ErrorCode_NoError.Enum()

	var (
		friendUid string
		err       error
	)
	friendUid, errStruct = api.getFriendUid(uid, friendSid, source)
	if friendUid != "" && errStruct.GetCode() == battery.ErrorCode_NoError {
		waitTime, err = api.subApplyStamina(uid, friendUid)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "subApplyStamina failed : %v, waitTime : %d", err, waitTime)
			errStruct.Code = battery.ErrorCode_ApplyStaminaToFriendError.Enum()
			return
		}

		//修改玩家的好友邮件数目
		api.changeFriendMailCount(friendUid)
	}

	return
}

//赠送好友体力（从好友列表）
// uid string 玩家id
// friendUid string 好友id
//return:
// waitTime int64 剩余等待时间
func (api *XYAPI) subGiveStamina(uid string, friendUid string) (waitTime int64, err error) {
	waitTime = -1
	var (
		staminaGiveLogItem battery.StaminaGiveLogItem
		bFound             = true
	)

	staminaGiveLogItem, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_STAMINAGIVEAPPLYLOG).GetStaminaGiveLogItem(uid, friendUid)
	if err != xyerror.ErrNotFound && err != xyerror.ErrOK {
		xylog.Debug(uid, "GetStaminaGiveLogItem failed : %v", err)
		return
	} else if err == xyerror.ErrNotFound {
		bFound = false
	}

	now := time.Now().Unix()
	waitTime = api.subGetWaitTime(staminaGiveLogItem.StaminaGiveLastTime, DefConfigCache.Configs().DefaultGiftGiveCooldown, now)

	if waitTime <= c_FailoverTime { //剩余请求时间必须是小于这个值才能允许请求（为了防止时间差，即前端时间小于后端时间，提早了10秒钟放开申请），防止恶意刷体力
		staminaGiveLogItem.StaminaGiveLastTime = &now
		staminaGiveLogItem.Uid = &uid
		staminaGiveLogItem.FriendUid = &friendUid

		if !bFound {
			staminaGiveLogItem.StaminaApplyLastTime = proto.Int64(0)
		}

		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_STAMINAGIVEAPPLYLOG).UpsertStaminaGiveLogItem(&staminaGiveLogItem)
		xylog.Debug(uid, "UpsertStaminaGiveLogItem (%v) err : %v", &staminaGiveLogItem, err)
		if err == xyerror.ErrOK {
			err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).AddFriendMail(friendUid, uid, battery.FriendMailType_FriendMailType_StaminaGive) //添加邮件
			if err == xyerror.ErrOK {
				waitTime = -DefConfigCache.Configs().DefaultGiftGiveCooldown
				//刷新体力赠送的相关任务状态
				quotas := []*battery.Quota{&battery.Quota{Id: battery.QuotaEnum_Quota_GiveStamina.Enum(), Value: proto.Uint64(1)}}
				missionTypes := []battery.MissionType{battery.MissionType_MissionType_Study, battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_MainLine}
				api.updateUserMissionsQuotas(uid, missionTypes, quotas, time.Now().Unix(), MissionQuotasNoNeedFinish)

				//触发好友邮件推送
				//go api.TrigerFriendMailNotify(friendUid)

			}
		}

	} else {
		err = xyerror.ErrTimeLimitError
		return
	}

	return
}

//获取剩余等待时间
// waitTimeSrc *int64 到期时间戳
// coolDownSecs int64 唤醒时间
// now int64 当前时间戳
func (api *XYAPI) subGetWaitTime(waitTimeSrc *int64, coolDownSecs, now int64) (waitTime int64) {
	if nil != waitTimeSrc {
		nextGiveTime := *waitTimeSrc + coolDownSecs
		waitTime = nextGiveTime - now
		if waitTime < 0 { //已经过了刷新时间
			waitTime = 0
		}
	}
	return
}

func (api *XYAPI) subApplyStamina(uid string, friendUid string) (waitTime int64, err error) {
	waitTime = -1
	var (
		staminaGiveLogItem battery.StaminaGiveLogItem
		bFound             = true
	)

	staminaGiveLogItem, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_STAMINAGIVEAPPLYLOG).GetStaminaGiveLogItem(uid, friendUid)
	if err != xyerror.ErrNotFound && err != xyerror.ErrOK {
		xylog.Debug(uid, "GetStaminaGiveLogItem failed : %v", err)
		return
	} else if err == xyerror.ErrNotFound {
		bFound = false
	}

	now := time.Now().Unix()
	waitTime = api.subGetWaitTime(staminaGiveLogItem.StaminaApplyLastTime, DefConfigCache.Configs().DefaultGiftAskCooldown, now)

	if waitTime <= c_FailoverTime { //剩余请求时间必须是小于这个值才能允许请求（为了防止时间差，即前端时间小于后端时间，提早了10秒钟放开申请），防止恶意刷体力
		staminaGiveLogItem.StaminaApplyLastTime = &now
		staminaGiveLogItem.Uid = &uid
		staminaGiveLogItem.FriendUid = &friendUid

		if !bFound {
			staminaGiveLogItem.StaminaGiveLastTime = proto.Int64(0)
		}

		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_STAMINAGIVEAPPLYLOG).UpsertStaminaGiveLogItem(&staminaGiveLogItem)
		if err == xyerror.ErrOK {
			err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).AddFriendMail(friendUid, uid, battery.FriendMailType_FriendMailType_StaminaApply) //添加邮件
			if err == xyerror.ErrOK {
				waitTime = -DefConfigCache.Configs().DefaultGiftAskCooldown
				//触发好友邮件推送
				//go api.TrigerFriendMailNotify(friendUid)
			}
		}
	} else {
		err = xyerror.ErrTimeLimitError
		return
	}

	return
}

//从邮件中赠送好友体力
// uid string 玩家id
// friendSid string 好友sid
// source battery.ID_SOURCE 第三方好友类型
// createTime int64 创建时间
func (api *XYAPI) giveStaminaFromMail(uid, friendSid string, source battery.ID_SOURCE, createTime int64) (errStruct *battery.Error) {
	var friendUid string
	friendUid, errStruct = api.getFriendUid(uid, friendSid, source)
	if friendUid != "" && errStruct.GetCode() == battery.ErrorCode_NoError {
		if createTime != 0 {
			err := api.MailStaminaGive(uid, friendUid, createTime)
			if err != xyerror.ErrOK {
				xylog.Error(uid, "give stamina from mail failed : %v", err)
				errStruct.Code = battery.ErrorCode_GiveStaminaToFriendError.Enum()
				return
			}
		}
	}

	return
}

//从邮件领取体力
// uid string 玩家id
// friendSid string 玩家sid
// createTime int64 邮件创建时间
func (api *XYAPI) collectStaminaFromMail(uid, friendSid string, createTime int64, source battery.ID_SOURCE) (staminatNum int32, errStruct *battery.Error) {
	var (
		friendUid string
	)

	errStruct = new(battery.Error)
	errStruct.Code = battery.ErrorCode_NoError.Enum()

	friendUid, errStruct = api.getFriendUid(uid, friendSid, source)
	if createTime != 0 && friendUid != "" && errStruct.GetCode() == battery.ErrorCode_NoError {
		err := api.StaminaGet(uid, friendUid, createTime)
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_CollectStaminaMailFromMailError.Enum()
			return
		} else {
			//获取玩家当前体力
			staminatNum, _, err = api.GetCurrentStaminaDirect(uid)
			if err != xyerror.ErrOK {
				errStruct.Code = battery.ErrorCode_QueryStaminaError.Enum()
				return
			}
		}
	} else {
		errStruct.Code = battery.ErrorCode_BadInputData.Enum()
		return
	}

	return
}

//一键确认所有邮件
//
func (api *XYAPI) confirmAllMail(uid string) (staminatNum int32, errStruct *battery.Error) {
	errStruct = new(battery.Error)
	errStruct.Code = battery.ErrorCode_NoError.Enum()

	//领取所有的好友赠送体力
	err := api.StaminaGetAll(uid)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_CollectStaminaMailFromMailError.Enum()
		return
	}

	//确认所有的好友体力请求
	var friendUids []string
	friendUids, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).GiveStaminaToAllFriendApplyMailList(uid)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_GiveStaminaToFriendError.Enum()
		return
	} else {
		//刷新所有好友的邮件数目
		for _, friendUid := range friendUids {
			api.changeFriendMailCount(friendUid)
		}

		//查询最新的体力值，返回前端
		staminatNum, _, err = api.GetCurrentStaminaDirect(uid)
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_QueryStaminaError.Enum()
			return
		}
	}
	return
}

//邮件确认赠送体力
// uid string 玩家id
// friendUid string 好友uid
// createTime int64 邮件创建时间
func (api *XYAPI) MailStaminaGive(uid string, friendUid string, createTime int64) (err error) {
	if api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).IsFriendMailExisting(uid, friendUid, battery.FriendMailType_FriendMailType_StaminaApply, createTime) { //这个校验必须是有的，防止恶意刷体力
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).AddFriendMail(friendUid, uid, battery.FriendMailType_FriendMailType_StaminaGive) //在好友的邮箱里添加邮件
		if err == xyerror.ErrOK {
			xylog.Error(uid, "[app RemoveFriendMail] Invalid user")
			err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).RemoveFriendMail(uid, friendUid, battery.FriendMailType_FriendMailType_StaminaApply, createTime) //删除邮件
		}

		//修改玩家的好友邮件数目
		api.changeFriendMailCount(uid)
		api.changeFriendMailCount(friendUid)

	} else {
		err = xyerror.ErrNotFound
	}
	return
}

//邮件领取体力
// uid string 玩家id
// friendUid string 好友uid
// createTime int64 邮件创建时间
func (api *XYAPI) StaminaGet(uid string, friend_uid string, createTime int64) (err error) {
	if api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).IsFriendMailExisting(uid, friend_uid, battery.FriendMailType_FriendMailType_StaminaGive, createTime) {
		//获取体力
		var propItem battery.PropItem
		var propType battery.PropType = battery.PropType_PROP_STAMINA
		propItem.Type = &propType
		propItem.Amount = proto.Uint32(1)
		err = api.GainProp(uid, nil, &propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
		//删除邮件
		if err == nil {
			err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).RemoveFriendMail(uid, friend_uid, battery.FriendMailType_FriendMailType_StaminaGive, createTime)
		}

		//修改玩家的好友邮件数目
		api.changeFriendMailCount(uid)

	} else {
		err = xyerror.ErrNotFound
	}
	return err
}

//一键领取所有赠送体力
// uid string 玩家id
func (api *XYAPI) StaminaGetAll(uid string) (err error) {

	var (
		staminaGetCount int
	)

	//获取所有体力赠送
	staminaGetCount, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).GetStaminaGiveMailCount(uid)
	if xyerror.ErrOK != err {
		xylog.Error(uid, "GetStaminaGiveMailCount failed : %v", err)
	} else if staminaGetCount > 0 {
		propItem := &battery.PropItem{
			Type:   battery.PropType_PROP_STAMINA.Enum(),
			Amount: proto.Uint32(uint32(staminaGetCount)),
		}
		xylog.Debug(uid, "StaminaGetAll Stamina propItem %v", propItem)
		//发放邮件里面的所有体力
		err = api.GainProp(uid, nil, propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
		if err == xyerror.ErrOK { //删除所有邮件
			err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).RemoveAllStaminaGiveMail(uid)
		}
	}

	//修改玩家的好友邮件数目
	api.changeFriendMailCount(uid)

	return
}

// changeFriendMailCount 修改玩家的好友邮件计数
//uid string    玩家标识
//change int32  邮件变化数
func (api *XYAPI) changeFriendMailCount(uid string) {
	//查询玩家当前邮件数目
	count, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).GetFriendMailCount(uid)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "GetFriendMailCount from DB failed : %v", err)
		return
	}

	//刷新玩家当前邮件数目
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAILCOUNT).UpsertFriendMailCount(uid, int32(count))
	if err != xyerror.ErrOK {
		xylog.Error(uid, "UpsertFriendMailCount failed : %v", err)
		return
	}

	return
}

////根据好友邮件的数目判断是否需要发送推送信息
//func (api *XYAPI) TrigerFriendMailNotify(uid string) {
//	//获取玩家当前的好友邮件数
//	count, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_FRIENDMAIL).GetFriendMailCount(uid)
//	if err != xyerror.ErrOK {
//		xylog.Error(uid, "[%s] GetFriendMailCount from DB failed : %v", uid, err)
//		return
//	}

//	//判断邮件数是否达到阈值，达到的话，就发送推送信息
//	xylog.Debug(uid, "[%s] FriendMailCount(%d) DefaultGiftNotifyCount(%d)", uid, count, DefConfigCache.Configs().DefaultGiftNotifyCount)

//	if count >= DefConfigCache.Configs().DefaultGiftNotifyCount {

//		//获取玩家的device token
//		deviceToken, err := api.GetDeviceToken(uid, mgo.Strong)
//		if err != xyerror.ErrOK {
//			xylog.Error(uid, "[%s] GetDeviceToken failed : %v", uid, err)
//			return
//		}

//		notification := &xyapn.APNNotification{
//			Cmd:         xyapn.APNNotificationCMD_Notification.Enum(),
//			Uid:         &uid,
//			DeviceToken: &deviceToken,
//			Content:     proto.String(NOTIFICATION_CONTENT_FRIEND_CALL),
//		}
//		xyapn.Send(notification)
//	}
//}
