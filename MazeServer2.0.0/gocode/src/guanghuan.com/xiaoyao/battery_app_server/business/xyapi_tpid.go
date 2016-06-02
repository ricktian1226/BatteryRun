package batteryapi

import (
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// 绑定账户消息处理
func (api *XYAPI) OperationBind(req *battery.BindRequest, resp *battery.BindResponse) (err error) {
	var (
		uid     = req.GetUid()
		errCode battery.ErrorCode
	)

	resp.Error = xyerror.DefaultError()
	resp.Uid = req.Uid

	errCode, err = api.SendToTransaction(uid, xybusiness.BusinessCode_Bind, req, resp)
	if errCode != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(errCode)
	}
	return
}
