package main

import (
	"fmt"
	xyversion "guanghuan.com/xiaoyao/common/version"
	"log"
	"strconv"
)

type config struct {
	Host          string
	Port          int
	Log_dir       string
	Id            int
	Name          string
	Cfgname       string
	UseStdOutput  bool
	Security      int
	DBUrl         string
	DBName        string
	NatsUrl       string
	ServerVersion xyversion.Version
}

var (
	default_config config
)

func init() {
	default_config.Host = ""
	default_config.Port = 12356
	default_config.Log_dir = "."
	default_config.Id = 1
	default_config.Name = "default"
	default_config.UseStdOutput = false
	default_config.Security = -1
	default_config.Cfgname = "default"
	default_config.DBUrl = DefaultDBURL
	default_config.DBName = DefaultDB
	default_config.NatsUrl = DefaultNatsUrl
}

func processConfigFile() {

}

func (cfg *config) ProcessCommandLine(args []string) (ok bool) {
	if len(args) <= 0 {
		return true
	}
	//
	if len(args)%2 != 0 {
		fmt.Printf("\nError: invalid args numbers: %d", len(args))
		return false
	}
	var err error
	ok = true

	for i := 0; i < int(len(args)/2); i++ {
		key := args[i*2]
		val := args[i*2+1]

		switch key {
		case "-port":
			cfg.Port, err = strconv.Atoi(val)
			if err != nil {
				log.Printf("error: %s", err.Error())
				ok = false
				break
			}
		case "-log":
			cfg.Log_dir = val
		case "-config":
			cfg.Cfgname = val
		case "-host":
			cfg.Host = val
		case "-name":
			cfg.Name = val
		case "-id":
			cfg.Id, err = strconv.Atoi(val)
			if err != nil {
				log.Printf("error: %s", err.Error())
				ok = false
				break
			}
		case "-stdout":
			v := 0
			v, err = strconv.Atoi(val)
			if err != nil {
				log.Printf("error: %s", err.Error())
				ok = false
				break
			}
			cfg.UseStdOutput = (v > 0)
		case "-security":
			v := 0
			v, err = strconv.Atoi(val)
			if err == nil {
				cfg.Security = v
				log.Printf("security = %d", cfg.Security)
			}
		case "-ns":
			cfg.NatsUrl = val
		case "-dburl":
			cfg.DBUrl = val
		case "-db":
			cfg.DBName = val
		}
	}
	return ok
}

func (cfg *config) PrintUsage() {
	fmt.Printf("\n========= configuration =========")
	fmt.Printf("\n-name <server name> default: default")
	fmt.Printf("\n-id   <server id>   default: 1")
	fmt.Printf("\n-host <server ip>   default: null")
	fmt.Printf("\n-port <server port> default: 12356")
	fmt.Printf("\n-log  <log path>    default: .")
	fmt.Printf("\n-stdout  <write to stdout only>    default: 0=false")
	fmt.Println()
	fmt.Printf("\n example:  -name test -id 1 -host localhost -port 8765 -log ./log ")
	fmt.Printf("\n=================================\n")
}

func (cfg *config) Print() {
	log.Printf("========= configuration =========")
	log.Printf("Server Version = %s", ServerVersion.String())
	log.Printf("Server ID   = %d", cfg.Id)
	log.Printf("Server Name = %s", cfg.Name)
	log.Printf("Host        = %s", cfg.Host)
	log.Printf("Port        = %d", cfg.Port)
	log.Printf("Log path    = %s", cfg.Log_dir)
	log.Printf("Log to stdout only = %v", cfg.UseStdOutput)
	log.Printf("Security Enabled = %d", cfg.Security)
	log.Printf("dburl = %s", cfg.DBUrl)
	log.Printf("db = %s", cfg.DBName)
	log.Printf("nats = %s", cfg.NatsUrl)
	log.Printf("=================================")
}

func (cfg *config) PrintBanner() {
	fmt.Println("*********************************************")
	fmt.Println("*          Battery Run Main Server          *")
	fmt.Println("*                                           *")
	fmt.Println("* input -usage for help                     *")
	fmt.Println("*********************************************")
}
