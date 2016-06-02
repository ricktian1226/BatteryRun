package main

import (
    proto "code.google.com/p/goprotobuf/proto"
    "fmt"
    batteryapi "guanghuan.com/xiaoyao/battery_maintenance_server/bussiness"
    xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

var DefMsgHandlerMap xynatsservice.MsgHandlerMap = make(xynatsservice.MsgHandlerMap, 20)

func OperationCDkeyExchange(req proto.Message, resp proto.Message) (err error) {
    fmt.Println("test")
    return batteryapi.NewXYAPI().CDkeyExchange(req.(*battery.CDkeyExchangeRequest), resp.(*battery.CDkeyExchangeResponse), DefConfig.TestEnv)
}

func initMsgHandlerMap() {
    DefMsgHandlerMap.AddHandler("cdkey", battery.CDkeyExchangeRequest{}, battery.CDkeyExchangeResponse{}, OperationCDkeyExchange)
}
