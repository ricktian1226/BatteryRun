// A gateway server for HTTP <-> nats message routing
package main

import (
	"fmt"

	batteryapi "guanghuan.com/xiaoyao/battery_apns_server/business"
	batterydb "guanghuan.com/xiaoyao/battery_apns_server/db"
	"guanghuan.com/xiaoyao/common/apn"
	"guanghuan.com/xiaoyao/common/log"
	xyperf "guanghuan.com/xiaoyao/common/performance"
	xyserver "guanghuan.com/xiaoyao/common/server"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	"guanghuan.com/xiaoyao/common/service/timer"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xybusinesscache "guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"os"
	"runtime"
)

func banner() {
	fmt.Println("*********************************************")
	fmt.Println("*        Battery Run Apns Server            *")
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

	//fmt.Printf("DefIniConfigs : %v", DefIniConfigs)

	//性能统计模块初始化
	xyperf.InitPerf()
	xyperf.StartPerfTimer()
	xyperf.InitTimeRange()
	xyperf.StartTimeRangeTimer()

	//设置使用多cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	//启动服务器管理器
	opts := xyserver.Options{
		StopOnErr: true, //服务启动失败就退出
	}
	server = xyserver.New(DefConfig.ServerName, &opts)

	//初始化业务数据库
	initDBService(server, businessCollections)

	//初始化业务服务
	initBatteryService(server)

	//初始化nats相关的服务
	initNatsService(server)

	//初始化pprof
	xyperf.InitPProf(DefConfig.LogConfig.Path, "apns", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId)

	//启动服务器失败，直接退出
	var err error
	err = server.Start()
	if xyerror.ErrOK != err {
		fmt.Printf("server.Start failed : %v", err)
		os.Exit(-1)
		return
	}

	go server.Run()

	//服务启动，发个告警
	xybusiness.SendAlert("", "", xybusiness.ALERT_EVENT_SERVER_START, xylog.DefConfig.NodeIdentity)

	//初始化业务配置项（来源:brcommondb.apiconfig）
	if !apiService.LoadConfig() {
		xylog.ErrorNoId("LoadApiConfig failed : %v", err)
		os.Exit(-1)
	}

	initCache()

	//初始化定时器
	initTimer()

	//退出goroutine
	runtime.Goexit()
}

//初始化nats服务
func initNatsService(server *xyserver.Server) {
	//业务服务
	natsService = xynatsservice.NewNatsService("Battery api nats", DefConfig.NatsUrl)
	natsService.AddSubscriber(xybusinesscache.RES_RELOAD_API, NatsHandlerConfig)
	natsService.AddSubscriber("pprof", NatsHandlerPProf)                                                                               //所有的服务节点侦听
	natsService.AddSubscriber("pprof_apns", NatsHandlerPProf)                                                                          //所有的app节点侦听
	natsService.AddSubscriber(fmt.Sprintf("pprof_apns_%d_%d", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId), NatsHandlerPProf) //特定的app节点侦听

	xynatsservice.Nats_service = natsService //global xynatsservice.Nats_service refer to this nats_service
	server.QuickRegService(natsService)

	//apns服务
	apnNatsService = xynatsservice.NewNatsService("Battery apns nats(client app)", DefConfig.ApnNatsUrl)
	apnNatsService.AddSubscriber(xyapn.SubjectPushNotification, NatsPushNotification)
	server.QuickRegService(apnNatsService)
	xyapn.InitNatsService(apnNatsService)

	//alert服务
	alertNatsService = xynatsservice.NewNatsService("Battery alert nats(client app)", DefConfig.AlertNatsUrl)
	server.QuickRegService(alertNatsService)
	xybusiness.InitAlertNats(alertNatsService)
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

// initTimer 初始化定时器
func initTimer() {

	//抽奖机会推送定时任务
	{
		//4 test
		//batteryapi.NewXYAPI().LottoChanceNotify()

		timerOption := xytimer.TimerOption{
			Type: xytimer.TIMER_TYPE_FIXED,
		}

		moment := xytimer.TimerMoment{
			Hour:   20,
			Minute: 00,
			Second: 00,
		}
		timerOption.Moments = append(timerOption.Moments, moment)

		xytimer.InitTimer(timerOption, func() { batteryapi.NewXYAPI().LottoChanceNotify() })
	}

	//好友邮件数目推送定时任务
	{
		timerOption := xytimer.TimerOption{
			Type: xytimer.TIMER_TYPE_FIXED,
		}

		moment := xytimer.TimerMoment{
			Hour:   12,
			Minute: 00,
			Second: 00,
		}
		timerOption.Moments = append(timerOption.Moments, moment)

		////4 test
		//batteryapi.NewXYAPI().FriendMailCountNotify()

		//timerOption := xytimer.TimerOption{
		//	Type:     xytimer.TIMER_TYPE_INTERVAL,
		//	Interval: 10,
		//}

		xytimer.InitTimer(timerOption, func() { batteryapi.NewXYAPI().FriendMailCountNotify() })
	}
}

func initBatteryService(server *xyserver.Server) {
	apiService = batteryapi.NewBatteryService("Battery Service", DefConfig.Name)
	server.QuickRegService(apiService)
}

//初始化业务静态配置信息缓存
// ldb *batterydb.BatteryDB 缓存数据库指针
func initCache() {
	//获取缓存需要的操作指针
	var ldb *xybusiness.XYBusinessDB
	if dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TIP_CONFIG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN); dbInterface != nil {
		ldb = &((dbInterface.(*batterydb.BatteryDB)).XYBusinessDB)
	} else { //如果查找失败直接退出进程
		xylog.ErrorNoId("Get dbInterface for config cache failed")
		os.Exit(-1)
	}

	//初始化缓存数据库操作指针
	xybusinesscache.Init(ldb, nil)

	//初始化广告管理器
	xybusinesscache.DefTipManager.InitWhileStart()
}
