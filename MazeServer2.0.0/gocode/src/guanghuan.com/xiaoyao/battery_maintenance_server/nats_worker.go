package main

import (
	nats "github.com/nats-io/nats"

	xylog "guanghuan.com/xiaoyao/common/log"

	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"

	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

// Echo service
func NatsHandlerEcho(m *nats.Msg) {
	xylog.DebugNoId("Receive msg: %s", string(m.Data))
	nats_service.Publish(m.Reply, m.Data)
}
func ProcessBusinessRequest(h *xynatsservice.MsgHandler, req_data []byte, subj string, reply_subj string) {

	var (
		err       error
		resp_data []byte
	)

	resp_data, err = ProcessMessage(subj, req_data, h.ReqType, h.RespType, h.Handler)

	if err != xyerror.ErrOK {
		xylog.ErrorNoId("Error processing message: %s", err.Error())
		// TODO: 需要添加错误返回
	}
	h.Nats.Publish(reply_subj, resp_data)
}

func NatsDispatcher(m *nats.Msg) {
	var (
		subj      string = m.Subject
		req_data  []byte = m.Data
		resp_data []byte
		h         *xynatsservice.MsgHandler
	)

	xylog.DebugNoId("call message processor")
	h = DefMsgHandlerMap.GetHandler(subj)
	if h == nil {
		// TODO: 需要有办法通知 gateway，处理失败，否则客户端会一直收到没有错误的空数据包
		// err = errors.New("subject has no handler:" + subj)
		nats_service.Publish(m.Reply, resp_data)
	} else {

		go ProcessBusinessRequest(h, req_data, subj, m.Reply)

	}
}
