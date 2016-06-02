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

//配置项结构体定义
type Config struct {
	IniFile             string
	Name                string
	ServerName          string
	NatsUrl             string
	ApnNatsUrl          string
	AlertNatsUrl        string
	IapStatisticNatsUrl string
	LogConfig           *xylog.LoggerConfig
}

func (cfg *Config) String() (str string) {
	return fmt.Sprintf(`
	IniFile    = %s
	ServerName = %s
	ConfigName = %s
	---- Nats Service ----
	NatsUrl = %s
	ApnNatsUrl = %s
	AlertNatsUrl = %s
	IapStatisticNatsUrl = %s
	---- Logging ---
	%s
	`, cfg.IniFile,
		cfg.ServerName,
		cfg.Name,
		cfg.NatsUrl,
		cfg.ApnNatsUrl,
		cfg.AlertNatsUrl,
		cfg.IapStatisticNatsUrl,
		cfg.LogConfig.String(),
	)
}

var (
	//配置项信息实例
	DefConfig = Config{
		Name:                "default",
		ServerName:          DefaultName,
		NatsUrl:             DefNatsUrl,
		AlertNatsUrl:        DefAlertNatsUrl,
		IapStatisticNatsUrl: DefIapStatisticNatsUrl,
		LogConfig:           xylog.DefConfig,
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
func ParseConfigFile(businessCollections *xybusiness.BusinessCollections) (err error) {

	//解析文件
	DefIniConfigs, err = beegoconf.NewConfig("ini", DefConfig.IniFile)
	if err != xyerror.ErrOK {
		fmt.Printf("load config file %s failed : %v\n", DefConfig.IniFile, err)
		return
	}

	DefConfig.Name = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_CONFIG)                              //服务器名称
	DefConfig.NatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_NATSURL)                          //主业务nats url
	DefConfig.IapStatisticNatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_IAPSTATISTIC_NATSURL) //iapstatistic nats url
	DefConfig.AlertNatsUrl = DefIniConfigs.String(xybusiness.INI_CONFIG_ITEM_SERVER_ALERT_NATSURL)               //alert nats url

	//解析业务数据库信息并且初始化数据库会话管理器
	xybusiness.ParseBusinessCollections(DefIniConfigs, businessCollections)

	//设置日志相关参数
	err = DefConfig.LogConfig.ProcessIniConfig(DefIniConfigs)
	if err != xyerror.ErrOK {
		fmt.Printf("DefConfig.LogConfig.ProcessIniConfig failed : %v\n", err)
		return
	}

	return
}

// initServerConfig 初始化服务器配置项
// businessCollections *xybusiness.BusinessCollections 保存业务数据库表信息
func initServerConfig(businessCollections *xybusiness.BusinessCollections) bool {

	//解析命令行参数
	ProcessCmd()

	//解析配置文件
	if err := ParseConfigFile(businessCollections); err != xyerror.ErrOK {
		return false
	}

	ApplyConfig()

	return true
}
