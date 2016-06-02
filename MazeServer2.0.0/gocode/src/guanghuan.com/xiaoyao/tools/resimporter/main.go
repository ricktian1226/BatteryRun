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

const API_URI_RES_OP = "/v2/res/res_op" //抽奖资源请求
const OP_NUM_PER_REQ = 150              //每次请求的

// var Default_Config = Config{
// 	Url: "192.168.1.205:10003",
// }

func test_op(req proto.Message, resp proto.Message, s string, rwtimeout time.Duration) (err error) {
	var (
		data    []byte
		api_url string = "http://" + DefConfig.Url + s
	)
	data, err = xyencoder.PbEncode(req)
	if err != nil {
		xylog.ErrorNoId("Error Encoding: %s", err.Error())
		return
	}
	//xylog.Debug("reqdata after PbEncode: [%d][%v]\n", len(data), data)

	data, err = crypto.Encrypt(data)
	if err != nil {
		xylog.ErrorNoId("Error Encrypt: %s", err.Error())
		return
	}
	//xylog.Debug("reqdata after Encrypt: [%d][%v]\n", len(data), data)

	fmt.Println(api_url)
	request := httplib.Post(api_url)
	request.Body(data)
	request.SetTimeout(60*time.Second, rwtimeout*time.Microsecond)

	data, err = request.Bytes()
	if err != nil {
		xylog.ErrorNoId("Error Send: %s", err.Error())
		return err
	} else if len(data) <= 0 {
		xylog.ErrorNoId("Error resp len(data)(%d) <= 0.", len(data))
		return err
	}
	//xylog.Debug("respdata after Send: [%d][%v]\n", len(data), data)

	data, err = crypto.Decrypt(data)
	if err != nil {
		xylog.ErrorNoId("Error Decrypt: %s", err.Error())
		return err
	}
	//xylog.Debug("respdata after Decrypt: [%d][%v]\n", len(data), data)

	err = xyencoder.PbDecode(data, resp)
	if err != nil {
		xylog.ErrorNoId("Error Decoding: %s", err.Error())
		return err
	}

	if err == nil {
		xylog.DebugNoId("response : %v", resp)
	}

	return err
}

func main() {

	//flag.StringVar(&Default_Config.Url, "url", "192.168.1.205:10003", "gateway server url")
	flag.BoolVar(&DefConfig.Check, "check", false, "check the content of files")
	//xylog.ProcessCmdAndApply()
	initLogConfig()
	var err error

	req := &battery.ResRequest{}
	ops := make([]*battery.ResOpItem, 0)

	flag := GetProperty(&ops)

	if !(DefConfig.Check) {
		//如果资源配置信息太多，分拆成多个消息
		var rest, loopTime, sum int
		sum = len(ops)
		rest = (sum % OP_NUM_PER_REQ)
		if rest > 0 {
			loopTime = (sum / OP_NUM_PER_REQ) + 1
		} else {
			loopTime = (sum / OP_NUM_PER_REQ)
		}

		xylog.DebugNoId("len(ops) : %d, loopTime : %d", len(ops), loopTime)

		for i := 0; i < loopTime && sum > 0; i++ {
			//第一条消息需要删除老数据
			if i == 0 {
				req.RemoveOld = proto.Bool(true)
				req.Flag = proto.Int32(flag)
			} else {
				req.RemoveOld = proto.Bool(false)
			}

			begin := i * OP_NUM_PER_REQ
			if sum > OP_NUM_PER_REQ {
				req.Ops = ops[begin : begin+OP_NUM_PER_REQ]
			} else {
				req.Ops = ops[begin : begin+sum]
			}

			xylog.DebugNoId("len(req.Ops) : %v", len(req.Ops))

			resp := &battery.ResResponse{}

			err = test_op(req, resp, API_URI_RES_OP, 1000*1000*1000)

			if err != nil {
				xylog.ErrorNoId("======== res %d import failed =========", i)
				break
			} else {
				xylog.DebugNoId("======== res %d import succeed =========", i)
			}
			sum -= OP_NUM_PER_REQ
		}
	}

	time.Sleep(5 * time.Second)

}
