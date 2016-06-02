package batteryapi

import (
    "fmt"

    proto "code.google.com/p/goprotobuf/proto"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationGetGlobalRankList(req *battery.QueryGlobalRankRequest, resp *battery.QueryGlobalRankResponse) (err error) {
    var (
        uid        = req.GetUid()
        failReason battery.ErrorCode
    )

    // 获取请求的终端平台类型
    platform := req.GetPlatform()
    api.SetDB(platform)

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_GlobalRankList, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }
    return
}

// 玩家起名
func (api *XYAPI) OperationCreatName(req *battery.CreatNameRequest, resp *battery.CreatNameResponse) (err error) {
    var (
        uid        = req.GetUid()
        failReason battery.ErrorCode
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    if !api.isUidValid(uid) {
        errStr := fmt.Sprintf("[%s] uid invalid", uid)
        xylog.ErrorNoId(errStr)
        resp.Error = xyerror.Resp_BadInputData
        resp.Error.Desc = proto.String(errStr)
        err = xyerror.ErrBadInputData
        goto ErrHandle
    }

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_CreatName, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }

ErrHandle:
    return
}
