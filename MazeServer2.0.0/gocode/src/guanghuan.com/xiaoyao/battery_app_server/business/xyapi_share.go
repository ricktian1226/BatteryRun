package batteryapi

import (
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

// 分享查询
func (api *XYAPI) OperationSharedQuery(req *battery.SharedQueryRequest, resp *battery.SharedQueryResponse) (err error) {
    var (
        uid        = req.GetUid()
        failReason battery.ErrorCode
    )
    platform := req.GetPlatform()
    api.SetDB(platform)

    if !api.isUidValid(uid) {
        err = xyerror.ErrGetAccountByUidError
        return
    }
    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_ShareQuery, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }
    return
}

// 分享请求
func (api *XYAPI) OperationSharedRequest(req *battery.SharedRequest, resp *battery.SharedResponse) (err error) {
    var (
        uid        = req.GetUid()
        failReason battery.ErrorCode
    )
    platform := req.GetPlatform()
    api.SetDB(platform)

    if !api.isUidValid(uid) {
        err = xyerror.ErrGetAccountByUidError
        return
    }
    failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_ShareRequest, req, resp)
    if failReason != battery.ErrorCode_NoError {
        resp.Error = xyerror.ConstructError(failReason)
    }

    return
}
