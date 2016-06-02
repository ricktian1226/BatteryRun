package batteryapi

import (
    "time"

    "code.google.com/p/goprotobuf/proto"

    "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationSharedQuery(req *battery.SharedQueryRequest, resp *battery.SharedQueryResponse) (err error) {
    var (
        platform         = req.GetPlatform()
        now              int64
        uid              = req.GetUid()
        activeShareInfos = make([]*battery.DBUserShareInfo, 0) // 需要返回的分享活动信息
        mapID2exit       = make(Set, 0)
        counter          int32
        dailyCounter     int32
        id               uint32
    )
    api.SetDB(platform)
    resp.Error = xyerror.DefaultError()
    // 查询玩家分享活动信息
    userShareInfos, err := api.QueryAllShareInfo(uid)
    if err != xyerror.ErrOK {
        if err == xyerror.ErrNotFound {
            xylog.ErrorNoId("Query useShareInfo not find")
        } else {
            xylog.ErrorNoId("usershareinfo failed :%v", err)
            resp.Error.Code = battery.ErrorCode_QueryUserShareRecordError.Enum()
            return
        }
    }
    xylog.DebugNoId("shareinfos %v", userShareInfos)
    // 查找需要添加的分享活动
    for _, shareinfo := range userShareInfos {
        mapID2exit[shareinfo.GetId()] = empty{}
        if shareinfo.GetState() == battery.UserMissionState_UserMissionState_Active {
            activeShareInfos = append(activeShareInfos, shareinfo)
        }
    }

    now = time.Now().Unix()
    shareActivitys, err := api.queryShareActivity(uid, now, mapID2exit)
    if err != xyerror.ErrOK {
        resp.Error.Code = battery.ErrorCode_QueryShareActivityError.Enum()
        return
    }
    xylog.DebugNoId("shareActivity %v", shareActivitys)
    for _, shareActivity := range shareActivitys {
        newShareInfo := api.defaultShareInfo(uid, now, shareActivity)
        err = api.addUserShareInfo(newShareInfo)
        if err == xyerror.ErrOK {
            activeShareInfos = append(activeShareInfos, newShareInfo)
        }
    }
    xylog.DebugNoId("userShareInfo %v", activeShareInfos)
    for _, shareInfo := range activeShareInfos {
        counter = shareInfo.GetCounter()
        dailyCounter = shareInfo.GetDailyCounter()
        lastTimeStamp := shareInfo.GetLastResetTimestamp()
        diff := xyutil.DayDiff(lastTimeStamp, now)
        if diff > 0 {
            // 隔天重置每日分享次数
            dailyCounter = 0
        }
        id = shareInfo.GetId()
        // 查询分享奖励信息
        giftlist := api.QuerySharedGift(uid, counter+1, id) // +1获取下次分享的礼包
        // 目前自有一个活动，直接覆盖掉，以后有需要在列表进行保存
        resp.Units = giftlist
        resp.Counter = &counter
        resp.DailyCounter = &dailyCounter
        resp.DailyLimit = shareInfo.DailyLimit
    }

    resp.Uid = req.Uid
    resp.Error = xyerror.DefaultError()
    return
}

func (api *XYAPI) addUserShareInfo(share *battery.DBUserShareInfo) (err error) {
    return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSHARE).AddUserInfo(share)
}
func (api *XYAPI) defaultShareInfo(uid string, now int64, shareinfo *xybusinesscache.ShareActivity) *battery.DBUserShareInfo {
    return &battery.DBUserShareInfo{
        Uid:                proto.String(uid),
        Id:                 proto.Uint32(shareinfo.Id),
        Type:               shareinfo.ShareType.Enum(),
        DailyCounter:       proto.Int32(0),
        Counter:            proto.Int32(0),
        LastResetTimestamp: proto.Int64(0),
        State:              battery.UserMissionState_UserMissionState_Active.Enum(),
        DailyLimit:         proto.Int32(shareinfo.DailyLimit),
        Limit:              proto.Int32(shareinfo.GoalValue),
        BeginTime:          proto.Int64(shareinfo.AwardsStartTime),
        EndTime:            proto.Int64(shareinfo.AwardsEndTime),
        Restart:            proto.Bool(shareinfo.Restart),
    }

}
func (api *XYAPI) OperationSharedRequest(req *battery.SharedRequest, resp *battery.SharedResponse) (err error) {
    var (
        platform       = req.GetPlatform()
        uid            = req.GetUid()
        now            = time.Now().Unix()
        counter        int32
        dailyCounter   int32
        restart        bool
        id             uint32
        userShareInfos []*battery.DBUserShareInfo
    )
    xylog.DebugNoId("%s share request", uid)
    api.SetDB(platform)
    resp.Uid = req.Uid
    resp.Error = xyerror.DefaultError()
    // // 加载所有活动信息
    // activitys := xybusinesscache.DefSharedActivityManager.GetAllShareActivity()
    // if activitys == nil {
    //     xylog.Error(uid, "DefSharedGiftManager .Uints is nil for %v", counter)
    //     resp.Error.Code = battery.ErrorCode_QueryShareActivityError.Enum()
    //     return
    // }
    // 遍历活动发放奖励
    userShareInfos, err = api.QueryAllShareInfo(uid)
    if err != xyerror.ErrOK {
        if err == xyerror.ErrNotFound {
            xylog.ErrorNoId("Query useShareInfo not find")
        } else {
            xylog.ErrorNoId("usershareinfo failed :%v", err)
            resp.Error.Code = battery.ErrorCode_QueryUserShareRecordError.Enum()
            return
        }
    }
    for _, userShareInfo := range userShareInfos {
        if userShareInfo.GetBeginTime() > now || (userShareInfo.GetEndTime() < now) {
            xylog.DebugNoId("activity not in active time skip it", userShareInfo)
            continue
        }
        restart = userShareInfo.GetRestart()
        id = userShareInfo.GetId()
        // 活动非激活状态直接跳过
        if userShareInfo.GetState() != battery.UserMissionState_UserMissionState_Active {
            continue
        }
        counter = userShareInfo.GetCounter()
        dailyCounter = userShareInfo.GetDailyCounter()
        lastResetTimeStamp := userShareInfo.GetLastResetTimestamp()
        diff := xyutil.DayDiff(lastResetTimeStamp, now)

        // 每日活动
        if userShareInfo.GetType() == battery.SHARE_TYPE_SHARE_TYPE_DAILY {
            if diff > 0 {
                userShareInfo.Counter = proto.Int32(1)
                userShareInfo.LastResetTimestamp = proto.Int64(now)
            } else {
                userShareInfo.Counter = proto.Int32(counter + 1)
            }

        }

        // 累计分享
        if userShareInfo.GetType() == battery.SHARE_TYPE_SHARE_TYPE_TOTAL {

            if diff > 0 {
                userShareInfo.Counter = proto.Int32(counter + 1)
                // 跨天重置当日分享次数和时间戳
                userShareInfo.DailyCounter = proto.Int32(1)
                userShareInfo.LastResetTimestamp = proto.Int64(now)
            } else {
                if dailyCounter >= userShareInfo.GetDailyLimit() {
                    continue
                } else {
                    userShareInfo.Counter = proto.Int32(counter + 1)
                    userShareInfo.DailyCounter = proto.Int32(dailyCounter + 1)
                }
            }
            // 活动是否可以重置
            if userShareInfo.GetCounter() == userShareInfo.GetLimit() {
                if restart {
                    userShareInfo.Counter = proto.Int32(0)
                } else {
                    userShareInfo.State = battery.UserMissionState_UserMissionState_Reset.Enum()
                }

            }
        }
        api.GainGift(uid, counter+1, id, resp.Error)
        api.updateShareInfo(uid, id, userShareInfo)

    }

    return
}

func (api *XYAPI) QuerySharedGift(uid string, counter int32, id uint32) (giftlist []*battery.SharedUnit) {
    // 从缓存加载配置信息
    giftlist = make([]*battery.SharedUnit, 0)
    activity := xybusinesscache.DefSharedActivityManager.GetShareActivityByID(id)
    if activity == nil {
        xylog.Error(uid, "DefSharedGiftManager .Uints is nil for %v", id)
        return
    }

    items := activity.Items
    for _, item := range items {
        if counter == item.GetCounter() {
            gift := &battery.SharedUnit{
                Counter: item.Counter,
                Items:   item.Award,
            }
            giftlist = append(giftlist, gift)
        }

    }

    xylog.Debug(uid, "QuerySharedGift %v", giftlist)
    return
}

func (api *XYAPI) QueryAllShareInfo(uid string) (userShareInfo []*battery.DBUserShareInfo, err error) {
    // 查询玩家分享次数
    userShareInfo = make([]*battery.DBUserShareInfo, 0)
    err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSHARE).QueryAllShareInfo(uid, &userShareInfo)

    return
}

func (api *XYAPI) QueryShareInfo(uid string, id uint32) (userShareInfo *battery.DBUserShareInfo, err error) {
    userShareInfo = new(battery.DBUserShareInfo)
    err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSHARE).QueryShareInfo(uid, id, userShareInfo)
    return
}
func (api *XYAPI) GainGift(uid string, counter int32, id uint32, errStruct *battery.Error) {
    var (
        propItems []*battery.PropItem
    )
    shareActivity := xybusinesscache.DefSharedActivityManager.GetShareActivityByID(id)
    items := shareActivity.Items
    for _, shareItem := range items {
        if shareItem.GetCounter() == counter {
            propItems = shareItem.GetAward()
            break
        }
    }
    switch shareActivity.DispenseType {
    case battery.DISPENSE_TYPE_DISPENSE_TYPE_DIRECT:
        for _, item := range propItems {
            api.maintenanceAddProp(
                uid,
                item.GetId(),
                item.GetType(),
                int32(item.GetAmount()),
                errStruct)
        }
    case battery.DISPENSE_TYPE_DISPENSE_TYPE_SYSTEMMAIL:
        for _, item := range propItems {
            api.maintenanceAddSystemmail(uid,
                item.GetId(),
                item.GetType(),
                int32(item.GetAmount()),
                shareActivity.MailId,
                errStruct)
        }
    }

    return
}

func (api *XYAPI) updateShareInfo(uid string, id uint32, shareInfo *battery.DBUserShareInfo) (err error) {

    return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSHARE).UpsertUserInfo(uid, id, shareInfo)
}

func (api *XYAPI) queryShareActivity(uid string, now int64, expectid Set) (shareActivitys []*xybusinesscache.ShareActivity, err error) {
    shareActivitys = make([]*xybusinesscache.ShareActivity, 0)
    activitys := xybusinesscache.DefSharedActivityManager.GetAllShareActivity()
    if 0 == len(activitys) {
        xylog.DebugNoId("query shareActivity from cache fail")
        err = xyerror.ErrQueryShareActivityFromCacheError
        return
    }
    for _, activity := range activitys {
        xylog.DebugNoId("expectid %v %v", expectid, activity.Id)
        if _, ok := expectid[activity.Id]; ok {
            xylog.DebugNoId("activity %v exit,skip it", activity)
            continue
        }
        if activity.StartTime > now || (activity.EndTime < now) {
            xylog.DebugNoId("activitys %v not in active time", activity)
            continue
        }

        shareActivitys = append(shareActivitys, activity)
    }
    return
}
