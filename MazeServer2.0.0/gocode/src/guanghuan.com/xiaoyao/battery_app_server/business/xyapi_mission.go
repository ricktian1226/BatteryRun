// xyapi_mission
package batteryapi

//app任务系统模块

import (
	//"code.google.com/p/goprotobuf/proto"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//查询玩家任务消息
func (api *XYAPI) OperationQueryUserMission(req *battery.QueryUserMissionRequest, resp *battery.QueryUserMissionResponse) (err error) {
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

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_QueryUserMission, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrHandle:

	go xybusiness.AddOperationLog(uid, req.String(), resp.String(), api.GetLogXYDB(), xybusiness.BusinessCode_QueryUserMission, 0, resp.Error)

	return
}

//领取任务奖励接口
func (api *XYAPI) OperationConfirmUserMission(req *battery.ConfirmUserMissionRequest, resp *battery.ConfirmUserMissionResponse) (err error) {

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

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_ConfirmUserMission, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrHandle:

	go xybusiness.AddOperationLog(uid, req.String(), resp.String(), api.GetLogXYDB(), xybusiness.BusinessCode_ConfirmUserMission, 0, resp.Error)

	return
}
