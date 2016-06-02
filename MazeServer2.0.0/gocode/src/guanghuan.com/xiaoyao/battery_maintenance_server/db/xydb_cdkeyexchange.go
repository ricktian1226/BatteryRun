package db

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

type CDkeyRequest struct {
	Appid              string // 游戏后台id
	Channel            string // 渠道
	Platform           string // 来源标识
	Code               string // 兑换码
	Timestamp          int64  // 请求时间戳
	Userip, Userid     string // 玩家ip，区服，玩家id
	Zoneid             int64
	UserName, NickName string // 账号，昵称
	Sig                string // 加密串
}

// 由uid获取账户信息并构造运营兑换消息
func (db *BatteryDB) GetAccountInfor(req *battery.CDkeyExchangeRequest) (resp CDkeyRequest, err error) {
	uid := req.GetUid()
	platform := req.GetPlatform()

	resp.Timestamp = time.Now().Unix()
	resp.Channel = req.GetChannelId()
	resp.Code = req.GetCdkey()
	resp.Userip = "127.0.0.1"
	resp.Zoneid = 1 // 不分区服设置为1

	switch platform {
	// 如果是通过网页进行兑换，需要将uid转换成真实的用户uid
	case battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN:
		resp.Platform = "web"
		err = getUidByIdentity(uid)
		if err != nil {
			xylog.ErrorNoId("get uid by uid identity:%v", err.Error())
			return
		}
	case battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS:
		resp.Platform = "ios"
	case battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID:
		resp.Platform = "android"
	default:
	}

	resp.Userid = uid
	err = db.GetAccountBuyUid(uid, &resp)
	if err != nil {
		return
	}
	xylog.DebugNoId("resp :%v", resp)

	return
}

func (db *BatteryDB) GetAccountBuyUid(uid string, resp *CDkeyRequest) (err error) {

	selector, tpid := bson.M{"sid": 1, "note": 1}, &battery.IDMap{}
	err = db.XYBusinessDB.QueryTpidByUid(selector, uid, tpid, mgo.Strong)
	if err != nil {
		xylog.ErrorNoId("get account by uid error:%v", err.Error())
		return
	}
	resp.NickName = tpid.GetNote()
	if resp.NickName == "" {
		resp.NickName = resp.Userid //没有用户名则设置为id
	}
	resp.UserName = resp.NickName

	return

}

// 通过网页兑换是将uid转换为真实用户uid
func getUidByIdentity(uid string) (err error) {
	return
}
