package batteryapi

import (
    "code.google.com/p/goprotobuf/proto"
    "fmt"
    xylog "guanghuan.com/xiaoyao/common/log"
    xyutil "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//// 处理单次游戏数据
//func (api *XYAPI) OperationAddGameData(req *battery.GameDataRequest, resp *battery.GameDataResponse) (err error) {
//	var (
//		uid        = req.GetUid()
//		gameId     = req.GetGameId()
//		gameType   = req.GetType()
//		failReason = xyerror.Resp_NoError.GetCode()
//		errStr     string
//	)

//	//获取请求的终端平台类型
//	platform := req.GetPlatformType()
//	api.SetDB(platform)

//	if !api.isUidValid(uid) || battery.GameType_GameType_Unkown == gameType {
//		errStr = fmt.Sprintf("[%s] uid or gametype invalid %v", uid, gameType)
//		xylog.ErrorNoId(errStr)
//		failReason = xyerror.Resp_BadInputData.GetCode()
//		resp.Error = xyerror.Resp_BadInputData
//		resp.Error.Desc = proto.String(errStr)
//		goto ErrHandle
//	}

//	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_GameResult, req, resp)
//	if failReason != battery.ErrorCode_NoError {
//		resp.Error = xyerror.ConstructError(failReason)
//	}

//ErrHandle:
//	go api.AddGameLog(uid, gameId, 0, gameType, GAME_OP_UPLOAD, failReason, resp.Error.GetDesc(), true)

//	return
//}

// 处理游戏结算数据提交请求
func (api *XYAPI) OperationGameResultCommit(req *battery.GameResultCommitRequest, resp *battery.GameResultCommitResponse) (err error) {

    var (
        uid          = req.GetUid()
        gameId       = req.GetGameId()
        gameType     = req.GetType()
        checkPointId = req.GetCheckPointId()
        isFinish     = req.GetIsFinish()
        failReason   battery.ErrorCode
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    //校验请求的合法性
    isValid, errStr := api.isGameResultCommitValid(req)
    if !isValid {
        resp.Error = xyerror.Resp_BadInputData
        resp.Error.Desc = proto.String(errStr)
        err = xyerror.ErrBadInputData
        goto ErrHandle
    }

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_GameResult2, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }

ErrHandle:
    go api.AddGameLog(uid, gameId, checkPointId, gameType, GAME_OP_UPLOAD, resp.Error.GetCode(), resp.Error.GetDesc(), isFinish)

    return
}

//校验游戏结算数据有效性
func (api *XYAPI) isGameResultCommitValid(req *battery.GameResultCommitRequest) (isValid bool, errStr string) {
    var (
        uid          = req.GetUid()
        gameId       = req.GetGameId()
        gameType     = req.GetType()
        checkPointId = req.GetCheckPointId()
        quotas       = req.GetGameResult().GetQuotas()
    )
    //玩家非法
    if !api.isUidValid(uid) {
        errStr = fmt.Sprintf("[%s] uid invalid", uid)
        goto ErrHandle
    }
    //游戏类型非法
    if battery.GameType_GameType_Unkown == gameType {
        errStr = fmt.Sprintf("[%s] gametype invalid", uid)
        goto ErrHandle
    }
    //记忆点标识非法
    if 0 == checkPointId {
        errStr = fmt.Sprintf("[%s] checkPointId invalid", uid)
        goto ErrHandle
    }
    //游戏标识非法
    if "" == gameId || "0" == gameId {
        errStr = fmt.Sprintf("[%s] gameId invalid", uid, gameId)
        goto ErrHandle
    }
    //游戏结算结果非法
    if nil == req.GetGameResult() || 0 >= len(req.GetGameResult().GetQuotas()) {
        errStr = fmt.Sprintf("[%s] gameresult is nil", uid)
        goto ErrHandle
    }
    //游戏时长非法
    if 0 >= req.GetDuration() {
        errStr = fmt.Sprintf("[%s] game duration %d is nil", uid, req.GetDuration())
        goto ErrHandle
    }
    if !api.isQuotaValid(quotas) {
        errStr = fmt.Sprintf("[%s] game quotas  %d invalid", quotas)
        goto ErrHandle
    }
    isValid = true

    return

ErrHandle:
    isValid = false
    xylog.ErrorNoId(errStr)
    return
}
func (api *XYAPI) isQuotaValid(quotas []*battery.Quota) (isvalid bool) {
    isvalid = true
    for _, quota := range quotas {

        if quota.GetId() == battery.QuotaEnum_Quota_Score && quota.GetValue() >= DefConfigCache.Configs().InvalidScore {
            isvalid = false
            xylog.ErrorNoId("invalid score :%d	", quota.GetValue())
            break
        }
    }
    return
}
func (api *XYAPI) OperationNewGame(req *battery.NewGameRequest, resp *battery.NewGameResponse) (err error) {
    var (
        uid          = req.GetUid()
        gameType     = req.GetType()
        gameId       = InvalidGameID
        checkPointId = req.GetCheckPointId()
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    failReason := xyerror.Resp_NoError.GetCode()
    if !api.isUidValid(uid) || battery.GameType_GameType_Unkown == gameType {
        xylog.Error(uid, "uid or gametype(%v) invalid ", gameType)
        failReason = xyerror.Resp_BadInputData.GetCode()
        resp.Error = xyerror.Resp_BadInputData
        goto ErrHandle
    }

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_NewGame, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }

ErrHandle:

    gameId = resp.GetGameId()

    go api.AddGameLog(uid, gameId, checkPointId, gameType, GAME_OP_START, failReason, resp.Error.GetDesc(), false)

    return
}

func (api *XYAPI) AddGameLog(uid, gameId string, checkPointId uint32, gameType battery.GameType, gameOp int32, fail_reason battery.ErrorCode, data string, isFinish bool) (err error) {
    gameLog := battery.GameLog{
        Uid:          proto.String(uid),
        GameId:       proto.String(gameId),
        CheckPointId: proto.Uint32(checkPointId),
        Type:         gameType.Enum(),
        OpType:       proto.Int32(gameOp),
        FailReason:   fail_reason.Enum(),
        OpDate:       proto.Int64(xyutil.CurTimeSec()),
        OpDateStr:    proto.String(xyutil.CurTimeStr()),
        Result:       proto.String(data),
        IsFinish:     proto.Bool(isFinish),
    }

    api.GetLogDB().AddGameLog(gameLog)

    return
}
