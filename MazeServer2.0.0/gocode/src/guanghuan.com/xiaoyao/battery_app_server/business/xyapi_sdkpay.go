package batteryapi

import (
    //"guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationSDKOrderOp(req *battery.SDKOrderOperationRequest, resp *battery.SDKOrderOperationResponse) (err error) {
    var (
        uid        = req.GetUid()
        failReason battery.ErrorCode
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_SDKOrderOp, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }
    return

}

func (api *XYAPI) OperationSDKOrderQuery(req *battery.SDKOrderRequest, resp *battery.SDKOrderResponse) (err error) {
    var (
        uid        = req.GetUid()
        failReason battery.ErrorCode
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_SDKOrderQuery, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }
    return
}

func (api *XYAPI) OperationSDKAddOrder(req *battery.SDKAddOrderRequest, resp *battery.SDKAddOrderResponse) (err error) {
    var (
        uid        = req.GetUid()
        failReason battery.ErrorCode
    )

    //获取请求的终端平台类型
    xylog.DebugNoId("sdk callback")
    platform := req.GetPlatformType()
    api.SetDB(platform)

    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_SDKAddOrder, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }
    return
}
