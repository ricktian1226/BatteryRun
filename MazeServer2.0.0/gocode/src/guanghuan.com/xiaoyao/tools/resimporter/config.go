// config
package main

import (
	"flag"
	"fmt"
	beegoconf "github.com/astaxie/beego/config"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/superbman_server/error"
)

type Config struct {
	Url       string
	Check     bool
	IniFile   string
	LogConfig *xylog.LoggerConfig
}

var (
	DefConfig = Config{
		LogConfig: xylog.DefConfig,
	}
)

func (conf *Config) String() (str string) {
	str = fmt.Sprintf(`
        IniFile= %s
        Url=%s 
        Chech=%v
        `,
		conf.IniFile,
		conf.Url,
		conf.Check)
	str += fmt.Sprintf(`----Logging----
     %s`, conf.LogConfig.String())
	return
}
func ProcessCmd() bool {
	flag.StringVar(&DefConfig.IniFile, "config", "resimporter.ini", "ini file")
	DefConfig.LogConfig.ProcessCmd()
	flag.Parse()
	return true
}

func PareseConfigFile() (err error) {
	defInfConfig, err := beegoconf.NewConfig("ini", DefConfig.IniFile)
	DefConfig.Url = defInfConfig.String("Server::url")
	if err != xyerror.ErrOK {
		fmt.Printf("load config file %s fail,%v\n ", DefConfig.IniFile, err)
	}
	err = DefConfig.LogConfig.ProcessIniConfig(defInfConfig)
	return
}

func ApplyConfig() {
	xylog.ApplyConfig(DefConfig.LogConfig)
	xylog.Info("Config: %v", DefConfig.String())
}

func initLogConfig() bool {
	if !ProcessCmd() {
		return false
	}
	if err := PareseConfigFile(); err != nil {
		return false
	}
	ApplyConfig()
	return true
}
