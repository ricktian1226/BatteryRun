package main

import (
	"code.google.com/p/goprotobuf/proto"
	nats "github.com/nats-io/nats"

	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"

	business "guanghuan.com/xiaoyao/battery_app_server/business"
)

// Echo service
func NatsHandlerEcho(m *nats.Msg) {
	xylog.Debug("Receive msg: %s", string(m.Data))
	natsService.Publish(m.Reply, m.Data)
}

func NatsHandlerProfile(m *nats.Msg) {
	var op string = string(m.Data)
	switch op {
	case "start":
		pm.Start()
	case "stop":
		pm.Stop()
	case "reset":
		pm.Reset()
	default:
		xylog.InfoNoId("Profiling: %s", pm.String())
	}
}

// pprof消息处理函数
func NatsHandlerPProf(m *nats.Msg) {
	var op string = string(m.Data)
	xyperf.OperationPProf(op)
}

// 重新加载配置
func NatsHandlerConfig(m *nats.Msg) {
	apiService.LoadConfig()
}

//重新加载封号玩家信息
func NatsHandlerBannedUser(m *nats.Msg) {
	//business.DefBannedUserManager.Load()
	banUnit := &battery.MaintenanceBanUnit{}
	err := proto.Unmarshal(m.Data, banUnit)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("proto.Unmarshal failed : %v", err)
		return
	}

	business.DefBannedUserManager.ProcessBanUnit(banUnit)
}

//debuguser重载消息处理函数
func NatsHandlerDebugUsers(m *nats.Msg) {
	xybusiness.LoadDebugUsers(business.NewXYAPI().GetCommonXYBusinessDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEBUGUSER))
}

//重载公告配置信息
func NatsHandlerAdvertisementResReload(m *nats.Msg) {
	xybusinesscache.DefAdvertisementManager.Reload()
}

//重载所有资源配置信息
func NatsHandlerAllResReload(m *nats.Msg) {

	xybusiness.LoadDebugUsers(business.NewXYAPI().GetCommonXYBusinessDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEBUGUSER))

	xybusinesscache.DefAdvertisementManager.Reload()
}

func ProcessBusinessRequest(h *xynatsservice.MsgHandler, req_data []byte, subj string, reply_subj string) {

	var (
		err       error
		resp_data []byte
	)

	resp_data, err = ProcessMessage(subj, req_data, h.ReqType, h.RespType, h.Handler)

	if err != xyerror.ErrOK {
		xylog.ErrorNoId("Error processing message: %v", err)
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
		natsService.Publish(m.Reply, resp_data)
	} else {
		if DefConfig.Concurrent {
			go ProcessBusinessRequest(h, req_data, subj, m.Reply)
		} else {
			ProcessBusinessRequest(h, req_data, subj, m.Reply)
		}

	}
}
