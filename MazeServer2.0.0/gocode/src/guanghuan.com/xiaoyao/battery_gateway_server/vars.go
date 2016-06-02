package main

import (
	//	"sync"
	xyprofiler "guanghuan.com/xiaoyao/common/profiler"
	xyserver "guanghuan.com/xiaoyao/common/server"
	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
)

var (
	server *xyserver.Server
)

var (
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
