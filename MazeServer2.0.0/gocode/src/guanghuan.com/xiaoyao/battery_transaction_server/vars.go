package main

import (
	//	"sync"
	batteryapi "guanghuan.com/xiaoyao/battery_transaction_server/business"
	xyprofiler "guanghuan.com/xiaoyao/common/profiler"
	xyserver "guanghuan.com/xiaoyao/common/server"
	//xydbservice "guanghuan.com/xiaoyao/common/service/db"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
)

var (
	server *xyserver.Server
)

var (
	nats_service            *xynatsservice.NatsService
	alertNatsService        *xynatsservice.NatsService
	apnNatsService          *xynatsservice.NatsService
	iapStatisticNatsService *xynatsservice.NatsService
	api_service             *batteryapi.BatteryService
)

var (
	pm *xyprofiler.ProfilerMap
)

//保存消息的channel map
type ChannelMap map[uint32]*chan []byte

var DefChannelMap = make(ChannelMap, 0)
