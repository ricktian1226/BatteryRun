// xyapi_lotto
package batteryapi

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	proto "code.google.com/p/goprotobuf/proto"

	"guanghuan.com/xiaoyao/common/idgenerate"
	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xycache "guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/money"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//抽奖id生成器
var DefLottoIdGenerater *xyidgenerate.IdGenerater

//游戏后抽奖删除次数和商品id的对应关系
//var DefAfterGameLottoDeleteCount2MallItem = make(map[int]uint64, 0)

//抽奖相关的配置项的初始化
// 游戏后抽奖删除格子删除次数与商品id的对应关系
func LottoInit() {

	cache := DefConfigCache.Slave()

	var NOMALLITEM = uint64(0)
	s := strings.Split(cache.Configs.AfterGameLottoMallItems, ";")
	var mallId uint64
	for i, v := range s {
		tmp := strings.TrimSpace(v)
		if len(tmp) <= 0 {
			mallId = NOMALLITEM
		} else {
			mallIdTmp, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				mallId = NOMALLITEM
			} else {
				mallId = uint64(mallIdTmp)
			}
		}
		cache.AfterGameLottoDeleteCount2MallItem[i+1] = mallId
	}
}

//操作入口
func (api *XYAPI) OperationLottoRequest(req *battery.LottoRequest, resp *battery.LottoResponse) (err error) {
	var (
		uid        = req.GetUid()
		cmd        = req.GetCmd()
		failReason = battery.ErrorCode_NoError
		now        = time.Now().Unix()
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//xylog.Debug(uid, "[%s] [LottoRequest] cmd %v start", uid, cmd)
	//defer xylog.Debug(uid, "[%s] [LottoRequest] cmd %v done", uid, cmd)

	//初始化resp
	resp.Cmd = req.Cmd
	resp.Uid = req.Uid
	resp.Lottoid = req.Lottoid
	resp.Parentlottoid = req.Parentlottoid
	resp.Error = xyerror.DefaultError()

	switch cmd {
	case battery.LottoCmd_Lotto_Initial:
		failReason, err = api.initialSysLotto(req, resp, now)
	case battery.LottoCmd_Lotto_Commit, battery.LottoCmd_Lotto_AfterGame_Commit:
		failReason, err = api.commitLotto(req, resp)
	case battery.LottoCmd_Lotto_AfterGame_Initial:
		failReason, err = api.initialAfterGameLotto(req, resp)
	case battery.LottoCmd_Lotto_AfterGame_NoInitial:
		failReason, err = api.noInitialAfterGameLotto(req, resp)
	default:
		xylog.Error(uid, "Unkown Lotto Cmd<%v>", cmd)
		failReason = xyerror.LOTTO_UNKOWN_CMD
		goto ErrHandle
	}

	if cmd == battery.LottoCmd_Lotto_Commit || cmd == battery.LottoCmd_Lotto_Commit {
		if err == xyerror.ErrOK {
			resp.Wallet, err = xymoney.QueryWallet(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT))
		}
	}

ErrHandle:
	resp.Error.Code = failReason.Enum()
	xylog.Debug(uid, "failReason : %d, err : %v", failReason, err)

	return
}

//系统抽奖请求
func (api *XYAPI) initialSysLotto(req *battery.LottoRequest, resp *battery.LottoResponse, now int64) (failReason battery.ErrorCode, err error) {
	var (
		uid              = req.GetUid()
		cmd              = req.GetCmd()
		forceRefreshSlot = req.GetForceRefreshSlot() //是否强制刷新奖池信息
		serialNum        = req.GetSerialNum()        //获取抽奖指定序列号
		info             = &battery.SysLottoInfo{
			Uid: &uid,
		}
		restFreshSec                     = new(int64)
		lottoid                          uint64
		lottoTransaction                 *battery.LottoTransaction
		restFreeCount, restGainFreeCount int32
		selected                         = new(uint32)
		stuff                            *battery.LottoStuff
	)

	if serialNum > 0 { //前端指定了抽奖标识
		failReason, err = api.getSysSerialLottoInfo(uid, serialNum, &info, now, selected)
		xylog.Debug(uid, "getSysSerialLottoInfo(%d) : info(%v) selected(%d) err : %v", serialNum, info, *selected, err)
		if err != xyerror.ErrOK {
			goto ErrHandle
		}

	} else {
		failReason, err = api.getSysLottoInfo(uid, &info, restFreshSec, now, forceRefreshSlot)
		if err != xyerror.ErrOK {
			goto ErrHandle
		}

		xylog.Debug(uid, "info : %v, value : %d", *info, info.GetValue())

		failReason, err = api.getSysSelectedSlot(uid, info.GetValue(), selected)
		if failReason != battery.ErrorCode_NoError || err != xyerror.ErrOK {
			goto ErrHandle
		}

	}

	//生成lottoid
	lottoid = DefLottoIdGenerater.NewID()
	restFreeCount = info.GetFreeCount()
	//if info.GetFreecountAdditionalExpiredTimestamp() > time.Now().Unix() { //如果有加成，需要加上加成的
	//	restFreeCount += info.GetFreecountAdditional()
	//}
	restGainFreeCount = info.GetGainFreeCount()
	stuff = &battery.LottoStuff{
		Lottoid:       proto.Uint64(lottoid),
		Parentlottoid: proto.Uint64(lottoid),          //初始抽奖内容的lottoid和parentlottoid相同
		RestFreeNum:   proto.Int32(restFreeCount),     //免费抽奖次数
		RestGainNum:   proto.Int32(restGainFreeCount), //抽奖券数
		RestFreshSec:  restFreshSec,
		//（刷新时间只有非抽奖标识的抽奖需要）
		Selected: proto.Uint32(*selected), //本次抽奖的结果提前算出来
		Slots:    info.Slots,
	}

	//添加事务信息
	lottoTransaction = api.defaultLottoTransaction(uid, battery.DrawAwardType_DrawAwardType_System, stuff, now)
	lottoTransaction.SerialNum = proto.Int32(serialNum)
	err = api.addLottoTransaction(lottoTransaction)
	if err != nil {
		xylog.Error(uid, "add lottoTransaction failed : %v", err)
		goto ErrHandle
	} else {
		xylog.Debug(uid, "add lottoTransaction  : %v", *lottoTransaction)
	}

	resp.Stuff = stuff

ErrHandle:
	//记录日志
	l := api.defaultLottoLog(uid, cmd, info.GetValue(), battery.DrawAwardType_DrawAwardType_System, stuff, resp.Error)

	go api.addLottoLog(l)

	return
}

//构建默认的用户系统抽奖信息
func (api *XYAPI) defaultSysLottoInfo(uid string, now int64) (info *battery.SysLottoInfo) {
	//获取玩家系统抽奖次数
	info = &battery.SysLottoInfo{
		Uid:                       proto.String(uid),
		Timestamp:                 proto.Int64(now),
		Value:                     proto.Int32(DefConfigCache.Configs().LottoInitUserValue),
		GainFreeCount:             proto.Int32(0),
		FreecountRefreshTimestamp: proto.Int64(now),
	}

	freeCount, freeCountLimitation, expiredTimestamp := api.getDefaultSysLottoFreeCount(uid)
	info.FreeCount, info.FreecountLimitation, info.FreecountLimitationExpiredTimestamp = proto.Int32(freeCount), proto.Int32(freeCountLimitation), proto.Int64(expiredTimestamp)

	api.resetLottoInfoSlots(info)

	return
}

func (api *XYAPI) defaultLottoTransaction(uid string, adType battery.DrawAwardType, stuff *battery.LottoStuff, now int64) (lottoTransaction *battery.LottoTransaction) {

	lottoTransaction = &battery.LottoTransaction{
		Uid:           proto.String(uid),
		Lottoid:       stuff.Lottoid,
		Parentlottoid: stuff.Parentlottoid,
		Type:          adType.Enum(),
		Slots:         stuff.Slots,
		Selected:      stuff.Selected,
		State:         battery.LottoState_LottoState_Initial.Enum(),
		Timestamp:     proto.Int64(now),
	}

	stateEntry := &battery.LottoStateEntry{
		State:  battery.LottoState_LottoState_Initial.Enum(),
		Opdate: proto.String(xyutil.CurTimeStr()),
	}

	lottoTransaction.States = append(lottoTransaction.States, stateEntry)

	return
}

func (api *XYAPI) defaultLottoLog(uid string, cmd battery.LottoCmd, value int32, adType battery.DrawAwardType, stuff *battery.LottoStuff, err *battery.Error) (lottoLog *battery.LottoLog) {
	lottoLog = &battery.LottoLog{
		Uid:           proto.String(uid),
		Lottoid:       stuff.Lottoid,
		Parentlottoid: stuff.Parentlottoid,
		Type:          adType.Enum(),
		Slots:         stuff.Slots,
		Selected:      stuff.Selected,
		Value:         proto.Int32(value),
		Opdate:        proto.String(xyutil.CurTimeStr()),
		Timestamp:     proto.Int64(time.Now().Unix()),
		Cmd:           cmd.Enum(),
		Error:         err,
	}

	return
}

//游戏后抽奖（初始）
func (api *XYAPI) initialAfterGameLotto(req *battery.LottoRequest, resp *battery.LottoResponse) (failReason battery.ErrorCode, err error) {
	var (
		uid     = req.GetUid()
		quota   = req.GetQuota()
		lottoid = DefLottoIdGenerater.NewID()
		stuff   = &battery.LottoStuff{
			Lottoid:       proto.Uint64(lottoid),
			Parentlottoid: proto.Uint64(lottoid), //初始抽奖内容的lottoid和parentlottoid相同
		}
		stage      uint32
		quotaId    battery.QuotaEnum
		quotaValue uint64
		errStr     string
		cmd        = req.GetCmd()
		lt         *battery.LottoTransaction
	)

	if nil != quota {
		quotaId = quota.GetId()
		quotaValue = quota.GetValue()
	} else {
		errStr = fmt.Sprintf("[%s] quota is nil", uid)
		xylog.Error(uid, errStr)
		err = xyerror.ErrBadInputData
		failReason = xyerror.Resp_BadInputData.GetCode()
		resp.Error = xyerror.Resp_BadInputData
		goto ErrHandle
	}

	//获取游戏阶段号
	failReason, err, stage = api.getAfterGameStage(uid, quotaId, quotaValue)
	if err != xyerror.ErrOK {
		errStr = fmt.Sprintf("[%s] quotaId (%d) quotaValue (%d) getAfterGameStage failed", uid)
		xylog.Error(uid, errStr)
		failReason = xyerror.Resp_QueryAfterGameStageError.GetCode()
		resp.Error = xyerror.Resp_QueryAfterGameStageError
		goto ErrHandle
	}

	//获取格子奖品内容
	failReason, err = api.constructLottoSlots(uid, stage, battery.DrawAwardType_DrawAwardType_GameFinish, &(stuff.Slots))
	if err != nil {
		resp.Error = xyerror.Resp_QueryAfterGameStageError
		return
	}

	//获取抽奖结果
	failReason, err = api.getAfterGameSelectedSlot(uid, stage, stuff.Selected, stuff.GetSlots(), false)
	if err != nil {
		return
	}

	resp.Stuff = stuff

	//记录抽奖事务
	lt = api.defaultLottoTransaction(uid, battery.DrawAwardType_DrawAwardType_GameFinish, stuff, time.Now().Unix())
	err = api.addLottoTransaction(lt)
	if err != nil {
		xylog.Error(uid, "lottoTransactionLog failed : %v", err)
		return
	}

ErrHandle:
	//记录日志
	l := api.defaultLottoLog(uid, cmd, 0, battery.DrawAwardType_DrawAwardType_GameFinish, stuff, resp.Error)
	go api.addLottoLog(l)
	return
}

//游戏后抽奖（删除抽奖格子）
func (api *XYAPI) noInitialAfterGameLotto(req *battery.LottoRequest, resp *battery.LottoResponse) (failReason battery.ErrorCode, err error) {
	var (
		uid               = req.GetUid()
		cmd               = req.GetCmd()
		lottoId           uint64
		lastLottoId       = req.GetLottoid()
		parentLottoId     = req.GetParentlottoid()
		lottoTransactions = make([]*battery.LottoTransaction, 0)
		lastTransaction   *battery.LottoTransaction
		reLottoCount      int
		errStr            string
		stuff             *battery.LottoStuff
		lottoTransaction  *battery.LottoTransaction
		kickedSlotId      = req.GetKickedSlotId()
		quota             = req.GetQuota()
	)

	//校验，根据parentLottoId获取transaction
	err = api.queryLottoTransactionsByParentLottoId(uid, parentLottoId, &lottoTransactions)
	if err != nil {
		errStr = fmt.Sprintf("[%s] queryLottoTransactionsByParentLottoId(%d) failed.", uid, parentLottoId)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_QueryLottoTransactionError
		failReason = xyerror.Resp_QueryLottoTransactionError.GetCode()
		goto ErrHandle
	}

	//校验事务们是否合法
	if !api.checkLottoTransactions(lastLottoId, parentLottoId, lottoTransactions, &lastTransaction) {
		errStr = fmt.Sprintf("[%s] checkLottoTransactions(%v) failed.", uid, lottoTransactions)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_QueryLottoTransactionError
		failReason = xyerror.Resp_QueryLottoTransactionError.GetCode()
		err = xyerror.ErrQueryPropsFromDBError
		goto ErrHandle
	}

	//校验是否超过了删除格子次数的阈值
	reLottoCount = len(lottoTransactions)
	if reLottoCount >= DefConfigCache.Configs().AfterGameLottoDeleteSlotLimit {
		errStr = fmt.Sprintf("[%s] AfterGameLottoDeleteSlot(%d)times >= AfterGameLottoDeleteSlotLimit(%d) failed.", uid, reLottoCount, DefConfigCache.Configs().AfterGameLottoDeleteSlotLimit)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_AfterGameNotEnoughDeleteSlotChanceError
		failReason = xyerror.Resp_AfterGameNotEnoughDeleteSlotChanceError.GetCode()
		err = xyerror.ErrAfterGameNotEnoughDeleteSlotChanceError
		goto ErrHandle
	}

	//查找删除次数对应的商品id
	if goodId, ok := DefConfigCache.Master().AfterGameLottoDeleteCount2MallItem[reLottoCount]; ok {
		api.BuyGoods(uid, "", goodId, resp.Error)
		if resp.Error.GetCode() != xyerror.Resp_NoError.GetCode() {
			failReason = resp.Error.GetCode()
			err = xyerror.ErrBuyGoodsError
			goto ErrHandle
		}
	}

	//生成新的抽奖id
	lottoId = DefLottoIdGenerater.NewID()

	stuff = &battery.LottoStuff{
		Lottoid:       proto.Uint64(lottoId),
		Parentlottoid: proto.Uint64(parentLottoId),
	}

	stuff.Slots = lastTransaction.GetSlots()
	//如果格子号非法或者已经删除了，则报错+返回
	if kickedSlotId >= uint32(len(stuff.Slots)) || false == stuff.Slots[kickedSlotId].GetValid() {
		errStr = fmt.Sprintf("[%s] kickedSlotId is no correct", uid, kickedSlotId)
		failReason = xyerror.Resp_BadInputData.GetCode()
		err = xyerror.ErrBadInputData
		resp.Error = xyerror.Resp_BadInputData
		//resp.Error.Desc = proto.String(errStr)
		goto ErrHandle
	}

	//删除格子
	stuff.Slots[kickedSlotId].Valid = proto.Bool(false)

	//获取抽奖结果，游戏后抽奖的权重值是通过剩余格子数来获取格子权重列表
	if nil != quota {
		quotaId := quota.GetId()
		quotaValue := quota.GetValue()
		var stage uint32
		failReason, err, stage = api.getAfterGameStage(uid, quotaId, quotaValue)
		if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
			errStr = fmt.Sprintf("[%s] quotaId(%d) quotaValue(%d) getAfterGameStage failed.", uid, quotaId, quotaValue)
			resp.Error = xyerror.Resp_QueryAfterGameStageError
			goto ErrHandle
		}

		failReason, err = api.getAfterGameSelectedSlot(uid, stage, stuff.Selected, stuff.GetSlots(), true)
		if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
			errStr = fmt.Sprintf("[%s] getAfterGameSelectedSlot for stage(%d) failed.", uid, stage)
			resp.Error = xyerror.Resp_GetSelectedSlotError
			goto ErrHandle
		}
	}

	resp.Stuff = stuff

	//记录抽奖记录，设置状态为initial
	lottoTransaction = api.defaultLottoTransaction(uid, battery.DrawAwardType_DrawAwardType_GameFinish, stuff, time.Now().Unix())

	err = api.addLottoTransaction(lottoTransaction)
	if err != nil {
		xylog.Error(uid, "lottoTransactionLog failed : %v", err)
		return
	}

ErrHandle:

	//记录日志
	l := api.defaultLottoLog(uid, cmd, 0, battery.DrawAwardType_DrawAwardType_GameFinish, stuff, resp.Error)
	go api.addLottoLog(l)

	return
}

//提交抽奖结果
func (api *XYAPI) commitLotto(req *battery.LottoRequest, resp *battery.LottoResponse) (failReason battery.ErrorCode, err error) {

	var (
		uid           = req.GetUid()
		lottoid       = req.GetLottoid()
		cmd           = req.GetCmd()
		parentlottoid = req.GetParentlottoid()
		errStr        string
		lottoType     battery.DrawAwardType
		lottoValue    int32
		deductValue   int32
	)

	//
	switch cmd {
	case battery.LottoCmd_Lotto_Commit:
		lottoType = battery.DrawAwardType_DrawAwardType_System
	case battery.LottoCmd_Lotto_AfterGame_Commit:
		lottoType = battery.DrawAwardType_DrawAwardType_GameFinish
	default:
		errStr = fmt.Sprintf("[%s] Unkown lotto commit cmd (%d)", uid, cmd)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_BadInputData
		return
	}

	//查询lottotransaction
	lt := &battery.LottoTransaction{}
	err = api.getLottoTransaction(uid, lottoid, parentlottoid, lt)
	if err != nil {
		errStr = fmt.Sprintf("[%s] GetLottoTransaction lottoid %d parentlottoid %d failed : %v", uid, lottoid, parentlottoid, err)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_QueryLottoTransactionError
		return
	}

	//校验是否提交状态是否正确
	if len(lt.States) != 1 || lt.States[0].GetState() != battery.LottoState_LottoState_Initial {
		errStr = fmt.Sprintf("[%s] lotto[%d] state is wrong, can't commit", uid, lottoid)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_BadInputData
		return
	}

	//追加抽奖事务commit状态信息
	state := battery.LottoState_LottoState_Commit
	stateEntry := &battery.LottoStateEntry{
		State:  &state,
		Opdate: proto.String(xyutil.CurTimeStr()),
	}
	err = api.pushLottoTransactionState(lt, stateEntry)
	if err != nil {
		errStr = fmt.Sprintf("[%s] pushLottoTransactionState failed : %v", uid, err)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_PushLottoTransactionStateError
		return
	}

	//修改lottoinfo用户权值
	//查询抽中道具的权值
	selected := lt.GetSelected()
	if selected >= DefConfigCache.Configs().LottoSlotCount || len(lt.Slots) != int(DefConfigCache.Configs().LottoSlotCount) {
		errStr = fmt.Sprintf("[%s] selected slotid %s is out of range len(Slots) %d", uid, selected, len(lt.Slots))
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_BadInputData
		return
	}

	propId := lt.Slots[selected].Items.GetId()

	propStruct := xycache.DefPropCacheManager.Prop(propId)
	if nil == propStruct {
		xylog.Error(uid, "get propStruct of %d failed", propId)
		return
	}
	xylog.Debug(uid, "PropStruct : %v", *propStruct)

	xylog.Debug(uid, "SubstractUserValue %d", propStruct.LottoValue)

	slots := lt.GetSlots()
	slots[selected].Valid = proto.Bool(false)

	//如果是系统抽奖，刷新一下玩家系统抽奖信息
	if lottoType == battery.DrawAwardType_DrawAwardType_System {

		api.calculateLottoValue(int32(propStruct.LottoValue), &lottoValue, &deductValue)
		sysLottoInfo := &battery.SysLottoInfo{}
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).QuerySysLottoInfo(uid, sysLottoInfo)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "QuerySysLottoInfo failed : %v", err)
			return
		}

		//系统抽奖次数减一
		freeCount, gainFreeCount := sysLottoInfo.GetFreeCount(), sysLottoInfo.GetGainFreeCount()
		if freeCount > 0 { //优先扣系统赠送抽奖次数
			sysLottoInfo.FreeCount = proto.Int32(freeCount - 1)
		} else if gainFreeCount > 0 { //扣背包中的抽奖次数
			sysLottoInfo.GainFreeCount = proto.Int32(gainFreeCount - 1)
		}

		//玩家的内部价值重新赋值
		sysLottoInfo.Value = proto.Int32(sysLottoInfo.GetValue() - lottoValue)

		//重新设置一下奖池信息
		sysLottoInfo.Slots = slots

		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).UpdateSysLottoInfo(uid, sysLottoInfo)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "SubstractUserValue -%d failed", propStruct.LottoValue)
			return
		}
	}

	//发放奖品
	for _, item := range propStruct.Items {
		err = api.GainProp(uid, nil, item, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "GainProp failed : %v", err)
			return
		}
	}

	//修改抽奖done状态信息
	state = battery.LottoState_LottoState_Done
	stateEntry = &battery.LottoStateEntry{
		State:  &state,
		Opdate: proto.String(xyutil.CurTimeStr()),
	}

	err = api.pushLottoTransactionState(lt, stateEntry)
	if err != nil {
		errStr = fmt.Sprintf("[%s] pushLottoTransactionState failed : %v", uid, err)
		xylog.Error(uid, errStr)
		resp.Error = xyerror.Resp_PushLottoTransactionStateError
		return
	}

	//记录日志
	l := &battery.LottoLog{
		Uid:           &uid,
		Lottoid:       lt.Lottoid,
		Parentlottoid: lt.Parentlottoid,
		Type:          lottoType.Enum(),
		Slots:         lt.Slots,
		Selected:      lt.Selected,
		Value:         proto.Int32(-1 * lottoValue),
		Opdate:        proto.String(xyutil.CurTimeStr()),
		Timestamp:     proto.Int64(time.Now().Unix()),
		Cmd:           &cmd,
		Deduct:        proto.Int32(-1 * deductValue),
		Error:         resp.Error,
	}

	go api.addLottoLog(l)

	//刷新系统抽奖相关任务
	quotas := []*battery.Quota{&battery.Quota{Id: battery.QuotaEnum_Quota_Lotto.Enum(), Value: proto.Uint64(1)}}
	missionTypes := []battery.MissionType{battery.MissionType_MissionType_Study, battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_MainLine}
	now := time.Now().Unix()
	err = api.updateUserMissionsQuotas(uid, missionTypes, quotas, now, MissionQuotasNoNeedFinish)

	return
}

//获取系统抽奖内容
// uid string 玩家id
// info **battery.SysLottoInfo 系统抽奖内容
// restFreshSec *int64 剩余刷新时间
// now int64 当前刷新时间戳
// forceRefreshSlot bool 是否强制刷新（强制刷新是玩家通过花钻石刷新奖池内容）
func (api *XYAPI) getSysLottoInfo(uid string, info **battery.SysLottoInfo, restFreshSec *int64, now int64, forceRefreshSlot bool) (failReason battery.ErrorCode, err error) {

	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).QuerySysLottoInfo(uid, *info)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			//没有找到，则增加一个
			xylog.Error(uid, "GetSysLottoInfo failed : %v， something wrong", err)
			*restFreshSec = int64(DefConfigCache.Configs().SysLottoRefreshTime)
			err = api.addNewSysLottoInfo(uid, info, now)
			if err != xyerror.ErrOK {
				xylog.Error(uid, "addNewSysLottoInfo failed : %v", err)
				return
			}
		} else {
			return
		}
	} else {
		timestamp := (*info).GetTimestamp()
		xylog.Debug(uid, "now %d , timestamp %d", now, timestamp)
		if now > timestamp+DefConfigCache.Configs().SysLottoRefreshTime || forceRefreshSlot { //抽奖内容过期或者购买刷新奖池，重置slots
			//重新构造抽奖内容
			api.resetLottoInfoSlots(*info)
			(*info).Timestamp = proto.Int64(now)
			*restFreshSec = int64(DefConfigCache.Configs().SysLottoRefreshTime) //倒计时重置
		} else {
			*restFreshSec = int64(timestamp + DefConfigCache.Configs().SysLottoRefreshTime - now) //倒计时重置为剩余时间
		}

		freecountRefreshTimestamp := (*info).GetFreecountRefreshTimestamp()

		//修复12点后在抽奖内容刷新前，可以无限抽奖的bug。
		//以下对修复bug前的玩家数据做兼容。
		if freecountRefreshTimestamp == 0 {
			(*info).FreecountRefreshTimestamp = proto.Int64(now)
		}

		//如果跨天需要重置freecount
		xylog.Debug(uid, "freecountRefreshTimestamp : %d, now : %d", freecountRefreshTimestamp, now)
		if xyutil.DayDiff(freecountRefreshTimestamp, now) > 0 {
			xylog.Debug(uid, "FreecountRefresh refresh")
			freeCount := api.getSysLottoFreeCount(uid, *info)
			(*info).FreeCount = proto.Int32(freeCount)
			(*info).FreecountRefreshTimestamp = proto.Int64(now)
		}

		failReason, err = api.constructLottoSlots(uid, 0, battery.DrawAwardType_DrawAwardType_System, &((*info).Slots))
		if err != xyerror.ErrOK {
			return
		}

		//刷新lottoinfo
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).UpdateSysLottoInfo(uid, *info)
		if err != nil {
			return
		}
	}

	xylog.Debug(uid, "restFreshSec : %d, syslottoinfo : %v", *restFreshSec, *info)

	return

}

//获取系统抽奖内容
// uid string 玩家id
// info **battery.SysLottoInfo 系统抽奖内容
// restFreshSec *int64 剩余刷新时间
// now int64 当前刷新时间戳
// forceRefreshSlot bool 是否强制刷新（强制刷新是玩家通过花钻石刷新奖池内容）
func (api *XYAPI) getSysSerialLottoInfo(uid string, serialNum int32, info **battery.SysLottoInfo, now int64, selected *uint32) (failReason battery.ErrorCode, err error) {

	//获取序号对应的抽奖信息
	var serialNumSlots *battery.LottoSerialNumSlot
	failReason, serialNumSlots = xycache.DefLottoCacheManager.SpecificSysSerialNumSlots(serialNum)
	if failReason != battery.ErrorCode_NoError {
		xylog.Error(uid, "get SpecificSysSerialNumSlots(%d) failed : %d", serialNum, failReason)
		return
	}

	xylog.Debug(uid, "get SpecificSysSerialNumSlots(%d)  : %v", serialNum, serialNumSlots)

	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).QuerySysLottoInfo(uid, *info)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			//没有找到，则增加一个
			xylog.Error(uid, "GetSysLottoInfo failed : %v， something wrong", err)
			err = api.addNewSysLottoInfo(uid, info, now)
			if err != xyerror.ErrOK {
				xylog.Error(uid, "addNewSysLottoInfo failed : %v", err)
				return
			}
		} else {
			return
		}
	} else {
		//特殊抽奖，不需要刷新奖池信息刷新事件
		timestamp := (*info).GetTimestamp()
		//如果跨天需要重置freecount
		if xyutil.DayDiff(timestamp, now) > 0 {
			freeCount := api.getSysLottoFreeCount(uid, *info)
			(*info).FreeCount = proto.Int32(freeCount)
			(*info).FreecountRefreshTimestamp = proto.Int64(now)
		}

		//刷新lottoinfo
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).UpdateSysLottoInfo(uid, *info)
		if err != nil {
			return
		}
	}

	//填充抽奖格子内容
	props := serialNumSlots.GetPropList()
	slots := make([]*battery.Slot, 0)
	for i, prop := range props {
		slot := &battery.Slot{
			Slotid: proto.Uint32(uint32(i)),
			Valid:  proto.Bool(true),
			Items: &battery.PropItem{
				Id: proto.Uint64(prop),
			},
		}
		slots = append(slots, slot)
	}

	(*info).Slots = slots

	//设置抽奖选中格子号
	*selected = serialNumSlots.GetSelected()

	return
}

//增加新的系统抽奖信息
// info **battery.SysLottoInfo 新的系统抽奖信息内存指针
func (api *XYAPI) addNewSysLottoInfo(uid string, info **battery.SysLottoInfo, now int64) (err error) {

	//begin := time.Now()
	//defer xyperf.Trace(LOGTRACE_ADDSYSLOTTOINFO, &begin)

	//生成默认
	*info = api.defaultSysLottoInfo(uid, now)
	//构建抽奖内容
	_, err = api.constructLottoSlots(uid, 0, battery.DrawAwardType_DrawAwardType_System, &((*info).Slots))
	if err != xyerror.ErrOK {
		return
	}

	//新的抽奖内容入库
	xylog.Debug(uid, "New SysLottoInfo : %v", *info)
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).AddLottoInfo(*info)
	return
}

//获取系统抽奖免费抽奖次数
//如果玩家拥有了免费抽奖上限符文，则按照符文对应的免费抽奖次数设置freecount
func (api *XYAPI) getDefaultSysLottoFreeCount(uid string) (freeCount, freeCountLimitation int32, expiredTimestamp int64) {
	freeCount = DefConfigCache.Configs().SysLottoFreeCount //系统默认的免费抽奖次数

	//如果拥有抽奖上限符文，则设置当前的呃抽奖免费次数为
	userRune, err := api.queryUserRune(uid, RUNE_LOTTO_ADDITIONAL)
	if err == xyerror.ErrOK {
		freeCountLimitation = api.RuneConfigValue(RUNE_LOTTO_ADDITIONAL)
		freeCount = freeCountLimitation
		expiredTimestamp = userRune.GetExpiredTimestamp()
	}

	return
}

// 获取玩家当前的抽奖机会值
// uid string 玩家标识
// sysLotto *battery.SysLottoInfo 玩家的系统抽奖信息
//returns:
// freeCount int32 系统默认的免费抽奖次数
func (api *XYAPI) getSysLottoFreeCount(uid string, sysLotto *battery.SysLottoInfo) (freeCount int32) {
	if sysLotto.GetFreecountLimitationExpiredTimestamp() > 0 { //如果加成有效，则设置当前值为加成后的值
		freeCount = sysLotto.GetFreecountLimitation()
	} else {
		freeCount = DefConfigCache.Configs().SysLottoFreeCount //系统默认的免费抽奖次数
	}
	return
}

//重置抽奖格子
func (api *XYAPI) resetLottoInfoSlots(info *battery.SysLottoInfo) {
	info.Slots = make([]*battery.Slot, 0)
	var i uint32 = 0
	for ; i < DefConfigCache.Configs().LottoSlotCount; i++ {
		slot := new(battery.Slot)
		slot.Slotid = proto.Uint32(i)
		slot.Valid = proto.Bool(false)
		info.Slots = append(info.Slots, slot)
	}
}

//构建抽奖内容，每个格子填一个奖品
// uid string 玩家id
// stage uint32 游戏阶段(游戏后抽奖)
// dAType battery.DrawAwardType 抽奖类型
// slots *[]*battery.Slot 格子列表
func (api *XYAPI) constructLottoSlots(uid string, stage uint32, dAType battery.DrawAwardType, slots *[]*battery.Slot) (failReason battery.ErrorCode, err error) {

	setNeedConstruct := make(Set, 0)

	//找出需要构造的格子
	for i, s := range *slots {
		if nil == s {
			xylog.Error(uid, "Slots[%d] is nil", i)
			continue
		}

		if nil == s.Valid {
			xylog.Error(uid, "Slots[%d].Valid is nil", i)
			continue
		}

		if *s.Valid { //s.Valid == true 表示对应的格子非空
			continue
		}
		setNeedConstruct[uint32(i)] = empty{}
	}

	//没有需要构建的格子，直接返回
	if len(setNeedConstruct) <= 0 {
		xylog.Debug(uid, "No need to construct any slot")
		return
	}
	xylog.Debug(uid, "slots need to construct : %v", setNeedConstruct)

	var (
		mapTmpSlots = make(map[uint32]uint64, 0)
	)

	switch dAType {
	case battery.DrawAwardType_DrawAwardType_System:
		//获取系统抽奖的奖池
		failReason, err = api.getSysSlots(uid, setNeedConstruct, mapTmpSlots)
		if err != xyerror.ErrOK {
			return
		}

	case battery.DrawAwardType_DrawAwardType_GameFinish:
		//获取游戏后抽奖的奖池
		failReason, err = api.getAfterGameSlots(uid, stage, setNeedConstruct, mapTmpSlots)
		if err != xyerror.ErrOK {
			return
		}

	default:
		xylog.Error(uid, "wrong DrawAwardType %d", dAType)
		err = errors.New(fmt.Sprintf("wrong DrawAwardType %d", dAType))
		failReason = xyerror.Resp_UnkownLottoTypeError.GetCode()
		return
	}

	//填充格子
	for k, v := range mapTmpSlots {
		(*slots)[k].Valid = proto.Bool(true)
		item := &battery.PropItem{}
		item.Id = proto.Uint64(v)
		(*slots)[k].Items = item
	}

	return
}

//获取系统抽奖结果
// uid string 玩家id
// uvalue int32 玩家的内部价值
// slotid **uint32 抽奖结果指针
func (api *XYAPI) getSysSelectedSlot(uid string, uvalue int32, slotid *uint32) (failReason battery.ErrorCode, err error) {
	var (
		weights      = xycache.DefLottoCacheManager.SysWeights()
		weightStruct *xycache.MAPWeightStruct
		isFound      = false
	)

	for _, w := range weights.S {
		if uvalue < int32(w) {
			weightStruct = weights.M[int32(w)]
			isFound = true
			break
		}
	}
	//记下权重日志
	xylog.Debug(uid, "uvalue %d get weight list : %v ", uvalue, weightStruct)

	if isFound {
		failReason, err = api.getSelectedSlot(uid, weightStruct, slotid)
	} else { //没找到内部价值对应的权重列表
		xylog.Error(uid, "get sys slot weightlist failed for value %d", uvalue)
		failReason = xyerror.Resp_GetSelectedSlotError.GetCode()
		err = xyerror.ErrGetSelectedSlotError
		return
	}

	return
}

//获取系统抽奖格子信息
// uid string 玩家id
// setNeedConstruct Set 需要构建的格子列表
// mapTmpSlots map[uint32]uint64 重新构建的格子奖品列表
func (api *XYAPI) getSysSlots(uid string, setNeedConstruct Set, mapTmpSlots map[uint32]uint64) (failReason battery.ErrorCode, err error) {
	m := xycache.DefLottoCacheManager.SysSlots()
	for slotid, slotMapStruct := range *m {
		//非空格子，跳过
		if _, ok := setNeedConstruct[slotid-1]; !ok {
			xylog.Debug(uid, "getSysSlots : No Need to construct slot %d", slotid-1)
			continue
		}
		randNum := rand.Intn(slotMapStruct.S[len(slotMapStruct.S)-1])
		isFound := false
		for _, v := range slotMapStruct.S {
			if randNum < v {
				mapTmpSlots[slotid-1] = slotMapStruct.M[uint32(v)]
				isFound = true
				break
			}
		}

		if !isFound {
			xylog.Error(uid, "getSysSlots : weight of randNum %d nofound slot %d", randNum, slotid)
		}
	}

	xylog.Debug(uid, "getSysSlots slots %v", mapTmpSlots)

	return
}

//获取游戏后抽奖结果
// uid string 玩家id
// stage uint32 阶段编码
// slotid **uint32 抽奖结果指针
// slots []*battery.Slot 抽奖格子信息
// reCalcWeight bool 是否需要重算权重(游戏后抽奖删除格子后，需要重算格子权重)
func (api *XYAPI) getAfterGameSelectedSlot(uid string, stage uint32, slotid *uint32, slots []*battery.Slot, reCalcWeight bool) (failReason battery.ErrorCode, err error) {
	var (
		weights      = xycache.DefLottoCacheManager.AfterGameWeights()
		weightStruct *xycache.MAPWeightStruct
		ok           bool
	)

	//没找到阶段对应的权重列表，直接返回错误
	if weightStruct, ok = (*weights)[stage]; !ok {
		xylog.Error(uid, "bad stage %d", stage)
		failReason = xyerror.Resp_GetSelectedSlotError.GetCode()
		err = xyerror.ErrGetSelectedSlotError
		return
	}

	//如果是删除格子的抽奖，需要重新算一下权重
	if reCalcWeight {
		err = api.reCalcWeight(stage, weightStruct, slots)
	}

	//记下权重日志
	xylog.Debug(uid, "stage %d get weight list : %v ", stage, weightStruct)

	failReason, err = api.getSelectedSlot(uid, weightStruct, slotid)

	return
}

//重算格子权重
func (api *XYAPI) reCalcWeight(stage uint32, weightStruct *xycache.MAPWeightStruct, slots []*battery.Slot) (err error) {

	var (
		weightList *xycache.MAPSlot2Weight
		ok         bool
		weightTmp  int
	)

	weightStruct.Clear()

	afterGameWeightOriginal := xycache.DefLottoCacheManager.AfterGameWeightsOriginal()

	if weightList, ok = (*afterGameWeightOriginal)[stage]; !ok {
		err = xyerror.ErrQueryAfterGameStageError
		return
	}

	for slotid, weight := range *weightList {
		if slots[slotid].GetValid() {
			weightTmp += int(weight)
			weightStruct.M[uint32(weightTmp)] = slotid
			weightStruct.S = append(weightStruct.S, weightTmp)
		}
	}

	if len(weightStruct.S) > 0 {
		sort.Sort(weightStruct.S)
	}

	return
}

func (api *XYAPI) getAfterGameSlots(uid string, stage uint32, setNeedConstruct Set, mapTmpSlots map[uint32]uint64) (failReason battery.ErrorCode, err error) {

	m := xycache.DefLottoCacheManager.AfterGameSlots()

	if _, ok := (*m)[stage]; !ok {
		xylog.Error(uid, "[%s] AfterGameSlots for stage %d nofound", uid, stage)
		failReason = xyerror.Resp_QueryAfterGameLottoWeightError.GetCode()
		err = xyerror.ErrQueryAfterGameLottoWeightError
		return
	}

	for slotid, slotMapStruct := range *((*m)[stage]) {
		//非空格子，跳过
		if _, ok := setNeedConstruct[slotid]; !ok {
			xylog.Debug(uid, "getAfterGameSlots : No Need to construct slot %d", slotid)
			continue
		}
		randNum := rand.Intn(slotMapStruct.S[len(slotMapStruct.S)-1])
		isFound := false
		for _, v := range slotMapStruct.S {
			if randNum < v {
				mapTmpSlots[slotid] = slotMapStruct.M[uint32(v)]
				isFound = true
				break
			}
		}

		if !isFound {
			xylog.Error(uid, "AfterGameSlots : weight of randNum %d nofound slot %d", randNum, slotid)
			failReason = xyerror.Resp_QueryAfterGameLottoWeightError.GetCode()
			err = xyerror.ErrQueryAfterGameLottoWeightError
			return
		}
	}

	xylog.Debug(uid, "getAfterGameSlots slots %v", mapTmpSlots)

	return
}

//获取游戏结束对应的阶段
func (api *XYAPI) getAfterGameStage(uid string, quotaId battery.QuotaEnum, quotaValue uint64) (failReason battery.ErrorCode, err error, stage uint32) {
	quotaId2Stages := xycache.DefLottoCacheManager.AfterGameQuotaId2Stages()
	var (
		stageMapStruct *xycache.StageMapStruct
		ok             bool
	)

	if stageMapStruct, ok = (*quotaId2Stages)[quotaId]; !ok {
		xylog.Error(uid, "[%s] AfterGameQuotaId2Stages for quotaId %v error", uid, quotaId)
		failReason = xyerror.Resp_GetSelectedSlotError.GetCode()
		err = xyerror.ErrGetSelectedSlotError
		return
	}

	isFound := false
	for _, v := range stageMapStruct.S {
		if quotaValue < uint64(v) {
			stage = stageMapStruct.M[uint64(v)]
			isFound = true
			break
		}
	}

	if !isFound {
		xylog.Error(uid, "AfterGameQuotaId2Stages for quotaId %v  quotaValue %d error", quotaId, quotaValue)
		failReason = xyerror.Resp_GetSelectedSlotError.GetCode()
		err = xyerror.ErrGetSelectedSlotError
		return
	}

	return
}

//获取抽中的格子号
func (api *XYAPI) getSelectedSlot(uid string, weightStruct *xycache.MAPWeightStruct, slotid *uint32) (failReason battery.ErrorCode, err error) {

	if len(weightStruct.S) <= 0 {
		xylog.Error(uid, "len(weightStruct.S) <= 0")
		failReason = xyerror.Resp_GetSelectedSlotError.GetCode()
		err = xyerror.ErrGetSelectedSlotError
		return
	}

	randNum := rand.Intn(weightStruct.S[len(weightStruct.S)-1])

	xylog.Debug(uid, "syslotto randNum : %d ", randNum)
	selected := false
	for _, v := range weightStruct.S {
		if randNum < v {
			*slotid = weightStruct.M[uint32(v)]
			selected = true
			break
		}
	}

	if !selected { //没找到，杯了个具
		xylog.Error(uid, "selected slot no found, bad randNum %d ", randNum)
		failReason = xyerror.Resp_GetSelectedSlotError.GetCode()
		err = xyerror.ErrGetSelectedSlotError
		return
	}

	xylog.Debug(uid, "syslotto select slot %d", *slotid)

	return
}

//获取抽奖事务信息
func (api *XYAPI) getLottoTransaction(uid string, lottoid, parentlottoid uint64, lt *battery.LottoTransaction) (err error) {
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOTRANSACTION).GetLottoTransaction(uid, lottoid, parentlottoid, lt)
	return
}

//增加抽奖事务信息
func (api *XYAPI) addLottoTransaction(lt *battery.LottoTransaction) (err error) {
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOTRANSACTION).AddLottoTransaction(lt)
	return
}

//增加事务状态信息
func (api *XYAPI) pushLottoTransactionState(lt *battery.LottoTransaction, state *battery.LottoStateEntry) (err error) {
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOTRANSACTION).PushLottoTransactionState(lt, state)
	return

}

//增加抽奖日志
//func (api *XYAPI) addLottoLog(l *battery.LottoLog) (fail_reason int32, err error) {
func (api *XYAPI) addLottoLog(l *battery.LottoLog) (fail_reason int32, err error) {
	if err = api.GetLogDB().AddLottoLog(l); err != xyerror.ErrOK {
		fail_reason = xyerror.LOTTO_DB_ERR_LOTTO_LOG
	}

	return
}

//抽奖后价值和抽水价值计算
// propValue int32 奖品的内部价值
// lottoValue *int32 奖品价值
// deductValue *int32 抽水价值
func (api *XYAPI) calculateLottoValue(propValue int32, lottoValue, deductValue *int32) {

	//抽中奖品的价值N   抽奖消耗的价值M   奖池金额X    返奖比例Y
	//若N - M >= 0 	奖池金额变为X-(N-M)
	//若N - M < 0 	奖池金额变为X-(N-M)*Y

	var diff int32 = propValue - DefConfigCache.Configs().LottoCostPerTime
	*deductValue = 0 //抽水值初始化为0
	if diff >= 0 {
		*lottoValue = diff
	} else {
		*lottoValue = (diff * (100 - DefConfigCache.Configs().LottoDeduct)) / 100
		*deductValue = (diff * DefConfigCache.Configs().LottoDeduct) / 100
	}

	return
}

func (api *XYAPI) checkLottoTransactions(lastLottoId, parentLottoId uint64, lottoTransactions []*battery.LottoTransaction, lastTransaction **battery.LottoTransaction) bool {

	isFound := false

	for _, lottoTransaction := range lottoTransactions {
		//任务状态校验，必须是initial
		if lottoTransaction.GetState() != battery.LottoState_LottoState_Initial {
			return false
		}

		if lastLottoId == lottoTransaction.GetLottoid() {
			*lastTransaction = lottoTransaction
			isFound = true
		}
	}

	//没找到上次的抽奖id
	if !isFound {
		return false
	}

	return true
}

func (api *XYAPI) queryLottoTransactionsByParentLottoId(uid string, parentLottoId uint64, lottoTransactions *[]*battery.LottoTransaction) (err error) {
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOTRANSACTION).QueryLottoTransactionsByParentLottoId(uid, parentLottoId, lottoTransactions)
	return
}

//将抽奖券更新到玩家账户
// uid string   玩家id
// amount uint32 抽奖券数目
func (api *XYAPI) updateUserDataWithLottoTicket(uid string, amount uint32) error {
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).IncreaseUserLottoTicket(uid, amount)
}

//设置玩家免费抽奖上限
// uid string   玩家id
// freeAmount int32 玩家免费抽奖次数
func (api *XYAPI) updateUserDataWithSysLottoFreeCount(uid string, freeAmount int32) error {
	sysLottoInfo := battery.SysLottoInfo{
		Uid:       proto.String(uid),
		FreeCount: proto.Int32(freeAmount),
	}
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).UpdateSysLottoInfo(uid, &sysLottoInfo)
}

//设置玩家免费抽奖上限加成信息
// uid string   玩家id
// freeCountAddtional int32 玩家免费抽奖次数加成
// expiredTimestamp int64 过期时间戳
//func (api *XYAPI) updateUserDataWithSysLottoFreeCountAdditional(uid string, freeCountAdditional int32, expiredTimestamp int64) error {
func (api *XYAPI) updateUserDataWithSysLottoFreeCountLimitation(uid string, freeCountLimitation int32, expiredTimestamp int64) error {
	sysLottoInfo := battery.SysLottoInfo{
		Uid:                                 proto.String(uid),
		FreecountLimitation:                 proto.Int32(freeCountLimitation),
		FreecountLimitationExpiredTimestamp: proto.Int64(expiredTimestamp),
	}
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO).UpdateSysLottoInfo(uid, &sysLottoInfo)
}
