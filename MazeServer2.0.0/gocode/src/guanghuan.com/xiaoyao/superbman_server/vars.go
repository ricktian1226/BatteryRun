package main

import (
	//	xyservice "guanghuan.com/xiaoyao/common/service"
	//	batteryapi "guanghuan.com/xiaoyao/superbman_server/api/v2"
	xydbservice "guanghuan.com/xiaoyao/common/service/db"
	//	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
)

var dbservice *xydbservice.DBService
var nats_service *xynatsservice.NatsService
var battery_service *BatteryHttpService
