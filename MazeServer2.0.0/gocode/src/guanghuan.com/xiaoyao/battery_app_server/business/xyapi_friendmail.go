package batteryapi

// xyapi_rolelist

//购买角色(返回购买确认信息（成功则扣除相关消耗：碎片、金币、宝石），失败则提示错误信息)

//返回游戏结算结果(在游戏结算中应该有了加进去就行)

import (
	//"code.google.com/p/goprotobuf/proto"
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationFriendMailListRequest(req *battery.FriendMailListRequest, resp *battery.FriendMailListResponse) (err error) {

	var (
		uid        = req.GetUid()
		errStr     string
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	resp.Uid = req.Uid
	resp.Cmd = req.Cmd

	// 查询用户是否可用
	if !api.isUidValid(uid) {
		errStr = fmt.Sprintf("[%s] invalid user.", uid)
		xylog.ErrorNoId(errStr)
		err = xyerror.ErrGetAccountByUidError
		resp.Error = xyerror.Resp_GetAccountByUidError
		goto ErrHandle
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_FriendMailInfoList, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}
ErrHandle:
	return err
}
