package main

const (
	DefaultName     = "App server"
	DefNatsUrl      = "nats://localhost:5555"
	DefApnNatsUrl   = "nats://localhost:5556"
	DefAlertNatsUrl = "nats://localhost:5557"
	DefHttpHost     = ""
	DefHttpPort     = 8080
)
const (
	DefaultDBURL        = "mongodb://localhost:27017/briosdb"
	DefaultDB           = "briosdb"
	DefaultAndroidDBURL = "mongodb://localhost:27017/brandroiddb"
	DefaultAndroidDB    = "brandroiddb"
	DefaultCommonDBURL  = "mongodb://localhost:27018/brcommondb"
	DefaultCommonDB     = "brcommondb"
	DefaultLogDBURL     = "mongodb://localhost:27018/brlogdb"
	DefaultLogDB        = "brlogdb"
)

const (
	//开发版本
	APN_CERT_DEVELOPMENT             = "aps_superbman_development.cer"
	APN_KEY_DEVELOPMENT              = "aps_superbman_dev.pem"
	APN_GATEWAY_DEVELOPMENT          = "gateway.sandbox.push.apple.com:2195"
	APN_GATEWAY_FEEDBACK_DEVELOPMENT = "feedback.sandbox.push.apple.com:2196"

	//生产环境
	APN_CERT_PRODUCTION             = "aps_superbman_production.cer"
	APN_KEY_PRODUCTION              = "aps_superbman_production.pem"
	APN_GATEWAY_PRODUCTION          = "gateway.push.apple.com:2195"
	APN_GATEWAY_FEEDBACK_PRODUCTION = "feedback.push.apple.com:2196"
)

const (
	//APN相关配置项
	INI_CONFIG_ITEM_APN_CERT_PRODUCTION  = "Apn::certproduction"
	INI_CONFIG_ITEM_APN_KEY_PRODUCTION   = "Apn::keyproduction"
	INI_CONFIG_ITEM_APN_CERT_DEVELOPMENT = "Apn::certdevelopment"
	INI_CONFIG_ITEM_APN_KEY_DEVELOPMENT  = "Apn::keydevelopment"
)
