// announcement
package main

import (
	httplib "github.com/astaxie/beego/httplib"
	//xycrypto "guanghuan.com/xiaoyao/common/crypto"
	//"github.com/lxn/walk"
	//. "github.com/lxn/walk/declarative"
	"bufio"
	proto "code.google.com/p/goprotobuf/proto"
	"flag"
	xyencoder "guanghuan.com/xiaoyao/common/encoding"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
	"io"
	"os"
	"strings"
	"time"
)

type Server_Config struct {
	Url         string //gateway 服务url
	ConfigFile  string //配置文件
	ItemsPerReq int    //每次请求的通告数(压力测试时可配置大于1）
	ReqCount    int    //请求次数(压力测试时可配置大于1）
}

var Default_Server_Config = NewServerConfig()

func NewServerConfig() *Server_Config {
	return &Server_Config{
		Url:        "192.168.93.129:10003",
		ConfigFile: "announcement.conf",
	}
}

const (
	TimeFormat           = "2006-01-02 15:04:05"
	API_URI_ANNOUNCEMENT = "/v2/announcement/:token" // 通告
	TIMEZONE_SECONDS     = 8 * 60 * 60               //中国时区，要将描述减去8小时
)

func NewAnnouncement() *battery.Announcement {
	return &battery.Announcement{
		Id:         proto.Uint64(0),
		SubmitTime: proto.Int64(time.Now().Unix()),
	}
}

func Process() {

	cmd := battery.OP_CMD_ADD
	req := &battery.AnnouncementRequest{
		Uid: proto.String("ricktian"),
		Cmd: &cmd,
	}

	item := NewAnnouncement()
	for i := 0; i < Default_Server_Config.ItemsPerReq; i++ {
		file, err := os.OpenFile(Default_Server_Config.ConfigFile, os.O_RDONLY, 0644)
		if nil != err {
			panic(err)
		}
		buf := bufio.NewReader(file)

		for {
			line, _, err := buf.ReadLine()
			//xylog.Debug("line : [%s]", line)
			sline := string(line)
			if err != nil || io.EOF == err {
				break
			} else if len(line) == 0 || "" == sline || sline[0] == '#' {
				continue
			}

			s := strings.Split(sline, "::")
			//xylog.Debug("%s split to %v", sline, s)
			if len(s) != 2 {
				continue
			}

			//xylog.Debug("config[%s] %s", s[0], s[1])

			switch {
			case s[0] == "begin_time":
				t, _ := time.Parse(TimeFormat, s[1])
				timeBegin := t.Unix()
				timeBegin -= TIMEZONE_SECONDS
				timeBegin += int64(i)
				item.BeginTime = &timeBegin
				//xylog.Debug("begin_time : [%s][%d]", s[1], *(item.BeginTime))
			case s[0] == "end_time":
				t, _ := time.Parse(TimeFormat, s[1])
				timeEnd := t.Unix()
				timeEnd -= TIMEZONE_SECONDS
				timeEnd += int64(i)
				item.EndTime = &timeEnd
				//xylog.Debug("end_time : [%s][%d]", s[1], *(item.EndTime))
			case s[0] == "title":
				item.Title = &s[1]
				//xylog.Debug("title : [%s]", *item.Title)
			case s[0] == "content":
				item.Content = &s[1]
				req.Items = append(req.Items, item)
				//xylog.Debug("item : %v %v %v %v %v %v", *item.Id, *item.SubmitTime, *item.BeginTime, *item.EndTime, *item.Title, *item.Content)
				//xylog.Debug("Shitems : %v", req.Items)
				//xylog.Debug("content : [%s]", *item.Content)
				item = NewAnnouncement() //重新分配一个对象
			}
		}
		file.Close()
	}

	for j := 0; j < Default_Server_Config.ReqCount; j++ {
		if len(req.Items) > 0 {
			xylog.Debug("Items : %v", req.Items)

			resp := &battery.AnnouncementResponse{}

			var (
				data    []byte
				api_url string = "http://" + Default_Server_Config.Url + API_URI_ANNOUNCEMENT
			)
			data, err := xyencoder.PbEncode(req)
			if err != nil {
				xylog.Error("Error Encoding: %s", err.Error())
				return
			}
			//xylog.Debug("reqdata after PbEncode: [%d][%v]\n", len(data), data)

			data, err = crypto.Encrypt(data)
			if err != nil {
				xylog.Error("Error Encrypt: %s", err.Error())
				return
			}
			//xylog.Debug("reqdata after Encrypt: [%d][%v]\n", len(data), data)

			request := httplib.Post(api_url)
			request.Body(data)

			data, err = request.Bytes()
			xylog.Debug("After request.Bytes()")
			if err != nil {
				xylog.Error("Error Send: %s", err.Error())
				return
			}
			//xylog.Debug("respdata after Send: [%d][%v]\n", len(data), data)

			data, err = crypto.Decrypt(data)
			if err != nil {
				xylog.Error("Error Decrypt: %s", err.Error())
				return
			}
			//xylog.Debug("respdata after Decrypt: [%d][%v]\n", len(data), data)

			err = xyencoder.PbDecode(data, resp)
			if err != nil {
				xylog.Error("Error Decoding: %s", err.Error())
				return
			} else {
				xylog.Debug("response : %v", resp)
			}
		}

		xylog.Debug("#%d finish", j)
	}
}

func main() {
	//获取配置参数
	flag.StringVar(&Default_Server_Config.Url, "url", Default_Server_Config.Url, "gateway server url")
	flag.IntVar(&Default_Server_Config.ItemsPerReq, "ipr", 1, "how many items per request")
	flag.IntVar(&Default_Server_Config.ReqCount, "rc", 1, "how many requests")
	flag.StringVar(&Default_Server_Config.ConfigFile, "c", Default_Server_Config.ConfigFile, "server config path")
	xylog.ProcessCmdAndApply()

	//解析配置文件，获取输入数据
	Process()

	time.Sleep(2 * time.Second)
}
