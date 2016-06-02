// xyapi_announcement
package batteryapi

import (
	//"code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	//"time"
	"fmt"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationAnnouncement(req *battery.AnnouncementRequest, resp *battery.AnnouncementResponse) (err error) {

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
		goto ErrHandle
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_Announcement, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrHandle:
	return
}
