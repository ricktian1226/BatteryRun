// xyapi_rune
package batteryapi

import (
	"time"

	"code.google.com/p/goprotobuf/proto"

	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//符文枚举值
var (
	RUNE_STAMINA_LIMIT_10   = uint64(110010000) //体力上限10
	RUNE_STAMINA_LIMIT_20   = uint64(110010001) //体力上限20
	RUNE_COIN_ADDITIONAL    = uint64(110020000) //金币增加系数
	RUNE_LOTTO_ADDITIONAL   = uint64(110030000) //抽奖次数上限
	RUNE_RESOLVE_ADDITIONAL = uint64(110040000) //道具分解增加系数
	RUNE_BASE_FACTOR        = int32(100)        //符文基础系数(100代表100%)
)

//符文操作消息接口
func (api *XYAPI) OperationRune(req *battery.RuneRequest, resp *battery.RuneResponse) (err error) {

	var (
		uid       = req.GetUid()
		errStruct = xyerror.DefaultError()
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//初始化返回消息
	resp.Uid, resp.Cmd, resp.Itemid = req.Uid, req.Cmd, req.Itemid

	switch req.GetCmd() {
	case battery.RuneCmd_RuneCmd_Query:
		resp.Runeidlist, resp.RuneList = api.queryUserRunes(uid, errStruct)
	}

	return
}

// 符文的永久时效
const IMMORTAL_RUNE_LIMITATION int64 = 10 * 360 * 60 * 60 * 24

//查询玩家拥有的符文信息
func (api *XYAPI) queryUserRunes(uid string, errStruct *battery.Error) (runeIdList []uint64, runeUnitList []*battery.RuneUnit) {
	runeList, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_RUNE).GetRuneList(uid)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_QueryUserRuneInfoError.Enum()
		return
	} else {
		var now = time.Now().Unix()
		for _, runeItem := range runeList {
			//旧版本的符文是没有时效限制的，需要修改成有时效的
			//对于这部分玩家，我们给他们的符文10年的时效时间，如果10年后他们来找我们申诉，那我们该好好请他们吃顿饭 :=)
			if runeItem.GetExpiredTimestamp() == 0 {
				runeItem.Uid = proto.String(uid)
				runeItem.ExpiredTimestamp = proto.Int64(now + IMMORTAL_RUNE_LIMITATION)
				err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_RUNE).UpsertRune(runeItem)
				if err != xyerror.ErrOK {
					xylog.Error(uid, "UpsertRune(%v) failed : %v", runeItem, err)
					continue
				}

				err = api.runeEffect(uid, runeItem.GetId(), runeItem.GetExpiredTimestamp(), nil, ACCOUNT_UPDATE_NO_DELAY)
				if err != xyerror.ErrOK {
					xylog.Error(uid, "runeEffect(%v) failed : %v", runeItem, err)
					continue
				}

			}

			//符文信息保存到符文列表
			runeTmp := &battery.RuneUnit{
				Id:               proto.Uint64(runeItem.GetId()),
				ExpiredTimestamp: proto.Int64(runeItem.GetExpiredTimestamp()),
			}
			runeUnitList = append(runeUnitList, runeTmp)

			//为了兼容旧版本必须传符文id列表。勿删！！！
			//将符文id加入id列表
			runeIdList = append(runeIdList, runeTmp.GetId())
		}
	}

	return
}

const (
	STAMINA_STOP_UPDATE_TIME int64 = -1
)

//增加玩家的符文数据
// uid string 玩家id
// propItem *battery.PropItem 符文信息
func (api *XYAPI) updateUserDataWithRune(uid string, accountWithFlag *AccountWithFlag, propItem *battery.PropItem, delay bool) (err error) {

	var (
		runeid   = propItem.GetId()
		amount   = propItem.GetAmount()
		userRune *battery.Rune
	)

	now := time.Now().Unix()
	expiredLimitation := int64(60 * 60 * amount) //每个符文的时效是1hour

	//获取玩家当前的对应符文信息
	userRune, err = api.queryUserRune(uid, runeid)
	if err == xyerror.ErrOK { //玩家已经拥有该符文
		xylog.Debug(uid, "rune(%d) existes, so let's set it longer. %v", runeid, userRune)

		expiredTimestamp := userRune.GetExpiredTimestamp()
		if expiredTimestamp > now { //如果未过期，累加上时效
			userRune.ExpiredTimestamp = proto.Int64(expiredTimestamp + expiredLimitation)
		} else { //如果已过期，以当前时间为起点，加上时效
			userRune.ExpiredTimestamp = proto.Int64(now + expiredLimitation)
		}

	} else if err == xyerror.ErrNotFound { //玩家未拥有该符文
		//则增加该符文，以当前时间为起点，加上时效
		userRune = &battery.Rune{
			Uid:              proto.String(uid),
			Id:               proto.Uint64(runeid),
			ExpiredTimestamp: proto.Int64(now + expiredLimitation),
		}

		xylog.Debug(uid, "rune(%d) does't exist, add it to db and set expired time to %d", runeid, userRune.GetExpiredTimestamp())
	} else {
		xylog.Error(uid, "queryUserRune failed : %v", err)
		return
	}

	//刷新玩家的符文信息
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_RUNE).UpsertRune(userRune)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "UpsertRune(%v) failed : %v", userRune, err)
		return
	}

	//刷新符文效果
	expiredTimestamp := userRune.GetExpiredTimestamp()
	err = api.runeEffect(uid, runeid, expiredTimestamp, accountWithFlag, delay)

	return
}

// 符文生效
// uid string 玩家标识
// runeid uint64 符文标识
// expiredTimestamp int64 符文过期时间
// accountWithFlag *AccountWithFlag 玩家账户信息指针
// delay bool 是否延迟刷新
func (api *XYAPI) runeEffect(uid string, runeid uint64, expiredTimestamp int64, accountWithFlag *AccountWithFlag, delay bool) (err error) {
	//获取符文的配置信息
	var runeConfig *battery.RuneConfig
	runeConfig, err = api.RuneConfigDetail(runeid)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "get RuneConfigDetail(%d) failed : %v", runeid, err)
		return
	}
	// 获取符文对应的效果值
	value := runeConfig.GetValue()
	if value > 0 {
		switch runeid {
		// 免费抽奖次数上限符文
		case RUNE_LOTTO_ADDITIONAL:
			err = api.updateUserDataWithSysLottoFreeCountLimitation(uid, value, expiredTimestamp)

		// 体力上限符文
		case RUNE_STAMINA_LIMIT_10, RUNE_STAMINA_LIMIT_20:
			maxStamina := DefConfigCache.Configs().DefaultMaxRegStamina
			if value > maxStamina {
				err = api.fillAccountWithFlag(uid, &accountWithFlag)
				if err == xyerror.ErrOK {
					account := accountWithFlag.account
					// 如果之前达到了体力上限值，则修改上次体力刷新时间，放开刷新
					if STAMINA_STOP_UPDATE_TIME == account.GetStaminaLastUpdateTime() {
						account.StaminaLastUpdateTime = proto.Int64(xyutil.CurTimeSec())
						xylog.Debug(uid, "account.StaminaLastUpdateTime set to %d", account.GetStaminaLastUpdateTime())
					}
					if account.StaminaLimitation == nil { //只有在第一次购买的时候
						account.StaminaLimitation = proto.Int32(value)
					}
					account.StaminaLimitationExpiredTimestamp = proto.Int64(expiredTimestamp)
					accountWithFlag.SetChange()
					xylog.Info(uid, "account info ,updatetime:%v,limit:%v", account.GetStaminaLastUpdateTime(), account.GetStaminaLimitation())

					if !delay {
						err = api.UpdateAccountWithFlag(accountWithFlag)
					}
				}
			}

		// 结算金币加成符文
		case RUNE_COIN_ADDITIONAL:
			err = api.fillAccountWithFlag(uid, &accountWithFlag)
			if err == xyerror.ErrOK {
				account := accountWithFlag.account
				account.CoinAddtional = proto.Int32(value)
				account.CoinAddtionalExpiredTimestamp = proto.Int64(expiredTimestamp)
				accountWithFlag.SetChange()

				if !delay {
					err = api.UpdateAccountWithFlag(accountWithFlag)
				}
			}

		default: //do nothing
		}
	}

	return
}

//判断玩家是否已经拥有了对应符文
// uid string 玩家id
// runeid uint64 符文标识
//return:
// isExisting bool true 已拥有，false 未拥有
func (api *XYAPI) IsRuneExisting(uid string, runeid uint64) (isExisting bool) {
	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_RUNEEXIST, &begin)

	isExisting, _ = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_RUNE).IsRuneExisting(uid, runeid)
	return
}

//判断玩家是否已经拥有了对应的有效符文
// uid string 玩家id
// runeid uint64 符文标识
//return:
// result bool true 已拥有，false 未拥有
func (api *XYAPI) IsRuneValid(uid string, runeid uint64) (result bool) {
	userRune, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_RUNE).GetRune(uid, runeid)
	if err == xyerror.ErrOK {
		if userRune.GetExpiredTimestamp() == 0 || //老玩家买的符文是没有时效限制的，所有没有这一项就认为是永久拥有
			userRune.GetExpiredTimestamp() > time.Now().Unix() { //如果找到符文信息，并且符文未过期
			result = true
		}
	}
	return
}

//查询符文道具详细信息
// runeid uint64 符文标识
//return:
// runeConfig *battery.RuneConfig 符文配置信息
func (api *XYAPI) RuneConfigDetail(runeid uint64) (runeConfig *battery.RuneConfig, err error) {
	runeConfig = xybusinesscache.DefRuneConfigCacheManager.RuneConfig(runeid)
	if nil == runeConfig { //没找到
		err = xyerror.ErrRuneConfigsFromCacheError
	}
	return
}

//查询符文道具对应的效果值
// runeid uint64 符文标识
//return:
// value int32 效果值
func (api *XYAPI) RuneConfigValue(runeid uint64) (value int32) {
	var runeconfig *battery.RuneConfig
	runeconfig, _ = api.RuneConfigDetail(runeid)
	if nil != runeconfig { //没找到
		value = runeconfig.GetValue()
	}
	return
}

// 查询玩家背包中是否有某个符文
// uid string 玩家标识
// id uint64 符文标识
//returns:
// userRune *battery.Rune 玩家符文信息
// err error 返回错误信息
func (api *XYAPI) queryUserRune(uid string, id uint64) (userRune *battery.Rune, err error) {
	userRune, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_RUNE).GetRune(uid, id)
	return
}
