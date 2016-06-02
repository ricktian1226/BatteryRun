package batteryapi

import (
	"time"

	"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
)

type NOTIFY_FUNC func(*XYAPI, string, string)

var NOTIFY_FUNC_MAP = map[battery.PLATFORM_TYPE]NOTIFY_FUNC{
	battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS:     (*XYAPI).NotifyIOS,
	battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID: (*XYAPI).NotifyAndroid,
}

// Notify 处理业务发起的推送消息请求
// platform battery.PLATFORM_TYPE 平台类型
// deviceTokens string 推送的目标deviceTokens
// alert string 提示信息
func (api *XYAPI) Notify(platform battery.PLATFORM_TYPE, deviceTokens, alert string) {
	if function, ok := NOTIFY_FUNC_MAP[platform]; ok {
		xylog.DebugNoId("send alert(%s) to deviceTokens(%s) UMeng server", alert, deviceTokens)
		function(api, deviceTokens, alert)
	} else {
		xylog.ErrorNoId("search func pointer for %v failed", platform)
		return
	}

}

// Notify 向友盟服务器推送消息
// platform battery.PLATFORM_TYPE 平台类型
// accounts []*battery.DBAccount 待推送的玩家信息列表(devicetoken)
// alert string 提示信息
func (api *XYAPI) NotifyWithAccounts(platform battery.PLATFORM_TYPE, accounts []*battery.DBAccount, alert string) {

	rest, begin, end := len(accounts), 0, 0

	for {

		if rest <= 0 { //如果剩余的推送玩家数为0，则跳出循环
			break
		}

		//受友盟接口限制，每次请求的devicetoken数不能超过500，但是实际上为了通信包不会过大，控制在200以内
		begin = end
		if rest > DefConfigCache.Configs().ApnNotifyDevicePerReq {
			end += DefConfigCache.Configs().ApnNotifyDevicePerReq
		} else {
			end += rest
		}

		if function, ok := NOTIFY_FUNC_MAP[platform]; ok {
			xylog.DebugNoId("send accounts[%d:%d] to UMeng server", begin, end)
			function(api, api.getDeviceTokensFromAccounts(accounts[begin:end]), alert)
			rest -= (end - begin)
		} else {
			xylog.ErrorNoId("search func pointer for %v failed", platform)
			return
		}
	}
}

// getDeviceTokensFromAccounts 根据account信息列表获取devicetoken列表
// accounts []*battery.DBAccount 玩家账户信息列表
//return:
// deviceTokens string devicetoken列表
func (api *XYAPI) getDeviceTokensFromAccounts(accounts []*battery.DBAccount) (deviceTokens string) {
	var i = 0
	for _, account := range accounts {
		if i == 0 {
			deviceTokens += account.GetDeviceid()
			i++
		} else {
			deviceTokens += "," + account.GetDeviceid()
		}
	}
	return
}

// NotifyIOS 向ios设备推送消息
// accounts []*battery.DBAccount 待推送的玩家信息列表(devicetoken)
// alert string 提示信息
func (api *XYAPI) NotifyIOS(deviceTokens, alert string) {
	notification := &UMengIOSNotification{
		AppKey:       DefConfigCache.Configs().AppKey,
		Timestamp:    time.Now().Unix(),
		Type:         "listcast", //都是采用listcast的方式
		DeviceTokens: deviceTokens,
		Payload: UMengIOSPayload{
			Aps: UMengIOSAps{
				Alert: alert,
				Badge: 1,
			},
		},
		ProductionMode: DefConfigCache.Configs().IsApnsProduction,
	}

	NewXYAPI().NotifyUMeng(notification)
}

// NotifyAndroid 想android设备推送消息
// accounts []*battery.DBAccount 待推送的玩家信息列表(devicetoken)
// alert string 提示信息
func (api *XYAPI) NotifyAndroid(deviceTokens, alert string) {
	notification := &UMengAndroidNotification{
		AppKey:       DefConfigCache.Configs().AppKey,
		Timestamp:    time.Now().Unix(),
		Type:         "listcast", //都是采用listcast的方式
		DeviceTokens: deviceTokens,
		Payload: UMengAndroidPayload{
			Body: UMengPayloadBody{
				Text: alert,
			},
		},
		ProductionMode: DefConfigCache.Configs().IsApnsProduction,
	}

	NewXYAPI().NotifyUMeng(notification)
}

// 对deviceid进行去重处理
// accounts []*battery.DBAccount 待去重的玩家账户信息列表
// distinctAccounts []*battery.DBAccount 去重后的玩家账户信息列表
func (api *XYAPI) distinctDeviceToken(accounts []*battery.DBAccount) (distinctAccounts []*battery.DBAccount) {
	deviceTokens := make(xybusinesscache.Set, 0)
	for _, account := range accounts {
		deviceToken := account.GetDeviceid()

		if deviceToken == "" { //devicetoke 为空，跳过
			continue
		}

		if _, ok := deviceTokens[deviceToken]; ok { //devicetoken已经存在，跳过
			continue
		}
		distinctAccounts = append(distinctAccounts, account)
		deviceTokens[deviceToken] = xybusinesscache.Empty{}
	}

	return
}
