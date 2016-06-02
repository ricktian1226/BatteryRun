package main

import (
	//"sync"

	//apns "github.com/timehop/apns"

	xyserver "guanghuan.com/xiaoyao/common/server"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"

	batteryapi "guanghuan.com/xiaoyao/battery_apns_server/business"
)

var (
	server *xyserver.Server

	apiService *batteryapi.BatteryService //业务服务对象指针

	natsService      *xynatsservice.NatsService //gateway nats,用于pprof消息的转发
	apnNatsService   *xynatsservice.NatsService //apn nats,用于apn消息的转发
	alertNatsService *xynatsservice.NatsService //alert nats,用于告警消息的转发

	//apnClient apns.Client
)
