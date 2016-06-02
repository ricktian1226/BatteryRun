package main

import (
	"flag"
	"fmt"
	beegoconf "github.com/astaxie/beego/config"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"os"
)

type Config struct {
	Name              string
	ServerName        string
	IniFile           string
	NatsUrl           string
	ApnNatsUrl        string
	AlertNatsUrl      string
	TestEnv           bool
	HttpHost          string
	HttpPort          int
	HttpDocRoot       string
	MaxRequestSize    int64
	MaxRequest        int
	MaxTimeoutRequest int
	MaxRequestTimeout int
	LogConfig         *xylog.LoggerConfig
}

const (
	INI_CONFIG_ITEM_SERVER_HOST                = "Server::host"
	INI_CONFIG_ITEM_SERVER_PORT                = "Server::port"
	INI_CONFIG_ITEM_SERVER_DIRECTORY           = "Server::directory"
	INI_CONFIG_ITEM_SERVER_MAX_REQUEST         = "Server::maxrequest"
	INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT         = "Server::maxtimeout"
	INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT_REQUEST = "Server::maxtimeoutrequest"
)

func (cfg *Config) String() string {
	return fmt.Sprintf(`
	ServerName = %s
	ConfigName = %s
	IniFile    = %s
	TestEnv    = %t
	---- HTTP Service -----
	Host             = %s
	Port             = %d
	DocRoot          = %s
	Max Request Size = %d (bytes)
	FlowControl:
		Max Request        = %d
		Max Timeout        = %d (ms)
		Max Timout request = %d
	---- Logging ---%s
	`, cfg.ServerName,
		cfg.Name,
		cfg.IniFile,
		cfg.TestEnv,
		cfg.HttpHost,
		cfg.HttpPort,
		cfg.HttpDocRoot,
		cfg.MaxRequestSize,
		cfg.MaxRequest,
		cfg.MaxRequestTimeout,
		cfg.MaxTimeoutRequest,
		xylog.String())
}

var (
	DefConfig = Config{
		Name:           "default",
		ServerName:     "battery file server",
		TestEnv:        true,
		HttpHost:       "",
		HttpPort:       80,
		HttpDocRoot:    "httpdoc",
		MaxRequestSize: 1024 * 20,
		LogConfig:      xylog.DefConfig,
	}
	//配置项容器
	DefIniConfigs beegoconf.ConfigContainer
)

//解析命令行参数
//[caution!]只有各个服务节点不同的配置项才放在进程命令行中，不同服务节点相同的配置项放在配置文件中
func ProcessCmd() {
	flag.StringVar(&DefConfig.IniFile, "config", os.Args[0]+".ini", "ini file")
	DefConfig.LogConfig.ProcessCmd()
	flag.Parse()
}

//应用配置
func ApplyConfig() {
	xylog.ApplyConfig(DefConfig.LogConfig)
	xylog.InfoNoId("Config: %s", DefConfig.String())
}

//解析配置文件
//[caution!!!]对于所有服务节点相同的配置项放在配置文件中
func ParseConfigFile() (err error) {

	//解析文件
	DefIniConfigs, err = beegoconf.NewConfig("ini", DefConfig.IniFile)
	if err != xyerror.ErrOK {
		fmt.Printf("load config file %s failed : %v\n", DefConfig.IniFile, err)
		return
	}

	DefConfig.Name = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_CONFIG)                //服务器名称
	DefConfig.NatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_NATSURL)            //主业务nats url
	DefConfig.ApnNatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_APNNATSURL)      //apns nats url
	DefConfig.AlertNatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_ALERT_NATSURL) //alert nats url

	DefConfig.HttpHost = DefIniConfigs.String(INI_CONFIG_ITEM_SERVER_HOST)
	DefConfig.HttpDocRoot = DefIniConfigs.String(INI_CONFIG_ITEM_SERVER_DIRECTORY)
	if DefConfig.HttpPort, err = DefIniConfigs.Int(INI_CONFIG_ITEM_SERVER_PORT); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_PORT, err)
		return xyerror.ErrReadIniFile
	}
	if DefConfig.MaxRequestTimeout, err = DefIniConfigs.Int(INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT, err)
		return xyerror.ErrReadIniFile
	}
	if DefConfig.TestEnv, err = DefIniConfigs.Bool(xybusiness.INI_CONFIG_ITEM_SERVER_TEST); err != nil {
		fmt.Printf("Get %s failed : %v", xybusiness.INI_CONFIG_ITEM_SERVER_TEST, err)
		return xyerror.ErrReadIniFile
	}

	if DefConfig.MaxRequest, err = DefIniConfigs.Int(INI_CONFIG_ITEM_SERVER_MAX_REQUEST); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAX_REQUEST, err)
		return xyerror.ErrReadIniFile
	}

	if DefConfig.MaxTimeoutRequest, err = DefIniConfigs.Int(INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT_REQUEST); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT_REQUEST, err)
		return xyerror.ErrReadIniFile
	}

	if DefConfig.MaxRequestTimeout, err = DefIniConfigs.Int(INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAX_TIMEOUT, err)
		return xyerror.ErrReadIniFile
	}

	//是否是测试环境
	DefConfig.TestEnv, err = DefIniConfigs.Bool(xybusiness.INI_CONFIG_ITEM_SERVER_TEST)
	if err != xyerror.ErrOK { //未设置，默认为非测试环境
		DefConfig.TestEnv = false
	}

	//设置日志相关参数
	err = DefConfig.LogConfig.ProcessIniConfig(DefIniConfigs)
	if err != xyerror.ErrOK {
		fmt.Printf("DefConfig.LogConfig.ProcessIniConfig failed : %v\n", err)
		return
	}

	return
}

//初始化服务器配置项
// businessCollections *xybusiness.BusinessCollections 保存业务数据库表信息
func initServerConfig() bool {

	//解析命令行参数
	ProcessCmd()

	ParseConfigFile()

	ApplyConfig()

	return true
}
