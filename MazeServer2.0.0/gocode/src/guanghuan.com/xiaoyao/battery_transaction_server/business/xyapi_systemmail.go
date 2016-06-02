package batteryapi

// xyapi_rolelist

//购买角色(返回购买确认信息（成功则扣除相关消耗：碎片、金币、宝石），失败则提示错误信息)

//返回游戏结算结果(在游戏结算中应该有了加进去就行)

import (
    "time"

    "code.google.com/p/goprotobuf/proto"

    "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/cache"
    xycache "guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/money"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//const c_StartChangeID int32 = 1000000

//系统邮件相关消息接口
func (api *XYAPI) OperationSystemMailInfoList(req *battery.SystemMailListRequest, resp *battery.SystemMailListResponse) (err error) {

    var (
        uid       = req.GetUid()
        errStruct = xyerror.DefaultError()
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    //初始化返回消息
    resp.Uid = req.Uid
    resp.Cmd = req.Cmd
    resp.MailId = req.MailId

    switch req.GetCmd() {

    case battery.SystemMailCmd_SystemMailCmd_Maillist: //查询玩家系统邮件
        resp.MailInfoItemList = api.queryUserSystemMails(uid, errStruct)

    case battery.SystemMailCmd_SystemMailCmd_GiftGet: //确认玩家系统邮件
        resp.Wallet = api.confirmSystemMail(uid, req.GetMailId(), errStruct)

    case battery.SystemMailCmd_SystemMailCmd_MailRead: //阅读玩家系统邮件
        api.readUserSystemMail(uid, req.GetMailId(), errStruct)
    }

    resp.Error = errStruct

    return
}

//查询玩家的系统邮件
// uid string 玩家id
//return:
// sysMailInfos []*battery.SystemMailInfoConfig 玩家系统邮件列表
func (api *XYAPI) queryUserSystemMails(uid string, errStruct *battery.Error) (sysMailInfos []*battery.SystemMailInfoConfig) {

    var err error

    sysMailInfos, err = api.GetAndCreateSystemMailListItem(uid)

    if err != xyerror.ErrOK {
        xylog.Error(uid, "GetAndCreateSystemMailListItem failed : %v")
        errStruct.Code = battery.ErrorCode_QueryUserSysMailError.Enum()
        return
    }

    return
}

//确认邮件
// uid string 玩家id
// mailId int32 邮件id
func (api *XYAPI) confirmSystemMail(uid string, mailId int32, errStruct *battery.Error) (wallet []*battery.MoneyItem) {

    err := api.ConfirmGiftEmail(uid, mailId)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_ConfirmUserSysMailError.Enum()
        xylog.Error(uid, "ConfirmGiftEmail failed : %v", err)
        return
    } else {
        wallet, err = xymoney.QueryWallet(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT))
        if err != xyerror.ErrOK {
            errStruct.Code = battery.ErrorCode_QueryWalletError.Enum()
            xylog.Error(uid, "QueryWallet failed : %v", err)
            return
        }
    }

    return
}

//阅读系统邮件
// uid string 玩家id
// mailId int32 邮件id
func (api *XYAPI) readUserSystemMail(uid string, mailId int32, errStruct *battery.Error) {

    err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).ReadSystemMail(uid, mailId)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "ReadSystemMail failed : %v", err)
        errStruct.Code = battery.ErrorCode_ReadUserSysMailError.Enum()
        return
    }

    return
}

//获取玩家的系统邮件信息
// uid string 玩家id
func (api *XYAPI) GetAndCreateSystemMailListItem(uid string) (sysMailInfos []*battery.SystemMailInfoConfig, err error) {
    var (
        resultSystemMailList = battery.DBSystemMailListItem{
            Uid: &uid,
        }
        newSystemMailList = battery.DBSystemMailListItem{
            Uid: &uid,
        }
        timeNow = time.Now()
    )

    //查询玩家当前的系统邮件信息
    var (
        tempSystemMailList *battery.DBSystemMailListItem
        change             = false
    )
    tempSystemMailList, change, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).GetSystemMailListWithoutChange(uid)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "GetSystemMailList failed : %v", err)
        return
    }

    newSystemMailList.DynamicCounter = tempSystemMailList.DynamicCounter

    ExistConfigMailIDList := make(Set, 0)
    for _, mailItem := range tempSystemMailList.MailInfoList {
        id := mailItem.GetSourceMailId()
        xylog.Debug(uid, "mailItem : %v", mailItem)
        config := xybusinesscache.DefMailConfigCacheManager.MailConfig(id)
        if nil == config {
            xylog.Error(uid, "failed to get MailConfig of %d", id)
            continue
        }

        var (
            endTimeStamp = mailItem.GetEndTime()
            mailEndTime  = time.Unix(endTimeStamp, 0)
            propId       = config.GetPropID()
            sourceType   = mailItem.GetSourceType()
        )

        //如果是动态邮件，奖励的道具id必须从邮件信息中取，而不是从静态配置信息中取
        if mailItem.GetSourceType() == battery.SystemMailSourceType_SystemMailSourceType_Dynamic {
            propId = mailItem.GetPropId()
        }

        //记录当前已存在的邮件哪些是要返回客户端或保存的
        if propId <= xybusinesscache.INVALID_PROPID && endTimeStamp != 0 && mailEndTime.After(timeNow) { //非礼包邮件，只要未过期，就保存数据库中并返回客户端
            resultSystemMailList.MailInfoList = append(resultSystemMailList.MailInfoList, mailItem)
            newSystemMailList.MailInfoList = append(newSystemMailList.MailInfoList, mailItem)
        } else if propId > xybusinesscache.INVALID_PROPID { //礼包邮件
            if mailItem.GetIsConfirm() && endTimeStamp != 0 && mailEndTime.After(timeNow) { //已经确认，只保存在数据库中，不返回客户端
                newSystemMailList.MailInfoList = append(newSystemMailList.MailInfoList, mailItem)
            } else if !mailItem.GetIsConfirm() { //只要是未确认的，不管是否过期，都保存数据库中，并且返回客户端
                resultSystemMailList.MailInfoList = append(resultSystemMailList.MailInfoList, mailItem)
                newSystemMailList.MailInfoList = append(newSystemMailList.MailInfoList, mailItem)
            }
        }

        //配置邮件,记录该邮件已存在邮箱中
        if sourceType == battery.SystemMailSourceType_SystemMailSourceType_Configuration {
            ExistConfigMailIDList[id] = empty{}
        }

        change = true
    }

    xylog.Debug(uid, "ExistConfigMailIDList : %v", ExistConfigMailIDList)

    //遍历邮件配置信息，添加需要新增的配置邮件
    mapMailConfig := xybusinesscache.DefMailConfigCacheManager.MailConfigs()
    if mapMailConfig == nil {
        xylog.Error(uid, "get mapMailConfig from cache failed")
        return
    } else {
        for id, mailConfig := range *mapMailConfig {

            //跳过动态邮件配置信息
            if id == SYSMAIL_TYPE_DYNAMIC*SYSMAIL_SEGMENT || //运营动态邮件
                id == SYSMAIL_TYPE_WEIBO_FIRSTLOGIN_DYNAMIC*SYSMAIL_SEGMENT ||
                id == SYSMAIL_TYPE_GUEST_FIRSTLOGIN_DYNAMIC*SYSMAIL_SEGMENT { //首次登录动态邮件
                continue
            }

            mailConfigTmp := mailConfig

            mailBeginTime, mailEndTime := time.Unix(mailConfig.GetStartTime(), 0), time.Unix(mailConfig.GetEndTime(), 0)
            if mailBeginTime.Before(timeNow) && mailEndTime.After(timeNow) {
                if _, ok := ExistConfigMailIDList[id]; !ok { //如果是未存在的配置邮件，就增加一条
                    dbUserMail := &battery.DBUserSystemMail{
                        MailId:       proto.Int32(id),
                        SourceMailId: proto.Int32(id),
                        MailType:     mailConfigTmp.Mailtype,
                        SourceType:   battery.SystemMailSourceType_SystemMailSourceType_Configuration.Enum(),
                        BeginTime:    proto.Int64(mailConfig.GetStartTime()),
                        EndTime:      proto.Int64(mailConfig.GetEndTime()),
                        IsRead:       proto.Bool(false),
                        IsConfirm:    proto.Bool(false),
                        PropId:       proto.Uint64(mailConfig.GetPropID()),
                    }
                    resultSystemMailList.MailInfoList = append(resultSystemMailList.MailInfoList, dbUserMail)
                    newSystemMailList.MailInfoList = append(newSystemMailList.MailInfoList, dbUserMail)
                    change = true
                }
            }
        }

    }

    xylog.Debug(uid, "newSystemMailList : %v", newSystemMailList)

    //如果系统邮件信息有变更，则upsert
    if change {
        err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).UpsertMailInfo(&newSystemMailList)
    }

    //拼装请求消息返回结果
    if err == xyerror.ErrOK {
        for _, userMail := range resultSystemMailList.MailInfoList {

            id := userMail.GetSourceMailId()
            config := xybusinesscache.DefMailConfigCacheManager.MailConfig(id)
            if nil == config {
                xylog.Error(uid, "failed to get MailConfig of %d", id)
                continue
            }

            sysMailInfo := &battery.SystemMailInfoConfig{
                MailID:      userMail.MailId,
                Mailtype:    userMail.MailType,
                Title:       config.Title,
                Message:     config.Message,
                Description: config.Description,
                StartTime:   userMail.BeginTime,
                EndTime:     userMail.EndTime,
                IsReaded:    userMail.IsRead,
                //PropItems:   propitems,
            }

            // 根据礼包id获取道具列表
            propId := userMail.GetPropId()
            if propId != xybusinesscache.INVALID_PROPID {
                propStruct := xycache.DefPropCacheManager.Prop(propId)
                if nil == propStruct {
                    xylog.Error(uid, "get propStruct of %d failed", propId)
                    return
                }
                //propitems := propStruct.Items
                sysMailInfo.PropItems = propStruct.Items
            }

            xylog.DebugNoId("sysmail propitems:%v", sysMailInfo.PropItems)

            //动态邮件的propid非配置的
            if userMail.GetSourceType() == battery.SystemMailSourceType_SystemMailSourceType_Dynamic {
                sysMailInfo.PropID = proto.Uint64(userMail.GetPropId())
            } else {
                sysMailInfo.PropID = proto.Uint64(config.GetPropID())
            }

            sysMailInfos = append(sysMailInfos, sysMailInfo)
        }
    }

    return
}

//确认邮件
// uid string 玩家id
// mailID int32 邮件id
func (api *XYAPI) ConfirmGiftEmail(uid string, mailID int32) (err error) {

    var (
        sourceId  int32
        isDynamic = api.isDynamic(uid, mailID)
    )
    xylog.Debug(uid, "isDynamic %v", isDynamic)
    if isDynamic {
        sourceId = (mailID / SYSMAIL_SEGMENT) * SYSMAIL_SEGMENT
    } else {
        sourceId = mailID
    }

    //有物品的话先领取物品
    var config *battery.DBMailInfoConfig
    config, err = api.MailConfigDetail(uid, sourceId)
    if err == xyerror.ErrOK && config != nil {
        //获取系统邮件信息，判断该邮件是否已经确认过
        var tempSystemMailList *battery.DBSystemMailListItem
        tempSystemMailList, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).GetSystemMailList(uid)
        if err != xyerror.ErrOK {
            xylog.Error(uid, "GetSystemMailList failed : %v", err)
            return
        }

        var mailItemTmp *battery.DBUserSystemMail
        found := false
        for _, mailItem := range tempSystemMailList.MailInfoList {
            if mailItem.GetMailId() == mailID { //找到对应的邮件
                mailItemTmp = mailItem
                found = true
                //if *mailItem.EndTime < *mailItem.BeginTime { //邮件已经确认过
                if mailItem.GetIsConfirm() { //邮件已经确认过
                    xylog.Error(uid, "mail(%d) already be removed", mailID)
                    err = xyerror.ErrBadInputData
                    return
                }
                break
            }
        }

        //没找到，就有点扯了。。。
        if !found {
            xylog.Error(uid, "mailId(%d) no found in systemmailbox", mailID)
            err = xyerror.ErrBadInputData
            return
        }

        //发放奖励
        propItem := &battery.PropItem{
            Type: battery.PropType_PROP_PACKAGE.Enum(),
            //Id:     config.PropID,
            Amount: proto.Uint32(1),
        }

        if isDynamic {
            propItem.Id = mailItemTmp.PropId
            err = api.GainProp(uid, nil, propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
            if err != xyerror.ErrOK {
                xylog.Error(uid, "GainProp failed : %v", err)
                return
            }
        } else {
            if config.GetPropID() != xybusinesscache.INVALID_PROPID {
                propItem.Id = config.PropID
                err = api.GainProp(uid, nil, propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
                if err != xyerror.ErrOK {
                    xylog.Error(uid, "GainProp failed : %v", err)
                    return
                }
            }
        }

        //然后再将邮件设置为已确认
        mailItemTmp.IsConfirm = proto.Bool(true) //设置为邮件已经确认
        err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).UpsertMailInfo(tempSystemMailList)
    }
    return
}

// isDynamic 判断邮件是否是动态邮件
func (api *XYAPI) isDynamic(uid string, id int32) bool {
    sysmailBase := id / SYSMAIL_SEGMENT
    xylog.Debug(uid, "sysmailBase %d", sysmailBase)
    if sysmailBase == SYSMAIL_TYPE_DYNAMIC || //运营动态邮件
        sysmailBase == SYSMAIL_TYPE_WEIBO_FIRSTLOGIN_DYNAMIC || //微博首次登录动态邮件
        sysmailBase == SYSMAIL_TYPE_GUEST_FIRSTLOGIN_DYNAMIC { //游客首次登录动态邮件
        return true
    }
    return false
}

////直接删除邮件 //如果是静态配置表里的系统通知和公告，不允许被删除
////func (api *XYAPI) RemoveSystemMail(uid string, mailID int32) (err error) {
////	err = api.GetDB().RemoveSystemMail(uid, mailID)
////	return
////}

//// 增加系统登录奖励邮件
//func (api *XYAPI) AddLoginGiftMail(uid string) (err error) {

//	timeNow := time.Now()

//	// 查询当天的登陆奖励配置信息
//	var configs []*battery.DBMailInfoConfig
//	configs, err = api.getTodayLoginMailConfig(uid, timeNow)
//	if err != xyerror.ErrOK {
//		return
//	}

//	if len(configs) > 0 { //存在当天的登录奖励配置
//		change := false

//		//查找玩家系统邮件信息
//		var tempSystemMailList *battery.DBSystemMailListItem
//		tempSystemMailList, change, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).GetSystemMailListWithoutChange(uid)
//		if err != xyerror.ErrOK {
//			return
//		}

//		//当天没有登录奖励记录，增加当天登录奖励的邮件和公告
//		if !api.IsHasTodayLoginGift(uid, tempSystemMailList) {

//			xylog.Debug(uid, "[%s] IsHasTodayLoginGift : false", uid)
//			//增加公告信息
//			//announcements := make([]interface{}, 0)
//			for _, config := range configs {

//				//增加玩家系统邮件
//				userSystemMail := &battery.DBUserSystemMail{
//					MailId:    config.MailID,
//					MailType:  config.Mailtype,
//					BeginTime: config.StartTime,
//					EndTime:   config.EndTime,
//					IsRead:    proto.Bool(false),
//				}
//				tempSystemMailList.MailInfoList = append(tempSystemMailList.MailInfoList, userSystemMail)
//				change = true
//			}

//			//更新玩家的系统邮件信息
//			if change {
//				err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).UpsertMailInfo(tempSystemMailList)
//				if err != xyerror.ErrOK {
//					return
//				}
//			} else {
//				xylog.Debug(uid, "[%s] SystemMailList not change", uid)
//			}

//		} else {
//			xylog.Debug(uid, "[%s] IsHasTodayLoginGift : true", uid)
//		} //if !hasGetConfigMailIDList(uid, tempSystemMailList)

//	} //if len(configs) > 0

//	return
//}

//// IsHasTodayLoginGift 判断当天是否已经有登录奖励通知
//// uid string 玩家id
//// sysMails *battery.DBSystemMailListItem 玩家的系统邮件信息
////return:
//// true 已经有登录奖励邮件，false 没有当天登录奖励邮件
//func (api *XYAPI) IsHasTodayLoginGift(uid string, sysMails *battery.DBSystemMailListItem) bool {

//	for _, mailItem := range sysMails.MailInfoList {
//		id := mailItem.GetMailId()

//		//先判断邮件类型是否为登录奖励邮件
//		if id/1000 != 111 || mailItem.GetMailType() != battery.SystemMailType_SystemMailType_Gift {
//			continue
//		}

//		//获取邮件的配置信息
//		config := xybusinesscache.DefMailConfigCacheManager.MailConfig(id)
//		if nil == config || config.GetMailtype() != battery.SystemMailType_SystemMailType_Gift {
//			xylog.Debug(uid, "[%s] failed to find MailConfig of %d mailType %v,pls check.", uid, id, config.GetMailtype())
//			continue
//		}

//		//判断是否是当天的邮件
//		startTimeValue := mailItem.GetBeginTime()
//		startTime := time.Unix(startTimeValue, 0)
//		nowTime := time.Now()
//		if config.GetPropID() >= 0 && startTime.Year() == nowTime.Year() && startTime.Month() == nowTime.Month() && startTime.Day() == nowTime.Day() {
//			return true
//		}
//	}

//	return false
//}

////查找当天的登录奖励配置
//func (api *XYAPI) getTodayLoginMailConfig(uid string, timeNow time.Time) (configs []*battery.DBMailInfoConfig, err error) {

//	//从缓存中获取当天的登录奖励配置信息
//	mapMailConfig := xybusinesscache.DefMailConfigCacheManager.MailConfigs()
//	if mapMailConfig == nil {
//		xylog.Error(uid, "[%s] mapMailConfig from cache is nil", uid)
//		err = xyerror.ErrQueryMailConfigsFromCacheError
//		return
//	} else {

//		for _, mailConfig := range *mapMailConfig {
//			mailId := mailConfig.GetMailID()
//			startTime := time.Unix(mailConfig.GetStartTime(), 0)
//			xylog.Debug(uid, "startTime : %v", startTime)
//			// ToDo: 下面这两个魔数需要去掉，或者这个过滤条件需要去掉
//			if mailId/1000 == 111 &&
//				mailConfig.GetPropID() >= xybusinesscache.INVALID_PROPID &&
//				startTime.Year() == timeNow.Year() &&
//				startTime.Month() == timeNow.Month() &&
//				startTime.Day() == timeNow.Day() &&
//				mailConfig.GetEndTime() > timeNow.Unix() &&
//				mailConfig.GetStartTime() <= timeNow.Unix() { //这么一堆组成了登录奖励的查询条件 :=()
//				configs = append(configs, mailConfig)
//			}
//		}

//		xylog.Debug(uid, "[%s] GetTodayLoginMailConfig : %v", uid, configs)
//	}

//	return
//}

//查询邮件配置
func (api *XYAPI) MailConfigDetail(uid string, mailid int32) (mailconfig *battery.DBMailInfoConfig, err error) {
    mailconfig = xybusinesscache.DefMailConfigCacheManager.MailConfig(mailid)
    if nil == mailconfig { //没找到
        xylog.Error(uid, "MailConfigDetail of %d failed", mailid)
        err = xyerror.ErrQueryMailConfigsFromCacheError
    }
    return
}

// 动态系统邮件类型
const (
    SYSMAIL_TYPE_DYNAMIC                  int32 = 8  //动态邮件类型
    SYSMAIL_TYPE_WEIBO_FIRSTLOGIN_DYNAMIC int32 = 9  //微博首次登录动态邮件类型
    SYSMAIL_TYPE_GUEST_FIRSTLOGIN_DYNAMIC int32 = 10 //游客首次登录动态邮件类型
)

const (
    SYSMAIL_SEGMENT int32 = 100000 //动态邮件号段
)

// 增加运营系统邮件
// uid string 玩家标识
// propId uint64 道具标识
// propType battery.PropType 道具类型
// amount int32 数目
// err error 返回错误
func (api *XYAPI) addMaintenanceSysmail(uid string,
    propId uint64,
    propType battery.PropType,
    amount int32,
    systemMailBaseId int32) (err error) {

    xylog.Debug(uid, "addMaintenanceSysmail %v %v %v", propId, propType, amount)

    // 查找当前玩家对应的sysmail文档
    var systemMailList *battery.DBSystemMailListItem
    systemMailList, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).GetSystemMailList(uid)
    if err != xyerror.ErrOK {
        if err == xyerror.ErrNotFound { // 如果没有找到，则创建一条
            systemMailList = &battery.DBSystemMailListItem{
                Uid:          proto.String(uid),
                MailInfoList: make([]*battery.DBUserSystemMail, 0, 1),
            }
            err = xyerror.ErrOK
        } else { // 如果是数据库其他错误，直接返回
            xylog.Error(uid, "GetSystemMailList failed : %v", err)
            return
        }
    }

    xylog.Debug(uid, "systemMailList : %v", systemMailList)

    // 填入新邮件
    now := xyutil.CurTimeSec()
    //如果计数超过上限，则清0重算
    if systemMailList.GetDynamicCounter() >= SYSMAIL_SEGMENT {
        systemMailList.DynamicCounter = proto.Int32(0)
    }
    systemMailList.DynamicCounter = proto.Int32(systemMailList.GetDynamicCounter() + 1) //计数器加一

    systemMail := &battery.DBUserSystemMail{
        MailId:       proto.Int32(systemMailBaseId + systemMailList.GetDynamicCounter()),
        SourceMailId: proto.Int32(systemMailBaseId),
        MailType:     battery.SystemMailType_SystemMailType_Gift.Enum(),
        SourceType:   battery.SystemMailSourceType_SystemMailSourceType_Dynamic.Enum(),
        BeginTime:    proto.Int64(now),
        EndTime:      proto.Int64(now + int64(time.Hour*24*7/time.Second)), //7天后过期
        IsRead:       proto.Bool(false),
        IsConfirm:    proto.Bool(false),
        PropId:       proto.Uint64(propId),
    }

    systemMailList.MailInfoList = append(systemMailList.MailInfoList, systemMail)

    xylog.Debug(uid, "systemMailList : %v", systemMailList)

    //更新玩家系统邮件
    err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST).UpsertMailInfo(systemMailList)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "UpsertMailInfo failed : %v", err)
        return
    }

    return
}
