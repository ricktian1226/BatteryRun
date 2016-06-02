package main

import (
	//	"sync"
	batteryapi "guanghuan.com/xiaoyao/battery_maintenance_server/bussiness"
	xyprofiler "guanghuan.com/xiaoyao/common/profiler"
	xyserver "guanghuan.com/xiaoyao/common/server"
	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
)

var (
	server *xyserver.Server
)

var (
	apiService       *batteryapi.BatteryService
	nats_service     *xynatsservice.NatsService
	alertNatsService *xynatsservice.NatsService
	martini_service  *xyhttpservice.MartiniService
)

var (
	pm = xyprofiler.NewProfilerMap(20)
)

//var (
//	//	ApnMutex  sync.Mutex
//	MsgMutex        sync.Mutex
//	LastMsgId       MessageId = 0
//	NotificationMap NoteMap   = make(NoteMap)
//)
