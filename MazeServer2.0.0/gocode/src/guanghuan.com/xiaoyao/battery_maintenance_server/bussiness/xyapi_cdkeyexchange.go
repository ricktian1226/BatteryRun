package bussiness

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	batterydb "guanghuan.com/xiaoyao/battery_maintenance_server/db"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

const (
	MAINTENANCESUCESS = 1
)

const (
	MAINTENANCEURL = "http://api.open.737.com/v1/game/cdkey"
)

// cdkey 兑换返回码
const (
	CDKEYOK           = 1     // 兑换成功
	CDKEYEXCHANGEFAIL = 0     // 兑换失败
	CDKEYNOTEXITED    = 32004 // 兑换劵不存在
	CDKEYUSED         = 32003 // 兑换劵已使用
	CDKEYEXPIRED      = 32005 // 兑换劵过期
	CDKEYGETED        = 32007 // 已领取
	CDKEYERROR        = 32001 // 兑换码错误
)

var CDkeyRequestArgs = []string{
	"appid",
	"channel",
	"pf",
	"code",
	"ts",
	"userip", "zoneid", "userid",
	"username", "nickname",
}

func (api *XYAPI) CDkeyExchange(req *battery.CDkeyExchangeRequest, respMsg *battery.CDkeyExchangeResponse, debug bool) (err error) {
	var (
		errStruct = xyerror.DefaultError()
		states    int
		text      string
		urlStr    string
	)

	CDkeyinfor, err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetAccountInfor(req)
	if err != nil {
		// 返回查找用户失败
		errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
		goto ErrorHandler
	}
	CDkeyinfor.Appid = DefConfigCache.Configs().Appid
	urlStr = constructUrlstr(&CDkeyinfor)

	// 用于通信测试
	if debug {
		xylog.DebugNoId("url:%s", urlStr)
		goto ErrorHandler
	}
	states, text = sendCDkeyMsg(urlStr)
	switch states {
	case CDKEYOK:
		battery.ErrorCode_NoError.Enum()
	case CDKEYNOTEXITED, CDKEYERROR:
		errStruct.Code = battery.ErrorCode_CDkeyInvalidError.Enum()
	case CDKEYUSED:
		errStruct.Code = battery.ErrorCode_CDkeyUsedError.Enum()
	case CDKEYEXPIRED:
		errStruct.Code = battery.ErrorCode_CDkeyExpiredError.Enum()
	case CDKEYGETED:
		errStruct.Code = battery.ErrorCode_CDkeyUpperLimitError.Enum()

	default:
		errStruct.Code = battery.ErrorCode_CDkeyExchangeFail.Enum()
	}
	errStruct.Desc = &text
	xylog.InfoNoId("cdkey exchange errDescribe:%s", text)
ErrorHandler:
	respMsg.Error = errStruct
	respMsg.Uid = &CDkeyinfor.Userid
	xylog.DebugNoId("PbMsg:%v", respMsg)

	return
}

// 运营响应消息体
type MaintenanceResp struct {
	Status int
	Text   string
}

// 发送兑换请求
// status 为1表示兑换成功，text为状态描述
func sendCDkeyMsg(uri string) (status int, text string) {
	var (
		respData = MaintenanceResp{}
	)
	resp, err := http.Get(uri)
	xylog.DebugNoId("resp:%v", resp)
	if err != nil {
		xylog.ErrorNoId("Error:%v", err.Error())
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&respData)
	status = respData.Status
	text = respData.Text
	xylog.DebugNoId("status :%v,text:%s,", status, text)
	return
}

func constructUrlstr(resp *batterydb.CDkeyRequest) (urlStr string) {
	var (
		argMap url.Values
		uri    *url.URL
		sig    string
	)
	urlStr = fmt.Sprintf("%s%s", MAINTENANCEURL, "?")
	// 构造请求url
	urlStr = fmt.Sprintf("%sappid=%s&channel=%s&pf=%s&code=%s&ts=%d&userip=%s&zoneid=%d&userid=%s&username=%s&nickname=%s",
		urlStr, resp.Appid, resp.Channel, resp.Platform, resp.Code, resp.Timestamp, resp.Userip, resp.Zoneid, resp.Userid, resp.UserName, resp.NickName)
	xylog.InfoNoId("urlStr:%s", urlStr)
	uri, _ = url.Parse(urlStr)
	argMap = uri.Query() // 获取参数urlmap
	sig = constructSig(argMap)
	urlStr = fmt.Sprintf("%s&sig=%s", urlStr, sig)
	xylog.InfoNoId("uriStr :%s", urlStr)
	return
}

// 构造sig加密串
// sig 生成说明
// 1、按URL去掉sig按参数名升序排再 urlencode 后 跟上 &和配置的加密串
// 2、md5(urlencode(coin=100&ts=1426669167&user=111) . '&'. 配置加密串)
// 3、md5串采用小写形式
func constructSig(argMap url.Values) (sig string) {
	var uriStr string
	xylog.InfoNoId("argMap:%s", argMap)
	//uriStr = fmt.Sprintf("%s%s", MAINTENANCEURL, "?")
	// 对url参数进行升序排序
	sort.Strings(CDkeyRequestArgs)
	for index, arg := range CDkeyRequestArgs {
		uriStr = fmt.Sprintf("%s%s=%s", uriStr, arg, argMap.Get(arg))
		if index != len(CDkeyRequestArgs)-1 {

			uriStr = fmt.Sprintf("%s&", uriStr)
		}
	}
	xylog.InfoNoId("uriStr :%s", uriStr)
	// 参数进行url编码
	sig = url.QueryEscape(uriStr)

	// md5加密
	sig = fmt.Sprintf("%s&%s", sig, DefConfigCache.Configs().Appkey)
	xylog.DebugNoId("sig before md5:%s", sig)
	h := md5.New()
	io.WriteString(h, sig)
	sig = fmt.Sprintf("%x", h.Sum(nil))

	sig = strings.ToLower(sig)
	xylog.InfoNoId("urlEncode:%s", sig)
	return sig
}
