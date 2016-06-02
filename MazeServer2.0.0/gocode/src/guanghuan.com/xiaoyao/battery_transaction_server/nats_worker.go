package main

import (
    "encoding/binary"
    nats "github.com/nats-io/nats"
    batteryapi "guanghuan.com/xiaoyao/battery_transaction_server/business"
    xylog "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/performance"
    xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
    xycache "guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/server"
    //xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "fmt"
    //"sync"
    "time"
    "unsafe"
)

// Echo service
func NatsHandlerEcho(m *nats.Msg) {
    xylog.DebugNoId("Receive msg: %s", string(m.Data))
    nats_service.Publish(m.Reply, m.Data)
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

//重载系统配置
func NatsHandlerConfig(m *nats.Msg) {
    api_service.LoadConfig()
}

//重载商品配置信息
func NatsHandlerGoodsResReload(m *nats.Msg) {
    xycache.DefGoodsCacheManager.ReLoad()
}

//重载抽奖配置信息
func NatsHandlerLottoResReload(m *nats.Msg) {
    lottoConfig := xycache.LottoConfig{
        LottoSlotCount:     api_service.ConfigCache.Configs().LottoSlotCount,
        LottoInitUserValue: api_service.ConfigCache.Configs().LottoInitUserValue,
        LottoCostPerTime:   api_service.ConfigCache.Configs().LottoCostPerTime,
        LottoDeduct:        api_service.ConfigCache.Configs().LottoDeduct,
        SysLottoFreeCount:  api_service.ConfigCache.Configs().SysLottoFreeCount,
    }
    xycache.DefLottoCacheManager.ReLoad(lottoConfig)
}

//重载道具配置信息
func NatsHandlerPropResReload(m *nats.Msg) {
    xycache.DefPropCacheManager.ReLoad()
}

//重载符文配置信息
func NatsHandlerRuneResReload(m *nats.Msg) {
    xycache.DefRuneConfigCacheManager.ReLoad()
}

//重载签到配置信息
func NatsHandlerSignInResReload(m *nats.Msg) {
    xycache.DefSignInCacheManager.ReLoad()
}

//重载游戏前商城配置信息
func NatsHandlerBeforeGameResReload(m *nats.Msg) {
    xycache.DefBeforeGameWeightCacheManager.ReLoad()
}

//重载收集物配置信息
func NatsHandlerPickUpResReload(m *nats.Msg) {
    xycache.DefPickUpCacheManager.ReLoad()
}

//重载角色配置信息
func NatsHandlerRoleInfoResReload(m *nats.Msg) {
    xycache.DefRoleInfoCacheManager.ReLoad()
}

//重载角色加成配置信息
func NatsHandlerRoleLevelBonusResReload(m *nats.Msg) {
    xycache.DefRoleLevelBonusCacheManager.ReLoad()
}

//重载任务配置信息
func NatsHandlerMissionResReload(m *nats.Msg) {
    xycache.DefMissionCacheManager.ReLoad()
}

//重载公告配置信息
func NatsHandlerAnnouncementResReload(m *nats.Msg) {
    xycache.DefAnnouncementsCacheManager.ReLoad(true)
}

//重载登录礼包配置信息
func NatsHandlerNewAccountPropResReload(m *nats.Msg) {
    xycache.DefNewAccountPropManager.Reload()
}

//重载所有资源配置信息
func NatsHandlerAllResReload(m *nats.Msg) {

    //商城
    xycache.DefGoodsCacheManager.ReLoad()

    //抽奖信息
    lottoConfig := xycache.LottoConfig{
        LottoSlotCount:     api_service.ConfigCache.Configs().LottoSlotCount,
        LottoInitUserValue: api_service.ConfigCache.Configs().LottoInitUserValue,
        LottoCostPerTime:   api_service.ConfigCache.Configs().LottoCostPerTime,
        LottoDeduct:        api_service.ConfigCache.Configs().LottoDeduct,
        SysLottoFreeCount:  api_service.ConfigCache.Configs().SysLottoFreeCount,
    }
    xycache.DefLottoCacheManager.ReLoad(lottoConfig)

    //道具信息
    xycache.DefPropCacheManager.ReLoad()

    //签到
    xycache.DefSignInCacheManager.ReLoad()

    //符文信息
    xycache.DefRuneConfigCacheManager.ReLoad()

    //游戏前商城
    xycache.DefBeforeGameWeightCacheManager.ReLoad()

    //邮件
    xycache.DefMailConfigCacheManager.ReLoad()

    //拼图
    xycache.DefJigsawConfigCacheManager.ReLoad()

    //角色配置信息
    xycache.DefRoleInfoCacheManager.ReLoad()

    //角色等级加成配置信息
    xycache.DefRoleLevelBonusCacheManager.ReLoad()

    //任务配置信息
    xycache.DefMissionCacheManager.ReLoad()

    //公告配置信息
    xycache.DefAnnouncementsCacheManager.ReLoad(true)

    //登录礼包配置信息
    xycache.DefNewAccountPropManager.Reload()

    // 分享礼包配置信息
    xycache.DefSharedActivityManager.Reload()

    // 关卡解锁商品配置
    xycache.DefCheckPointUnlockGoodsManager.Reload()

}

// pprof消息处理函数
func NatsHandlerPProf(m *nats.Msg) {
    var op string = string(m.Data)
    xyperf.OperationPProf(op)
}

//debuguser重载消息处理函数
func NatsHandlerDebugUsers(m *nats.Msg) {
    xybusiness.LoadDebugUsers(batteryapi.NewXYAPI().GetCommonXYBusinessDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEBUGUSER))
}

func ProcessBusinessRequest(h *xynatsservice.MsgCodeHandler, reqData []byte, businessCode uint32, reply string) {
    var (
        err      error
        respData []byte
    )
    //4 test, to test if goruntine will block...
    /*if testCount == 0 {
      	xylog.Debug("process block testCount(%d)...", testCount)
      	testCount++
      	time.Sleep(1200 * time.Second)
      } else {
      	xylog.Debug("process continue  testCount(%d)...", testCount)
      }*/

    respData, err = ProcessMessage(businessCode, reqData, h.ReqType, h.RespType, h.Handler)

    if err != nil {
        xylog.ErrorNoId("Error processing message: %s", err.Error())
        // TODO: 需要添加错误返回
    }
    h.Nats.Publish(reply, respData)

}

//消息分发函数
func NatsDispatcher(m *nats.Msg) {

    reqData := m.Data

    //获取消息channel
    channel := getChannel(reqData)
    if nil == channel {
        xylog.ErrorNoId("failed to get channel.")
        return
    }

    //拼出队列消息
    sizeOfReply := uint32(len(m.Reply))
    sizeOfReplyBytes := make([]byte, unsafe.Sizeof(sizeOfReply))
    binary.LittleEndian.PutUint32(sizeOfReplyBytes, sizeOfReply)
    msgHeader := append(sizeOfReplyBytes, ([]byte(m.Reply))...)
    msg := append(msgHeader, reqData...)

    //放入队列中
    *channel <- msg

    //xylog.Debug("Push msg to channel")
}

//初始化消息队列
func initChannels() {
    //根据配置的goruntine数，初始化channel列表
    for i := uint32(0); i < api_service.ConfigCache.Configs().ChannelCount; i++ {
        initChannel(i)
    }
}

//4test
func printChanLen(index uint32, c *chan []byte) {
    for {
        time.Sleep(time.Millisecond)
        fmt.Printf("len(channel%d):%d\n", index, len(*c))
    }

}

//初始化单个消息队列
func initChannel(index uint32) {
    xylog.DebugNoId("Create channel for index[%d]", index)
    c := make(chan []byte, api_service.ConfigCache.Configs().ChannelMaxMsg)
    DefChannelMap[index] = &c

    //4test
    //go printChanLen(index, DefChannelMap[index])

    //create a go rountine to process msg
    go ProcessChannelMsg(DefChannelMap[index])
}

//根据请求的消息头获取消息队列指针
func getChannel(reqData []byte) (channel *chan []byte) {
    //get uid suffix
    suffix := binary.LittleEndian.Uint32(reqData[:4])
    //取uid后9位十进制整型的低4位作为channel node的路由依据
    index := (suffix % xybusiness.BASE_UID_BAND) % (api_service.ConfigCache.Configs().ChannelCount)
    //4test
    xylog.DebugNoId("suffix : %d, index : %d\n", suffix, index)

    //如果没找到，就初始化一个
    if _, ok := DefChannelMap[index]; !ok {
        initChannel(index)
    }

    channel = DefChannelMap[index]

    return
}

//轮训处理通道内的消息
//消息协议如下：
// replyLen    reply     uidSuffix   businessCode   reqData
//     4         n           4           4            m
//var testCount = 0

func ProcessChannelMsg(channel *chan []byte) {

    var respData []byte
    for msg := range *channel {

        //获取nats msg reply
        reply, length := getReply(msg)

        //获取业务码
        businessCode := getBusinessCode(msg[4+length+4 : 4+length+4+4])

        xylog.DebugNoId("Nats to reply : %s , businessCode : %d", reply, businessCode)

        //根据业务码获取处理函数handle信息
        h := DefMsgHandlerMap.GetHandler(businessCode)
        if h == nil {
            // TODO: 需要有办法通知 gateway，处理失败，否则客户端会一直收到没有错误的空数据包
            // err = errors.New("subject has no handler:" + subj)、
            xylog.ErrorNoId("no processor for businessCode %v.", businessCode)
            nats_service.Publish(reply, respData)
        } else {
            //go ProcessBusinessRequest(h, msg[4+length+4+4:], businessCode, reply)
            ProcessBusinessRequest(h, msg[4+length+4+4:], businessCode, reply)
        }
    }
}

func getReply(msg []byte) (reply string, length uint32) {
    length = binary.LittleEndian.Uint32(msg[:4])
    reply = string(msg[4 : 4+length])

    return
}

func getBusinessCode(businessCodeBytes []byte) (businessCode uint32) {
    businessCode = binary.LittleEndian.Uint32(businessCodeBytes)
    return
}
