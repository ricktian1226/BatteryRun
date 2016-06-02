package main

import (
	xyserver "guanghuan.com/xiaoyao/common/server"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"

	batteryapi "guanghuan.com/xiaoyao/battery_statistic_server/business"
)

var (
	server *xyserver.Server

	apiService *batteryapi.BatteryService //业务服务对象指针

	natsService             *xynatsservice.NatsService //gateway nats,用于pprof消息的转发
	iapStatisticNatsService *xynatsservice.NatsService //iapstatistic nats,用于内购统计信息的转发
	alertNatsService        *xynatsservice.NatsService //alert nats,用于告警消息的转发

)
