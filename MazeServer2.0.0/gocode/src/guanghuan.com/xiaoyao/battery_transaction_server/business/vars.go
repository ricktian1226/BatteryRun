package batteryapi

import (
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
)

var nats_svc *xynatsservice.NatsService

func SetNs(ns *xynatsservice.NatsService) {
	nats_svc = ns
}

const (
	InvalidGameID                = "0"
	OneDaySeconds                = 24 * 60 * 60 //一天的秒数
	DefaultStaminaLastUpdateTime = int64(-1)
)

const (
	MSG_VERSION_NOT_SUPPORT = "Sorry, client version %s is not support now, please update to the latest version"
)
