// main
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	httplib "github.com/astaxie/beego/httplib"
	xyencoder "guanghuan.com/xiaoyao/common/encoding"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
	"time"
)

const API_URI_LOTTO_RES_OP = "/v2/lotto/lotto_res_op" //抽奖资源请求

type Config struct {
	Url string
}

var Default_Config = Config{
	Url: "192.168.1.205:10003",
}

func test_op(req proto.Message, resp proto.Message, s string, rwtimeout time.Duration) (err error) {
	var (
		data    []byte
		api_url string = "http://" + Default_Config.Url + s
	)
	data, err = xyencoder.PbEncode(req)
	if err != nil {
		xylog.Error("Error Encoding: %s", err.Error())
		return err
	}
	//xylog.Debug("reqdata after PbEncode: [%d][%v]\n", len(data), data)

	data, err = crypto.Encrypt(data)
	if err != nil {
		xylog.Error("Error Encrypt: %s", err.Error())
		return err
	}
	//xylog.Debug("reqdata after Encrypt: [%d][%v]\n", len(data), data)

	fmt.Println(api_url)
	request := httplib.Post(api_url)
	request.Body(data)
	request.SetTimeout(60*time.Second, rwtimeout*time.Microsecond)

	data, err = request.Bytes()
	if err != nil {
		xylog.Error("Error Send: %s", err.Error())
		return err
	}
	//xylog.Debug("respdata after Send: [%d][%v]\n", len(data), data)

	data, err = crypto.Decrypt(data)
	if err != nil {
		xylog.Error("Error Decrypt: %s", err.Error())
		return err
	}
	//xylog.Debug("respdata after Decrypt: [%d][%v]\n", len(data), data)

	err = xyencoder.PbDecode(data, resp)
	if err != nil {
		xylog.Error("Error Decoding: %s", err.Error())
		return err
	}

	if err == nil {
		xylog.Debug("response : %v", resp)
	}

	return err
}

func main() {

	flag.StringVar(&Default_Config.Url, "url", "192.168.1.205:10003", "gateway server url")
	xylog.ProcessCmdAndApply()

	var err error

	req := &battery.LottoResRequest{}
	ops := make([]*battery.LottoResOpItem, 0)

	GetProperty(&ops)

	req.Ops = ops

	xylog.Debug("req.Ops : %v", req.GetOps())

	resp := &battery.LottoResResponse{}

	err = test_op(req, resp, API_URI_LOTTO_RES_OP, 1000*1000)

	if err != nil {
		xylog.Error("======== lotto import failed =========")
	} else {
		xylog.Debug("======== lotto import succeed =========")
	}

	time.Sleep(5 * time.Second)

}
