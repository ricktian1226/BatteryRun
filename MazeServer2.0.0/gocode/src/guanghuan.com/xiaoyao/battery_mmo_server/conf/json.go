package xyconf

import (
	beegoConf "github.com/astaxie/beego/config"
	//xylogs "guanghuan.com/xiaoyao/common/log"
	"fmt"
)

var GJsonConf *beegoConf.JsonConfigContainer

func NewJSonConfig(file string) bool {
	conf, err := beegoConf.NewConfig("json", file)
	if err != nil {
		fmt.Println("NewJSonConfig failed, pls check : " + err.Error())
		return false
	}
	fmt.Println("NewJSonConfig OK.")
	fmt.Printf("NewJSonConfig %v.", conf)

	GJsonConf = conf.(*beegoConf.JsonConfigContainer)

	return true
}
