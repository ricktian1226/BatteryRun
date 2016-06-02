// xyapi_signin
package batteryapi

//app任务系统模块

import (
	//"code.google.com/p/goprotobuf/proto"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//签到查询请求
func (api *XYAPI) OperationQuerySignIn(req *battery.QuerySignInRequest, resp *battery.QuerySignInResponse) (err error) {
	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if !api.isUidValid(uid) {
		err = xyerror.ErrGetAccountByUidError
		goto ErrHandle
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_QuerySignIn, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrHandle:

	return
}

//新版本签到查询请求
func (api *XYAPI) OperationQuerySignIn2(req *battery.NewQuerySignInRequest, resp *battery.NewQuerySignInResponse) (err error) {
	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if !api.isUidValid(uid) {
		err = xyerror.ErrGetAccountByUidError
		goto ErrHandle
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_QuerySignIn2, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrHandle:

	return
}

//签到上报请求
func (api *XYAPI) OperationSignIn(req *battery.SignInRequest, resp *battery.SignInResponse) (err error) {
	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if !api.isUidValid(uid) {
		err = xyerror.ErrGetAccountByUidError
		goto ErrHandle
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_SignIn, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrHandle:

	return
}
