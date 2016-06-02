package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	//	"github.com/codegangsta/martini"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//获取好友排行榜消息
func (api *XYAPI) OperationGetFriendData(req *battery.QueryFriendsDataRequest, resp *battery.QueryFriendsDataResponse) (err error) {
	var (
		uid       = req.GetUid()
		src       = req.GetSource()
		sids      = req.GetSids()
		maxReturn int
		counter   int
		users     = make([]*battery.FriendData, 0)
		acts      = make([]*battery.GiftActivity, 0)
		uids      = make([]string, 0)
		tpids     = make([]*battery.IDMap, 0)
		//dbAccounts                      = make([]*battery.DBAccount, 0)
		userAccomplishments = make([]*battery.DBUserAccomplishment, 0)
		//mapUid2Sid                      = make(map[string]string, 0)
		mapUid2Tpid                     = make(map[string]*battery.IDMap, 0)
		mapFriendUid2StaminaGiveLogItem = make(map[string]*battery.StaminaGiveLogItem, 0)
		now                             = time.Now().Unix()
		errStruct                       = new(battery.Error)
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	*errStruct = *(xyerror.Resp_NoError)
	if !api.isUidValid(uid) {
		xylog.Error(uid, "OperationGetFriendData invalid uid")
		errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
		err = xyerror.ErrGetAccountByUidError
		goto ErrHandle
	}

	maxReturn = len(sids)
	if maxReturn > api.Config.Configs().MaxFriendsRequestCount {
		maxReturn = api.Config.Configs().MaxFriendsRequestCount
	}

	xylog.Debug(uid, "Sid count: %d", len(sids))

	//查询好友的tpid信息列表
	{
		selector := bson.M{"gid": 1, "sid": 1, "note": 1, "iconurl": 1}
		err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).XYBusinessDB.QueryTpidsBySids(selector, sids, src, &tpids)
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_QueryUserFriendsUidError.Enum()
			err = xyerror.ErrQueryUserFriendsUidError
			goto ErrHandle
		}
	}

	for _, tpid := range tpids {
		gid := tpid.GetGid()
		//sid := tpid.GetSid()
		if len(gid) > 0 {
			mapUid2Tpid[gid] = tpid
			uids = append(uids, gid)
		}
	}

	//查询玩家的账户信息
	//err = api.GetDBAccountsDirect(uids, &dbAccounts, mgo.Monotonic) //不需要实时信息，可以从备节点中查询
	err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERACCOMPLISHMENT).GetUsersAccoplishment(uids, &userAccomplishments, mgo.Monotonic)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
		err = xyerror.ErrGetAccountByUidError
		goto ErrHandle
	}

	//查询玩家所有好友的体力赠送日志
	api.GetUserStaminaGiveItems(uid, uids, &mapFriendUid2StaminaGiveLogItem, errStruct)
	if errStruct.GetCode() != battery.ErrorCode_NoError {
		goto ErrHandle
	}

	//for _, account := range dbAccounts {
	//xylog.Debug("[%s] userAccomplishments: %v", uid, userAccomplishments)
	for _, userAccomplishment := range userAccomplishments {
		var (
			gid = userAccomplishment.GetUid()
			//sid string
			tpid                        *battery.IDMap
			ok                          bool
			askCountDown, giveCountDown int64
		)

		//user := api.getFriendDataFromAccount(account)
		user := api.getFriendDataFromAccomplishment(userAccomplishment)
		//将好友id转换为第三方账户id
		if tpid, ok = mapUid2Tpid[gid]; ok {
			user.Uid = proto.String(tpid.GetSid())
			user.IconUrl = proto.String(tpid.GetIconUrl())
			user.Name = proto.String(tpid.GetNote())
		}
		users = append(users, user)

		xylog.Debug(uid, "gid : %v, mapFriendUid2StaminaGiveLogItem : %v", gid, mapFriendUid2StaminaGiveLogItem)

		if item, ok := mapFriendUid2StaminaGiveLogItem[gid]; ok {
			askCountDown = api.GetStaminaWaitTime(now, item.GetStaminaApplyLastTime(), api.Config.Configs().DefaultGiftAskCooldown)
			giveCountDown = api.GetStaminaWaitTime(now, item.GetStaminaGiveLastTime(), api.Config.Configs().DefaultGiftGiveCooldown)
		}

		act := &battery.GiftActivity{
			FriendSid:         proto.String(tpid.GetSid()),
			GiftGiveCountDown: proto.Int64(giveCountDown),
			GiftAskCountDown:  proto.Int64(askCountDown),
		}
		acts = append(acts, act)

		counter++
		if counter > maxReturn {
			// 最多返回10条记录
			break
		}
	}

	xylog.Debug(uid, "Find friends count: %d", counter)

	resp.Data = users
	resp.GiftActivities = make([]*battery.GiftActivity, len(acts))
	resp.GiftActivities = acts
	xylog.Debug(uid, "GiftActivities : %v", resp.GiftActivities)
	resp.Source = src.Enum()
	resp.SystemTime = proto.Int64(xyutil.CurTimeSec())

ErrHandle:

	resp.Error = errStruct

	xylog.Debug(uid, "errStruct : %v, resp.Error : %v", errStruct, resp.Error)

	return
}

//查询玩家相关的所有好友体力赠送日志
// uid string 玩家id
// friendUids []string 玩家好友id列表
// mapFriendUid2StaminaGiveLogItem
//return:
// errStruct *battery.Error 返回的错误信息
func (api *XYAPI) GetUserStaminaGiveItems(uid string, friendUids []string, mapFriendUid2StaminaGiveLogItem *map[string]*battery.StaminaGiveLogItem, errStruct *battery.Error) {
	*errStruct = *(xyerror.Resp_NoError)
	staminaGiveLogItems := make([]*battery.StaminaGiveLogItem, 0)
	err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_STAMINAGIVEAPPLYLOG).GetStaminaGiveLogItems(uid, friendUids, &staminaGiveLogItems)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		errStruct.Code = battery.ErrorCode_QueryUserStaminaGiveLogError.Enum()
		return
	}

	for _, item := range staminaGiveLogItems {
		friendUid := item.GetFriendUid()
		pItem := item
		(*mapFriendUid2StaminaGiveLogItem)[friendUid] = pItem
	}

	xylog.Debug(uid, "mapFriendUid2StaminaGiveLogItem : %v", *mapFriendUid2StaminaGiveLogItem)

	return
}

//获取剩余等待时间
// now int64      当前时间戳，单位：秒
// lastTime int64 上次更新时间戳，单位：秒
// defautInterval int64 更新间隔，单位：秒
//return:
// waitTime int64 剩余等待时间，单位：秒
func (api *XYAPI) GetStaminaWaitTime(now, lastTime, defautInterval int64) (waitTime int64) {
	nextTime := lastTime + defautInterval
	waitTime = nextTime - now
	if waitTime < 0 { //如果已经超出了时间间隔，则等待时间为0
		waitTime = 0
	}
	return
}
