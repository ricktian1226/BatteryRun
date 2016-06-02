package main

import (
	xydb "guanghuan.com/xiaoyao/common/db"
	//	xyservice "guanghuan.com/xiaoyao/common/service"
	xydbservice "guanghuan.com/xiaoyao/common/service/db"
	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	batteryapi "guanghuan.com/xiaoyao/superbman_server/api/v2"
	batterydb "guanghuan.com/xiaoyao/superbman_server/db/v2"
)

type BatteryHttpService struct {
	xyhttpservice.MartiniService
	dbs         *xydbservice.DBService
	ns          *xynatsservice.NatsService
	config_name string
	apihandler  *batteryapi.XYAPI
	//	db         *batterydb.BatteryDB
}

func NewBatteryHttpService(name string,
	http_host string, http_port int, dburl string, dbname string,
	natsurl string, cfgname string) (svc *BatteryHttpService) {

	svc = &BatteryHttpService{
		MartiniService: *xyhttpservice.DefaultMartiniService(name, http_host, http_port),
		dbs:            xydbservice.NewXYDBService(name, dburl, dbname),
		ns:             xynatsservice.NewNatsService(name, natsurl),
		config_name:    cfgname,
	}
	return
}

func (svc *BatteryHttpService) Init() (err error) {
	svc.dbs.Init()
	svc.ns.Init()

	db := svc.dbs.GetDB()

	svc.apihandler = batteryapi.NewXYAPI(batterydb.NewBatteryDB(db.(*xydb.XYDB)), svc.config_name)

	svc.apihandler.LoadConfig()
	svc.apihandler.PrintConfig()
	svc.apihandler.SetNatsService(svc.ns)

	// login
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_LOGIN, svc.apihandler.HTTPHandlerLogin)
	//	r.Get(API_URI_HEARTBEAT, BatterRunHTTPAPI.RefreshSession)

	// 开始新游戏
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_NEWGAME, svc.apihandler.HTTPHandlerNewGame)
	// 提交游戏结果
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_ADD_GAMEDATA, svc.apihandler.HTTPHandlerAddGameData)
	// 查询当前玩家数据 (uid)
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_GET_USER_GAMEDATA, svc.apihandler.HTTPHandlerGetUserData)
	// 查询好友数据(facebook id)
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_GET_FRIEND_GAMEDATA, svc.apihandler.HTTPHandlerGetFriendData)

	// 查询体力请求服务 ((邮件服务的特例)
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_GIFT_QUERY, svc.apihandler.HTTPHandlerQueryStaminaGiftService)
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_GIFT_OP, svc.apihandler.HTTPHandlerStaminaGiftService)

	// 查询体力
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_STAMINA, svc.apihandler.HTTPHandlerQueryStaminaService)

	// 查询/购买 商品
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_GOODS_QUERY, svc.apihandler.HTTPHandlerQueryGoods)
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_GOODS_BUY, svc.apihandler.HTTPHandlerBuyGoods)
	// IAP服务端验证
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_ORDER_NUM_REQUEST, svc.apihandler.HTTPHandlerRequestOrderNum)
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_ORDER_VERIFY, svc.apihandler.HTTPHandlerVerifyOrder)

	// 设置对应设备ID
	svc.AddRouter(xyhttpservice.HttpPost, API_URI_DEVICE_ID_SUBMIT, svc.apihandler.HTTPHandlerDeviceIdSubmit)

	// 重新加载设置
	svc.AddRouter(xyhttpservice.HttpGet, API_URI_LOAD_CONFIG, svc.apihandler.HTTPHandlerReloadConfig)
	svc.MartiniService.Init()

	return
}
func (svc *BatteryHttpService) Start() (err error) {
	svc.dbs.Start()
	svc.ns.Start()
	svc.MartiniService.Start()
	return
}
func (svc *BatteryHttpService) Stop() (err error) {
	svc.MartiniService.Stop()
	svc.ns.Stop()
	svc.dbs.Stop()
	return
}
