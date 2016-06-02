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

const (
	//	INI_CONFIG_ITEM_SERVER_CONFIG            = "Server::config"
	//	INI_CONFIG_ITEM_SERVER_NATSURL           = "Server::natsurl"
	INI_CONFIG_ITEM_SERVER_HOST              = "Server::host"
	INI_CONFIG_ITEM_SERVER_PORT              = "Server::port"
	INI_CONFIG_ITEM_SERVER_NATSTIMEOUT       = "Server::natstimeout"
	INI_CONFIG_ITEM_SERVER_TESTENV           = "Server::test"
	INI_CONFIG_ITEM_SERVER_MAXREQUEST        = "Server::maxrequest"
	INI_CONFIG_ITEM_SERVER_MAXTIMEOUTREQUEST = "Server::maxtimeoutrequest"
	INI_CONFIG_ITEM_SERVER_MAXREQUESTTIMEOUT = "Server::maxrequesttimeout"
	INI_CONFIG_ITEM_SERVER_MAXGOROUTINE      = "Server::maxgoroutine"
)

type Config struct {
	IniFile           string
	Name              string
	ServerName        string
	TestEnv           bool
	NatsUrl           string
	AlertNatsUrl      string
	NatsTimeout       int
	HttpHost          string
	HttpPort          int
	MaxRequestSize    int64
	MaxRequest        int
	MaxTimeoutRequest int
	MaxRequestTimeout int
	MaxGoroutine      int
}

func (cfg *Config) ToString() (str string) {
	str = fmt.Sprintf(`
	IniFile    = %s
	ServerName = %s
	ConfigName = %s
	TestEnv    = %t
	---- Nats Service ---
	NatsUrl    = %s
	AlertNatsUrl = %s
	Timeout    = %d sec
	---- HTTP Service -----
	Host             = %s
	Port             = %d
	Max Request Size = %d (bytes)
	FlowControl:
		Max Request        = %d
		Max Timeout        = %d (ms)
		Max Timout request = %d
		Max Goroutine = %d
	---- Logging ---%s
	`, cfg.IniFile,
		cfg.ServerName,
		cfg.Name,
		cfg.TestEnv,
		cfg.NatsUrl,
		cfg.AlertNatsUrl,
		cfg.NatsTimeout,
		cfg.HttpHost,
		cfg.HttpPort,
		cfg.MaxRequestSize,
		cfg.MaxRequest,
		cfg.MaxRequestTimeout,
		cfg.MaxTimeoutRequest,
		cfg.MaxGoroutine,
		xylog.String())

	return
}

var (
	DefConfig = Config{
		Name:           "default",
		ServerName:     DefaultName,
		TestEnv:        true,
		NatsUrl:        DefNatsUrl,
		NatsTimeout:    10,
		HttpHost:       DefHttpHost,
		HttpPort:       DefHttpPort,
		MaxRequestSize: 1024 * 20,
		MaxGoroutine:   10000,
	}
)

func ProcessCmd() (should_continue bool) {
	//[caution!]只有各个服务节点不同的配置项才放在进程命令行中，不同服务节点相同的配置项放在配置文件中
	flag.StringVar(&DefConfig.IniFile, "config", os.Args[0]+".ini", "ini file")
	xylog.ProcessCmd()
	flag.Parse()
	should_continue = true
	return
}

func ApplyConfig() {
	xylog.ApplyConfig(nil)
	xylog.InfoNoId("Config: %s", DefConfig.ToString())

}

//解析配置文件
//[caution!!!]对于所有服务节点相同的配置项放在配置文件中
func ParseConfigFile() error {
	//fileName := os.Args[0] + ".ini"
	configs, err := beegoconf.NewConfig("ini", DefConfig.IniFile)
	if err != xyerror.ErrOK {
		fmt.Printf("load config file %s failed : %v\n", DefConfig.IniFile, err)
		return err
	}

	DefConfig.Name = configs.String(xybusiness.INI_CONFIG_ITEM_SERVER_CONFIG)
	DefConfig.NatsUrl = configs.String(xybusiness.INI_CONFIG_ITEM_SERVER_NATSURL)
	DefConfig.AlertNatsUrl = configs.String(xybusiness.INI_CONFIG_ITEM_SERVER_ALERT_NATSURL)
	DefConfig.HttpHost = configs.String(INI_CONFIG_ITEM_SERVER_HOST)
	if DefConfig.HttpPort, err = configs.Int(INI_CONFIG_ITEM_SERVER_PORT); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_PORT, err)
		return xyerror.ErrReadIniFile
	}
	if DefConfig.NatsTimeout, err = configs.Int(INI_CONFIG_ITEM_SERVER_NATSTIMEOUT); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_NATSTIMEOUT, err)
		return xyerror.ErrReadIniFile
	}
	if DefConfig.TestEnv, err = configs.Bool(INI_CONFIG_ITEM_SERVER_TESTENV); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_TESTENV, err)
		return xyerror.ErrReadIniFile
	}

	if DefConfig.MaxRequest, err = configs.Int(INI_CONFIG_ITEM_SERVER_MAXREQUEST); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAXREQUEST, err)
		return xyerror.ErrReadIniFile
	}

	if DefConfig.MaxTimeoutRequest, err = configs.Int(INI_CONFIG_ITEM_SERVER_MAXTIMEOUTREQUEST); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAXTIMEOUTREQUEST, err)
		return xyerror.ErrReadIniFile
	}

	if DefConfig.MaxRequestTimeout, err = configs.Int(INI_CONFIG_ITEM_SERVER_MAXREQUESTTIMEOUT); err != nil {
		fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAXREQUESTTIMEOUT, err)
		return xyerror.ErrReadIniFile
	}

	//if DefConfig.MaxGoroutine, err = configs.Int(INI_CONFIG_ITEM_SERVER_MAXGOROUTINE); err != nil {
	//	fmt.Printf("Get %s failed : %v", INI_CONFIG_ITEM_SERVER_MAXGOROUTINE, err)
	//	return xyerror.ErrReadIniFile
	//}

	err = xylog.DefConfig.ProcessIniConfig(configs)
	if err != xyerror.ErrOK {
		fmt.Printf("DefConfig.LogConfig.ProcessIniConfig failed : %v\n", err)
		return err

	}

	fmt.Printf("DefConfig : %v", DefConfig)

	return err
}

//初始化服务器配置项
// businessCollections *xybusiness.BusinessCollections 保存业务数据库表信息
func initServerConfig() bool {

	//从进程输入参数中读取配置项
	if !ProcessCmd() {
		return false
	}

	if err := ParseConfigFile(); err != nil {
		return false
	}

	ApplyConfig()

	return true
}
