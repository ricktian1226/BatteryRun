package batteryapi

import (
	"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

////////////////////////////////////
// 商品相关操作
// 客户端：查询，购买
////////////////////////////////////

func (api *XYAPI) OperationQueryGoods(req *battery.QueryGoodsRequest, resp *battery.QueryGoodsResponse) (err error) {
	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_QueryGoods, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}
	return
}

func (api *XYAPI) OperationBuyGoods(req *battery.BuyGoodsRequest, resp *battery.BuyGoodsResponse) (err error) {
	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if !api.isUidValid(uid) {
		resp.Error.Code = battery.ErrorCode_BadInputData.Enum()
		xylog.Error(uid, "OperationBeforeGameProp invalid uid")
		err = xyerror.ErrGetAccountByUidError
		return
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_GoodsBuy, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}
	return
}
