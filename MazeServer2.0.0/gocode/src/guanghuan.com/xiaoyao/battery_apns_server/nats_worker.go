package main

import (
	"code.google.com/p/goprotobuf/proto"
	nats "github.com/nats-io/nats"

	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"

	batteryapi "guanghuan.com/xiaoyao/battery_apns_server/business"
)

//pprof消息处理函数
func NatsHandlerPProf(m *nats.Msg) {
	var op string = string(m.Data)
	xyperf.OperationPProf(op)
}

//重置数据库配置项消息处理函数
func NatsHandlerConfig(m *nats.Msg) {
	apiService.LoadConfig()
}

//处理业务推送请求处理函数
func NatsPushNotification(m *nats.Msg) {
	var (
		reqData []byte = m.Data
	)

	req := &battery.APNNotification{}

	err := proto.Unmarshal(reqData, req)
	if err != xyerror.ErrOK { //解码失败，直接返回
		xylog.ErrorNoId("NatsPushNotification proto.Unmarshal failed : %v", err)
		return
	}

	xylog.DebugNoId("APNNotification : %v", req)

	switch req.GetCmd() {
	case battery.APNNotificationCMD_Notification: //推送消息
		batteryapi.NewXYAPI().Notify(req.GetPlatform(), req.GetDeviceToken(), req.GetContent())
		//case xyapn.APNNotificationCMD_EnableDeviceToken: //使能设备token
		//	batteryapi.NewXYAPI().EnableDeviceToken(req.GetUid(), req.GetDeviceToken())
		//case xyapn.APNNotificationCMD_DisableDeviceToken: //去使能设备token
		//	batteryapi.NewXYAPI().DisableDeviceToken(req.GetUid(), req.GetDeviceToken())
	}
}
