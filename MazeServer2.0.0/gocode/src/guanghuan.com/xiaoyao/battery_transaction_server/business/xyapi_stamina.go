package batteryapi

import (
	"math"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"

	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	//"guanghuan.com/xiaoyao/superbman_server/server"
)

//查询当前玩家体力信息
// uid string 玩家uid
//return:
// amount int32 玩家体力数目
// timeleft int32 玩家体力刷新剩余时间
func (api *XYAPI) GetCurrentStamina(uid string, change int32) (amount int32, timeleft int32, err error) {

	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_GETCURRENTSTAMINA, &begin)

	account := &battery.DBAccount{}
	err = api.GetDBAccountDirect(uid, account, mgo.Strong)
	if err != xyerror.ErrOK {
		return
	}

	accountWithFlag := &AccountWithFlag{
		account: account,
		bChange: false,
	}

	amount, timeleft, err = api.getCurrentStamina(uid, accountWithFlag, change)
	if err == xyerror.ErrOK {
		err = api.UpdateAccountWithFlag(accountWithFlag)
	}

	return
}

//查询玩家当前账户信息（带刷新体力信息）
// uid string 玩家id
// account *battery.DBAccount 玩家账户信息
// change int32 体力的改变值（例如购买xx体力或者游戏消耗xx体力）
//return:
// amount int32 体力数目
// timeleft int32 剩余刷新时间
//func (api *XYAPI) getCurrentStamina(uid string, account *battery.DBAccount, change int32) (amount int32, timeleft int32, err error) {
func (api *XYAPI) getCurrentStamina(uid string, accountWithFlag *AccountWithFlag, change int32) (amount int32, timeleft int32, err error) {

	timeleft = int32(DefaultStaminaLastUpdateTime)

	//获取体力上限
	maxAmount := DefConfigCache.Master().Configs.DefaultMaxRegStamina //默认玩家体力上限
	if accountWithFlag.account.GetStaminaLimitationExpiredTimestamp() > xyutil.CurTimeSec() {
		maxAmount = accountWithFlag.account.GetStaminaLimitation()
	}

	xylog.Debug(uid, "stamina maxAmount : %d ", maxAmount)
	//获取当前的体力数
	account := accountWithFlag.account
	if len(account.Wallet) < int(battery.MoneyType_MoneyType_max) {
		err = xyerror.ErrBadInputData
		xylog.Error(uid, "len(account.Wallet)(%d) < battery.MoneyType_MoneyType_max(%d)", len(account.Wallet), battery.MoneyType_MoneyType_max)
		return
	}
	amountChange := false
	amount = int32(account.Wallet[int(battery.MoneyType_stamina)].GetGainamount())
	currentTime := xyutil.CurTimeSec()
	startTime := account.GetStaminaLastUpdateTime()
	regInterval := DefConfigCache.Configs().DefaultStaminRegIntervalSec // 自动增长体力的间隔，秒
	changeDiff := int32(0)
	if amount < maxAmount { // 当前体力值没到上限，则根据时间差，计算最新体力值
		if startTime == DefaultStaminaLastUpdateTime {
			// 这种异常情况出现在玩家购买了体力上限符文后，但是因为某种原因更新StaminaLastUpdateTime失败
			// 在此修复一下
			xylog.Error(uid, "DefaultStaminaLastUpdateTime(%d) error status amount(%d) limitAmount(%d), set StaminaLastUpdateTime to (%d)", DefaultStaminaLastUpdateTime, amount, maxAmount, currentTime)
			account.StaminaLastUpdateTime = proto.Int64(currentTime)
			accountWithFlag.SetChange()
		} else { //上次体力未达到上限，则根据时间差，计算出当前的体力值
			if startTime <= currentTime { //时间差要校验一下  设置为《=，购买符文时会出现更新时间和请求时间相等的情况
				changeDiff = int32(math.Floor((float64(currentTime - startTime)) / float64(regInterval)))
				xylog.Debug(uid, "stamina amount %d changeDiff: %d", amount, changeDiff)

				if changeDiff > 0 {
					amount += changeDiff //加上时间差产生的体力
					amountChange = true

					//超过玩家体力上限，则设置为上限值
					if amount > maxAmount {
						amount = maxAmount
					} else if amount < 0 {
						xylog.Error(uid, "stamina Error status : amount(%d) < 0", amount)
						err = xyerror.ErrNotEnoughStamina
						return
					}
				}
			} else { //出现当前时间戳小于上次更新时间戳的状态是非法的，直接返回错误
				xylog.Error(uid, "stamina Error status : currentTime(%d) > StaminaLastUpdateTime(%d)", currentTime, startTime)
				err = xyerror.ErrQueryStaminaError
				return
			}
		}
	}

	//结合体力的变化值，算出当前体力值
	if change != 0 {
		amount += change
		amountChange = true
		if amount < 0 { //在消费体力的时候，可能出现当前值小于需要的体力值
			err = xyerror.ErrNotEnoughStamina
			return
		}
	}

	//如果体力值变更，需要刷新体力值和lastUpdateTime
	if amountChange { //
		if amount >= maxAmount {
			account.StaminaLastUpdateTime = proto.Int64(DefaultStaminaLastUpdateTime)
		} else {
			if startTime == DefaultStaminaLastUpdateTime { //之前时间戳是默认值的，设置为当前值；如果不是默认值，则保持之前的值
				account.StaminaLastUpdateTime = proto.Int64(currentTime)
			} else {
				//设置一下lastUpdateTime,如果changeDiff==0就保持不变
				account.StaminaLastUpdateTime = proto.Int64(startTime + regInterval*(int64(changeDiff)))
			}
		}

		account.Wallet[int(battery.MoneyType_stamina)].Gainamount = proto.Uint32(uint32(amount))
		accountWithFlag.SetChange()
	}

	//根据体力值设置返回的timeleft
	if amount >= maxAmount {
		timeleft = int32(DefaultStaminaLastUpdateTime)
	} else {
		timeleft = int32(account.GetStaminaLastUpdateTime() + regInterval - currentTime)
	}

	xylog.Debug(uid, "[getCurrentStamina] return amountChange= %t stamina=%d, timeleft=%d", amountChange, amount, timeleft)

	if timeleft < -1 || timeleft > int32(DefConfigCache.Configs().DefaultStaminRegIntervalSec) {
		xylog.Error(uid, "WARNING: invalid timeleft : %d", timeleft)
	}

	return
}

//直接更新玩家当前体力，没有其他体力变更相关的行为，如玩家开始新游戏 -1, 购买体力 +x
// uid string 玩家id
//return:
// amount int32 剩余的体力数目
// timeleft int32 体力刷新剩余时间，单位：秒
func (api *XYAPI) GetCurrentStaminaDirect(uid string) (amount int32, timeleft int32, err error) {
	amount, timeleft, err = api.GetCurrentStamina(uid, 0)
	return
}

// 更新玩家当前体力，同事有其他体力变更相关的行为，如玩家开始新游戏 -1, 购买体力 +x
// uid string 玩家id
// change int32 体力变更数目
//return:
// amount int32 剩余的体力数目
// timeleft int32 体力刷新剩余时间，单位：秒
func (api *XYAPI) UpdateStamina(uid string, change int32) (amount int32, timeleft int32, err error) {
	amount, timeleft, err = api.GetCurrentStamina(uid, change)
	return
}

//更新玩家体力数据
// uid string 玩家id
// accountWithFlag *AccountWithFlag 保存玩家数据的临时数据结构
// amount uint32 体力数目
// delay bool 是否刷新延迟
func (api *XYAPI) updateUserDataWithStamina(uid string, accountWithFlag *AccountWithFlag, amount uint32, delay bool) (err error) {
	if delay {
		_, _, err = api.getCurrentStamina(uid, accountWithFlag, int32(amount))
	} else {
		_, _, err = api.UpdateStamina(uid, int32(amount))
	}
	return
}

//// updateUserDataWithStaminaLastUpdateTime 刷新玩家的上次体力刷新时间
//// 在玩家体力上线发生变化的时候，如果之前玩家的上次体力刷新时间是-1，则需要刷新为当前时间。这样重启体力计算。
//// uid string 玩家标识
//// limitAddtional int32 玩家体力上限加成
//func (api *XYAPI) updateUserDataWithStaminaLastUpdateTimeAndLimitAddtional(uid string, addtional int32) (err error) {
//	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).UpdateStaminaLastUpdateTimeAndLimitAddtional(uid, addtional)
//}

//// updateUserDataWithStaminaLimitAddtional 刷新玩家的上次体力刷新时间
//// 用处：在玩家体力上线发生变化的时候，如果之前玩家的上次体力刷新时间是-1，则需要刷新为当前时间。这样重启体力计算。
//// uid string 玩家标识
//// limitAddtional int32
//func (api *XYAPI) updateUserDataWithCoinAddtional(uid string, addtional int32) (err error) {
//	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).UpdateCoinAddtional(uid, addtional)
//}

//// updateUserDataWithStaminaLimitAddtional 刷新玩家的上次体力刷新时间
//// 用处：在玩家体力上线发生变化的时候，如果之前玩家的上次体力刷新时间是-1，则需要刷新为当前时间。这样重启体力计算。
//// uid string 玩家标识
//// limitAddtional int32
//func (api *XYAPI) updateUserDataWithResolveAddtional(uid string, addtional int32) (err error) {
//	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).UpdateResolveAddtional(uid, addtional)
//}
