package main

import (
	nats "github.com/nats-io/nats"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyperf "guanghuan.com/xiaoyao/common/performance"
)

// Echo service
func NatsHandlerEcho(m *nats.Msg) {
	xylog.DebugNoId("Receive msg: %s", string(m.Data))
	nats_service.Publish(m.Reply, m.Data)
}

// pprof消息处理函数
func NatsHandlerPProf(m *nats.Msg) {
	var op string = string(m.Data)
	xyperf.OperationPProf(op)
}
