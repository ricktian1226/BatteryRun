// A gateway server for HTTP <-> nats message routing
package main

import (
	"fmt"
	//	xylog "guanghuan.com/xiaoyao/common/log"
	xyserver "guanghuan.com/xiaoyao/common/server"
	"runtime"
)

func banner() {
	fmt.Println("*********************************************")
	fmt.Println("*      Battery Run File Server              *")
	fmt.Println("*                                           *")
	fmt.Println("* input -usage for help                     *")
	fmt.Println("*********************************************")
}

func main() {
	banner()
	initServerConfig()

	// starting server
	opts := xyserver.Options{}
	server := xyserver.New(DefConfig.ServerName, &opts)

	// create a martini service
	martini_service := server.EnableMartiniService(DefConfig.HttpHost, DefConfig.HttpPort)
	if DefConfig.MaxRequest > 0 {
		martini_service.EnableFlowControl(DefConfig.MaxRequest, DefConfig.MaxRequestTimeout, DefConfig.MaxTimeoutRequest)
	}
	// 添加静态文件目录
	martini_service.SetStaticFilePath(DefConfig.HttpDocRoot)

	server.Start()
	go server.Run()
	runtime.Goexit()
}
