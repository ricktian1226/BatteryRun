// xystatistic 交易成功统计模块
package xyiapstatistic

import (
	"code.google.com/p/goprotobuf/proto"

	"guanghuan.com/xiaoyao/common/log"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//"guanghuan.com/xiaoyao/superbman_server/iapstatistics"
)

var (
	DefNats *xynatsservice.NatsService
)

const (
	SubjectSendIapStatistic = "iapstatistic"
)

//初始化iapstatistic nats服务指针
func InitNatsService(natsService *xynatsservice.NatsService) {
	DefNats = natsService
}

// Send发送内容上报请求
//  iapStatistic *battery.IapStatistic 上报信息结构体指针
func Send(iapStatistic *battery.IapStatistic) (err error) {

	var data []byte
	data, err = proto.Marshal(iapStatistic)
	if err != nil {
		xylog.ErrorNoId("IapStatistic proto.Marshal failed : %v", err)
		return
	}

	err = DefNats.Publish(SubjectSendIapStatistic, data)
	if err != nil {
		xylog.ErrorNoId("Nats.Publish failed : %v", err)
		return
	}

	return
}
