package xylog

import (
	"flag"
	"fmt"
	beegoconf "github.com/astaxie/beego/config"
	"os"
	"strings"
)

//日志相关配置项定义
const (
	INI_CONFIG_ITEM_LOG_APP        = "Log::app"
	INI_CONFIG_ITEM_LOG_LOGPATH    = "Log::logpath"
	INI_CONFIG_ITEM_LOG_LOGMAXLINE = "Log::logmaxline"
	INI_CONFIG_ITEM_LOG_LOGMAXSIZE = "Log::logmaxsize"
	INI_CONFIG_ITEM_LOG_LOGDAILY   = "Log::logdaily"
	INI_CONFIG_ITEM_LOG_LOGMAXDAYS = "Log::logmaxdays"
	INI_CONFIG_ITEM_LOG_LOGROTATE  = "Log::logrotate"
	INI_CONFIG_ITEM_LOG_STDOUT     = "Log::stdout"
	INI_CONFIG_ITEM_LOG_LEVEL      = "Log::loglevel"
	INI_CONFIG_ITEM_LOG_VERBOSE    = "Log::verbose"
)

type LoggerConfig struct {
	AppName      string
	NodeIdentity string
	DCId         int
	NodeId       int
	LogId        int
	Path         string
	Filename     string
	Maxlines     int
	Maxsize      int
	Daily        bool
	Maxdays      int
	Rotate       bool
	Level        LogLevel
	Stdout       bool
	Verbose      bool // enable FuncCallDepth?
}

func (l *LoggerConfig) String() string {
	return fmt.Sprintf(`
	    app      : %s
		nodeidentity:%s
	    dcid     : %d
		nodeid   : %d
		path     : %s
		filename : %s
		maxlines : %d
		maxsize  : %d
		daily    : %t
		maxdays  : %d
		rotate   : %t
		level    : %s
		stdout   : %t
		verbose  : %t
	   `,
		l.AppName,
		l.NodeIdentity,
		l.DCId,
		l.NodeId,
		l.Path,
		l.Filename,
		l.Maxlines,
		l.Maxsize,
		l.Daily,
		l.Maxdays,
		l.Rotate,
		l.Level.String(),
		l.Stdout,
		l.Verbose)
}

//解析命令行中的配置项
func (l *LoggerConfig) ProcessCmd() {
	flag.IntVar(&l.DCId, "dcid", l.DCId, "dc id")
	flag.IntVar(&l.NodeId, "nodeid", l.NodeId, "node id")
}

//解析配置文件中的配置项
// configs beegoconf.ConfigContainer ini文件配置项管理器
func (l *LoggerConfig) ProcessIniConfig(configs beegoconf.ConfigContainer) (err error) {
	l.AppName = configs.String(INI_CONFIG_ITEM_LOG_APP)
	l.Path = configs.String(INI_CONFIG_ITEM_LOG_LOGPATH)
	if l.Maxlines, err = configs.Int(INI_CONFIG_ITEM_LOG_LOGMAXLINE); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_LOGMAXLINE, err)
		return
	}

	if l.Maxsize, err = configs.Int(INI_CONFIG_ITEM_LOG_LOGMAXSIZE); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_LOGMAXSIZE, err)
		return
	}

	if l.Daily, err = configs.Bool(INI_CONFIG_ITEM_LOG_LOGDAILY); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_LOGDAILY, err)
		return
	}

	if l.Maxdays, err = configs.Int(INI_CONFIG_ITEM_LOG_LOGMAXDAYS); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_LOGMAXDAYS, err)
		return
	}

	if l.Rotate, err = configs.Bool(INI_CONFIG_ITEM_LOG_LOGROTATE); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_LOGROTATE, err)
		return
	}

	if l.Stdout, err = configs.Bool(INI_CONFIG_ITEM_LOG_STDOUT); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_STDOUT, err)
		return
	}

	if l.Verbose, err = configs.Bool(INI_CONFIG_ITEM_LOG_VERBOSE); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_VERBOSE, err)
		return
	}

	var tmp int
	if tmp, err = configs.Int(INI_CONFIG_ITEM_LOG_LEVEL); err != nil {
		fmt.Printf("Get %s failed : %v\n", INI_CONFIG_ITEM_LOG_LEVEL, err)
		return
	} else {
		l.Level = LogLevel(tmp)
	}

	l.NodeIdentity = fmt.Sprintf("%s_%d_%d", l.AppName, l.DCId, l.NodeId)
	return
}

func NewLoggerConfig() *LoggerConfig {
	app := strings.ToLower(os.Args[0])
	app = strings.TrimRight(app, ".exe")
	idx := strings.LastIndex(app, "\\")
	if idx >= 0 {
		app = app[idx+1:]
	}
	idx = strings.LastIndex(app, "/")
	if idx >= 0 {
		app = app[idx+1:]
	}

	return &LoggerConfig{
		AppName:  app,
		DCId:     0,
		NodeId:   0,
		Path:     ".",
		Maxlines: 100000,
		Maxsize:  10 * 1024 * 1024,
		Daily:    true,
		Maxdays:  30,
		Rotate:   true,
		Level:    def_log_level,
		Stdout:   true,
		Verbose:  false,
	}
}

var DefConfig = NewLoggerConfig()

func ProcessCmd() {
	DefConfig.ProcessCmd()
}

func String() string {
	return DefConfig.String()
}

func ProcessCmdAndApply() {
	ProcessCmd()
	flag.Parse()
	ApplyConfig(nil)
}
