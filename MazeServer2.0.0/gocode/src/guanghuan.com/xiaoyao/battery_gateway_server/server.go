// A gateway server for HTTP <-> nats message routing
package main

import (
	"fmt"
	"os"
	"runtime"
	//"time"

	xylog "guanghuan.com/xiaoyao/common/log"
	xyperf "guanghuan.com/xiaoyao/common/performance"
	//"guanghuan.com/xiaoyao/common/pool"
	xyserver "guanghuan.com/xiaoyao/common/server"
	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func banner() {
	fmt.Println("*********************************************")
	fmt.Println("*      Battery Run Gateway Server           *")
	fmt.Println("*                                           *")
	fmt.Println("* input -usage for help                     *")
	fmt.Println("*********************************************")
}

func init() {
	InitRouteTable()
}
func main() {
	banner()
	//ProcessCmd()
	//ApplyConfig()

	//设置使用多cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始化服务配置项
	initServerConfig()

	//perf trace init
	xyperf.InitPerf()
	xyperf.StartPerfTimer()
	xyperf.InitTimeRange()
	xyperf.StartTimeRangeTimer()

	// starting server
	opts := xyserver.Options{}
	server = xyserver.New(DefConfig.ServerName, &opts)

	// init nats service
	initNatsService(server)

	// init a martini service
	initMartiniService(server)

	if xyerror.ErrOK != server.Start() {
		os.Exit(-1)
		return
	}

	xybusiness.SendAlert("", "", xybusiness.ALERT_EVENT_SERVER_START, xylog.DefConfig.NodeIdentity)

	//初始化pprof
	xyperf.InitPProf(xylog.DefConfig.Path, "gateway", xylog.DefConfig.DCId, xylog.DefConfig.NodeId)

	go server.Run()
	runtime.Goexit()
}

// initNatsService 初始化业务nats服务
func initNatsService(server *xyserver.Server) {
	//业务 nats service
	nats_service = server.EnableNatsService(DefConfig.NatsUrl)
	nats_service.AddSubscriber("pprof", NatsHandlerPProf)                                                                          //所有的服务节点侦听
	nats_service.AddSubscriber("pprof_gateway", NatsHandlerPProf)                                                                  //所有的app节点侦听
	nats_service.AddSubscriber(fmt.Sprintf("pprof_gateway_%d_%d", xylog.DefConfig.DCId, xylog.DefConfig.NodeId), NatsHandlerPProf) //特定的app节点侦听

	//告警 nats service
	alertNatsService = server.EnableNatsService(DefConfig.AlertNatsUrl)
	xybusiness.InitAlertNats(alertNatsService)
}

// initMartiniService
func initMartiniService(server *xyserver.Server) {
	martini_service = server.EnableMartiniService(DefConfig.HttpHost, DefConfig.HttpPort)
	if DefConfig.MaxRequest > 0 {
		martini_service.EnableFlowControl(DefConfig.MaxRequest, DefConfig.MaxRequestTimeout, DefConfig.MaxTimeoutRequest)
	}

	martini_service.AddRouter(xyhttpservice.HttpGet, "/config/reload", HttpApiConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/lotto/reload", HttpLottoConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/prop/reload", HttpPropConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/rune/reload", HttpRuneConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/mission/reload", HttpMissionConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/signin/reload", HttpSigninConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/goods/reload", HttpGoodsConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/beforegame/reload", HttpBeforeGameConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/pickup/reload", HttpPickUpConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/roleinfo/reload", HttpRoleInfoConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/rolelevelbonus/reload", HttpRoleLevelBonusConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/announcement/reload", HttpAnnouncementConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/advertisement/reload", HttpAdvertisementConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/tip/reload", HttpTipConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/newaccountprop/reload", HttpNewAccountPropConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/allres/reload", HttpAllConfigReload)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/profile/:op", HttpConfigProfile)
	//pprof调试开关消息接口，格式为 /pprof/:appname/:dcid/:nodeid/:op
	// appname 服务名称，[app|transaction]
	// dcid    dc标识，[0...N],如果全部打开，用all
	// nodeid  服务节点标识，[0...N],如果全部打开，用all
	// op      操作标识，[lookupheap|lookupthreadcreate|lookupblock|mem|cpustart|cpustop]
	martini_service.AddRouter(xyhttpservice.HttpGet, "/pprof/:appname/:dcid/:nodeid/:op", HttpPProf)
	martini_service.AddRouter(xyhttpservice.HttpGet, "/debugusers/reload", HttpDebugUsersReload)

	for _, route := range DefHttpPostTable {
		martini_service.AddRouter(xyhttpservice.HttpPost, route.GetHttpUri(), HttpPostWorker)
		xylog.DebugNoId("Adding router: %s -> %s", route.GetHttpUri(), route.GetNatsSubject())
	}
	for _, route := range DefHttpPostNoTokenTable {
		martini_service.AddRouter(xyhttpservice.HttpPost, route.GetHttpUri(), HttpPostWorkerNoToken)
		xylog.DebugNoId("Adding router: %s -> %s", route.GetHttpUri(), route.GetNatsSubject())
	}
}
