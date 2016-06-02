package main

import (
    "fmt"
    _ "net/http/pprof"
    "os"
    "runtime"
    //"time"

    batteryapi "guanghuan.com/xiaoyao/battery_transaction_server/business"
    batterydb "guanghuan.com/xiaoyao/battery_transaction_server/db"
    "guanghuan.com/xiaoyao/common/apn"
    xyidgenerate "guanghuan.com/xiaoyao/common/idgenerate"
    xylog "guanghuan.com/xiaoyao/common/log"
    xyperf "guanghuan.com/xiaoyao/common/performance"
    xyserver "guanghuan.com/xiaoyao/common/server"
    xynatsservice "guanghuan.com/xiaoyao/common/service/nats"

    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xycache "guanghuan.com/xiaoyao/superbman_server/cache"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/iapstatistic"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func banner() {
    fmt.Println("*********************************************")
    fmt.Println("*     Battery Run Transaction Server            *")
    fmt.Println("*                                                                     *")
    fmt.Println("* input -usage for help                                *")
    fmt.Println("*********************************************")
}

func main() {

    banner()

    //初始化服务配置项
    businessCollections := xybusiness.NewBusinessCollections()
    initServerConfig(businessCollections)

    //初始化性能测试相关信息
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

    //初始化数据库相关的服务
    api_service = batteryapi.NewBatteryService("Battery Service", DefConfig.Name /*, dbs*/)
    server.QuickRegService(api_service)

    //初始化nats相关的服务
    initNatsService(server)

    //初始化业务数据库
    initDBService(server, businessCollections)

    err := server.Start()
    if err != xyerror.ErrOK { //启动失败，进程直接退出
        //xylog.ErrorNoId(" server.Start failed : %v", err)
        fmt.Printf(" server.Start failed : %v\n", err)
        os.Exit(-1)
        return
    }

    //xybusiness.DefBusinessDBSessionManager.Print()

    //初始化业务配置项（来源:brcommondb.apiconfig）
    if api_service.LoadConfig() {
        api_service.ApplyConfig()
        api_service.Defsvc.Init()
    }

    //初始化业务静态配置缓存信息
    initCache()

    //初始化id生成器
    initIDGenerator()

    //初始化业务消息处理器
    initMsgHandlers(server)

    //初始化消息队列
    initChannels()

    //初始化打开调试开关的玩家id信息
    xybusiness.LoadDebugUsers(batteryapi.NewXYAPI().GetCommonXYBusinessDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEBUGUSER))

    //初始化pprof工具接口
    xyperf.InitPProf(DefConfig.LogConfig.Path, "transaction", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId)

    xybusiness.SendAlert("", "", xybusiness.ALERT_EVENT_SERVER_START, xylog.DefConfig.NodeIdentity)

    go server.Run()

    runtime.Goexit()
}

func initNatsService(server *xyserver.Server) {
    nats_service = xynatsservice.NewNatsService("Battery route nats", DefConfig.NatsUrl)
    nats_service.AddSubscriber(xycache.RES_RELOAD_API, NatsHandlerConfig)
    nats_service.AddSubscriber(xycache.RES_RELOAD_LOTTO, NatsHandlerLottoResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_GOODS, NatsHandlerGoodsResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_PROP, NatsHandlerPropResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_RUNE, NatsHandlerRuneResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_SIGNIN, NatsHandlerSignInResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_BEFOREGAME, NatsHandlerBeforeGameResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_PICKUP, NatsHandlerPickUpResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_ROLE_INFO, NatsHandlerRoleInfoResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_ROLE_LEVEL_BONUS, NatsHandlerRoleLevelBonusResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_ANNOUNCEMENT, NatsHandlerAnnouncementResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_MISSION, NatsHandlerMissionResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_NEWACCOUNTPROP, NatsHandlerNewAccountPropResReload)
    nats_service.AddSubscriber(xycache.RES_RELOAD_ALL, NatsHandlerAllResReload)
    nats_service.AddSubscriber("profile", NatsHandlerProfile)
    nats_service.AddSubscriber("pprof", NatsHandlerPProf)                                                                                      //所有的服务节点侦听
    nats_service.AddSubscriber("pprof_transaction", NatsHandlerPProf)                                                                          //所有的app节点侦听
    nats_service.AddSubscriber(fmt.Sprintf("pprof_transaction_%d_%d", DefConfig.LogConfig.DCId, DefConfig.LogConfig.NodeId), NatsHandlerPProf) //特定的app节点侦听
    nats_service.AddSubscriber(xylog.DEBUGUSER_SUBJECT, NatsHandlerDebugUsers)
    server.QuickRegService(nats_service)

    alertNatsService = xynatsservice.NewNatsService("Battery alert nats", DefConfig.AlertNatsUrl)
    server.QuickRegService(alertNatsService)
    xybusiness.InitAlertNats(alertNatsService)

    apnNatsService = xynatsservice.NewNatsService("Battery apn nats", DefConfig.ApnNatsUrl)
    server.QuickRegService(apnNatsService)
    xyapn.InitNatsService(apnNatsService)

    iapStatisticNatsService = xynatsservice.NewNatsService("Battery iapstatistic nats", DefConfig.IapStatisticNatsUrl)
    server.QuickRegService(iapStatisticNatsService)
    xyiapstatistic.InitNatsService(iapStatisticNatsService)
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
    var ldb *xybusiness.XYBusinessDB
    if dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ANNOUNCEMENT_CONFIG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN); dbInterface != nil {
        ldb = &((dbInterface.(*batterydb.BatteryDB)).XYBusinessDB)
    } else { // 如果查找失败直接退出进程
        xylog.ErrorNoId("Get dbInterface for config cache failed")
        os.Exit(-1)
    }

    var ldbCheckPoint *xybusiness.XYBusinessDB
    if dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN); dbInterface != nil {
        ldbCheckPoint = &((dbInterface.(*batterydb.BatteryDB)).XYBusinessDB)
    } else { // 如果查找失败直接退出进程
        xylog.ErrorNoId("Get dbInterface(usercheckpoint) for cache failed")
        os.Exit(-1)
    }

    // 初始化缓存数据库操作指针
    xycache.Init(ldb, ldbCheckPoint)

    // 初始化抽奖缓存信息
    lottoConfig := xycache.LottoConfig{
        LottoSlotCount:     api_service.ConfigCache.Configs().LottoSlotCount,
        LottoInitUserValue: api_service.ConfigCache.Configs().LottoInitUserValue,
        LottoCostPerTime:   api_service.ConfigCache.Configs().LottoCostPerTime,
        LottoDeduct:        api_service.ConfigCache.Configs().LottoDeduct,
        SysLottoFreeCount:  api_service.ConfigCache.Configs().SysLottoFreeCount,
    }
    xycache.DefLottoCacheManager.InitWhileStart(lottoConfig)

    // 初始化道具缓存信息
    xycache.DefPropCacheManager.InitWhileStart()

    // 初始化黑名单信息
    xycache.DefRankBlacklistManager.InitWhileStart()

    // 初始化签到活动缓存信息
    xycache.DefSignInCacheManager.InitWhileStart()

    // 初始化商品缓存信息
    xycache.DefGoodsCacheManager.InitWhileStart()

    // 系统道具缓存信息
    xycache.DefRuneConfigCacheManager.InitWhileStart()

    // 初始化记忆点全局排行榜缓存信息
    checkPointConfig := &xycache.CheckPointConfig{
        IdNum:                api_service.ConfigCache.Configs().CheckPointIdNum,
        GlobalRankReLoadSecs: api_service.ConfigCache.Configs().CheckPointGlobalRankReLoadSecs,
        GlobalRankSize:       api_service.ConfigCache.Configs().CheckPointGlobalRankSize,
    }

    xycache.DefCheckPointGlobalRankManager.InitWhileStart(checkPointConfig)

    globalConfig := &xycache.RankListConfig{
        GlobalRankReLoadSecs: api_service.ConfigCache.Configs().GlobalRankReLoadSecs,
        GlobalRankSize:       api_service.ConfigCache.Configs().GlobalRankSize,
    }
    // 初始化全局排行榜
    xycache.DefGlobalRankListManager.InitWhileStart(globalConfig)
    // 初始化游戏前商店缓存信息
    xycache.DefBeforeGameWeightCacheManager.InitWhileStart()

    // 拼图
    xycache.DefJigsawConfigCacheManager.InitWhileStart()

    // 邮件
    xycache.DefMailConfigCacheManager.InitWhileStart()

    // 收集物缓存初始化
    xycache.DefPickUpCacheManager.InitWhileStart()

    // 任务缓存初始化
    xycache.DefMissionCacheManager.InitWhileStart()

    // 角色加成缓存初始化
    xycache.DefRoleLevelBonusCacheManager.InitWhileStart()

    // 角色信息缓存初始化
    xycache.DefRoleInfoCacheManager.InitWhileStart()

    // 公告信息缓存初始化
    xycache.DefAnnouncementsCacheManager.InitWhileStart()

    // 服务节点玩家计数初始化
    xycache.DefUserIdentityManager.InitWhileStart(xylog.DefConfig.DCId, xylog.DefConfig.NodeId)

    //登录礼包
    xycache.DefNewAccountPropManager.InitWhileStart()

    xycache.DefSharedActivityManager.InitWhileStart()
    xycache.DefCheckPointUnlockGoodsManager.InitWhileStart()
}

//初始化id生成器
func initIDGenerator() {
    //抽奖id生成器
    batteryapi.DefLottoIdGenerater = xyidgenerate.NewIdGenerater(xyidgenerate.DefIdGenerateBeginTimeStamp, int64(DefConfig.LogConfig.DCId), int64(DefConfig.LogConfig.NodeId), "lotto")
    //mission id生成器
    batteryapi.DefMissionIdGenerater = xyidgenerate.NewIdGenerater(xyidgenerate.DefIdGenerateBeginTimeStamp, int64(DefConfig.LogConfig.DCId), int64(DefConfig.LogConfig.NodeId), "mission")
    //购买票据id生成器
    batteryapi.DefReceiptIdGenerater = xyidgenerate.NewIdGenerater(xyidgenerate.DefIdGenerateBeginTimeStamp, int64(DefConfig.LogConfig.DCId), int64(DefConfig.LogConfig.NodeId), "receipt")
    //公告id生成器
    batteryapi.DefAnnouncementIdGenerater = xyidgenerate.NewIdGenerater(xyidgenerate.DefIdGenerateBeginTimeStamp, int64(DefConfig.LogConfig.DCId), int64(DefConfig.LogConfig.NodeId), "announcement")

}

//初始化业务处理消息处理器
func initMsgHandlers(server *xyserver.Server) {

    initMsgHandlerMap()

    subject := fmt.Sprintf("transaction%03d", DefConfig.LogConfig.NodeId)
    ns0 := xynatsservice.NewNatsService(fmt.Sprintf("Battery route nats-%s", subject), DefConfig.NatsUrl)
    ns0.AddQueueSubscriber(subject, NatsDispatcher)
    server.QuickRegService(ns0)

    for _, h := range DefMsgHandlerMap {
        h.Nats = ns0
    }
}
