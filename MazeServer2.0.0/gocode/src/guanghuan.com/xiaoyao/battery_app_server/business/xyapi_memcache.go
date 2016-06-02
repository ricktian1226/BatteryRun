package batteryapi

import (
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//玩家登录消息
func (api *XYAPI) OperationMemCache(req *battery.MemCacheRequest, resp *battery.MemCacheResponse) (err error) {
	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)
	//获取请求的终端平台类型
	api.SetDB(req.GetPlatformType())

	resp.Error = xyerror.DefaultError()

	if !api.isUidValid(uid) {
		resp.Error.Code = battery.ErrorCode_BadInputData.Enum()
		xylog.Error(uid, "OperationMemCache invalid uid")
		err = xyerror.ErrGetAccountByUidError
		return
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_MemCache, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

	xylog.Debug(uid, "OperationMemCache resp : %v", resp)

	return
}
