// A gateway server for HTTP <-> nats message routing
package main

import (
    batteryapi "guanghuan.com/xiaoyao/battery_maintenance_server/bussiness"
    batterydb "guanghuan.com/xiaoyao/battery_maintenance_server/db"
    xylog "guanghuan.com/xiaoyao/common/log"
    xyperf "guanghuan.com/xiaoyao/common/performance"
    xyserver "guanghuan.com/xiaoyao/common/server"
    xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
    xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"

    "fmt"
    "os"
    "runtime"
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

    //初始化服务配置项
    businessCollections := xybusiness.NewBusinessCollections()
    initServerConfig(businessCollections)

    //perf trace init
    xyperf.InitPerf()
    xyperf.StartPerfTimer()
    xyperf.InitTimeRange()
    xyperf.StartTimeRangeTimer()

    // starting server
    opts := xyserver.Options{}
    server = xyserver.New(DefConfig.ServerName, &opts)

    // create a nats service
    nats_service = server.EnableNatsService(DefConfig.NatsUrl)
    alertNatsService = server.EnableNatsService(DefConfig.AlertNatsUrl)
    initMsgHandlers(server)

    // create a martini service
    martini_service = server.EnableMartiniService(DefConfig.HttpHost, DefConfig.HttpPort)
    if DefConfig.MaxRequest > 0 {
        martini_service.EnableFlowControl(DefConfig.MaxRequest, DefConfig.MaxRequestTimeout, DefConfig.MaxTimeoutRequest)
    }

    //pprof调试开关消息接口，格式为 /pprof/:appname/:dcid/:nodeid/:op
    // appname 服务名称，[app|transaction]
    // dcid    dc标识，[0...N],如果全部打开，用all
    // nodeid  服务节点标识，[0...N],如果全部打开，用all
    // op      操作标识，[lookupheap|lookupthreadcreate|lookupblock|mem|cpustart|cpustop]
    martini_service.AddRouter(xyhttpservice.HttpGet, "/pprof/:appname/:dcid/:nodeid/:op", HttpPProf)

    //初始化运营数据接口，格式为 /maintenance/:content
    // content 请求内容
    martini_service.AddRouter(xyhttpservice.HttpGet, "/maintenance/:content", HttpGetWorker)

    martini_service.AddRouter(xyhttpservice.HttpPost, "/sdk/callback", HttpSDKCallBack) // 支付回调接口

    // 运营数据请求接口
    initMartiniRouter()

    // 初始化数据库接口
    initDBService(server, businessCollections)

    xybusiness.InitAlertNats(alertNatsService)

    if xyerror.ErrOK != server.Start() {
        os.Exit(-1)
        return
    }

    xybusiness.SendAlert("", "", xybusiness.ALERT_EVENT_SERVER_START, xylog.DefConfig.NodeIdentity)
    initBatteryService(server)

    //初始化业务配置项（来源:brcommondb.apiconfig）
    if apiService.LoadConfig() {
        apiService.Defsvc.Init()
    }
    go server.Run()
    runtime.Goexit()
}

func initMartiniRouter() {

    //	martini_service.AddRouter(xyhttpservice.HttpPost, "/maintenancereq/cdkey", HttpCDkeyExchange)

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

func initBatteryService(server *xyserver.Server) {
    apiService = batteryapi.NewBatteryService("Battery Service", DefConfig.Name)
    server.QuickRegService(apiService)
}
