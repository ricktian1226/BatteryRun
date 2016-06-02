package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	xypanic "guanghuan.com/xiaoyao/common/panic"
	xyperf "guanghuan.com/xiaoyao/common/performance"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
	"reflect"
	"time"
)

// 预处理
func Before(data_in []byte, pbmsg proto.Message) (err error) {
	data_in, err = crypto.Decrypt(data_in)
	if err != nil {
		xylog.ErrorNoId("Error decrypt: %s", err.Error())
	} else {
		err = proto.Unmarshal(data_in, pbmsg)
		if err != nil {
			xylog.ErrorNoId("Error unmarshal: %s", err.Error())
		} else {
			//			xylog.Debug("req: %s", pbmsg.String())
		}
	}
	return
}

// 后处理
func After(pbmsg proto.Message) (data_out []byte, err error) {
	//	xylog.Debug("resp: %s", pbmsg.String())
	data_out, err = proto.Marshal(pbmsg)
	if err != nil {
		xylog.ErrorNoId("Error marshal: %s", err.Error())
	} else {
		data_out, err = crypto.Encrypt(data_out)
		if err != nil {
			xylog.ErrorNoId("Error encrypt: %s", err.Error())
		}
	}

	return
}

// 利用relection来分发消息
//func ProcessMessage(id int32, subject string, req_data []byte, req_type reflect.Type, resp_type reflect.Type, handle xynatsservice.ApiHandler) (resp_data []byte, err error) {
func ProcessMessage(subject string, req_data []byte, req_type reflect.Type, resp_type reflect.Type, handle xynatsservice.ApiHandler) (resp_data []byte, err error) {

	begin := time.Now()
	defer xyperf.Trace(xyperf.DefLogId, &begin)

	xylog.InfoNoId("[%s] Process Message start", subject)
	defer xylog.InfoNoId("[%s] Process Message end", subject)

	defer xypanic.Crash()

	req := reflect.New(req_type)
	resp := reflect.New(resp_type)

	err = Before(req_data, req.Interface().(proto.Message))

	if err != nil {
		xylog.ErrorNoId("Error pre-processing : %s", err.Error())
		goto ErrorHandler
	}

	err = handle(req.Interface().(proto.Message), resp.Interface().(proto.Message))

	if err != nil {
		xylog.ErrorNoId("Error processing : %s", err.Error())
		goto ErrorHandler
	}

ErrorHandler:

	resp_data, err = After(resp.Interface().(proto.Message))

	if err != nil {
		xylog.ErrorNoId("Error post-processing : %s", err.Error())

	}

	return
}
