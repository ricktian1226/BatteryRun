// xyapi_rune
package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	"fmt"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationBeforeGameProp(req *battery.BeforeGamePropRequest, resp *battery.BeforeGamePropResponse) (err error) {

	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if !api.isUidValid(uid) {
		errStr := fmt.Sprintf("[%s]: OperationBeforeGameProp invalid uid", uid)
		resp.Error = xyerror.ConstructError(battery.ErrorCode_BadInputData)
		resp.Error.Desc = proto.String(errStr)
		err = xyerror.ErrGetAccountByUidError
		return
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_BeforeGameProp, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

	return
}
