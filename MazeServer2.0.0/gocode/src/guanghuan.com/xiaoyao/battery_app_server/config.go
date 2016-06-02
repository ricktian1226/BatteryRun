package main

import (
	"flag"
	"fmt"
	beegoconf "github.com/astaxie/beego/config"
	//"guanghuan.com/xiaoyao/common/conf"
	xylog "guanghuan.com/xiaoyao/common/log"
	//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"os"
)

type Config struct {
	IniFile      string
	Name         string
	ServerName   string
	TestEnv      bool
	NatsUrl      string
	ApnNatsUrl   string
	AlertNatsUrl string
	Concurrent   bool
	LogConfig    *xylog.LoggerConfig
}

func (cfg *Config) String() (str string) {
	str = fmt.Sprintf(`
	IniFile    = %s
	ServerName = %s
	ConfigName = %s
	TestEnv    = %t
	---- Nats Service ----
	NatsUrl = %s
	ApnNatsUrl = %s
	AlertNatsUrl = %s
	---- Logging ---
	%s
	`, cfg.IniFile,
		cfg.ServerName,
		cfg.Name,
		cfg.TestEnv,
		cfg.NatsUrl,
		cfg.ApnNatsUrl,
		cfg.AlertNatsUrl,
		cfg.LogConfig.String(),
	)

	return
}

var (
	DefConfig = Config{
		Name:         "default",
		ServerName:   DefaultName,
		TestEnv:      true,
		NatsUrl:      DefNatsUrl,
		AlertNatsUrl: DefAlertNatsUrl,
		Concurrent:   true,
		LogConfig:    xylog.DefConfig,
	}

	DefIniConfigs beegoconf.ConfigContainer
)

func ProcessCmd() (should_continue bool) {
	//[caution!]只有各个服务节点不同的配置项才放在进程命令行中，不同服务节点相同的配置项放在配置文件中
	flag.StringVar(&DefConfig.IniFile, "config", os.Args[0]+".ini", "ini file")
	DefConfig.LogConfig.ProcessCmd()
	flag.Parse()
	should_continue = true

	return
}

func ApplyConfig() {
	xylog.ApplyConfig(DefConfig.LogConfig)
	xylog.InfoNoId("Config: %s", DefConfig.String())
}

//解析配置文件
//[caution!!!]对于所有服务节点相同的配置项放在配置文件中
func ParseConfigFile(businessCollections *xybusiness.BusinessCollections) error {
	var err error
	DefIniConfigs, err = beegoconf.NewConfig("ini", DefConfig.IniFile)
	if err != xyerror.ErrOK {
		fmt.Printf("load config file %s failed : %v\n", DefConfig.IniFile, err)
		return err
	}

	DefConfig.Name = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_CONFIG)
	DefConfig.NatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_NATSURL)
	DefConfig.ApnNatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_APNNATSURL)
	DefConfig.AlertNatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_ALERT_NATSURL)

	//解析业务数据库信息并且初始化数据库会话管理器
	xybusiness.ParseBusinessCollections(DefIniConfigs, businessCollections)

	err = DefConfig.LogConfig.ProcessIniConfig(DefIniConfigs)
	if err != xyerror.ErrOK {
		fmt.Printf("DefConfig.LogConfig.ProcessIniConfig failed : %v\n", err)
		return err

	}

	fmt.Printf("DefConfig : %v", DefConfig)

	return err
}

//初始化服务器配置项
// businessCollections *xybusiness.BusinessCollections 保存业务数据库表信息
func initServerConfig(businessCollections *xybusiness.BusinessCollections) bool {

	//从进程输入参数中读取配置项
	if !ProcessCmd() {
		return false
	}

	if err := ParseConfigFile(businessCollections); err != nil {
		return false
	}

	ApplyConfig()

	return true
}
