// config
package main

import (
	"flag"
	//"fmt"
	"guanghuan.com/xiaoyao/battery_mmo_server/conf"
	"guanghuan.com/xiaoyao/battery_mmo_server/server"
	xylog "guanghuan.com/xiaoyao/common/log"
)

const (
	DEFAULT_CONFIG_FILE = "./mmo.conf"
)

func GetConfigFromFile(opts *server.Options) {

	if ok := xyconf.NewJSonConfig(DEFAULT_CONFIG_FILE); ok {
		opts.Port, _ = xyconf.GJsonConf.Int("port")
		opts.MaxConn, _ = xyconf.GJsonConf.Int("maxconnection")
		opts.PingInterval, _ = xyconf.GJsonConf.Int("pinginterval")
		opts.PingTimeOut, _ = xyconf.GJsonConf.Int("pingtimeout")
		opts.RoundTime, _ = xyconf.GJsonConf.Int("roundtime")
		opts.RoundDelay, _ = xyconf.GJsonConf.Int("rounddelay")
		opts.PidFile = xyconf.GJsonConf.String("pid")
		opts.NoSigs, _ = xyconf.GJsonConf.Bool("nosig")
		//fmt.Print(opts.String())
	}
}

func ProcessCmd(opts *server.Options) {
	flag.IntVar(&opts.Port, "port", opts.Port, "service port")
	flag.IntVar(&opts.MaxConn, "mcs", opts.MaxConn, "service maxconnections")
	flag.IntVar(&opts.PingInterval, "pi", opts.PingInterval, "service ping interval")
	flag.IntVar(&opts.PingTimeOut, "pto", opts.PingTimeOut, "service ping timeout")
	flag.IntVar(&opts.RoundTime, "rt", opts.RoundTime, "round time limit")
	flag.IntVar(&opts.RoundDelay, "rd", opts.RoundDelay, "round delay time ")
	flag.StringVar(&opts.PidFile, "pid", opts.PidFile, "pid file")
	flag.BoolVar(&opts.NoSigs, "nosig", opts.NoSigs, "no signal")
	xylog.DefConfig.ProcessCmd()
	flag.Parse()
}

func ApplyConfig() {
	xylog.ApplyConfig(xylog.DefConfig)
	xylog.Info("Config: %s", xylog.DefConfig.String())
}
