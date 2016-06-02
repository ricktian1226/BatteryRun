// xyapi_BeforeGameGoods
package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	xymoney "guanghuan.com/xiaoyao/superbman_server/money"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"time"
)

const (
	BEFOREGAME_RANDOM_GOODSID = uint64(160010000)
)

//赛前道具操作消息
func (api *XYAPI) OperationBeforeGameProp(req *battery.BeforeGamePropRequest, resp *battery.BeforeGamePropResponse) (err error) {

	var (
		uid       = req.GetUid()
		errStruct = xyerror.DefaultError()
		errStr    string
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//初始化resp
	resp.Uid = req.Uid
	resp.Cmd = req.Cmd

	switch req.GetCmd() {
	case battery.BeforeGamePropCmd_BeforeGamePropCmd_Query: //打开赛前道具面板查询信息

		var randItemId uint64 = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_CONSUMABLE).GetRandomConsumableID(uid)
		resp.RandItemId = proto.Uint64(randItemId)
		resp.BeforeGamePropList = make([]*battery.Consumable, 0) //玩家已拥有的赛前道具列表
		resp.Goods = make([]*battery.MallItem, 0)                //赛前道具商品列表
		api.queryBeforeGameGoods(uid, &(resp.BeforeGamePropList), &(resp.Goods), errStruct)
		if errStruct.GetCode() != battery.ErrorCode_NoError {
			goto ErrHandle
		}

	case battery.BeforeGamePropCmd_BeforeGamePropCmd_Buy: //购买赛前道具商品

		resp.ItemId = req.ItemId
		goodsId := req.GetItemId()
		api.buyBeforeGameGoods(uid, goodsId, errStruct)
		if errStruct.GetCode() != battery.ErrorCode_NoError {
			//resp.Error = errStruct
			goto ErrHandle
		} else {
			if goodsId == BEFOREGAME_RANDOM_GOODSID {
				randomGood := xybusinesscache.DefBeforeGameWeightCacheManager.RandomGood()
				if randomGood == xybusinesscache.InvalidRandomGood {
					xylog.Error(uid, "randomGood failed")
				} else {
					//查询商品信息
					mallItem := xybusinesscache.DefGoodsCacheManager.Good(randomGood)
					if mallItem == nil { //没找到，报错
						errStr = fmt.Sprintf("[%s] get mallItem for %d failed.", uid, randomGood)
						errStruct.Code = battery.ErrorCode_QueryGoodsError.Enum()
						//errStruct.Desc = proto.String(errStr)
						goto ErrHandle
					}

					//发放道具
					err = api.GainProps(uid, nil, mallItem.GetItems(), ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
					if err != xyerror.ErrOK {
						errStr = fmt.Sprintf("[%s] GainProps failed : %v", uid, err)
						errStruct.Code = battery.ErrorCode_GainPropsError.Enum()
						//errStruct.Desc = proto.String(errStr)
						goto ErrHandle
					}

					//刷新购买随机道具的相关任务状态
					quotas := []*battery.Quota{&battery.Quota{Id: battery.QuotaEnum_Quota_BuyRandomGoods.Enum(), Value: proto.Uint64(1)}}
					missionTypes := []battery.MissionType{battery.MissionType_MissionType_Study, battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_MainLine}
					api.updateUserMissionsQuotas(uid, missionTypes, quotas, time.Now().Unix(), MissionQuotasNoNeedFinish)
				}
				resp.RandItemId = proto.Uint64(randomGood)
			}

			//返回玩家钱包信息
			if err == xyerror.ErrOK {
				resp.Wallet, err = xymoney.QueryWallet(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT))
			}
		}
	case battery.BeforeGamePropCmd_BeforeGamePropCmd_Use: //使用赛前道具
		resp.ItemId = req.ItemId
		id := resp.GetItemId()
		err = api.useConsumable(uid, id)
		if err != xyerror.ErrOK {
			errStr = fmt.Sprintf("[%s] useConsumable %d failed : %v", uid, id, err)
			errStruct.Code = battery.ErrorCode_BuyGoodUpdateUserDataError.Enum()
			//errStruct.Desc = proto.String(errStr)
			goto ErrHandle
		}
	}

ErrHandle:
	if errStruct.GetCode() != battery.ErrorCode_NoError {
		xylog.Error(uid, errStr)
	}

	resp.Error = errStruct
	return
}

//查询赛前道具商城
// uid string 玩家id
// consumables *[]*battery.Consumable 消耗品列表
// goods *[]*battery.MallItem 赛前商品列表
func (api *XYAPI) queryBeforeGameGoods(uid string, consumables *[]*battery.Consumable, goods *[]*battery.MallItem, errStruct *battery.Error) {

	var (
		errStr string
	)

	//查询基础的赛前道具商品列表(目前的实现是前端直接从配置中拉商店数据，不需要后端传回商品列表，以下代码暂时注释，在需要后端传回商品列表的时候重新启用)
	//*goods = api.QueryGoods(uid, battery.MallType_Mall_BeforeGame, battery.MallSubType_MallSubType_BeforeGame_Common)

	//查询玩家的消耗品背包信息
	err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_CONSUMABLE).GetNoRandomConsumables(uid, consumables)
	if err != xyerror.ErrOK {
		errStr = fmt.Sprintf("[%s] GetNoRandomConsumables failed : %v", uid, err)
		xylog.Error(uid, errStr)
		errStruct.Code = battery.ErrorCode_QueryUserConsumableError.Enum()
		//errStruct.Desc = proto.String(errStr)
		return
	}

	return
}

//购买赛前道具
// uid string 玩家id
// goodsId uint64 商品
func (api *XYAPI) buyBeforeGameGoods(uid string, goodsId uint64, errStruct *battery.Error) {
	var (
		errStr string
	)

	//判断是否是赛前道具
	if !api.isBeforeGameGoods(goodsId) {
		errStr = fmt.Sprintf("[%s] %d is not BeforeGameGoods", uid, goodsId)
		errStruct.Code = battery.ErrorCode_BadInputData.Enum()
		//errStruct.Desc = proto.String(errStr)
		xylog.Error(uid, errStr)
		return
	}

	//购买商品
	api.BuyGoods(uid, "", goodsId, errStruct)

	return
}

//判断商品是否合法
// goodsId uint64 商品id
func (api *XYAPI) isBeforeGameGoods(goodsId uint64) bool {
	if _, ok := (*(xybusinesscache.DefBeforeGameWeightCacheManager.GoodsSet()))[goodsId]; ok {
		return true
	}
	return false
}

//判断商品是否合法
// goodsId uint64 商品id
func (api *XYAPI) isBeforeGameRandomGoods(goodsId uint64) bool {
	if _, ok := (*(xybusinesscache.DefBeforeGameWeightCacheManager.RandomGoodsSet()))[goodsId]; ok {
		return true
	}
	return false
}

//使用消耗品道具
// uid string 玩家id
// id uint64 消耗品道具id
func (api *XYAPI) useConsumable(uid string, id uint64) (err error) {
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_CONSUMABLE).DecreaseConsumable(uid, id, 1)
}
