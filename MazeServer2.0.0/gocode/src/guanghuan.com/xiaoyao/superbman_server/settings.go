package main

import (
	xyversion "guanghuan.com/xiaoyao/common/version"
)

const (
	DefaultDBURL = "mongodb://117.25.150.86:20143/brdb01"
	DefaultDB    = "brdb01"
)

const (
	API_URI_AUTH                = "/v1/auth/:key"
	API_URI_GET_CONFIG          = "/config"
	API_URI_ECHO                = "/echo"
	API_URI_LOGIN               = "/v1/login/:token"
	API_URI_GET_USER_GAMEDATA   = "/v1/user/:token"
	API_URI_GET_FRIEND_GAMEDATA = "/v1/friend/:token"
	API_URI_NEWGAME             = "/v1/newgame/:token"
	API_URI_ADD_GAMEDATA        = "/v1/gameresult/:token"
	API_URI_STAMINA             = "/v1/stamina/:token"
	API_URI_GIFT_QUERY          = "/v1/gift/query/:token"  // 查询
	API_URI_GIFT_OP             = "/v1/gift/op/:token"     // 确认
	API_URI_GOODS_QUERY         = "/v1/goods/query/:token" // 查询商品列表
	API_URI_GOODS_BUY           = "/v1/goods/buy/:token"   // 购买商品
	API_URI_ORDER_NUM_REQUEST   = "/iap_verify/order_request"
	API_URI_ORDER_VERIFY        = "/iap_verify/order_verify"
	API_URI_DEVICE_ID_SUBMIT    = "/v1/device/device_id"
	API_URI_LOAD_CONFIG         = "/config/reload"
)

const (
	DefaultNatsUrl = "nats://localhost:5555"
)

var (
	ServerVersion = *xyversion.New(1, 0, 0) // 2014.5.12
)
