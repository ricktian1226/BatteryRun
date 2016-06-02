package main

import _ "net/http/pprof"
import (
	"fmt"
	xyserver "guanghuan.com/xiaoyao/common/server"
	xyutil "guanghuan.com/xiaoyao/common/util"
	"log"
	"os"
)

func banner() {
	fmt.Println("*********************************************")
	fmt.Println("*          Battery Run Main Server          *")
	fmt.Println("*                                           *")
	fmt.Println("* input -usage for help                     *")
	fmt.Println("*********************************************")
}

func main() {
	banner()

	args := os.Args[1:]
	pid := os.Getpid()
	//	default_config.PrintBanner()
	ok := default_config.ProcessCommandLine(args)
	if !ok {
		default_config.PrintUsage()
		return
	}
	default_config.Print()
	//batteryapi.DefaultConfig.Load(default_config.Cfgname)

	if !default_config.UseStdOutput {
		logfile := fmt.Sprintf("%s/battery_server_%s_%d.%s.log",
			default_config.Log_dir, default_config.Name, default_config.Id, xyutil.CurTimeStr())

		f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		//		defer f.Close()
		if err != nil {
			log.Printf("error opening file: %v", err)
			log.Printf("Server stopped")
			return
		}

		log.Printf("all log writes to %s", logfile)
		log.SetOutput(f)
	}

	log.Printf("Starting Server %s-%d on (%s:%d) ... (pid=%d)",
		default_config.Name, default_config.Id, default_config.Host, default_config.Port, pid)

	opts := xyserver.Options{}
	server := xyserver.New(default_config.Name, &opts)

	// create db service
	//dbservice = xyservice.NewDBService(DefaultDBServiceName, *db_url, *db_name)
	//server.RegisterService(dbservice.Name(), dbservice)
	//dbservice = server.EnableDBService(default_config.DBUrl, default_config.DBName)

	// create nats service
	// nats_service = xyservice.NewNatsService(DefaultNatsServiceName, *server_url)
	// server.RegisterService(nats_service.Name(), nats_service)
	//nats_service = server.EnableNatsService(default_config.NatsUrl)
	//	nats_service.AddQueueSubscriber(*subject, HandleAuditMessage)

	// create a martini service
	//	martini_service = server.EnableMartiniService(default_config.Host, default_config.Port)
	//	db := dbservice.GetDB()
	//	btdb := batterydb.NewBatteryDB(db.(*xydb.XYDB))
	battery_service = NewBatteryHttpService("battery service", default_config.Host, default_config.Port,
		default_config.DBUrl, default_config.DBName, default_config.NatsUrl,
		default_config.Cfgname)

	// 添加静态文件目录
	battery_service.SetStaticFilePath("httpdoc")
	//	m2 := xyservice.DefaultMartiniService("martini2", "", 8888)
	server.RegisterService(battery_service.Name(), battery_service)
	//	battery_service.AddRouter(xyservice.HttpPost, API_URI_LOGIN, )

	server.Start()
}
