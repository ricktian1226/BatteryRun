package main

import (
    "flag"
    "fmt"
    beegoconf "github.com/astaxie/beego/config"
    //batterydb "guanghuan.com/xiaoyao/battery_transaction_server/db"
    xylog "guanghuan.com/xiaoyao/common/log"
    //xyserver "guanghuan.com/xiaoyao/common/server"
    //battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
    "os"
)

type Config struct {
    IniFile                                                string
    Name                                                   string
    ServerName                                             string
    TestEnv                                                bool
    NatsUrl, AlertNatsUrl, ApnNatsUrl, IapStatisticNatsUrl string
    Concurrent                                             bool
    LogConfig                                              *xylog.LoggerConfig
}

func (cfg *Config) String() (str string) {
    str = fmt.Sprintf(`
	IniFile    = %s
	ServerName = %s
	ConfigName = %s
	TestEnv    = %t
	---- Nats Service ----
	NatsUrl = %s
	AlertNatsUrl = %s
	ApnNatsUrl = %s
	IapStatisticNatsUrl = %s
	`,
        cfg.IniFile,
        cfg.ServerName,
        cfg.Name,
        cfg.TestEnv,
        cfg.NatsUrl,
        cfg.AlertNatsUrl,
        cfg.ApnNatsUrl,
        cfg.IapStatisticNatsUrl)

    str += fmt.Sprintf(`---- Logging ---
	%s
	`,
        cfg.LogConfig.String(),
    )

    return
}

var (
    DefConfig = Config{
        //IniFile:    os.Args[0] + ".ini",
        Name:         "default",
        ServerName:   DefaultName,
        TestEnv:      true,
        NatsUrl:      DefNatsUrl,
        AlertNatsUrl: DefAlertNatsUrl,
        ApnNatsUrl:   DefApnNatsUrl,
        Concurrent:   true,
        LogConfig:    xylog.DefConfig,
    }

    defIniConfig beegoconf.ConfigContainer
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
func ParseConfigFile(businessCollections *xybusiness.BusinessCollections) (err error) {
    defIniConfig, err = beegoconf.NewConfig("ini", DefConfig.IniFile)
    if err != xyerror.ErrOK {
        fmt.Printf("load config file %s failed : %v\n", DefConfig.IniFile, err)
        return
    }

    DefConfig.Name = defIniConfig.String(xybusiness.INI_CONFIG_ITEM_SERVER_CONFIG)
    DefConfig.NatsUrl = defIniConfig.String(xybusiness.INI_CONFIG_ITEM_SERVER_NATSURL)
    DefConfig.AlertNatsUrl = defIniConfig.String(xybusiness.INI_CONFIG_ITEM_SERVER_ALERT_NATSURL)
    DefConfig.ApnNatsUrl = defIniConfig.String(xybusiness.INI_CONFIG_ITEM_SERVER_APNNATSURL)
    DefConfig.IapStatisticNatsUrl = defIniConfig.String(xybusiness.INI_CONFIG_ITEM_SERVER_IAPSTATISTIC_NATSURL)

    //解析业务数据库信息并且初始化数据库会话管理器
    xybusiness.ParseBusinessCollections(defIniConfig, businessCollections)

    err = DefConfig.LogConfig.ProcessIniConfig(defIniConfig)
    if err != xyerror.ErrOK {
        fmt.Printf("DefConfig.LogConfig.ProcessIniConfig failed : %v\n", err)
        return

    }

    fmt.Printf("DefConfig : %v", DefConfig)

    return
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
