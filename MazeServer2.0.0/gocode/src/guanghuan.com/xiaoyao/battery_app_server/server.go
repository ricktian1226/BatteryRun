// A gateway server for HTTP <-> nats message routing
package main

import (
	"fmt"
	"os"
	"runtime"

	batteryapi "guanghuan.com/xiaoyao/battery_app_server/business"
	batterydb "guanghuan.com/xiaoyao/battery_app_server/db"
	"guanghuan.com/xiaoyao/common/apn"
	"guanghuan.com/xiaoyao/common/log"
	xyperf "guanghuan.com/xiaoyao/common/performance"
	xyprofiler "guanghuan.com/xiaoyao/common/profiler"
	xyserver "guanghuan.com/xiaoyao/common/server"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xybusinesscache "guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func banner() {
	fmt.Println("*********************************************")
	fmt.Println("*        Battery Run App Server             *")
	fmt.Println("*                                           *")
	fmt.Println("* input -usage for help                     *")
	fmt.Println("*********************************************")
}

func main() {
	banner()

	//初始化服务配置项
	businessCollections := xybusiness.NewBusinessCollections()
	if !initServerConfig(businessCollections) {
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
	initBatteryService(server)

	//初始化nats相关的服务
	initNatsService(server)

	//初始化业务数据库
	initDBService(server, businessCollections)

	//初始化业务消息处理器
	initMsgHandlers(server)

	//初始化性能管理器
	pm = xyprofiler.NewProfilerMap(len(DefMsgHandlerMap) * 4)

	fmt.Printf("DefIniConfigs : %v", DefIniConfigs)

	//启动失败，直接退出
	err := server.Start()
	if xyerror.ErrOK != err {
		fmt.Printf("server.Start failed : %v", err)
		os.Exit(-1)
		return
	}

	//初始化业务配置项（来源:brcommondb.apiconfig）
	if apiService.LoadConfig() {
		apiService.Defsvc.Init()
	}

	//加载静态封号玩家信息
	err = batteryapi.DefBannedUserManager.Load()
	if xyerror.ErrOK != err {
		fmt.Printf("DefBannedUserManager.Load failed : %v", err)
		os.Exit(-1)
		return
	}

	//加载打开调试开关的玩家信息
	xybusiness.LoadDebugUsers(batteryapi.NewXYAPI().GetCommonXYBusinessDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEBUGUSER))

	//初始化pprof
	xyperf.InitPProf(DefConfig.LogConfig.Path, "app", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId)

	go server.Run()

	xybusiness.SendAlert("", "", xybusiness.ALERT_EVENT_SERVER_START, xylog.DefConfig.NodeIdentity)

	initCache()

	//xybusiness.DefBusinessDBSessionManager.Print()

	//退出协程
	runtime.Goexit()
}

func initNatsService(server *xyserver.Server) {
	//nats操作服务对象
	//ns := xynatsservice.NewNatsService("Battery api nats", DefConfig.NatsUrl)
	//server.QuickRegService(ns)
	//batteryapi.SetNs(ns)

	natsService = xynatsservice.NewNatsService("Battery route nats", DefConfig.NatsUrl)
	natsService.AddSubscriber(xybusinesscache.RES_RELOAD_API, NatsHandlerConfig)
	natsService.AddSubscriber("profile", NatsHandlerProfile)
	natsService.AddSubscriber("pprof", NatsHandlerPProf)                                                                              //所有的服务节点侦听
	natsService.AddSubscriber("pprof_app", NatsHandlerPProf)                                                                          //所有的app节点侦听
	natsService.AddSubscriber(fmt.Sprintf("pprof_app_%d_%d", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId), NatsHandlerPProf) //特定的app节点侦听
	natsService.AddSubscriber(batteryapi.RELOAD_BANNEDUSERS_SUBJECT, NatsHandlerBannedUser)
	natsService.AddSubscriber(xybusinesscache.RES_RELOAD_ADVERTISEMENT, NatsHandlerAdvertisementResReload)
	natsService.AddSubscriber(xybusinesscache.RES_RELOAD_ALL, NatsHandlerAllResReload)
	natsService.AddSubscriber(xylog.DEBUGUSER_SUBJECT, NatsHandlerDebugUsers)

	xynatsservice.Nats_service = natsService //global xynatsservice.Nats_service refer to this nats_service
	server.QuickRegService(natsService)

	alertNatsService = xynatsservice.NewNatsService("Battery alert nats(client app)", DefConfig.AlertNatsUrl)
	server.QuickRegService(alertNatsService)
	xybusiness.InitAlertNats(alertNatsService)

	apnNatsService = xynatsservice.NewNatsService("Battery apn nats(client app)", DefConfig.ApnNatsUrl)
	server.QuickRegService(apnNatsService)
	xyapn.InitNatsService(apnNatsService)

}

//初始化业务处理消息处理器
func initMsgHandlers(server *xyserver.Server) {
	//初始化消息处理路由信息
	initMsgHandlerMap()
	for _, h := range DefMsgHandlerMap {
		ns0 := xynatsservice.NewNatsService(fmt.Sprintf("Battery route nats-%s", h.Subject), DefConfig.NatsUrl)
		ns0.AddQueueSubscriber(h.Subject, NatsDispatcher)
		server.QuickRegService(ns0)
		h.Nats = ns0
	}
}

//根据配置项，初始化业务数据库会话信息
func initDBService(server *xyserver.Server, businessCollections *xybusiness.BusinessCollections) {
	for _, c := range *(businessCollections.Collections()) {
		dburl, dbname, platform, index := c.Detail()
		serverName := fmt.Sprintf("Battery DB %s %v", dburl, index)
		dbservice := batterydb.NewBatteryDBService(serverName, dburl, dbname)
		server.QuickRegService(dbservice)
		xybusiness.DefBusinessDBSessionManager.Insert(index, platform, dbservice.GetDB().(*batterydb.BatteryDB))
	}
}

//初始化业务静态配置信息缓存
// ldb *batterydb.BatteryDB 缓存数据库指针
func initCache() {
	//获取缓存需要的操作指针
	var ldb *xybusiness.XYBusinessDB
	if dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ADVERTISEMENT_CONFIG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN); dbInterface != nil {
		ldb = &((dbInterface.(*batterydb.BatteryDB)).XYBusinessDB)
	} else { //如果查找失败直接退出进程
		xylog.ErrorNoId("Get dbInterface for config cache failed")
		os.Exit(-1)
	}

	//初始化缓存数据库操作指针
	xybusinesscache.Init(ldb, nil)

	//初始化广告管理器
	xybusinesscache.DefAdvertisementManager.InitWhileStart()
}

func initBatteryService(server *xyserver.Server) {
	apiService = batteryapi.NewBatteryService("Battery Service", DefConfig.Name)
	server.QuickRegService(apiService)
}
