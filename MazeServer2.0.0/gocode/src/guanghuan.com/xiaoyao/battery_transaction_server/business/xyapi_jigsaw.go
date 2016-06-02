// xyapi_jigsaw
package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	xymoney "guanghuan.com/xiaoyao/superbman_server/money"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"time"
)

const (
	JIGSAW_NEED_UNLOCK_PROP           = true      //需要解锁道具
	JIGSAW_NO_NEED_UNLOCK_PROP        = false     //不需要解锁道具
	JIGSAW_FOR_NEWACCOUNT      uint64 = 120040001 //新账户拥有的拼图
)

//拼图相关消息接口
func (api *XYAPI) OperationJigsaw(req *battery.JigsawRequest, resp *battery.JigsawResponse) (err error) {

	var (
		uid       = req.GetUid()
		cmd       = req.GetCmd()
		errStruct = xyerror.DefaultError()
		money     *battery.Money
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//初始化返回消息
	resp.Uid = req.Uid
	resp.Cmd = req.Cmd
	resp.ItemId = req.ItemId

	switch cmd {
	case battery.JigsawCmd_JigsawCmd_Query: //查询玩家拼图信息
		resp.JigsawIdList, err = api.queryUserJigsaw(uid)
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_QueryUserJigsawError.Enum()
			goto ErrHandle
		}
	case battery.JigsawCmd_JigsawCmd_Buy: //玩家购买拼图
		api.buyJigsaw(uid, req.GetItemId(), errStruct)
		if errStruct.GetCode() != battery.ErrorCode_NoError {
			goto ErrHandle
		}
	}

	//碎片更新
	//获取当前账户的碎片数
	money, err = xymoney.Query(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_JIGSAW), battery.MoneyType_chip)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "xymoney.Query for chip failed : %v", err)
		goto ErrHandle
	} else {
		resp.ChipAmount = proto.Uint32(money.GetIapamount() + money.GetOapamount() + money.GetGainamount())
	}

ErrHandle:
	errStruct.Desc = nil //错误描述去掉去掉~
	resp.Error = errStruct

	return err
}

//查询玩家的拼图信息（保存在slice中）
// uid string 玩家id
//return:
// jigsawIdList []int64 玩家的拼图id列表
func (api *XYAPI) queryUserJigsaw(uid string) (jigsawIdList []uint64, err error) {
	jigsawIdList = make([]uint64, 0)
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_JIGSAW).GetJigsawList(uid, &jigsawIdList)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "GetJigsawList failed : %v", err)
		return
	}

	return
}

//查询玩家的拼图信息（保存在Set中）
// uid string 玩家id
//return:
// jigsawIdSet Set 玩家的拼图id列表
func (api *XYAPI) queryUserJigsaw2Set(uid string) (jigsawIdSet Set, err error) {
	var jigsawIdList []uint64
	jigsawIdList, err = api.queryUserJigsaw(uid)
	if err != xyerror.ErrOK {
		return
	}

	if jigsawIdSet == nil {
		jigsawIdSet = make(Set, 0)
	}

	for _, id := range jigsawIdList {
		jigsawIdSet[id] = empty{}
	}

	return
}

//购买拼图
// uid string 玩家id
// jigsawId uint64 拼图id
func (api *XYAPI) buyJigsaw(uid string, jigsawId uint64, errStruct *battery.Error) {
	var (
		jigsawIdSet Set
		err         error
	)

	//查询玩家的拼图信息
	jigsawIdSet, err = api.queryUserJigsaw2Set(uid)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_QueryUserJigsawError.Enum()
		return
	}

	//如果拼图已经拥有，报错并返回
	if _, ok := jigsawIdSet[jigsawId]; ok {
		xylog.Error(uid, "jigsaw(%d) already in your hands, don't buy again", jigsawId)
		errStruct.Code = battery.ErrorCode_BuyDuplicateJigsawError.Enum()
		return
	}

	if jigsawId%10000 != 0 { //小拼图
		api.BuyGoods(uid, "", jigsawId, errStruct)
		if errStruct.GetCode() == battery.ErrorCode_NoError { //碎片合成拼图块成功
			//刷新合成XX次拼图块相关任务状态
			quotas := []*battery.Quota{&battery.Quota{Id: battery.QuotaEnum_Quota_SyntheticJigsaw.Enum(), Value: proto.Uint64(1)}}
			missionTypes := []battery.MissionType{battery.MissionType_MissionType_Study, battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_MainLine}
			api.updateUserMissionsQuotas(uid, missionTypes, quotas, time.Now().Unix(), MissionQuotasNoNeedFinish)
		}
	} else { //大拼图
		//查询拼图配置信息
		jigsawConfig := xybusinesscache.DefJigsawConfigCacheManager.JigsawConfig(jigsawId)
		if jigsawConfig == nil {
			xylog.Error(uid, "invalid jigsaw(%d), kidding?", jigsawId)
			errStruct.Code = battery.ErrorCode_InvalidJigsawIdError.Enum()
			return
		}

		price := battery.MoneyItem{
			Type:   battery.MoneyType_chip.Enum(),
			Amount: proto.Uint32(0),
		}
		jigsawList := make([]uint64, 0)
		for _, goodsId := range jigsawConfig.Jigsawidlist {
			//商品是否存在
			var mallItem *battery.MallItem = xybusinesscache.DefGoodsCacheManager.Good(goodsId)
			if nil == mallItem { //只要有一个拼图子项不存在，就报错返回
				errStruct.Code = battery.ErrorCode_QueryGoodsError.Enum()
				xylog.Error(uid, "mallItem not exist : %v", goodsId)
				return
				//} else if _, ok := jigsawIdSet[jigsawId]; !ok { //只买玩家未拥有的拼图
			} else if _, ok := jigsawIdSet[goodsId]; !ok { //只买玩家未拥有的拼图
				*price.Amount += mallItem.GetPrice()[0].GetAmount()
				jigsawList = append(jigsawList, goodsId)
			}
		}

		xylog.Debug(uid, "jigsawList to buy for jigsawId(%d) : %v", jigsawId, jigsawList)

		if len(jigsawList) > 0 {

			// 计算价格，检查玩家是否有足够的货币
			account := new(battery.DBAccount)
			err = api.GetDBAccountDirect(uid, account, mgo.Strong)
			if err != xyerror.DBErrOK {
				errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
				return
			}

			var priceList = make([]*battery.MoneyItem, 0)
			priceList = append(priceList, &price)
			if !api.isUserHasEnoughCurrency(uid, account, priceList) {
				xylog.Error(uid, "wallet %v not enough for price %v", account.GetWallet(), priceList)
				err = xyerror.ErrNotEnoughCurrency
				errStruct.Code = battery.ErrorCode_NotEnoughCurrency.Enum()
				return
			}

			//购买所有的子拼图
			for _, goodsId := range jigsawList {
				api.BuyGoods(uid, "", goodsId, errStruct)
				if errStruct.GetCode() != battery.ErrorCode_NoError { //购买失败，跳出循环
					return
				}
			}
		}
	}

	return
}

//刷新玩家的拼图数据
// uid string 玩家id
// propItem *battery.PropItem 拼图信息
func (api *XYAPI) updateUserDataWithJigsaw(uid string, accountWithFlag *AccountWithFlag, propItem *battery.PropItem, delay bool) (err error) {
	//获取玩家的拼图数据
	var jigsawIdSet Set
	jigsawIdSet, err = api.queryUserJigsaw2Set(uid)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		return
	}

	if nil == jigsawIdSet {
		jigsawIdSet = make(Set, 0)
	}

	jigsawId := propItem.GetId()
	if api.isJigsawNeedResolve(uid, jigsawId, jigsawIdSet) { //如果需要分解，就分解
		err = api.ResolveProp(uid, accountWithFlag, propItem, delay)
	} else {
		//查找拼图所有子项是否已经集齐，如果已经集齐，则子项合并为拼图为大拼图（有角色解锁就解锁角色）
		parentJigsawId := api.getParentJigsawId(jigsawId)
		if parentJigsawId != jigsawId { //子拼图
			maxSubCount := api.GetMaxSubJigsawCount(uid, parentJigsawId)
			if maxSubCount <= 0 {
				err = xyerror.ErrBadInputData
				return
			}

			//子拼图不需要解锁道具
			xylog.Debug(uid, "add jigsaw(%d)", jigsawId)
			err = api.addJigsaw(uid, accountWithFlag, jigsawId, &jigsawIdSet, JIGSAW_NO_NEED_UNLOCK_PROP, delay)
			subJigsawCount := api.getSubJigsawCount(uid, parentJigsawId, jigsawIdSet)
			if subJigsawCount == maxSubCount {
				//集齐了所有拼图子项，增加拼图
				xylog.Debug(uid, "all subJigsaws collected, add parentJigsaw(%d)", parentJigsawId)
				err = api.addJigsaw(uid, accountWithFlag, parentJigsawId, &jigsawIdSet, JIGSAW_NEED_UNLOCK_PROP, delay)
				if err != xyerror.ErrOK {
					return
				}
			}
		} else { //父拼图
			err = api.addJigsaw(uid, accountWithFlag, parentJigsawId, &jigsawIdSet, JIGSAW_NEED_UNLOCK_PROP, delay)
		}
	}

	return
}

//获取父拼图id
// jigsawid uint64 拼图id
//return:
// uint64 父拼图id
func (api *XYAPI) getParentJigsawId(jigsawid uint64) uint64 {
	return jigsawid - (jigsawid % 10000)
}

//获取已经拥有的子拼图数目
// uid string 玩家id
// parentJigsawId uint64 父拼图id
// jigsawIdSet Set 玩家已拥有的拼图集合
//return:
// count int 已拥有的子拼图数目
func (api *XYAPI) getSubJigsawCount(uid string, parentJigsawId uint64, jigsawIdSet Set) (count int) {
	for id, _ := range jigsawIdSet {
		if id.(uint64)/10000*10000 == parentJigsawId {
			count++
		}
	}
	return
}

//增加拼图
// uid string 玩家id
// jigsawId uint64 拼图id
// needUnlockProp bool 是否需要解锁道具
func (api *XYAPI) AddJigsaw(uid string, accountWithFlag *AccountWithFlag, jigsawId uint64, needUnlockProp, delay bool) (err error) {
	//获取玩家当前的拼图列表
	var jigsawIdSet Set
	jigsawIdSet, err = api.queryUserJigsaw2Set(uid)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		return
	}

	//增加拼图
	err = api.addJigsaw(uid, accountWithFlag, jigsawId, &jigsawIdSet, needUnlockProp, delay)

	return
}

//增加拼图
// uid string 玩家id
// jigsawId uint64 拼图id
// jigsawIdSet *Set 玩家已拥有的拼图列表
// needUnlockProp bool 是否需要解锁道具
func (api *XYAPI) addJigsaw(uid string, accountWithFlag *AccountWithFlag, jigsawId uint64, jigsawIdSet *Set, needUnlockProp bool, delay bool) (err error) {

	//拼图id合法性校验
	if nil == xybusinesscache.DefPropCacheManager.Prop(jigsawId) {
		xylog.Error(uid, "invalid jigsawId %d", jigsawId)
		err = xyerror.ErrBadInputData
		return
	}

	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_JIGSAW).AddJigsaw(uid, jigsawId)
	if err == xyerror.ErrOK {
		(*jigsawIdSet)[jigsawId] = empty{}         //添加到玩家已拥有的拼图列表中
		if jigsawId%10000 == 0 && needUnlockProp { //之前不存在该大拼图,进行解锁
			xylog.Debug(uid, "addJigsaw : %d", jigsawId)
			var jigsawConfig *battery.JigsawConfig
			jigsawConfig, err = api.JigsawConfigDetail(jigsawId)
			if err == xyerror.ErrOK && jigsawConfig != nil {
				propItems := jigsawConfig.GetUnlockprops()
				for _, prop := range propItems {
					propID := prop.GetId()
					if propID/10000000 == uint64(battery.PropType_PROP_ROLE) { //如果是解锁角色，则玩家获取角色
						xylog.Debug(uid, "unlock role : %v", propID)
						err = api.updateUserDataWithRole(uid, accountWithFlag, propID, delay, UNLOCK_ROLE_NO_NEED_ADD_JIGSAW) //解锁角色
					} else { //如果是获取
						xylog.Debug(uid, "unlock prop : %v", prop)
						err = api.GainProp(uid, accountWithFlag, prop, ACCOUNT_UPDATE_DELAY, battery.MoneySubType_gain)
					}
				}
			}
		}
	}
	return
}

//玩家是否拥有拼图
// uid string 玩家id
// jigsawId uint64 拼图id
//return:
// isExisting bool true 玩家拥有拼图，false 玩家未拥有拼图
func (api *XYAPI) isJigsawExisting(uid string, jigsawId uint64) (isExisting bool, err error) {
	isExisting, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_JIGSAW).IsJigsawExisting(uid, jigsawId)
	return
}

//拼图是否需要分解
// uid string 玩家id
// jigsawId uint64 拼图id
//return:
// bool true 需要分解，false 不需要分解
func (api *XYAPI) isJigsawNeedResolve(uid string, jigsawId uint64, jigsawIdSet Set) bool {

	if len(jigsawIdSet) <= 0 { //未拥有任何拼图，不需要分解
		return false
	}

	//父拼图是否存在
	parentJigsawId := api.getParentJigsawId(jigsawId)
	if parentJigsawId != 0 {
		if _, ok := jigsawIdSet[parentJigsawId]; ok { //已经拥有父拼图
			xylog.Debug(uid, "already own parentJigsawId(%d), jigsawId(%d) will resolve", parentJigsawId, jigsawId)
			return true
		}
	}

	//子拼图是否存在
	if _, ok := jigsawIdSet[jigsawId]; ok { //已经拥有子拼图
		xylog.Debug(uid, "already own jigsawId(%d), will resolve", jigsawId)
		return true
	}

	return false
}

//查询拼图配置信息
func (api *XYAPI) JigsawConfigDetail(jigsawid uint64) (jigsawconfig *battery.JigsawConfig, err error) {
	jigsawconfig = xybusinesscache.DefJigsawConfigCacheManager.JigsawConfig(jigsawid)
	if nil == jigsawconfig { //没找到
		err = xyerror.ErrJigsawConfigsFromCacheError
	}
	return
}

//获取父拼图下的子拼图数目
// parentJigsawId uint64 父拼图数目
//return:
// count int 子拼图数目
func (api *XYAPI) GetMaxSubJigsawCount(uid string, parentJigsawId uint64) (count int) {
	jigsawConfig := xybusinesscache.DefJigsawConfigCacheManager.JigsawConfig(parentJigsawId)
	if jigsawConfig != nil {
		count = len(jigsawConfig.Jigsawidlist)
	} else {
		xylog.Error(uid, " invalid jigsaw(%d)", parentJigsawId)
	}
	return
}
