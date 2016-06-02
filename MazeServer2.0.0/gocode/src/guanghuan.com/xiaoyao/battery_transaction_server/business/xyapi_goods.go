package batteryapi

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "code.google.com/p/goprotobuf/proto"

    "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/money"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//查询商品信息
func (api *XYAPI) OperationQueryGoods(req *battery.QueryGoodsRequest, resp *battery.QueryGoodsResponse) (err error) {
    var (
        uid         = req.GetUid()
        mallType    = req.GetMallType()
        mallSubType = req.GetMallSubType()
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    goodslist := api.QueryGoods(uid, mallType, mallSubType, platform)
    resp.GoodsList = goodslist
    resp.MallType = req.MallType
    resp.MallSubType = req.MallSubType
    resp.Error = xyerror.DefaultError()
    resp.SystemTime = proto.Int64(xyutil.CurTimeSec())

    return
}

//购买商品消息
func (api *XYAPI) OperationBuyGoods(req *battery.BuyGoodsRequest, resp *battery.BuyGoodsResponse) (err error) {
    var (
        uid     = req.GetUid()
        goodsId = req.GetGoodsId()
        gameId  = req.GetGameId()
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    //初始化resp
    resp.Uid = req.Uid
    resp.GoodsId = req.GoodsId
    resp.GameId = req.GameId
    resp.Error = xyerror.DefaultError()

    //购买商品
    api.BuyGoods(uid, gameId, goodsId, resp.Error)
    if resp.Error.GetCode() != battery.ErrorCode_NoError {
        return
    }

    //返回玩家钱包数据
    resp.Wallet, err = xymoney.QueryWallet(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT))

    return
}

//根据商场类型查询商品列表
// uid string 玩家id
// mallType battery.MallType 商场类型
// mallSubType battery.MallSubType 商品类型
func (api *XYAPI) QueryGoods(uid string, mallType battery.MallType, mallSubType battery.MallSubType, platform battery.PLATFORM_TYPE) (goodsList []*battery.MallItem) {
    goodsList = xybusinesscache.DefGoodsCacheManager.SpecificMall(mallType, mallSubType, platform)
    xylog.Debug(uid, "QueryGoods mallType(%v) mallSubType(%v) result : %d", mallType, mallSubType, len(goodsList))
    return
}

//购买商品
// uid string 玩家id
// gameId string 游戏id（游戏中商城购买需要提供）
// goodsId uint64 商品id
func (api *XYAPI) BuyGoods(uid, gameId string, goodsId uint64, errStruct *battery.Error) {
    var (
        canBuy          = true
        mallItem        *battery.MallItem
        errStr          string
        accountWithFlag *AccountWithFlag
        receiptId       uint64
    )

    //errStruct = new(battery.Error)
    errStruct.Code = battery.ErrorCode_NoError.Enum()

    //查询玩家账户信息
    account := new(battery.DBAccount)
    err := api.GetDBAccountDirect(uid, account, mgo.Strong)
    if err != xyerror.ErrOK {
        errStr = fmt.Sprintf("[%s] GetDBAccount failed: %s", uid, err.Error())
        errStruct.Code = xyerror.Resp_GetAccountByUidError.GetCode().Enum()
        goto ErrHandle
    }

    //查询商品是否允许购买
    canBuy, mallItem, *errStruct = api.isGoodCanBuy(uid, account, goodsId, gameId)
    if !canBuy || xyerror.Resp_NoError.GetCode() != errStruct.GetCode() {
        errStr = errStruct.GetDesc()
        goto ErrHandle
    }

    // 购买该商品
    xylog.Debug(uid, "Ready to buy %v", *mallItem)
    accountWithFlag = &AccountWithFlag{
        account: account,
        bChange: false,
    }
    receiptId, *errStruct = api.TakeGoods(accountWithFlag, mallItem)
    if errStruct.GetCode() != xyerror.Resp_NoError.GetCode() {
        errStr = errStruct.GetDesc()
        goto ErrHandle
    }

    //刷新玩家数据
    err = api.UpdateAccountWithFlag(accountWithFlag)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_UpdateAccountError.Enum()
        errStr = fmt.Sprintf("[%s] UpdateAccountWithFlag failed : %v", uid, err)
        xylog.Error(uid, errStr)
        goto ErrHandle
    }

    //记录交易信息
    err = api.AddShoppingTransaction(uid, goodsId, gameId)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_AddUserShoppingTransactionError.Enum()
        errStr = fmt.Sprintf("[%s] UpdateAccountWithFlag failed : %v", uid, err)
        xylog.Error(uid, errStr)
        goto ErrHandle
    }

ErrHandle:
    //errStruct.Desc = proto.String(errStr)
    go api.AddShoppingLog(uid, goodsId, receiptId, gameId, errStruct.GetCode(), errStr)
    return
}

//判断商品是否可以购买
// uid string 玩家id
// account *battery.DBAccount 玩家账户信息
// goodId uint64 商品id
// gameId string 游戏id
func (api *XYAPI) isGoodCanBuy(uid string, account *battery.DBAccount, goodsId uint64, gameId string) (canBuy bool, mallItem *battery.MallItem, errStruct battery.Error) {
    var (
        errStr                      string
        expiredDate                 int64
        amountPerUser, amountPerDay uint32
    )

    //初始化返回值
    canBuy = true
    errStruct.Code = xyerror.Resp_NoError.GetCode().Enum()

    //商品是否存在
    mallItem = xybusinesscache.DefGoodsCacheManager.Good(goodsId)
    if nil == mallItem {
        errStr = fmt.Sprintf("[%s] goods(%d) doesn't exists", uid, goodsId)
        errStruct.Code = xyerror.Resp_QueryGoodsError.GetCode().Enum()
        xylog.Error(uid, "mallItem not exist : %v", goodsId)
        goto ErrHandle
    }

    //判断商品是否过期
    expiredDate = mallItem.GetExpiretimestamp()
    if expiredDate > 0 { //当expiredDate大于0时才有效
        curTime := xyutil.CurTimeSec()
        if expiredDate <= curTime {
            errStr = fmt.Sprintf("[%s] goods(%d) expired, can't buy", uid, goodsId)
            errStruct.Code = xyerror.Resp_QueryGoodsError.GetCode().Enum()
            xylog.Error(uid, "mallITem expired : %v", goodsId)
            goto ErrHandle
        }
    }

    //如果是游戏中购买，校验一下
    if mallItem.GetMallType() == battery.MallType_Mall_InGame {
        if gameId == "" {
            errStr = fmt.Sprintf("[%s] goods(%d) is only available in game, provide game id nil", uid, goodsId)
            errStruct.Code = xyerror.Resp_BuyGoodInvalidGame.GetCode().Enum()
            goto ErrHandle
        } else if !api.IsGameOnGoing(uid, gameId) { //游戏是否还在进行中
            errStr = fmt.Sprintf("[%s] goods(%d) is only available in game, game (%s) is not a valid game", uid, goodsId, gameId)
            errStruct.Code = xyerror.Resp_BuyGoodInvalidGame.GetCode().Enum()
            goto ErrHandle
        }

        //是否超过每局购买上限
        amountPerGame := mallItem.GetAmountpergame()
        if amountPerGame > 0 {
            curCount, _ := api.GetShoppingCountOfGame(uid, goodsId, gameId)
            xylog.Debug(uid, "goodsid(%d) gameId(%s) amountPerGame(%d) curCount(%d)", goodsId, gameId, amountPerGame, uint32(curCount))
            if amountPerGame <= uint32(curCount) {
                errStr = fmt.Sprintf("[%s] goods(%d) shopping count (%d) is over amountPerGame(%d), gameid(%s)", uid, goodsId, curCount, amountPerGame, gameId)
                errStruct.Code = battery.ErrorCode_BuyGoodOverAmountPerGame.Enum()
                xylog.Error(uid, "%d curCount(%d) over amountPerGame(%d)", goodsId, curCount, amountPerGame)
                goto ErrHandle
            }
        }
    }

    //是否超过玩家购买上限
    amountPerUser = mallItem.GetAmountperuser()
    if amountPerUser > 0 {
        curCount, _ := api.GetShoppingCountOfUser(uid, goodsId)
        if amountPerUser <= uint32(curCount) {
            //errStr = fmt.Sprintf("[%s] goods(%d) shopping count (%d) is over amountPerUser(%d)", uid, goodsId, curCount, amountPerUser)
            errStruct.Code = battery.ErrorCode_BuyGoodOverAmountPerUser.Enum()
            xylog.Error(uid, "%d curCount(%d) over amountPerUser(%d)", goodsId, curCount, amountPerUser)
            goto ErrHandle
        }
    }

    //是否超过玩家当日购买上限
    amountPerDay = mallItem.GetAmountperday()
    if amountPerDay > 0 {
        curCount, _ := api.GetShoppingCountOfDay(uid, goodsId)
        if amountPerDay <= uint32(curCount) {
            //errStr = fmt.Sprintf("[%s] goods(%d) shopping count (%d) is over amountPerUser(%d)", uid, goodsId, curCount, amountPerUser)
            errStruct.Code = battery.ErrorCode_BuyGoodOverAmountPerDay.Enum()
            xylog.Error(uid, "%d curCount(%d) over amountPerDay(%d)", goodsId, curCount, amountPerDay)
            goto ErrHandle
        }
    }

    // 计算价格，检查玩家是否有足够的货币
    if !api.isUserHasEnoughCurrency(uid, account, mallItem.GetPrice()) {
        errStr = fmt.Sprintf("[%s] user doesn't have enough currency(%v) for %v", uid, account.Wallet, mallItem.GetPrice())
        errStruct.Code = xyerror.Resp_NotEnoughCurrency.GetCode().Enum()
        goto ErrHandle
    }

    //没有错误，直接返回
    return

ErrHandle:
    if errStruct.GetCode() != battery.ErrorCode_NoError {
        xylog.Error(uid, errStr)
    }

    canBuy = false
    return
}

//检查玩家是否有足够的资源来购买商品
// uid string 玩家id
// account *battery.DBAccount 账户信息
// price []*battery.MoneyItem 商品价格
//return:
// isEnough bool 是否足够
func (api *XYAPI) isUserHasEnoughCurrency(uid string, account *battery.DBAccount, price []*battery.MoneyItem) (isEnough bool) {
    isEnough = true
    //只要有一个代币不够，就买不了
    for _, v := range price {
        amountHas := xymoney.Amount(uid, v.GetType(), account.GetWallet())
        if amountHas < v.GetAmount() {
            xylog.Warning(uid, "[%s] %v request %d but only %d ,no enough", uid, v.GetType(), v.GetAmount(), amountHas)
            isEnough = false
            return
        }
    }

    return
}

//生成新票据信息
// uid string 玩家id
// goodsId uint64 商品id
// price []*battery.MoneyItem 商品价格
//return:
// receiptId uint64 票据id
func (api *XYAPI) NewReceipt(uid string, goodsId uint64, price []*battery.MoneyItem) (receiptId uint64, err error) {
    receiptId = DefReceiptIdGenerater.NewID()
    receipt := battery.Receipt{
        Id:        proto.Uint64(receiptId),
        Uid:       proto.String(uid),
        GoodsId:   proto.Uint64(goodsId),
        Price:     price,
        Timestamp: proto.Int64(xyutil.CurTimeSec()),
        OpDateStr: proto.String(xyutil.CurTimeStr()),
    }
    //记录票据信息
    err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_RECEIPT).AddReceipt(receipt)
    if err != xyerror.DBErrOK {
        receiptId = 0
    }
    return
}

//将商品发放到玩家背包中
// uid string 玩家id
// goods *battery.MallItem 商品信息
func (api *XYAPI) updateUserDataWithGoods(uid string, accountWithFlag *AccountWithFlag, goods *battery.MallItem, delay bool) (err error) {
    err = api.GainProps(uid, accountWithFlag, goods.GetItems(), delay, battery.MoneySubType_gain)
    return
}

//玩家购买商品
// account *battery.DBAccount 玩家账户信息
// goods *battery.MallItem 商品信息
//return:
// receiptId uint64 购买票据
//func (api *XYAPI) TakeGoods(account *battery.DBAccount, goods *battery.MallItem) (receiptId uint64, errStruct battery.Error) {
func (api *XYAPI) TakeGoods(accountWithFlag *AccountWithFlag, goods *battery.MallItem) (receiptId uint64, errStruct battery.Error) {

    var (
        account = accountWithFlag.account
        uid     = account.GetUid()
        price   = goods.GetPrice()
        errStr  string
        err     error
    )

    errStruct = *(xyerror.Resp_NoError)

    //更新玩家背包
    err = api.updateUserDataWithGoods(uid, accountWithFlag, goods, ACCOUNT_UPDATE_DELAY)
    if err != xyerror.ErrOK {
        errStr = fmt.Sprintf("[%s] updateUserDataWithGoods failed : %v", uid, err.Error())
        errStruct.Code = xyerror.Resp_BuyGoodUpdateUserDataError.GetCode().Enum()
        goto ErrHandle
    }

    //消费代币
    err = xymoney.Consum(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT), account, price, false, ACCOUNT_UPDATE_DELAY)
    if err != xyerror.DBErrOK {
        errStr = fmt.Sprintf("[%s] Consum money failed : %v", uid, err.Error())
        errStruct.Code = xyerror.Resp_BuyGoodConsumMoneyError.GetCode().Enum()
        goto ErrHandle
    }
    accountWithFlag.SetChange()
    xylog.Debug(uid, "Consum %v for goods %d done", price, goods.GetId())

    return

ErrHandle:
    if errStruct.GetCode() != battery.ErrorCode_NoError {
        xylog.Error(uid, errStr)
    }

    errStruct.Desc = proto.String(errStr)
    return
}

//增加购买日志
// uid string 玩家id
// goodsId uint64 商品id
// receiptId uint64 票据
// gameId string 游戏id，如果是游戏内购买，该项有效
// failReason battery.ErrorCode 操作错误码
func (api *XYAPI) AddShoppingLog(uid string, goodsId uint64, receiptId uint64, gameId string, failReason battery.ErrorCode, errStr string) (err error) {
    gl := &battery.ShoppingLog{
        Uid:        proto.String(uid),
        GoodsId:    proto.Uint64(goodsId),
        GameId:     proto.String(gameId),
        ReceiptId:  proto.Uint64(receiptId),
        FailReason: failReason.Enum(),
        Error:      proto.String(errStr),
        OpDateStr:  proto.String(xyutil.CurTimeStr()),
    }

    return api.GetLogDB().AddShoppingLog(gl)
}

//记录商品购买交易信息（只记录交易成功的）
// uid string 玩家id
// goodsId uint64 商品id
// gameId string 游戏id
func (api *XYAPI) AddShoppingTransaction(uid string, goodsId uint64, gameId string) (err error) {
    st := &battery.ShoppingTransaction{
        Uid:     proto.String(uid),
        GoodsId: proto.Uint64(goodsId),
        GameId:  proto.String(gameId),
        OpDate:  proto.Int64(xyutil.CurTimeSec()),
    }

    return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SHOPPINGTRANSACTION).AddShoppingTransaction(st)
}

////获取玩家单局游戏内已经购买对应商品的数目（针对游戏内购买）
//// uid string 玩家id
//// goodsId uint64 商品id
//// gameId string 游戏id
//func (api *XYAPI) GetShoppingCount(uid string, goodsId uint64, gameId string) (count int, err error) {
//	count, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SHOPPINGTRANSACTION).GetShoppingCount(uid, goodsId, gameId)
//	return
//}

//获取玩家当天已经购买对应商品的数目
// uid string 玩家id
// goodsId uint64 商品id
func (api *XYAPI) GetShoppingCountOfDay(uid string, goodsId uint64) (count int, err error) {
    begin, end := xyutil.CurTimeRangeSec()
    condition := bson.M{"uid": uid, "goodsid": goodsId, "opdate": bson.M{"$gt": begin, "$lte": end}}
    count, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SHOPPINGTRANSACTION).GetShoppingCount(condition)
    return
}

//获取玩家当天已经购买对应商品的数目
// uid string 玩家id
// goodsId uint64 商品id
func (api *XYAPI) GetShoppingCountOfGame(uid string, goodsId uint64, gameId string) (count int, err error) {
    condition := bson.M{"uid": uid, "goodsid": goodsId, "gameid": gameId}
    count, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SHOPPINGTRANSACTION).GetShoppingCount(condition)
    return
}

//获取玩家当天已经购买对应商品的数目
// uid string 玩家id
// goodsId uint64 商品id
func (api *XYAPI) GetShoppingCountOfUser(uid string, goodsId uint64) (count int, err error) {
    condition := bson.M{"uid": uid, "goodsid": goodsId}
    count, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_SHOPPINGTRANSACTION).GetShoppingCount(condition)
    return
}
