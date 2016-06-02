// xyapn
// apple push notification 接口
// 定义apn推送的公共接口，供业务有推送需求时调用
package xyapn

import (
	"time"

	//nats "github.com/nats-io/nats"
	"code.google.com/p/goprotobuf/proto"

	"guanghuan.com/xiaoyao/common/log"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

const (
	SubjectPushNotification = "pushnotification"
)

var (
	DefNats    *xynatsservice.NatsService
	DefOptions = NewOptions()
)

type Options struct {
	NatsTimeOut time.Duration //nats请求超时时间
}

func NewOptions() *Options {
	return &Options{
		NatsTimeOut: (time.Second * 10),
	}
}

//初始化apn nats服务指针
func InitNatsService(natsService *xynatsservice.NatsService) {
	DefNats = natsService
}

// 发送apn推送消息 该接口供业务调用
// notification *APNNotification
func Send(notification *battery.APNNotification) (err error) {

	var data []byte
	data, err = proto.Marshal(notification)
	if err != nil {
		xylog.ErrorNoId("APNNotification proto.Marshal failed : %v", err)
		return
	}

	err = DefNats.Publish(SubjectPushNotification, data)
	if err != nil {
		xylog.ErrorNoId("Nats.Publish failed : %v", err)
		return
	}

	return
}
