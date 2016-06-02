package batteryapi

import (
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

var (
	Resp_NoError          = &battery.Error{Code: battery.ErrorCode_NoError.Enum()}
	Resp_NotEnoughStamina = &battery.Error{Code: battery.ErrorCode_NotEnoughStamina.Enum()}
	Resp_NotEnoughDiamond = &battery.Error{Code: battery.ErrorCode_NotEnoughDiamond.Enum()}
	Resp_ServerError      = &battery.Error{Code: battery.ErrorCode_ServerError.Enum()}
	Resp_BadInputData     = &battery.Error{Code: battery.ErrorCode_BadInputData.Enum()}
	Resp_NotSupport       = &battery.Error{Code: battery.ErrorCode_ClientVersionNotSupport.Enum()}
)

const (
	InvalidGameID = "0"
)

const (
	MSG_VERSION_NOT_SUPPORT = "Sorry, client version %s is not support now, please update to the latest version"
)

var nats_svc *xynatsservice.NatsService

func SetNs(ns *xynatsservice.NatsService) {
	nats_svc = ns
}
