package main

import (
	nats "github.com/nats-io/nats"

	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//处理告警请求
func NatsHandlerAlert(m *nats.Msg) {
	xylog.DebugNoId("receive alert msg: %v ", m)

	alert := &battery.BusinessAlert{}
	err := xybusiness.UnMarshalAlert(alert, m.Data)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("UnMarshalAlert failed")
		return
	}

	xylog.DebugNoId("[%s] alert : %v", alert.Uid, alert)

	go xybusiness.SendMail(alert.GetSubject(), alert.GetContent(), alert.GetNode(), false)
}

// Echo service
//func NatsHandlerEcho(m *nats.Msg) {
//	xylog.Debug("Receive msg: %s", string(m.Data))
//	nats_service.Publish(m.Reply, m.Data)
//}

//func NatsHandlerProfile(m *nats.Msg) {
//	var op string = string(m.Data)
//	switch op {
//	case "start":
//		pm.Start()
//	case "stop":
//		pm.Stop()
//	case "reset":
//		pm.Reset()
//	default:
//		xylog.Info("Profiling: %s", pm.String())
//	}
//}

// pprof消息处理函数
func NatsHandlerPProf(m *nats.Msg) {
	var op string = string(m.Data)
	xyperf.OperationPProf(op)
}

// 重新加载配置
//func NatsHandlerConfig(m *nats.Msg) {
//isSuccess := api_service.LoadConfig()
//xylog.Info("Reload config success: %t", isSuccess)
//}

//var (
//	cur_req   int32
//	job_id    int32
//	req_mutex sync.Mutex
//)

//type Resp struct {
//	Data []byte
//	Err  error
//}

//func ProcessBusinessRequest(h *xynatsservice.MsgHandler, req_data []byte, subj string, reply_subj string) {

//	defer xypanic.Crash()

//	var (
//		err       error
//		resp_data []byte
//	)

//	resp_data, err = ProcessMessage(subj, req_data, h.ReqType, h.RespType, h.Handler)

//	if err != xyerror.ErrOK {
//		xylog.Error("Error processing message: %s", err.Error())
//		// TODO: 需要添加错误返回
//	}
//	h.Nats.Publish(reply_subj, resp_data)
//}

//func NatsDispatcher(m *nats.Msg) {
//	var (
//		//		err       error
//		subj      string = m.Subject
//		req_data  []byte = m.Data
//		resp_data []byte
//		h         *xynatsservice.MsgHandler
//	)

//	xylog.Debug("call message processor")
//	h = DefMsgHandlerMap.GetHandler(subj)
//	if h == nil {
//		// TODO: 需要有办法通知 gateway，处理失败，否则客户端会一直收到没有错误的空数据包
//		// err = errors.New("subject has no handler:" + subj)
//		nats_service.Publish(m.Reply, resp_data)
//	} else {
//		if DefConfig.Concurrent {
//			go ProcessBusinessRequest(h, req_data, subj, m.Reply)
//		} else {
//			ProcessBusinessRequest(h, req_data, subj, m.Reply)
//		}

//	}
//}
