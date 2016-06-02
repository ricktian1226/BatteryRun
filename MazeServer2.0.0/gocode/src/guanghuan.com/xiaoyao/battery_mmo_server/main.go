// main
package main

import (
	//"flag"
	//"github.com/astaxie/beego/config"
	"guanghuan.com/xiaoyao/battery_mmo_server/server"
)

func main() {
	GetConfigFromFile(server.GOpts)

	ProcessCmd(server.GOpts)

	ApplyConfig()

	// Create the server with appropriate options.
	s := server.New(server.GOpts)
	// Start things up. Block here until done.
	s.Start()

}
