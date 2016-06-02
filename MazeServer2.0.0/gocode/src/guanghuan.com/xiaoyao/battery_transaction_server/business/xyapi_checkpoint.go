// xyapi_checkpoint
package batteryapi

import (
    "code.google.com/p/goprotobuf/proto"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    xylog "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/performance"
    xyutil "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xycache "guanghuan.com/xiaoyao/superbman_server/cache"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/money"
    "guanghuan.com/xiaoyao/superbman_server/server"

    //"sort"
    "fmt"
    "time"
)

//查询玩家记忆点信息
func (api *XYAPI) OperationQueryUserCheckPoints(req *battery.QueryUserCheckPointsRequest, resp *battery.QueryUserCheckPointsResponse) (err error) {
    var (
        uid               = req.GetUid()
        checkPointBeginId = req.GetBeginId()
        checkPointEndId   = req.GetEndId()
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    //初始化resp
    resp.Uid = req.Uid
    resp.BeginId = req.BeginId
    resp.EndId = req.EndId
    resp.Error = xyerror.DefaultError()

    //查询玩家最大的记忆点id
    dbUserCheckPoints := []*battery.DBUserCheckPoint{}
    err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPoints(uid, checkPointBeginId, checkPointEndId, &dbUserCheckPoints)
    if xyerror.ErrOK == err {
        for _, dbUserCheckPoint := range dbUserCheckPoints {
            checkPoint := xycache.GetCheckPointFromDBUserCheckPoint(dbUserCheckPoint)
            checkPoint.Sid = nil //把uid置为空，请求单个玩家的记忆点信息就不要每个都存一份uid了
            resp.CheckPoints = append(resp.CheckPoints, checkPoint)
        }

        xylog.Debug(uid, "QueryUserCheckPoints result : %v", resp.CheckPoints)

    } else {
        xylog.Error(uid, "QueryUserCheckPoints failed : %v", err)
        resp.Error.Code = battery.ErrorCode_QueryUserCheckPointError.Enum()
    }

    return
}

//查询记忆点排行版（好友排行榜或者全局排行榜）
func (api *XYAPI) OperationQueryCheckPointDetail(req *battery.QueryUserCheckPointDetailRequest, resp *battery.QueryUserCheckPointDetailResponse) (err error) {
    var (
        uid          = req.GetUid()
        checkPointId = req.GetCheckPointId()
        sids         = req.GetSids()
        idSource     = req.GetSource()
        rankType     = req.GetRankType()
        errStruct    = xyerror.DefaultError()
    )

    // 获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    //初始化resp
    resp.Uid = req.Uid
    resp.CheckPointId = req.CheckPointId
    resp.RankType = req.RankType

    //xylog.Debug(uid, "req : %v", req)

    switch rankType {
    case battery.CheckPointRankType_CheckPointRankType_Friend: //查询玩家记忆点对应的好友排行榜
        api.queryCheckPointFriendRank(uid, sids, idSource, checkPointId, &(resp.Rank), errStruct)
    case battery.CheckPointRankType_CheckPointRankType_Global: //查询玩家记忆点对应的全局排行榜
        api.queryCheckPointGlobalRank(uid, checkPointId, &(resp.Rank), platform, errStruct)
    }
    if errStruct.GetCode() != battery.ErrorCode_NoError {
        goto ErrHandle
    }

ErrHandle:

    resp.Error = errStruct

    //if resp.Error.GetCode() != battery.ErrorCode_NoError {
    //	xylog.Error(resp.Error.GetDesc())
    //} else {
    //	xylog.Debug(resp.Error.GetDesc())
    //}

    return
}

// //提交记忆点分数信息
// func (api *XYAPI) OperationCommitCheckPoint(req *battery.CommitCheckPointRequest, resp *battery.CommitCheckPointResponse) (err error) {
//     var (
//         uid          = req.GetUid()
//         checkPointId = req.GetCheckPointId()
//         gameId       = req.GetGameId()
//         score        = req.GetScore()
//         charge       = req.GetCharge()
//         coin         = req.GetCoin()
//         roleId       = req.GetRoleId()
//         bChange      = false
//         errStruct    = xyerror.DefaultError()
//     )

//     //获取请求的终端平台类型
//     platform := req.GetPlatformType()
//     api.SetDB(platform)

//     //初始化resp
//     resp.Uid = req.Uid
//     resp.CheckPointId = req.CheckPointId
//     resp.GameId = req.GameId

//     //查询玩家记忆点信息
//     dbDetail := &battery.DBUserCheckPoint{
//         Uid:          req.Uid,
//         CheckPointId: req.CheckPointId,
//         Score:        proto.Uint64(0),
//         Charge:       proto.Uint64(0),
//         Coin:         proto.Uint32(0),
//     }
//     err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPointDetail(uid, checkPointId, dbDetail)
//     if err != xyerror.ErrOK {
//         //没找到，说明是新的checkpoint
//         if err != xyerror.ErrNotFound {
//             errStruct.Code = battery.ErrorCode_QueryUserCheckPointError.Enum()
//             //errStruct.Desc = proto.String(fmt.Sprintf("[%s] QueryUserCheckPointDetail for checkPointId(%d) failed : %v", uid, checkPointId, err.Error()))
//             goto ErrHandle
//         }

//     }

//     if score > dbDetail.GetScore() {
//         dbDetail.Score = proto.Uint64(score)
//         dbDetail.RoleId = proto.Uint64(roleId)
//         bChange = true
//     }

//     if charge > dbDetail.GetCharge() {
//         dbDetail.Charge = proto.Uint64(charge)
//         bChange = true
//     }

//     if coin > dbDetail.GetCoin() {
//         dbDetail.Coin = proto.Uint32(coin)
//         bChange = true
//     }

//     //需要刷新的话，就刷新一把吧
//     if bChange {
//         err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).UpsertUserCheckPoint(dbDetail)
//         if err != xyerror.ErrOK {
//             errStruct.Code = battery.ErrorCode_CommitUserCheckPointDetailError.Enum()
//             goto ErrHandle
//         }
//     }

// ErrHandle:
//     resp.Error = errStruct
//     //记录记忆点提交日志
//     log := api.defaultCheckPointLog(uid, gameId, checkPointId, coin, score, charge, roleId, bChange)
//     go api.addCheckPointLog(log)

//     return
// }

// 钻石解锁关卡
func (api *XYAPI) OperationCheckPointUnlock(req *battery.CheckPointUnlockRequest, resp *battery.CheckPointUnlockResponse) (err error) {
    var (
        checkPointId      = req.GetCheckPointId()
        platform          = req.GetPlatform()
        uid               = req.GetUid()
        errStruct         = xyerror.DefaultError()
        missionTypes      = make([]battery.MissionType, 0)
        isFinish     bool = true
        collections       = make([]uint32, 0)
        goodid       uint64
        dbDetail     *battery.DBUserCheckPoint
    )

    resp.Platform = req.Platform
    resp.Uid = req.Uid
    resp.CheckPointId = req.CheckPointId
    resp.Error = xyerror.DefaultError()

    userAccomplishment := &battery.DBUserAccomplishment{}

    gameResult := make([]*battery.Quota, 0)
    farthestCheckpoint := new(battery.Quota)

    xylog.DebugNoId("unlock checkpoint %v", checkPointId)
    api.SetDB(platform)
    //  由解锁的关卡id获取需要购买的物品id

    checkPoint := &battery.DBUserCheckPoint{}
    err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPointDetail(uid, checkPointId, checkPoint)
    if err == xyerror.ErrOK {
        // 关卡已存在，解锁失败
        errStruct.Code = battery.ErrorCode_UnlockCheckPointError.Enum()
        xylog.ErrorNoId("unlock checkpoint %v fail,checkpoint has exit", checkPointId)
        goto ErrHandle
    }

    goodid = xycache.DefCheckPointUnlockGoodsManager.GoodsId(checkPointId)
    api.BuyGoods(uid, "", goodid, errStruct)
    if errStruct.GetCode() != battery.ErrorCode_NoError {
        xylog.ErrorNoId("goodsid not find")
        goto ErrHandle
    }

    // 增加玩家记忆点信息，解锁关卡只记录关卡id无游戏数据

    dbDetail = &battery.DBUserCheckPoint{
        Uid:              req.Uid,
        CheckPointId:     req.CheckPointId,
        Score:            proto.Uint64(0),
        Charge:           proto.Uint64(0),
        Coin:             proto.Uint32(0),
        RoleId:           proto.Uint64(0),
        Collections:      collections,
        CollectionsCount: proto.Uint32(0),
        Grade:            proto.Uint32(0),
        PlatformType:     req.Platform,
    }
    // 更新检查点
    err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).UpsertUserCheckPoint(dbDetail)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_CommitUserCheckPointDetailError.Enum()
        goto ErrHandle
    }

    // 更新玩家成就信息（只更新最远关卡信息）

    api.getUserAccomplishment(uid, userAccomplishment, errStruct)
    if errStruct.GetCode() != battery.ErrorCode_NoError {
        goto ErrHandle
    }

    farthestCheckpoint.Value = proto.Uint64(uint64(req.GetCheckPointId()))
    farthestCheckpoint.Id = battery.QuotaEnum_Quota_FarthestCheckPoint.Enum()
    gameResult = append(gameResult, farthestCheckpoint)

    if api.updateUserCheckPointTotalAccomplishment(uid, userAccomplishment, gameResult) { // 更新关卡指标
        api.updateUserAccomplishment(uid, userAccomplishment, gameResult, errStruct) // 更新所有指标
        userAccomplishment.Uid = proto.String(uid)
        xylog.DebugNoId("user accomplishment :%v", userAccomplishment)
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERACCOMPLISHMENT).UpdateUserAccomplishment(userAccomplishment)
        if err != xyerror.ErrOK {
            errStruct.Code = battery.ErrorCode_UpdateUserAccomplishmentError.Enum()
            goto ErrHandle
        }
    }

    missionTypes = append(missionTypes, battery.MissionType_MissionType_MainLine)
    missionTypes = append(missionTypes, battery.MissionType_MissionType_Daily)

    err = api.updateUserMissionsQuotas(uid, missionTypes, gameResult, time.Now().Unix(), isFinish)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "updateUserMissionsQuotas failed")
    }
ErrHandle:
    resp.Error = errStruct
    resp.Wallet, err = xymoney.QueryWallet(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT))
    return
}

//刷新玩家的记忆点信息
// uid string 玩家id
// roleId uint64 角色id
// checkPointId uint32 记忆点id
// gameResultQuotas []*battery.Quota 游戏结算指标列表
//return:
// err error 错误信息
func (api *XYAPI) updateUserCheckPoint(uid, gameId string, roleId uint64, checkPointId uint32, collections []uint32, gameResultQuotas, tmpGameResultQuotas *[]*battery.Quota, userAccomplishment *battery.DBUserAccomplishment, errStruct *battery.Error, platform battery.PLATFORM_TYPE) (isAccomplishmentChange bool) {

    begin := time.Now()
    defer xyperf.Trace(LOGTRACE_UPDATEUSERCHECKPOINT, &begin)

    var (
        score, charge     uint64
        coin, grade       uint32
        bChangeCheckpoint = false
        setCollections    = make(Set, 0)
    )

    //获取记忆点需要记录的任务指标数据
    for _, quota := range *tmpGameResultQuotas {
        switch quota.GetId() {
        case battery.QuotaEnum_Quota_Score:
            score = quota.GetValue()
        case battery.QuotaEnum_Quota_Coin:
            coin = uint32(quota.GetValue())
        case battery.QuotaEnum_Quota_Charge:
            charge = quota.GetValue()
        case battery.QuotaEnum_Quota_Grade:
            grade = uint32(quota.GetValue())
        }
    }

    //查询玩家记忆点信息
    dbDetail := &battery.DBUserCheckPoint{
        Uid:              &uid,
        CheckPointId:     &checkPointId,
        Score:            proto.Uint64(0),
        Charge:           proto.Uint64(0),
        Coin:             proto.Uint32(0),
        Grade:            proto.Uint32(0),
        Collections:      make([]uint32, 0),
        CollectionsCount: new(uint32),
    }
    xylog.DebugNoId("collection %p", dbDetail.CollectionsCount)
    err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPointDetail(uid, checkPointId, dbDetail)
    xylog.DebugNoId("dbdetail %v", dbDetail)
    if err != xyerror.ErrOK {
        //没找到，说明是新的checkpoint
        if err == xyerror.ErrNotFound {
            bChangeCheckpoint = true
        } else {
            errStruct.Code = battery.ErrorCode_QueryUserCheckPointError.Enum()
            errStruct.Desc = proto.String(fmt.Sprintf("[%s] QueryUserCheckPointDetail for checkPointId(%d) failed : %v", uid, checkPointId, err.Error()))
            goto ErrHandle
        }
    }
    dbDetail.PlatformType = platform.Enum()
    if score > dbDetail.GetScore() {
        dbDetail.Score = proto.Uint64(score)
        dbDetail.RoleId = proto.Uint64(roleId)
        bChangeCheckpoint = true
    }

    if charge > dbDetail.GetCharge() {
        dbDetail.Charge = proto.Uint64(charge)
        bChangeCheckpoint = true
    }

    if coin > dbDetail.GetCoin() {
        dbDetail.Coin = proto.Uint32(coin)
        bChangeCheckpoint = true
    }

    if grade > dbDetail.GetGrade() {
        dbDetail.Grade = proto.Uint32(grade)
        bChangeCheckpoint = true
    }
    //将新的收集物添加到记忆点的收集物列表中
    for _, c := range dbDetail.GetCollections() {
        setCollections[c] = empty{}
    }
    for _, collection := range collections {
        if _, ok := setCollections[collection]; ok {
            continue
        } else {
            dbDetail.Collections = append(dbDetail.Collections, collection)
            xylog.DebugNoId("collection %p", dbDetail.CollectionsCount)
            (*(dbDetail.CollectionsCount))++
            setCollections[collection] = empty{} //前台插入的要重新
            bChangeCheckpoint = true
        }
    }

    xylog.Debug(uid, "bChangeCheckpoint(%t)", bChangeCheckpoint)
    if bChangeCheckpoint { //刷新玩家对应的记忆点信息
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).UpsertUserCheckPoint(dbDetail)
        if err != xyerror.ErrOK {
            errStruct.Code = battery.ErrorCode_CommitUserCheckPointDetailError.Enum()
            xylog.Error(uid, "UpsertUserCheckPoint(%v) failed : %v", dbDetail, err)
            goto ErrHandle
        }
    }

    //至少有QuotaEnum_Quota_FinishGame指标需要刷新，只要记忆点成功，这个指标就需要刷新
    isAccomplishmentChange, err = api.setUserCheckPointAccomplishment(uid, userAccomplishment, gameResultQuotas, tmpGameResultQuotas, checkPointId, platform)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_CommitUserCheckPointDetailError.Enum()
        xylog.Error(uid, "setUserNOZeroCheckPoints(%v) failed : %v", dbDetail, err)
        goto ErrHandle
    }

ErrHandle:
    return
}

//生成记忆点日志结构体信息
// uid string 玩家id
// gameId string 游戏id
// checkPointId uint32 记忆点id
// score uint64 分数
// charge uint64 charge数
// best bool 是否最佳
func (api *XYAPI) defaultCheckPointLog(uid, gameId string, checkPointId, coin uint32, score, charge, roleId uint64, best bool) *battery.CheckPointLog {
    now := xyutil.CurTimeSec()
    return &battery.CheckPointLog{
        Uid:          proto.String(uid),
        CheckPointId: proto.Uint32(checkPointId),
        GameId:       proto.String(gameId),
        Score:        proto.Uint64(score),
        Charge:       proto.Uint64(charge),
        Coin:         proto.Uint32(coin),
        RoleId:       proto.Uint64(roleId),
        Best:         proto.Bool(best),
        Opdate:       proto.Int64(now),
        Opdatestr:    proto.String(xyutil.ToStrTime(now)),
    }
}

//增加记忆点提交日志
func (api *XYAPI) addCheckPointLog(log *battery.CheckPointLog) {
    api.GetLogDB().AddCheckPointLog(log)
}

//查询玩家的记忆点数据信息
// uid string 玩家id
// checkPointId uint32 记忆点id
// detail *battery.UserCheckPointDetail 玩家记忆点数据信息地址
func (api *XYAPI) queryCheckPointDetail(uid string, checkPointId uint32, detail *battery.UserCheckPointDetail) (errStruct *battery.Error) {
    errStruct = &battery.Error{}
    dbDetail := &battery.DBUserCheckPoint{}
    err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPointDetail(uid, checkPointId, dbDetail)
    if xyerror.ErrOK != err {
        if err == xyerror.ErrNotFound {
            errStruct.Code = battery.ErrorCode_DBNoFoundError.Enum()
            //errStruct.Desc = proto.String(fmt.Sprintf("[%s] QueryUserCheckPointDetail for checkPointId(%d) failed : %v", uid, checkPointId, err.Error))
        } else {
            errStruct.Code = battery.ErrorCode_QueryUserCheckPointError.Enum()
            //errStruct.Desc = proto.String(fmt.Sprintf("[%s] QueryUserCheckPointDetail for checkPointId(%d) failed : %v", uid, checkPointId, err.Error))
        }
    }

    return
}

//查询玩家所有记忆点（不包含零号记忆点）的各种参数和
// uid string 玩家id
// dbDetail *battery.DBUserCheckPoint 返回的记忆点
func (api *XYAPI) queryUserCheckPointsSum(uid string, dbDetail *battery.DBUserCheckPoint) (err error) {
    return api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPointsSum(uid, dbDetail)
}

//重新设置玩家0号记忆点的数据，0号记忆点保存的是玩家所有记忆点数据的和(零号就是关键，保存的是关键数据——纯属调侃)
// uid string 玩家id
// dbUserCheckPoint *battery.DBUserCheckPoint
func (api *XYAPI) setUserCheckPointAccomplishment(uid string, userAccomplishment *battery.DBUserAccomplishment, gameResultQuotas, tmpGameResultQuotas *[]*battery.Quota, checkPointId uint32, platform battery.PLATFORM_TYPE) (isChange bool, err error) {

    //查找玩家所有记忆点的指标和
    dbDetail := &battery.DBUserCheckPoint{}
    err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPointsSum(uid, dbDetail)
    if err != xyerror.ErrOK {
        if err != xyerror.ErrNotFound {
            xylog.Error(uid, "QueryUserCheckPointsSum failed : %v", err)
            return
        }
        err = xyerror.ErrOK
    }
    // 更新排行榜数据
    api.updateUserRankList(uid, dbDetail, platform)
    //将玩家所有记忆点的指标和刷新到玩家checkpoint成就信息
    checkPointQuotas := api.getCheckPointQuotas(uid, dbDetail, checkPointId)
    *gameResultQuotas = append(*gameResultQuotas, checkPointQuotas...)
    *tmpGameResultQuotas = append(*tmpGameResultQuotas, checkPointQuotas...)
    isChange = api.updateUserCheckPointTotalAccomplishment(uid, userAccomplishment, *tmpGameResultQuotas)

    return
}

//根据玩家所有checkpoint的指标和获取玩家的checkpoint相关的指标值
// dbUserCheckPointSum *battery.DBUserCheckPoint 玩家所有记忆点的指标和信息
// checkPointQuotas []*battery.Quota checkPoint相关的指标集合
func (api *XYAPI) getCheckPointQuotas(uid string, dbUserCheckPointSum *battery.DBUserCheckPoint, checkPointId uint32) (checkPointQuotas []*battery.Quota) {
    checkPointQuotas = []*battery.Quota{&battery.Quota{
        Id:    battery.QuotaEnum_Quota_AllCheckPointScore.Enum(), //分数
        Value: dbUserCheckPointSum.Score,
    }, &battery.Quota{
        Id:    battery.QuotaEnum_Quota_AllCheckPointCharge.Enum(), //charge
        Value: dbUserCheckPointSum.Charge,
    }, &battery.Quota{
        Id:    battery.QuotaEnum_Quota_AllCheckPointStar.Enum(), //有效星星
        Value: proto.Uint64(uint64(dbUserCheckPointSum.GetCollectionsCount())),
    },
    }

    //获取关卡最远指标信息
    if quota := api.getFarthestCheckPointQuota(checkPointId); nil != quota {
        checkPointQuotas = append(checkPointQuotas, quota)
    }

    xylog.Debug(uid, "getCheckPointQuotas : %v", checkPointQuotas)
    return
}

//每个关卡段预留10000个关卡
const (
    CHECKPOINTS_PER_SEGMENT    uint32 = 10000
    CHECKPOINTS_SEGMENT_NORMAL        = 0 //普通关卡阶段
    CHECKPOINTS_SEGMENT_GUIDE         = 1 //新手引导关卡阶段
)

//根据上报的关卡id获取是否需要记录关卡
func (api *XYAPI) getFarthestCheckPointQuota(checkPointId uint32) *battery.Quota {
    switch checkPointId / CHECKPOINTS_PER_SEGMENT {
    case CHECKPOINTS_SEGMENT_NORMAL:
        return &battery.Quota{
            Id:    battery.QuotaEnum_Quota_FarthestCheckPoint.Enum(), //完成的最远的记忆点编号
            Value: proto.Uint64(uint64(checkPointId)),
        }
    case CHECKPOINTS_SEGMENT_GUIDE: //do nothing，新手引导关卡暂时不需要记录，在memcache中保存
    default: //do nothing
    }
    return nil
}

//查询好友记忆点排行榜
// uid string 玩家id
// sids []string 好友sid列表
// idSource battery.ID_SOURCE 玩家好友来源
// checkPointId uint32 记忆点id
// friends *[]*battery.UserCheckPointDetail 好友记忆点排行榜列表指针
func (api *XYAPI) queryCheckPointFriendRank(uid string, sids []string, idSource battery.ID_SOURCE, checkPointId uint32, friends *[]*battery.UserCheckPointDetail, errStruct *battery.Error) {

    errStruct = &battery.Error{}

    //查询玩家好友的uid列表
    uids, mapUid2Sid, err := api.GetUidsFromSids(uid, sids, idSource)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_GetUidError.Enum()
        return
    }

    xylog.Debug(uid, "uids : %v", uids)

    //查询好友列表对应的记忆点信息
    dbUserCheckPoints := make([]*battery.DBUserCheckPoint, 0)
    err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryCheckPointFriendsRank(checkPointId, uids, &dbUserCheckPoints)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_QueryUserCheckPointFriendRankError.Enum()
        //errStruct.Desc = proto.String(fmt.Sprintf("[%s] QueryCheckPointFriendsRank for %v failed : %v", uid, uids, err.Error()))
        return
    }
    xylog.Debug(uid, "QueryCheckPointFriendsRank result: %v", dbUserCheckPoints)

    for _, dbUserCheckPoint := range dbUserCheckPoints {
        checkPoint := xycache.GetCheckPointFromDBUserCheckPoint(dbUserCheckPoint)
        checkPoint.CheckPointId = nil //将记忆点id清除，不需要每个checkpoint都保存一个记忆点id
        //将uid转换为sid
        uidTmp := checkPoint.GetSid()
        if sid, ok := mapUid2Sid[uidTmp]; ok {
            checkPoint.Sid = proto.String(sid)
        }
        //好友排行榜不需要这两个信息，前端在拉好友数据的时候可以获取
        checkPoint.Name = nil
        //checkPoint.IconUrl = nil
        *friends = append(*friends, checkPoint)
    }

    return
}

//查询记忆点全局排行榜
// uid string 玩家id
// sids []string 好友sid列表
// idSource battery.ID_SOURCE 玩家好友来源
// checkPointId uint32 记忆点id
// friends *[]*battery.UserCheckPointDetail 好友记忆点排行榜列表指针
func (api *XYAPI) queryCheckPointGlobalRank(uid string, checkPointId uint32, global *[]*battery.UserCheckPointDetail, platform battery.PLATFORM_TYPE, errStruct *battery.Error) {

    pGlobal := xycache.DefCheckPointGlobalRankManager.GlobalRank(checkPointId, platform)

    if nil == pGlobal {
        errStruct.Code = battery.ErrorCode_QueryUserCheckPointGlobalRankError.Enum()
        return
    } else {
        checkPoint := &battery.DBUserCheckPoint{}
        err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT).QueryUserCheckPointDetail(uid, checkPointId, checkPoint)
        if err != xyerror.ErrOK {
            xylog.ErrorNoId("query user checkpoint fail :%v", err)

        } else {
            dbdetail := api.getCheckPointFromDBUserCheckPoint(checkPoint)
            // 查询玩家tpid信息
            tpid := &battery.IDMap{}
            selector := bson.M{"sid": 1, "gid": 1, "note": 1, "iconurl": 1}
            err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).QueryTpidByUid(selector, uid, tpid, mgo.Monotonic)
            if err != nil {
                xylog.ErrorNoId("uid :%v, get tpid fail :%v ", uid, err)

            } else {
                dbdetail.Sid = tpid.Sid
                dbdetail.Name = tpid.Note
                dbdetail.IconUrl = tpid.IconUrl
            }

            // 玩家不在top排行榜，把最后一名替换为玩家自己
            ok, index := api.isUserInCheckPointRank(uid, *pGlobal)
            if !ok {

                xylog.DebugNoId("global %v, ", len(*pGlobal))
                if len(*pGlobal) < DefConfigCache.Configs().CheckPointGlobalRankSize {
                    *pGlobal = append(*pGlobal, dbdetail)
                } else {
                    (*pGlobal)[len(*pGlobal)-1] = dbdetail
                }

            } else {
                // 更新玩家实时数据
                (*pGlobal)[index] = dbdetail
            }

        }

        *global = *pGlobal
        xylog.Debug(uid, "GlobalRank of checkPointId(%d) : %v", checkPointId, *global)
    }

    return
}

// 角色id转化为初始等级角色id
func getRoleByRoleID(id uint64) (roleid uint64) {
    return id / 10000 * 10000
}

func (api *XYAPI) getCheckPointFromDBUserCheckPoint(dbUserCheckPoint *battery.DBUserCheckPoint) (detail *battery.UserCheckPointDetail) {

    detail = &battery.UserCheckPointDetail{
        Sid:          proto.String(dbUserCheckPoint.GetUid()),
        CheckPointId: proto.Uint32(dbUserCheckPoint.GetCheckPointId()),
        Score:        proto.Uint64(dbUserCheckPoint.GetScore()),
        Charge:       proto.Uint64(dbUserCheckPoint.GetCharge()),
        Coin:         proto.Uint32(dbUserCheckPoint.GetCoin()),
        RoleId:       proto.Uint64(getRoleByRoleID(dbUserCheckPoint.GetRoleId())),
        Grade:        proto.Uint32(dbUserCheckPoint.GetGrade()),
        Collections:  dbUserCheckPoint.Collections,
    }

    return
}

// 判断玩家是否在top排行榜
func (api *XYAPI) isUserInCheckPointRank(uid string, global []*battery.UserCheckPointDetail) (b bool, index int) {
    for index, checkpoint := range global {
        if checkpoint.GetSid() == uid {
            return true, index
        }
    }
    return
}

//判断指标是否是checkPointTotalAccomplishment相关的
// quotaId battery.QuotaEnum 指标枚举值
func (api *XYAPI) isCheckPointAccomplishmentQuotaId(quotaId battery.QuotaEnum) bool {
    if quotaId == battery.QuotaEnum_Quota_AllCheckPointScore ||
        quotaId == battery.QuotaEnum_Quota_AllCheckPointCharge ||
        quotaId == battery.QuotaEnum_Quota_AllCheckPointStar ||
        quotaId == battery.QuotaEnum_Quota_FarthestCheckPoint ||
        quotaId == battery.QuotaEnum_Quota_FinishGame {
        return true
    }
    return false
}
