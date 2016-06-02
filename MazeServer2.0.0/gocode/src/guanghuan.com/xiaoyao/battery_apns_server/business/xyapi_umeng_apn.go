package batteryapi

import (
	//"sync"
	"time"
	//"encoding/binary"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	//proto "code.google.com/p/goprotobuf/proto"
	//apns "github.com/timehop/apns"
	httplib "github.com/astaxie/beego/httplib"

	//"guanghuan.com/xiaoyao/common/apn"
	//xyconf "guanghuan.com/xiaoyao/common/conf"
	//"guanghuan.com/xiaoyao/common/db"
	"guanghuan.com/xiaoyao/common/log"
	//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	//"guanghuan.com/xiaoyao/superbman_server/server"

	//batterydb "guanghuan.com/xiaoyao/battery_apns_server/db"
)

const (
	UMENG_SEND_URL = "http://msg.umeng.com/api/send" //发送推送请求url
)

//友盟返回消息结构体定义
type UmengSendResponseData struct {
	MsgId        string `json:"msg_id"`        //当type为unicast、listcast或者customizedcast且alias不为空时
	TaskId       string `json:"task_id"`       //当type为于broadcast、groupcast、filecast、customizedcast
	ErrorCode    string `json:"error_code"`    // 当"ret"为"FAIL"时
	ThirdPartyId string `json:"thirdparty_id"` //如果开发者填写了thirdparty_id, 接口也会返回该值
}

type UmengSendResponse struct {
	Ret  string                `json:"ret"` //返回结果，SUCCESS/FAIL
	Data UmengSendResponseData `json:"data"`
}

//友盟请求消息结构体定义

type UMengPayloadBody struct {
	Ticker      string `json:"ticker"`       //必填，通知栏提示文字
	Title       string `json:"title"`        //必填，通知标题
	Text        string `json:"text"`         //必填，通知文字描述
	Icon        string `json:"icon"`         //可选，状态栏图标ID，默认使用应用图标
	LargeIcon   string `json:"largeIcon"`    //可选，通知栏拉开后左侧图标ID
	Img         string `json:"img"`          //可选，通知栏大图标的url链接，以http或https开头
	Sound       string `json:"sound"`        //可选
	BuilderId   string `json:"builder_id"`   //可选，通知采用的样式
	PlayVibrate bool   `json:"play_vibrate"` //可选，收到通知是否震动
	PlayLights  bool   `json:"play_lights"`  //可选，收到通知是否闪灯
	PlaySound   bool   `json:"play_sound"`   //可选，收到通知是否播放声音
	AfterOpen   string `json:"after_open"`   //可选，点击通知的后续行为，默认为打开app
	Url         string `json:"url"`          //after_open为go_url时此项必填
	Activity    string `json:"activity"`     //after_open为go_activity时此项必填
	Custom      string `json:"custom"`       //after_open为go_custom时此项必填
}

type UMengPolicy struct {
	StartTime  string `json:"start_time"`   //可选，格式"YYYY-MM-DD hh:mm:ss"，默认为立即发送
	ExpireTime string `json:"expire_time"`  //可选，格式"YYYY-MM-DD hh:mm:ss"
	MaxSendNum string `json:"max_send_num"` //可选，发送限速。发送的消息如果有请求自己服务器的资源，可以考虑设置此项
}

//android设备的推送消息体
type UMengAndroidPayload struct {
	DisplayType string            `json:"display_type"` //必填，消息类型
	Body        UMengPayloadBody  `json:"body"`         //消息体
	extra       map[string]string `json:"extra"`        //用户自定义的key-value对
}

//ios设备的推送消息体
type UMengIOSPayload struct {
	Aps UMengIOSAps `json:"aps"` //必填
	//extra map[string]string `json:"extra"` //用户自定义的key-value对
}

type UMengIOSAps struct {
	Alert            string `json:"alert"` //必填
	Badge            int32  `json:"badge"`
	Sound            string `json:"sound"`
	ContentAvailable string `json:"content-available"`
	Category         string `json:"category"` //ios8才支持
}

type UMengAndroidNotification struct {
	AppKey         string              `json:"appkey"`          //应用唯一标识
	Timestamp      int64               `json:"timestamp"`       //时间戳
	Type           string              `json:"type"`            //消息发送类型。unicast 单播；listcast 列播（不超过500个devicetoken）；filecast 文件播；groupcast 组播，按照filter条件筛选特定用户群；customizedcast 自定义的alias进行推送，alias 对单个或多个alias进行推送，file_id 将alias保存在文件中，根据file_id来推送
	DeviceTokens   string              `json:"device_tokens"`   //unicast 必填，指单个设备；listcast 时，必填，不超过500个，多个devicetoken间用逗号分隔
	AliasType      string              `json:"alias_type"`      //当type为customizedcast时为必填，alias的类型。 alias_type可由开发者自定义，开发者在SDK中调用setAlias(alias, alias_type)时所设置的alias_type
	Alias          string              `json:"alias"`           //可选，当type为customizedcast时，开发真填写自己的alias。要求不超过50个alias，多个alias以逗号分隔开
	FileId         string              `json:"file_id"`         //可选，当type为filecast时，file内容为多条devicetoken，devicetoken间以回车符分隔；当type为customizedcast是，file内容为多条alias，alias间以回车符分隔，注意同一个文件内的alias所对应的alias_type必须和接口参数的alias_type一致。注意，使用文件播前需要先调用文件上传接口获取file_id
	Filter         string              `json:"filter"`          //终端用户筛选条件
	Payload        UMengAndroidPayload `json:"payload"`         //payload
	Policy         UMengPolicy         `json:"policy"`          //发送策略
	ProductionMode bool                `json:"production_mode"` //是否是生产模式
	Description    string              `json:"description"`     //消息描述，建议设置
	ThirdPartyId   string              `json:"thirdparty_id"`
}

type UMengIOSNotification struct {
	AppKey         string          `json:"appkey"`          //应用唯一标识
	Timestamp      int64           `json:"timestamp"`       //时间戳
	Type           string          `json:"type"`            //消息发送类型。unicast 单播；listcast 列播（不超过500个devicetoken）；filecast 文件播；groupcast 组播，按照filter条件筛选特定用户群；customizedcast 自定义的alias进行推送，alias 对单个或多个alias进行推送，file_id 将alias保存在文件中，根据file_id来推送
	DeviceTokens   string          `json:"device_tokens"`   //unicast 必填，指单个设备；listcast 时，必填，不超过500个，多个devicetoken间用逗号分隔
	AliasType      string          `json:"alias_type"`      //当type为customizedcast时为必填，alias的类型。 alias_type可由开发者自定义，开发者在SDK中调用setAlias(alias, alias_type)时所设置的alias_type
	Alias          string          `json:"alias"`           //可选，当type为customizedcast时，开发真填写自己的alias。要求不超过50个alias，多个alias以逗号分隔开
	FileId         string          `json:"file_id"`         //可选，当type为filecast时，file内容为多条devicetoken，devicetoken间以回车符分隔；当type为customizedcast是，file内容为多条alias，alias间以回车符分隔，注意同一个文件内的alias所对应的alias_type必须和接口参数的alias_type一致。注意，使用文件播前需要先调用文件上传接口获取file_id
	Filter         string          `json:"filter"`          //终端用户筛选条件
	Payload        UMengIOSPayload `json:"payload"`         //payload
	Policy         UMengPolicy     `json:"policy"`          //发送策略
	ProductionMode bool            `json:"production_mode"` //是否是生产模式
	Description    string          `json:"description"`     //消息描述，建议设置
	ThirdPartyId   string          `json:"thirdparty_id"`
}

// NotifyUMeng 友盟消息推送
func (api *XYAPI) NotifyUMeng(notification interface{}) (err error) {
	//获取消息json字符串
	var data []byte
	data, err = json.Marshal(notification)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("json.Marshal(%v) failed : %v", notification, err)
		return
	}

	//获取md5校验串
	sign := api.Md5UMeng(string(data))

	xylog.DebugNoId("sign : %v", sign)

	//往友盟服务器发送推送请求
	req := httplib.Post(fmt.Sprintf("%s?sign=%s", UMENG_SEND_URL, sign))
	req.Body(data)
	req.SetTimeout(60*time.Second, 60*time.Second) //链接超时时间和读写超时时间都设置为1分钟
	data, err = req.Bytes()
	xylog.DebugNoId("error : %v", err)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("NotifyUMeng failed: %v", err)
		return
	}

	resp := UmengSendResponse{}
	err = json.Unmarshal(data, &resp)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("UmengSendResponse json.Unmarshal failed : %v", err)
		return
	}

	if resp.Ret != "SUCCESS" {
		xylog.ErrorNoId("UmengSendResponse : %v", resp)
	} else {
		xylog.DebugNoId("UmengSendResponse : %v", resp)
	}

	return
}

// Md5UMeng 按照友盟的规则生成md5校验串
//规则: Sign=MD5($http_method$url$post-body$app_master_secret);
//    $http_method: POST 全大写
//    $url: 请求url
//    $post-body: post请求的body
//    $app_master_secret: 1ygzzohcqzh8qzeqrjmhmojycmp3burk
func (api *XYAPI) Md5UMeng(body string) string {
	md5Ctx := md5.New()
	md5source := fmt.Sprintf("POST%s%s%s", UMENG_SEND_URL, body, DefConfigCache.Configs().AppMasterSecret)
	xylog.DebugNoId("md5source[%s]", md5source)
	md5Ctx.Write([]byte(md5source))
	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr)
}
