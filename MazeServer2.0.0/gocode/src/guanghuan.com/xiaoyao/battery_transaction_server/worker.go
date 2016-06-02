package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	xypanic "guanghuan.com/xiaoyao/common/panic"
	xyperf "guanghuan.com/xiaoyao/common/performance"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"reflect"
	"time"
)

// 预处理
func Before(data_in []byte, pbmsg proto.Message) (err error) {
	err = proto.Unmarshal(data_in, pbmsg)
	if err != nil {
		xylog.Error("Error unmarshal: %s", err.Error())
	}
	return
}

// 后处理
func After(pbmsg proto.Message) (data_out []byte, err error) {
	data_out, err = proto.Marshal(pbmsg)
	if err != nil {
		xylog.Error("Error marshal: %s", err.Error())
	}
	return
}

// 利用relection来分发消息
//func ProcessMessage(id int32, businessCode uint32, reqData []byte, reqType reflect.Type, respType reflect.Type, handle xynatsservice.ApiHandler) (respData []byte, err error) {
func ProcessMessage(businessCode uint32, reqData []byte, reqType reflect.Type, respType reflect.Type, handle xynatsservice.ApiHandler) (respData []byte, err error) {

	begin := time.Now()
	defer xyperf.Trace(0, &begin)

	xylog.InfoNoId("[%v] Process Message start", businessCode)
	defer xylog.InfoNoId("[%v] Process Message end", businessCode)

	defer xypanic.Crash()

	req := reflect.New(reqType)
	resp := reflect.New(respType)
	err = Before(reqData, req.Interface().(proto.Message))

	if err != xyerror.ErrOK {
		xylog.ErrorNoId("Error pre-processing : %s", err.Error())
		goto ErrorHandler
	}

	err = handle(req.Interface().(proto.Message), resp.Interface().(proto.Message))

	if err != nil {
		xylog.Error("Error processing : %s", err.Error())
		goto ErrorHandler
	}

ErrorHandler:
	xylog.DebugNoId("resp for After : %v", resp)

	respData, err = After(resp.Interface().(proto.Message))

	if err != nil {
		xylog.ErrorNoId("Error post-processing : %s", err.Error())
	}

	return
}
