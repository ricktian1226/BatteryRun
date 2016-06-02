package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//appstore receipt items
const (
	RECEIPT_STATUS               = "status"
	RECEIPT                      = "receipt"
	RECEIPT_INAPP                = "in_app"
	RECEIPT_TRACSACTION          = "transaction_id"
	RECEIPT_BUNDLE_ID            = "bundle_id"
	RECEIPT_BID                  = "bid"
	RECEIPT_PRODUCT_ID           = "product_id"
	RECEIPT_ORIGINAL_PUR_DATE_MS = "original_purchase_date_ms"
)

const (
	//电池人bundle id
	SB_BUNDLE_ID = "com.guanghuan.SuperBMan"
)

type empty struct{}
type Set map[interface{}]empty

//func (api *XYAPI) OperationIapOrder(reqData *battery.OrderNumRequest, resp *battery.OrderNumResponse) (err error) {
//	uid := reqData.GetUid()
//	xylog.Debug("[%s] [IapOrder] start", uid)
//	defer xylog.Debug("[%s] [IapOrder] end", uid)

//	xylog.Debug("[%s] Item id: %s", uid, reqData.GetItemId())

//	// 检查商品是否合法
//	_, goodErr := api.GetDB().GetGoodById(reqData.GetItemId())
//	if goodErr != nil {
//		return xyerror.ErrIapGoodNotFound
//	}

//	// 新建一条订单记录
//	var order_id string = api.GetDB().NewId()
//	resp.Uid = proto.String(uid)
//	resp.ItemId = proto.String(reqData.GetItemId())
//	resp.OrderId = proto.String(order_id)
//	return
//}

func (api *XYAPI) OperationIapValidate(req *battery.OrderVerifyRequest, resp *battery.OrderVerifyResponse) (err error) {
	var uid string = req.GetUid()
	//xylog.Debug(uid, "[%s] [IapValidate] start", uid)
	//defer xylog.Debug(uid, "[%s] [IapValidate] done", uid)

	var (
		failReason = battery.ErrorCode_NoError
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	// 查询用户是否存在
	if uid == "" || !api.IsUserExist(uid) {
		failReason = battery.ErrorCode_BadInputData
		goto ErrorHandler
	}

	// 初始化返回数据
	resp.Uid = proto.String(req.GetUid())
	resp.OrderId = proto.String(req.GetOrderId())
	resp.IsSucc = proto.Bool(false)
	resp.DiamondCount = proto.Int32(-1)

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_IapValidate, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrorHandler:
	xylog.Debug(uid, "respData : %v", resp)
	// 添加日志
	var strTids string
	sTids := req.GetTransactionId()
	for i, n := 0, len(sTids); i < n; i++ {
		strTids += sTids[i] + ","
	}

	l := battery.IapLog{
		Uid:           proto.String(uid),
		TransactionId: proto.String(strTids),
		FailReason:    failReason.Enum(),
		OpDate:        proto.Int64(xyutil.CurTimeSec()),
		OpDateStr:     proto.String(xyutil.CurTimeStr()),
	}
	l.IapReceipt = proto.String(resp.GetReceiptDetail())
	if err != nil {
		l.Error = proto.String(err.Error())
	}
	go api.GetLogDB().AddIapLog(l)

	//receiptdetail只是给iaplog用的，客户端信息返回信息中不需要，清空它
	resp.ReceiptDetail = nil

	return
}
