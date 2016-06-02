// A gateway server for HTTP <-> nats message routing
package main

import (
	"fmt"
	//batteryapi "guanghuan.com/xiaoyao/battery_app_server/business"
	//batterydb "guanghuan.com/xiaoyao/battery_app_server/db"
	xyperf "guanghuan.com/xiaoyao/common/performance"
	//xyprofiler "guanghuan.com/xiaoyao/common/profiler"
	xyserver "guanghuan.com/xiaoyao/common/server"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	//xycache "guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"os"
	"runtime"
)

func banner() {
	fmt.Println("*********************************************")
	fmt.Println("*        Battery Run Mail Server             *")
	fmt.Println("*                                           *")
	fmt.Println("* input -usage for help                     *")
	fmt.Println("*********************************************")
}

func main() {
	banner()

	//初始化服务配置项
	if !initServerConfig() {
		os.Exit(-1)
		return
	}

	//perf trace init
	xyperf.InitPerf()
	xyperf.StartPerfTimer()
	xyperf.InitTimeRange()
	xyperf.StartTimeRangeTimer()

	//设置使用多cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	opts := xyserver.Options{
		StopOnErr: true, //服务启动失败就退出
	}

	server = xyserver.New(DefConfig.ServerName, &opts)

	//业务服务基础api类
	//api_service = batteryapi.NewBatteryService("Battery Service", DefConfig.Name)
	//server.QuickRegService(api_service)

	//初始化nats相关的服务
	initNatsService(server)

	//初始化业务操作相关信息
	if xyerror.ErrOK != xybusiness.Init(&DefIniConfigs) {
		os.Exit(-1)
		return
	}

	xybusiness.InitAlertNats(nats_service)

	//启动服务
	if xyerror.ErrOK != server.Start() { //启动失败，直接退出
		os.Exit(-1)
		return
	}

	//初始化业务配置项（来源:brcommondb.apiconfig）
	//if api_service.LoadConfig() {
	//	api_service.Defsvc.Init()
	//}

	//初始化pprof
	xyperf.InitPProf(DefConfig.LogConfig.Path, "app", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId)

	go server.Run()

	//退出协程
	runtime.Goexit()
}

func initNatsService(server *xyserver.Server) {
	nats_service = xynatsservice.NewNatsService("Battery alert nats(client mail)", DefConfig.AlertNatsUrl)
	nats_service.AddSubscriber(xybusiness.ALERT_API, NatsHandlerAlert)
	nats_service.AddSubscriber("pprof", NatsHandlerPProf)                                                                               //所有的服务节点侦听
	nats_service.AddSubscriber("pprof_mail", NatsHandlerPProf)                                                                          //所有的app节点侦听
	nats_service.AddSubscriber(fmt.Sprintf("pprof_mail_%d_%d", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId), NatsHandlerPProf) //特定的app节点侦听

	xynatsservice.Nats_service = nats_service //global xynatsservice.Nats_service refer to this nats_service
	server.QuickRegService(nats_service)
}

//初始化业务处理消息处理器
//func initMsgHandlers(server *xyserver.Server) {
//	//初始化消息处理路由信息
//	//initMsgHandlerMap()
//	for _, h := range DefMsgHandlerMap {
//		ns0 := xynatsservice.NewNatsService(fmt.Sprintf("Battery route nats-%s", h.Subject), DefConfig.NatsUrl)
//		ns0.AddQueueSubscriber(h.Subject, NatsDispatcher)
//		server.QuickRegService(ns0)
//		h.Nats = ns0
//	}
//}

//根据配置项，初始化业务数据库会话信息
//func initDBService(server *xyserver.Server, businessCollections *xybusiness.BusinessCollections) {
//	for k, c := range *(businessCollections.Collections()) {
//		dburl, dbname, platform := c.Detail()
//		serverName := fmt.Sprintf("Battery DB %s %v", dburl, k)
//		dbservice := batterydb.NewBatteryDBService(serverName, dburl, dbname)
//		server.QuickRegService(dbservice)
//		xybusiness.DefBusinessDBSessionManager.Insert(k, platform, dbservice.GetDB().(*batterydb.BatteryDB))
//	}
//}
