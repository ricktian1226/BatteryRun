// xyapi_checkpoint
package batteryapi

import (
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//查询玩家区间记忆点信息
func (api *XYAPI) OperationQueryUserCheckPoints(req *battery.QueryUserCheckPointsRequest, resp *battery.QueryUserCheckPointsResponse) (err error) {

    var (
        uid               = req.GetUid()
        checkPointBeginId = req.GetBeginId()
        checkPointEndId   = req.GetEndId()
        failReason        battery.ErrorCode
    )

    // 获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    err = api.verifyCheckPointInput(uid, []uint32{checkPointBeginId, checkPointEndId})
    if err != xyerror.ErrOK {
        return
    }

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_QueryUserCheckPoints, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }

    return
}

//查询玩家记忆点排行榜信息
func (api *XYAPI) OperationQueryCheckPointDetail(req *battery.QueryUserCheckPointDetailRequest, resp *battery.QueryUserCheckPointDetailResponse) (err error) {

    var (
        uid          = req.GetUid()
        checkPointId = req.GetCheckPointId()
        failReason   battery.ErrorCode
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    err = api.verifyCheckPointInput(uid, []uint32{checkPointId})
    if err != xyerror.ErrOK {
        return
    }

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_QueryUserCheckPointDetail, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }

    return
}

// //提交玩家记忆点数据信息
// func (api *XYAPI) OperationCommitCheckPoint(req *battery.CommitCheckPointRequest, resp *battery.CommitCheckPointResponse) (err error) {

//     var (
//         uid          = req.GetUid()
//         checkPointId = req.GetCheckPointId()
//         failReason   battery.ErrorCode
//     )

//     //获取请求的终端平台类型
//     platform := req.GetPlatformType()
//     api.SetDB(platform)

//     err = api.verifyCheckPointInput(uid, []uint32{checkPointId})
//     if err != xyerror.ErrOK {
//         return
//     }

//     failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_CommitCheckPoint, req, resp)
//     if failReason != battery.ErrorCode_NoError {
//         resp.Error = xyerror.ConstructError(failReason)
//     }

//     return
// }

//校验请求的输入参数
func (api *XYAPI) verifyCheckPointInput(uid string, checkPointIds []uint32) (err error) {

    if !api.isUidValid(uid) {
        err = xyerror.ErrGetAccountByUidError
        return
    }

    //for _, checkPointId := range checkPointIds {
    //	if checkPointId > api.Config.CheckPointIdNum {
    //		err = xyerror.ErrBadInputData
    //		return
    //	}
    //}

    return
}

func (api *XYAPI) OperationCheckPointUnlock(req *battery.CheckPointUnlockRequest, resp *battery.CheckPointUnlockResponse) (err error) {
    var (
        uid          = req.GetUid()
        checkPointId = req.GetCheckPointId()
        failReason   battery.ErrorCode
    )

    //获取请求的终端平台类型
    platform := req.GetPlatform()
    api.SetDB(platform)

    err = api.verifyCheckPointInput(uid, []uint32{checkPointId})
    if err != xyerror.ErrOK {
        return
    }

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_CheckPointUnlock, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }

    return
}
