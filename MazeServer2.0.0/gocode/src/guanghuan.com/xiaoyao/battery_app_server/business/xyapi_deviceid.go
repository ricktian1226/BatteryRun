package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationSubmitDeviceId(req *battery.DeviceIdSubmitRequest, resp *battery.DeviceIdSubmitResponse) (err error) {
	// 获取对应的userdata
	uid := req.GetUid()
	dev_id := req.GetDeviceId()

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if uid != "" {
		err = api.UpdateAccountDeviceId(uid, dev_id)
	}
	resp.Uid = proto.String(uid)
	resp.DeviceId = proto.String(req.GetDeviceId())

	return
}

func (api *XYAPI) UpdateAccountDeviceId(uid string, devid string) (err error) {
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).UpdateAccountDeviceId(uid, devid)
	return
}
