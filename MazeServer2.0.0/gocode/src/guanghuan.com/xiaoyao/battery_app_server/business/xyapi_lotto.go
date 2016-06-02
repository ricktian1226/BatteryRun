// xyapi_lotto
package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	//xymoney "guanghuan.com/xiaoyao/superbman_server/money"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationLottoRequest(req *battery.LottoRequest, resp *battery.LottoResponse) (err error) {

	var (
		uid    = req.GetUid()
		cmd    = req.GetCmd()
		errStr string
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	xylog.Debug(uid, "[LottoRequest] cmd %v start", cmd)
	defer xylog.Debug(uid, "[LottoRequest] cmd %v done", cmd)

	//初始化resp
	resp.Uid = req.Uid
	resp.Cmd = req.Cmd

	// 查询用户是否存在
	if !api.isUidValid(uid) {
		errStr = fmt.Sprintf("[%s] [app lottorequest] Invalid user.", uid)
		xylog.ErrorNoId(errStr)
		resp.Error = xyerror.Resp_GetAccountByUidError
		resp.Error.Desc = proto.String(errStr)
		err = xyerror.ErrGetAccountByUidError
		goto ErrHandle
	}

	switch cmd {
	case battery.LottoCmd_Lotto_Initial:
		fallthrough
	case battery.LottoCmd_Lotto_Commit:
		fallthrough
	case battery.LottoCmd_Lotto_AfterGame_Initial:
		fallthrough
	case battery.LottoCmd_Lotto_AfterGame_NoInitial:
		fallthrough
	case battery.LottoCmd_Lotto_AfterGame_Commit:
		var failReason battery.ErrorCode
		failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_Lotto, req, resp)
		if failReason != battery.ErrorCode_NoError {
			resp.Error = xyerror.ConstructError(failReason)
		}
	default:
		errStr = fmt.Sprintf("[%s] [app lottorequest] Unkown cmd %v.", cmd)
		xylog.ErrorNoId(errStr)
		resp.Error = xyerror.Resp_BadInputData
		resp.Error.Desc = proto.String(errStr)
		err = xyerror.ErrBadInputData
		goto ErrHandle
	}

ErrHandle:
	xylog.Debug(uid, "cmd %d resp : %v", cmd, resp)
	//添加日志

	return
}
