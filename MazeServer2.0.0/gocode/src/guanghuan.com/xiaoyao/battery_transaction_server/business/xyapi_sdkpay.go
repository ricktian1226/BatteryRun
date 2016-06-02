package batteryapi

import (
    // "strconv"

    // "code.google.com/p/goprotobuf/proto"

    // "guanghuan.com/xiaoyao/common/idgenerate"
    "guanghuan.com/xiaoyao/common/log"
    // xyutil "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/money"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

// 订单操作
func (api *XYAPI) OperationSDKOrderOp(req *battery.SDKOrderOperationRequest, resp *battery.SDKOrderOperationResponse) (err error) {
    var (
        uid       string = req.GetUid()
        orderid   string = req.GetOrderId()
        errStruct        = xyerror.DefaultError()
        goodList         = make([]uint64, 0)
    )
    platform := req.GetPlatformType()
    api.SetDB(platform)

    resp.OrderId = req.OrderId
    resp.Uid = req.Uid

    sdkOrder := &battery.SDKOrderInfo{}
    err = api.getSDKOrder(uid, orderid, sdkOrder)
    if err != nil {
        xylog.ErrorNoId("get sdkorder :%v fail", orderid)
        errStruct.Code = battery.ErrorCode_SDKOrderInvalid.Enum()
        goto ErrHandle
    }
    if sdkOrder.GetState() != battery.SDKOrderState_SDKOrderState_Init {
        xylog.Error(uid, "orderid %d state wrong ", orderid)
        errStruct.Code = battery.ErrorCode_SDKOrderInvalid.Enum()
        goto ErrHandle
    }

    err = api.dealOrder(uid, sdkOrder, errStruct)
    if err != nil {
        goto ErrHandle
    }

    goodList = append(goodList, sdkOrder.GetGoodsId())
    resp.GoodsList = goodList

    resp.Wallet, err = xymoney.QueryWallet(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT))
ErrHandle:
    resp.Error = errStruct
    return
}

func (api *XYAPI) dealOrder(uid string, order *battery.SDKOrderInfo, errStruct *battery.Error) (err error) {
    var (
        goodsId uint64
    )
    // 发放商品
    goodsId = order.GetGoodsId()
    moneygoods := xybusinesscache.DefGoodsCacheManager.Good(goodsId)
    if nil == moneygoods {
        xylog.ErrorNoId("mallItem not exist : %v", goodsId)
        errStruct = xyerror.Resp_IapGoodNotFound
        return
    }
    err = api.GainProps(uid, nil, moneygoods.GetItems(), ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_iap)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("gain prop err : %v", goodsId)
        errStruct.Code = battery.ErrorCode_DBError.Enum()
        return
    }

    // 更新订单状态
    order.State = battery.SDKOrderState_SDKOrderState_Done.Enum()
    err = api.updateSDKOrder(uid, order.GetOrderId(), order)
    if err != nil {
        xylog.ErrorNoId("update order state fail:%v", order.GetOrderId())
        errStruct.Code = battery.ErrorCode_SDKUpdateOrderFail.Enum()
    }
    return
}

// 未完成订单查询
func (api *XYAPI) OperationSDKOrderQuery(req *battery.SDKOrderRequest, resp *battery.SDKOrderResponse) (err error) {
    var (
        uid = req.GetUid()
    )
    platform := req.GetPlatformType()
    api.SetDB(platform)

    resp.Error = xyerror.DefaultError()
    resp.Uid = req.Uid
    orders := make([]*battery.SDKOrderInfo, 0)
    err = api.getAllUnfinishedSDKOrder(uid, &orders)
    if err != nil {
        if err == xyerror.ErrNotFound {
            return
        } else {
            resp.Error.Code = battery.ErrorCode_DBError.Enum()
            return
        }
    }
    orderUints := make([]*battery.SDKOrderUint, 0)
    for _, order := range orders {
        orderUint := &battery.SDKOrderUint{
            OrderId: order.OrderId,
            GoodsId: []uint64{order.GetGoodsId()},
        }
        orderUints = append(orderUints, orderUint)
        // 未完成订单处理
        err = api.dealOrder(uid, order, resp.Error)
        if err != nil {
            xylog.ErrorNoId("order error :%v,order:%v", err, orderUint)
            // 有订单出错则退出等待下次请求再次处理
            break
        }
    }
    resp.Order = orderUints

    return
}

// 添加订单
func (api *XYAPI) OperationSDKAddOrder(req *battery.SDKAddOrderRequest, resp *battery.SDKAddOrderResponse) (err error) {

    resp.Error = xyerror.DefaultError()
    resp.Uid = req.Uid

    platform := req.GetPlatformType()
    api.SetDB(platform)

    // 订单防重放校验
    isExist, err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SDKORDER).IsOrderIdExist(req.GetOrderId())
    if isExist {
        resp.Error.Code = battery.ErrorCode_SDKOrderInvalid.Enum()
        return
    }
    sdkOrder := &battery.SDKOrderInfo{}
    sdkOrder.Uid = req.Uid
    sdkOrder.OrderId = req.OrderId
    sdkOrder.GoodsId = req.GoodsId
    sdkOrder.PayTime = req.PayTime
    sdkOrder.Sandbox = req.Sandbox
    sdkOrder.PayAmount = req.PayAmount
    if req.PlatformType == nil {
        sdkOrder.PlatformType = battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID.Enum()
    } else {
        sdkOrder.PlatformType = req.PlatformType
    }
    sdkOrder.State = battery.SDKOrderState_SDKOrderState_Init.Enum()

    xylog.DebugNoId("sdkorder :%v", sdkOrder)
    err = api.addSDKOrder(sdkOrder)
    if err != nil {
        resp.Error.Code = battery.ErrorCode_SDKAddOrderFail.Enum()
    }
    return
}

func (api *XYAPI) addSDKOrder(order *battery.SDKOrderInfo) (err error) {
    err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SDKORDER).AddOrder(order)
    return
}

func (api *XYAPI) getSDKOrder(uid, orderid string, order *battery.SDKOrderInfo) (err error) {
    return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SDKORDER).GetSDKOrder(uid, orderid, order)
}

func (api *XYAPI) getAllUnfinishedSDKOrder(uid string, orders *[]*battery.SDKOrderInfo) (err error) {
    return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SDKORDER).GetAllUnfinishedSDKOrder(uid, orders)
}

func (api *XYAPI) updateSDKOrder(uid, orderid string, order *battery.SDKOrderInfo) (err error) {
    return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SDKORDER).UpdateSDKOrder(uid, orderid, order)
}
