package main

import (
	//	"sync"
	xyserver "guanghuan.com/xiaoyao/common/server"
	//	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	batteryapi "guanghuan.com/xiaoyao/battery_app_server/business"
	//	batterydb "guanghuan.com/xiaoyao/battery_app_server/db"
	xyprofiler "guanghuan.com/xiaoyao/common/profiler"
	//xydbservice "guanghuan.com/xiaoyao/common/service/db"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
)

var (
	server *xyserver.Server
)

var (
	apiService       *batteryapi.BatteryService
	natsService      *xynatsservice.NatsService
	apnNatsService   *xynatsservice.NatsService
	alertNatsService *xynatsservice.NatsService
)

var (
	pm *xyprofiler.ProfilerMap
)
