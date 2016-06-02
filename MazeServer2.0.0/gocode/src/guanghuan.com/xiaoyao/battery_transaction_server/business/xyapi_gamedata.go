package batteryapi

import (
    proto "code.google.com/p/goprotobuf/proto"
    "gopkg.in/mgo.v2"
    xylog "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/performance"
    xyutil "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xycache "guanghuan.com/xiaoyao/superbman_server/cache"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
    "time"
)

//金币的道具id
var PROPID_COIN uint64 = 10010000

//游戏结算数据提交消息
func (api *XYAPI) OperationGameResultCommit(req *battery.GameResultCommitRequest, resp *battery.GameResultCommitResponse) (err error) {
    var (
        uid             = req.GetUid()
        gameId          = req.GetGameId()
        gameResult      = req.GetGameResult()
        gameType        = req.GetType()
        gameDuration    = req.GetDuration()
        checkPointId    = req.GetCheckPointId()
        collections     = req.GetCollections()
        roleId          = req.GetRoleId()
        pickUps         = req.GetPickups()
        isFinish        = req.GetIsFinish()
        game            battery.Game
        startTime       int64
        account         = new(battery.DBAccount)
        errStruct       = xyerror.DefaultError()
        awards          []*battery.PropItem
        missionTypes    = make([]battery.MissionType, 0)
        accountWithFlag *AccountWithFlag
    )
    xylog.DebugNoId("roleid :%v", roleId)
    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    //xylog.Debug("[%s] req : %v", uid, req)

    //初始化resp
    resp.GameId = req.GameId
    resp.Error = xyerror.Resp_NoError

    //查询游戏是否存在
    game, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_GAME).GetGame(gameId)
    if err != xyerror.ErrOK || game.GetUid() != uid { //查找失败或者玩家uid不符，则请求非法
        xylog.Error(uid, "[AddGameData]: game[%s] not found:%v", gameId, err)
        errStruct.Code = battery.ErrorCode_GameNotExistError.Enum()
        goto ErrHandle
    }

    //判断游戏状态是否正确（Uploadtime非，说明游戏已经被提交过了）
    if 0 != game.GetUploadTime() {
        xylog.Error(uid, "[AddGameData]: game[%s] already finish, uploadtime %d", gameId, game.GetUploadTime())
        err = xyerror.ErrBadInputData
        errStruct.Code = battery.ErrorCode_BadInputData.Enum()
        goto ErrHandle
    }

    //判断游戏数据是否有效
    startTime = game.GetStartTime()
    if gameResult == nil || !api.isGameResultValid(uid, startTime, gameDuration, gameResult) {
        xylog.Error(uid, "[AddGameData]: game[%s] not found:%v", gameId, err)
        err = xyerror.ErrGameResultInvalidError
        errStruct.Code = battery.ErrorCode_GameResultInvalidError.Enum()
        goto ErrHandle
    }

    // 查询数据结算前的玩家数据信息
    err = api.GetDBAccountDirect(uid, account, mgo.Strong)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "Error: user not existing!")
        errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
        goto ErrHandle
    }

    accountWithFlag = &AccountWithFlag{
        account: account,
        bChange: false,
    }

    //修改玩家数据，这个接口只更新玩家的checkpoint数据，不更新玩家的账户数据，账户数据在游戏后奖励后一把更新（避免同一请求多次更新玩家账户数据）
    api.updateGameResult2(uid, gameId, roleId, checkPointId, collections, accountWithFlag, gameType, gameResult, isFinish, errStruct, platform)
    if errStruct.GetCode() != battery.ErrorCode_NoError {
        goto ErrHandle
    }
    //获取游戏后奖励
    awards = api.afterGameAwards(uid, roleId, gameResult, accountWithFlag)
    xylog.Debug(uid, "afterGameAwards : %v", awards)

    //发放收集物
    api.GainPickUps(uid, checkPointId, pickUps, accountWithFlag)

    // 更新玩家数据
    err = api.UpdateAccountWithFlag(accountWithFlag)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_UpdateAccountError.Enum()
        goto ErrHandle
    }

    //刷新玩家任务状态
    switch gameType {
    case battery.GameType_GameType_Study:
        missionTypes = append(missionTypes, battery.MissionType_MissionType_Study)
    default:
        missionTypes = append(missionTypes, battery.MissionType_MissionType_Daily)
        missionTypes = append(missionTypes, battery.MissionType_MissionType_MainLine)
    }

    //if isFinish { //添加记忆点任务指标信息
    //	quota := &battery.Quota{
    //		Id:    battery.QuotaEnum_Quota_FarthestCheckPoint.Enum(),
    //		Value: proto.Uint64(uint64(checkPointId)),
    //	}
    //	gameResult.Quotas = append(gameResult.Quotas, quota)
    //}

    xylog.Debug(uid, "gameResult.Quotas : %v", gameResult.Quotas)

    err = api.updateUserMissionsQuotas(uid, missionTypes, gameResult.Quotas, time.Now().Unix(), isFinish)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "updateUserMissionsQuotas failed")
    }

    //更新本局游戏状态
    game.Duration = proto.Int64(gameDuration)
    game.UploadTime = proto.Int64(xyutil.CurTimeSec())
    game.Result = gameResult
    game.IsFinish = req.IsFinish

    {
        begin := time.Now()
        defer xyperf.Trace(LOGTRACE_UPDATEGAME, &begin)
    }

    err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_GAME).UpdateGame(game)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "[AddGameData] update game record failed: %v", err)
        errStruct.Code = battery.ErrorCode_UpdateGameError.Enum()
        goto ErrHandle
    }

    //设置一下返回信息
    resp.Data = account
    resp.Awards = awards

ErrHandle:
    resp.Error = errStruct

    return
}

func (api *XYAPI) OperationNewGame(req *battery.NewGameRequest, resp *battery.NewGameResponse) (err error) {
    var (
        uid             = req.GetUid()
        gameType        = req.GetType()
        gameId          = "0"
        startTime int64 = -1
        curAmount int32
        countDown int32 = -1
        game      *battery.Game
        errStruct = xyerror.DefaultError()
    )
    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    resp.Error = xyerror.Resp_NoError
    resp.CheckPointId = req.CheckPointId

    // 先检查是否有足够体力(非教学模式)
    // 更新体力值（非教学模式）
    if gameType != battery.GameType_GameType_Study {
        curAmount, countDown, err = api.UpdateStamina(uid, -1)
        if err != xyerror.ErrOK {
            xylog.Error(uid, "[NewGame] failed: %v", err)
            if err == xyerror.ErrNotEnoughStamina {
                errStruct.Code = battery.ErrorCode_NotEnoughStamina.Enum()
            } else {
                errStruct.Code = battery.ErrorCode_UpdateStaminaError.Enum()
            }
            goto ErrHandle
        } else {
            xylog.Debug(uid, "[NewGame] current stamina= %d", curAmount)

            //刷新消耗体力XX次相关任务
            quotas := []*battery.Quota{&battery.Quota{Id: battery.QuotaEnum_Quota_Stamina.Enum(), Value: proto.Uint64(1)}}
            missionTypes := []battery.MissionType{battery.MissionType_MissionType_Study, battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_MainLine}
            err = api.updateUserMissionsQuotas(uid, missionTypes, quotas, time.Now().Unix(), MissionQuotasNoNeedFinish)
        }
    }

    {
        begin := time.Now()
        defer xyperf.Trace(LOGTRACE_NEWGAME, &begin)
    }

    // 添加新游戏
    // 产生一个默认的游戏记录
    game = api.newDefaultGame(uid)
    game.Type = gameType.Enum()
    err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_GAME).AddNewGame(*game)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "[New Game] failed: %v", err)
        errStruct.Code = battery.ErrorCode_AddNewGameError.Enum()
        goto ErrHandle
    } else {
        startTime = game.GetStartTime()
        gameId = game.GetId()
    }

ErrHandle:
    //游戏日志中加上请求详情
    resp.Error = errStruct
    if errStruct.GetCode() == battery.ErrorCode_NoError {
        // 成功
        resp.GameId = proto.String(gameId)
        resp.StartTime = proto.Int64(startTime)
        resp.Stamina = proto.Int32(curAmount)
        resp.Timeleft = proto.Int32(countDown)
    } else {
        //错误返回 game_id=0, start_time=0
        resp.GameId = proto.String(InvalidGameID)
        resp.StartTime = proto.Int64(0)
    }

    return
}

func (api *XYAPI) newDefaultGame(uid string) *battery.Game {
    game := &battery.Game{}
    gameid := xyutil.NewId()
    game.Id = proto.String(gameid)
    game.Uid = proto.String(uid)
    game.StartTime = proto.Int64(xyutil.CurTimeSec())

    return game
}

//获取游戏后奖励
// distance uint64
func (api *XYAPI) afterGameAwards(uid string, roleId uint64, gameResult *battery.GameResult, accountWithFlag *AccountWithFlag) (awards []*battery.PropItem) {

    //begin := time.Now()
    //defer xyperf.Trace(LOGTRACE_AFTERGAMEREWARDS, &begin)

    awards = make([]*battery.PropItem, 0)

    for _, quota := range gameResult.GetQuotas() {
        //根据分数计算出游戏后奖励的金币数
        if quota.GetId() == battery.QuotaEnum_Quota_Score {
            bigValue := (uint32(quota.GetValue() * (DefConfigCache.Configs().AfterGameAwardFactor)))
            amount := bigValue / 10000
            if bigValue%10000 > 0 { //如果有余数，金币数加1
                amount += 1
            }

            factor := RUNE_BASE_FACTOR
            //如果金币加成符文有效，加上符文的金币加成
            if accountWithFlag.account.GetCoinAddtionalExpiredTimestamp() > xyutil.CurTimeSec() {
                factor += accountWithFlag.account.GetCoinAddtional()
            }

            //加上角色等级的金币加成
            roleLevelBonus := xycache.DefRoleLevelBonusCacheManager.Bonus(roleId)
            if nil != roleLevelBonus {
                factor += roleLevelBonus.GoldBonus
                //xylog.Debug("[%s] added roleLevelBonus.GoldBonus, factor : %d", uid, factor)
            }

            amount = (amount * uint32(factor)) / uint32(RUNE_BASE_FACTOR)

            xylog.Debug(uid, "afterGameAwards score %d coin amount %d", quota.GetValue(), amount)

            awards = append(awards, &battery.PropItem{
                Id:     proto.Uint64(PROPID_COIN),
                Type:   battery.PropType_PROP_COIN.Enum(),
                Amount: proto.Uint32(amount),
            })
            //账户钱包中加上金币
            accountWithFlag.account.Wallet[battery.MoneyType_coin].Gainamount = proto.Uint32(accountWithFlag.account.Wallet[battery.MoneyType_coin].GetGainamount() + uint32(amount))
            accountWithFlag.SetChange()
            break
        }
    }

    return
}

//to delete
// 按照玩家奔跑距离计算获取的金币数。该逻辑暂时无用。
type distanceMultiplerItem struct {
    distance uint64
    factor   float64
}

var distanceMultiplerSlice = []distanceMultiplerItem{
    {800, 0.448},  //[0,800)
    {600, 0.448},  //[800,1400)
    {800, 0.504},  //[1400,2200)
    {400, 0.504},  //[2200,2600)
    {200, 0.56},   //[2600,2800)
    {600, 0.56},   //[2800,3400)
    {200, 0.616},  //[3400,3600)
    {200, 0.616},  //[3600,3800)
    {600, 0.672},  //[3800,4400)
    {200, 0.672},  //[4400,4600)
    {400, 0.728},  //[4600,5000)
    {400, 0.728},  //[5000,5400)
    {800, 0.784},  //[5400,6200)
    {400, 0.784},  //[6200,6600)
    {400, 0.84},   //[6600,7000)
    {400, 0.84},   //[7000,7400)
    {400, 0.896},  //[7400,7800)
    {400, 0.896},  //[7800,8200)
    {1600, 0.896}, //[8200,9800)
    {0, 0.896},    //[9800,)
}

//根据米数获取奖励的金币数
func (api *XYAPI) getCoinsFromDistance(distance uint64) (amount uint32) {
    var i uint64
    for _, item := range distanceMultiplerSlice {
        if item.distance == 0 ||
            item.distance > distance {
            i = distance
        } else {
            i = item.distance
        }
        amount += uint32(float64(i) * item.factor)
        distance -= i
        if distance == 0 {
            break
        }
    }

    return
}

//func (api *XYAPI) isValidGameDataRequest(req battery.GameDataRequest) (isValid bool) {
//	if req.Uid == nil || req.GameId == nil || req.GameResult == nil || req.Duration == nil {
//		return false
//	}
//	game_id := req.GetGameId()
//	game_duration := req.GetDuration()

//	if !api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_GAME).IsValidGameId(game_id) || game_duration < 0 {
//		return false
//	}
//	return true
//}

func (api *XYAPI) isGameResultValid(uid string, start_time int64, duration int64, result *battery.GameResult) (isValid bool) {
    isValid = true
    return
}

//func (api *XYAPI) updateGameResult(uid string, account *battery.DBAccount, gameType battery.GameType, gameResult *battery.GameResult) (err error) {

//	//判断游戏account信息是否正确，必须有total和best结构
//	if len(account.Accomplishment) < 2 {
//		xylog.Warning("[%s] wrong account Accomplishment len(%d)", uid, len(account.Accomplishment))
//		account.Accomplishment = make([]*battery.Accomplishment, 2)
//		for i := 0; i < 2; i++ {
//			account.Accomplishment[i] = api.defaultAccomplishment()
//		}
//	}

//	gameResultQuotas := gameResult.GetQuotas()

//	xylog.Debug("[%s] GameResultQuotas : %v", uid, gameResultQuotas)

//	//游戏数据中没有指标，直接跳过
//	if len(gameResultQuotas) <= 0 {
//		return
//	}

//	//如果非教学模式，需要刷新玩家分数(account的accomplishment)
//	if gameType != battery.GameType_GameType_Study {
//		api.updateUserAccomplishment(uid, account, gameResultQuotas)
//	}

//	//刷新玩家任务状态
//	missionTypes := make([]battery.MissionType, 0)
//	switch gameType {
//	case battery.GameType_GameType_Study:
//		missionTypes = append(missionTypes, battery.MissionType_MissionType_Study)
//	default:
//		missionTypes = append(missionTypes, battery.MissionType_MissionType_Daily)
//		missionTypes = append(missionTypes, battery.MissionType_MissionType_MainLine)
//	}

//	now := time.Now().Unix()
//	err = api.updateUserMissionsQuotas(uid, missionTypes, gameResultQuotas, now)

//	return
//}

//结算数据处理
// uid string 玩家id
// gameId string 游戏id
// roleId uint64 角色id
// checkPointId uint32 记忆点id
// collections []uint32 收集物（大星星）
// account *battery.DBAccount 玩家帐号信息
// gameType battery.GameType 游戏类型
// gameResult *battery.GameResult 游戏结果集
// isFinish bool 是否完成
//return:
// err error 返回错误信息
func (api *XYAPI) updateGameResult2(uid, gameId string, roleId uint64, checkPointId uint32, collections []uint32, accountWithFlag *AccountWithFlag, gameType battery.GameType, gameResult *battery.GameResult, isFinish bool, errStruct *battery.Error, platform battery.PLATFORM_TYPE) {
    begin := time.Now()
    defer xyperf.Trace(LOGTRACE_GAMERESULT, &begin)

    //如果非教学模式，需要刷新玩家分数(account的accomplishment)
    if gameType != battery.GameType_GameType_Study {
        //游戏数据中没有指标，直接跳过
        gameResultQuotas := gameResult.GetQuotas()
        xylog.Debug(uid, "GameResultQuotas : %v", gameResultQuotas)
        if len(gameResultQuotas) <= 0 {
            return
        }

        //查找玩家的成就信息
        userAccomplishment := &battery.DBUserAccomplishment{}
        api.getUserAccomplishment(uid, userAccomplishment, errStruct)
        if errStruct.GetCode() != battery.ErrorCode_NoError {
            return
        }
        userAccomplishment.Uid = proto.String(uid) //查询时为了减少通信数据量，返回uid
        isAccomplishmentChange := false

        //如果记忆点完成，刷新玩家记忆点和checkPointAccomplishment数据
        if isFinish {
            //增加完成游戏指标
            //api.addQuota(battery.QuotaEnum_Quota_FinishGame, 1, &gameResultQuotas)

            isAccomplishmentChange = api.updateUserCheckPoint(uid, gameId, roleId, checkPointId, collections, &(gameResult.Quotas), &gameResultQuotas, userAccomplishment, errStruct, platform)
            if errStruct.GetCode() != battery.ErrorCode_NoError {
                xylog.Error("updateUserCheckPoint failed :%s", errStruct.GetDesc())
                return
            }
        }

        //刷新玩家totalAccomplishment和bestAccomplishment数据
        if api.updateUserAccomplishment(uid, userAccomplishment, gameResultQuotas, errStruct) {
            isAccomplishmentChange = true
        }

        //需要刷新成就的话，就刷新下
        if isAccomplishmentChange {
            err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERACCOMPLISHMENT).UpdateUserAccomplishment(userAccomplishment)
            if err != xyerror.ErrOK {
                errStruct.Code = battery.ErrorCode_UpdateUserAccomplishmentError.Enum()
                return
            }
        }
    }

    return
}

//根据上报的数据刷新玩家分数数据
// uid string 玩家id
// account *battery.Account 玩家账户信息
// gameResultQuotas []*battery.Quota 游戏数据指标信息
//return:
// isChange bool 玩家成绩数据是否有修改，true 有修改，false 无修改
func (api *XYAPI) updateUserAccomplishment(uid string, userAccomplishment *battery.DBUserAccomplishment, gameResultQuotas []*battery.Quota, errStruct *battery.Error) (isChange bool) {
    totalAccomplishment := userAccomplishment.Accomplishment[int(battery.AccomplishmentType_AccomplishmentType_Total)]
    bestAccomplishment := userAccomplishment.Accomplishment[int(battery.AccomplishmentType_AccomplishmentType_Best)]
    isChange = false

    for _, gameResultQuota := range gameResultQuotas {
        id := gameResultQuota.GetId()
        value := gameResultQuota.GetValue()

        //获取指标id对应的指标类型下标和指标下标
        quotaTypeIndex, quotaIndex := api.parseQuotaId(id)

        //修改玩家总计成绩
        if api.fixAccomplishment(uid, quotaTypeIndex, quotaIndex, totalAccomplishment) {
            isChange = true
        }

        totalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].Id = id.Enum()
        totalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].Value = proto.Uint64(totalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].GetValue() + value)
        isChange = true

        //修改玩家最好成绩
        if api.fixAccomplishment(uid, quotaTypeIndex, quotaIndex, bestAccomplishment) {
            isChange = true
        }

        if bestAccomplishment.List[quotaTypeIndex].Items[quotaIndex].GetId() == battery.QuotaEnum_Quota_Unkown || //如果是新建的指标节点，直接赋值
            bestAccomplishment.List[quotaTypeIndex].Items[quotaIndex].GetValue() < value { //如果新值大于旧值，赋值
            xylog.Debug(uid, "%v oldvalue:%d, newvalue:%d", id, bestAccomplishment.List[quotaTypeIndex].Items[quotaIndex].GetValue(), value)
            bestAccomplishment.List[quotaTypeIndex].Items[quotaIndex].Id = id.Enum()
            bestAccomplishment.List[quotaTypeIndex].Items[quotaIndex].Value = proto.Uint64(value)
            isChange = true
        }
    }

    return
}

//根据上报的数据刷新玩家的checkpoint成就信息信息
// uid string 玩家id
// account *battery.Account 玩家账户信息
// gameResultQuotas []*battery.Quota 游戏数据指标信息
//func (api *XYAPI) updateUserCheckPointTotalAccomplishment(uid string, accountWithFlag *AccountWithFlag, gameResultQuotas []*battery.Quota) {
func (api *XYAPI) updateUserCheckPointTotalAccomplishment(uid string, userAccomplishment *battery.DBUserAccomplishment, gameResultQuotas []*battery.Quota) (isChange bool) {
    checkPointTotalAccomplishment := userAccomplishment.Accomplishment[battery.AccomplishmentType_AccomplishmentType_CheckPoint_Total]
    isChange = false
    for _, gameResultQuota := range gameResultQuotas {

        id := gameResultQuota.GetId()
        if api.isCheckPointAccomplishmentQuotaId(id) { //如果不是checkpoint相关的指标，则跳过
            value := gameResultQuota.GetValue()

            //获取指标id对应的指标类型下标和指标下标
            quotaTypeIndex, quotaIndex := api.parseQuotaId(id)

            xylog.Debug(uid, "checkPointTotalAccomplishment quotaTypeIndex %d, quotaIndex %d", quotaTypeIndex, quotaIndex)

            //修改玩家总计成绩
            if api.fixAccomplishment(uid, quotaTypeIndex, quotaIndex, checkPointTotalAccomplishment) {
                isChange = true
            }
            //游戏完成总数需要特殊处理下，计算成总数
            if id == battery.QuotaEnum_Quota_FinishGame {
                value += checkPointTotalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].GetValue()
            }

            if checkPointTotalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].GetId() == battery.QuotaEnum_Quota_Unkown || //如果是新建的指标节点，直接赋值
                checkPointTotalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].GetValue() < value {
                checkPointTotalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].Id = id.Enum()
                checkPointTotalAccomplishment.List[quotaTypeIndex].Items[quotaIndex].Value = &value
                isChange = true
            }
        }
    }

    return
}

//判断指标是否是totalAccomplishment相关的
// quotaId battery.QuotaEnum 指标枚举值
func (api *XYAPI) isTotalAccomplishmentQuotaId(quotaId battery.QuotaEnum) bool {
    return false
}

//判断指标是否是bestAccomplishment相关的
// quotaId battery.QuotaEnum 指标枚举值
func (api *XYAPI) isBestAccomplishmentQuotaId(quotaId battery.QuotaEnum) bool {
    return false
}

func (api *XYAPI) parseQuotaId(quotaId battery.QuotaEnum) (quotaTypeIndex int, quotaIndex int) {
    quotaTypeIndex = (int(quotaId) / 1000)
    quotaIndex = (int(quotaId) % 1000)
    return
}

//修复玩家的成绩信息
// uid string 玩家id
// typeIndex int 指标类型
// index int 指标索引
// accomplishment *battery.Accomplishment 成绩信息的指针
//return:
// isFix bool 是否有修复 true 修复 false 未修复，可以通过isFix确认是否有修改
func (api *XYAPI) fixAccomplishment(uid string, typeIndex int, index int, accomplishment *battery.Accomplishment) (isFix bool) {

    isFix = false

    lengthList := len(accomplishment.List)
    //一级下标不存在，补一下
    if typeIndex >= lengthList {
        for i := 0; i < typeIndex+1-lengthList; i++ {
            accomplishment.List = append(accomplishment.List, new(battery.QuotaList))
        }
        isFix = true
    }

    //二级下标不存在，补一下
    lengthItems := len(accomplishment.List[typeIndex].Items)
    if index >= lengthItems {
        for i := 0; i < index+1-lengthItems; i++ {
            accomplishment.List[typeIndex].Items = append(accomplishment.List[typeIndex].Items, new(battery.Quota))
        }
        isFix = true
    }

    return
}

//判断游戏是否在进行中
func (api *XYAPI) IsGameOnGoing(uid string, game_id string) (isOnGoing bool) {
    game, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_GAME).GetGame(game_id)
    if err == xyerror.ErrOK && game.GetUid() == uid {
        if game.Result == nil && game.GetUploadTime() <= 0 {
            isOnGoing = true
        }
    }
    return
}
