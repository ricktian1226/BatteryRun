package batteryapi

import (
	//proto "code.google.com/p/goprotobuf/proto"
	//"encoding/binary"
	//"sync"
	//"time"

	//apns "github.com/timehop/apns"
	//batterydb "guanghuan.com/xiaoyao/battery_apns_server/db"
	//"guanghuan.com/xiaoyao/common/apn"
	//xyconf "guanghuan.com/xiaoyao/common/conf"
	//"guanghuan.com/xiaoyao/common/db"
	"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// LottoChanceNotify 抽奖机会推送消息
func (api *XYAPI) LottoChanceNotify() {

	platforms := []battery.PLATFORM_TYPE{battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS}

	for _, platform := range platforms {
		//查找当天抽奖券还有剩余的玩家
		api.SetDB(platform)
		lottoInfos, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).QueryUidsSysLottoFreeCountNoNull()
		if err != xyerror.ErrOK { //查询失败，跳过
			xylog.ErrorNoId("QueryUidsSysLottoFreeCountNoNull failed : %v", err)
			continue
		} else if len(lottoInfos) <= 0 { //查询没有符合条件的uid，跳过
			xylog.DebugNoId("QueryUidsSysLottoFreeCountNoNull no user")
			continue
		}

		uids := make([]string, 0)
		for _, lottoInfo := range lottoInfos {
			uids = append(uids, lottoInfo.GetUid())
		}

		//查找玩家的devicetoken
		var accounts []*battery.DBAccount
		accounts, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).QueryDeviceTokenByUids(uids)
		if err != xyerror.ErrOK { //查询失败，跳过
			xylog.ErrorNoId("QueryDeviceTokenByUids failed : %v", err)
			continue
		} else if len(accounts) <= 0 { //查询没有符合条件的accounts，跳过
			xylog.DebugNoId("QueryDeviceTokenByUids no user")
			continue
		}

		// 对devicei进行去重处理
		accounts = api.distinctDeviceToken(accounts)

		//4 test
		//accounts = []*battery.DBAccount{&battery.DBAccount{
		//	Deviceid: proto.String("752758d0f23966677e036df2ae9f681dcb93193e1255bfabdb33e38b06179a74"),
		//}, &battery.DBAccount{
		//	Deviceid: proto.String("f15d3410a815c3a813123026077cd2dd4e837ba156b00a5d9c8facc3e4ea819e"),
		//},
		//}
		var tip string
		tip, err = xybusinesscache.DefTipManager.Tip(battery.LANGUAGE_TYPE_LANGUAGE_TYPE_CHINESE, battery.TIP_IDENTITY_TIP_IDENTITY_LOTTOCHANGE)
		if err == xyerror.ErrOK {
			api.NotifyWithAccounts(platform, accounts, tip)
		} else {
			xylog.ErrorNoId("Get tip for %v %v failed : %v", battery.LANGUAGE_TYPE_LANGUAGE_TYPE_CHINESE, battery.TIP_IDENTITY_TIP_IDENTITY_LOTTOCHANGE, err)
		}
	}

}
