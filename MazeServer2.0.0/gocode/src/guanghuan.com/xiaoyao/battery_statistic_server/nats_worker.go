package main

import (
	"time"

	"code.google.com/p/goprotobuf/proto"
	nats "github.com/nats-io/nats"

	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"

	batteryapi "guanghuan.com/xiaoyao/battery_statistic_server/business"
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
func NatsIapStatistic(m *nats.Msg) {
	var (
		reqData []byte = m.Data
	)

	req := &battery.IapStatistic{}

	err := proto.Unmarshal(reqData, req)
	if err != xyerror.ErrOK { //解码失败，直接返回
		xylog.ErrorNoId("NatsIapStatistic proto.Unmarshal failed : %v", err)
		return
	}

	xylog.DebugNoId("IapStatistic : %v", req)

	batteryapi.NewXYAPI().TrySingleIapStatistic(req, time.Now().Unix())
}
