package batteryapi

import (
	//"code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationJigsawRequest(req *battery.JigsawRequest, resp *battery.JigsawResponse) (err error) {

	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if !api.isUidValid(uid) {
		err = xyerror.ErrGetAccountByUidError
		return
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_Jigsaw, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

	return err
}
