package batteryapi

import (
    "bytes"
    "code.google.com/p/goprotobuf/proto"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"

    "gopkg.in/mgo.v2"

    "guanghuan.com/xiaoyao/common/idgenerate"
    xylog "guanghuan.com/xiaoyao/common/log"
    xyutil "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/cache"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/iapstatistic"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

var DefReceiptIdGenerater *xyidgenerate.IdGenerater

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

func (api *XYAPI) OperationIapValidate(req *battery.OrderVerifyRequest, resp *battery.OrderVerifyResponse) (err error) {

    var (
        uid                                        = req.GetUid()
        fail_reason         int32                  = xyerror.IAP_SUCCESS
        iap_receipt                                = req.GetReceiptData()
        client_version      battery.CLIENT_VERSION = req.GetClientVer()
        moneygoods          *battery.MallItem
        isExist             bool = false
        receipt_str, opDate string
        now                 int64
    )

    //获取请求的终端平台类型
    platform := req.GetPlatformType()
    api.SetDB(platform)

    tids := req.GetTransactionId()

    // 初始化返回数据
    resp.Uid = proto.String(req.GetUid())
    resp.OrderId = proto.String(req.GetOrderId())
    resp.IsSucc = proto.Bool(false)
    resp.DiamondCount = proto.Int32(-1)
    resp.Error = xyerror.DefaultError()

    //为了方便后续的查找，把transactionid保存在set里面
    setTids := make(Set)

    //用于记录iaptransaction的切片
    iapTransactions := make([]battery.IapTransaction, 0)

    //用于单个transaction的标识定义
    var (
        bKick         bool  = false               //是否被过滤
        bOldStyle     bool  = false               //是否是1.00.00的客户端，判断标准是请求中的OrderId是否非""
        subFailReason int32 = xyerror.IAP_SUCCESS //过滤原因
    )

    //只有旧版本会上报orderid
    if resp.GetOrderId() != "" {
        bOldStyle = true
    }

    if len(iap_receipt) <= 10 {
        fail_reason = xyerror.IAP_INVALID_RECEIPT
        goto ErrorHandler
    }

    //1.00.00以后的客户端请求需要根据请求信息中transacationid列表进行防重放校验
    if !bOldStyle {
        for _, tid := range tids {
            isExist, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_IAPTRANSACTION).IsTransactionExsit(tid)
            if err != xyerror.ErrOK {
                subFailReason = xyerror.IAP_DB_ERROR
                bKick = true
                xylog.Error(uid, "verify receipt db error: %v", err)
            } else if isExist {
                subFailReason = xyerror.IAP_SUCCESS
                bKick = true
                xylog.Warning(uid, "verify receipt transactionid [%s] already exists in DB.iaptransaction", tid)
            }

            //如果交易被过滤掉了，添加到返回结果中
            if bKick {
                api.appendTransactionItem(tid, subFailReason, &(resp.Items))
                continue
            }

            setTids[tid] = empty{}
        }

        //待处理的列表为空的话直接返回成功，在resp的transation列表信息中会有相应的交易结果，客户端根据此结果信息判断是否要提交finishTransation
        //出现这种情况认为是因为iaptransation已经保存了相应transaction的信息
        if len(setTids) <= 0 {
            fail_reason = xyerror.IAP_SUCCESS
            resp.IsSucc = proto.Bool(true)
            goto ErrorHandler
        }
    }

    // 苹果验证
    // IAP_STEP_VERIFY_RECEIPT_BEFORE

    fail_reason, err = api.VerifyReceipt(uid, setTids, iap_receipt, &receipt_str, DefConfigCache.Configs().IsProduction, &iapTransactions, &(resp.Items), client_version, bOldStyle)
    xylog.Info(uid, "[IapValidate] verify receipt via apple done: %d , err[%v]", fail_reason, err)
    if err != xyerror.ErrOK || fail_reason != xyerror.IAP_SUCCESS {
        xylog.Error(uid, "verify receipt error: %v (%d)", err, fail_reason)
        goto ErrorHandler
    }

    // 验证成功
    resp.IsSucc = proto.Bool(true)

    xylog.Debug(uid, "result iapTransactions : %v", uid, &iapTransactions)

    //插入交易记录
    now, opDate = xyutil.CurTimeSec(), xyutil.CurTimeStr()
    for _, transaction := range iapTransactions {
        transaction.Uid = proto.String(uid)
        transaction.OpDate = proto.String(opDate)
        transaction.Timestamp = proto.Int64(now)
        xylog.Debug(uid, "transaction info : %v", &transaction)

        //1.00.00以前的客户端，需要在对收据内的transaction遍历时进行防重放校验
        tid := transaction.GetTransactionId()
        if bOldStyle {
            isExist, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_IAPTRANSACTION).IsTransactionExsit(tid)
            if err != xyerror.ErrOK || err != xyerror.ErrNotFound {
                subFailReason = xyerror.IAP_DB_ERROR
                bKick = true
                xylog.Error(uid, "verify receipt db error: %v", err)
            } else if isExist {
                subFailReason = xyerror.IAP_SUCCESS
                bKick = true
                xylog.Warning(uid, "verify receipt transactionid [%s] already exists in DB.iaptransaction", tid)
            }

            //如果交易被过滤掉了，添加到返回结果中
            if bKick {
                api.appendTransactionItem(tid, subFailReason, &(resp.Items))
                continue
            }
        }

        err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_IAPTRANSACTION).AddIapTransaction(&transaction)
        xylog.Debug(uid, "IapTransaction %v", &transaction)

        // 查询购买的商品属性
        moneygoods = xybusinesscache.DefGoodsCacheManager.IapGood(transaction.GetItemId())
        if moneygoods == nil {
            xylog.Error(uid, "[%s] [IapValidate] get iapgood(%s) info err: %v", uid, transaction.GetItemId(), err)
            api.appendTransactionItem(transaction.GetTransactionId(), xyerror.IAP_INVALID_GOODS, &(resp.Items))
            continue
        }

        // 扫尾工作：增加物品
        xylog.Debug(uid, "add MoneyGoods %v", moneygoods)

        //err = api.GainProps(uid, accountWithFlag, moneygoods.GetItems(), ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_iap)
        err = api.GainProps(uid, nil, moneygoods.GetItems(), ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_iap)
        if err == xyerror.ErrOK {
            api.SendIapStatistic(uid, transaction.GetItemId(), tid)
        }

        //记录该条交易结果
        api.appendTransactionItem(transaction.GetTransactionId(), xyerror.IAP_SUCCESS, &(resp.Items))
        resp.GoodsList = append(resp.GoodsList, moneygoods.GetId())
    }

    if err != xyerror.ErrOK {
        if fail_reason == xyerror.IAP_SUCCESS {
            fail_reason = xyerror.IAP_UPDATE_ACCOUNT_FAIL
        } else {
            err = errors.New("Verify success, fail to update account data")
        }
    }

    xylog.Debug(uid, "[IapValidate] current diamond: %d", resp.GetDiamondCount())

ErrorHandler:

    resp.ReceiptDetail = proto.String(receipt_str)

    //xylog.Debug("resp : %v", resp)

    return
}

func (api *XYAPI) SendIapStatistic(uid, productIdentity, transactionId string) {
    //如果商品发放成功，则往数据中心发送一条统计信息
    var sandbox int32 = 0
    if !DefConfigCache.Configs().IsProduction {
        sandbox = 1
    }

    //查询玩家的devicetoken
    //  devicetoken, identityString, err := api.GetDeviceToken(uid, mgo.Strong)
    devicetoken, identityString, err := api.GetStatisticsDevid(uid, mgo.Strong)
    if err == xyerror.ErrOK {
        var note string
        note, err = api.NoteByGid(uid)
        if err == xyerror.ErrOK {

            iapStatistic := &battery.IapStatistic{
                Uid:      &uid,
                Username: &note,
                Pid:      proto.String(productIdentity),
                Rname:    proto.String(identityString),
                Quantity: proto.Int32(1),
                Sandbox:  &sandbox,
                DeviceId: &devicetoken,
                Oid:      &transactionId,
            }
            xyiapstatistic.Send(iapStatistic)
        } else {
            xylog.Error(uid, "[%s] NoteByGid failed : %v", uid, err)
        }
    } else {
        xylog.Error(uid, "GetDeviceToken failed : %v", err)
    }
}

// 向app store 请求验证
func (api *XYAPI) VerifyReceipt(uid string, setTids Set, receipt_data string, receipt_str *string, isProduction bool, iapTransactions *[]battery.IapTransaction, items *[]*battery.TransactionItem, client_ver battery.CLIENT_VERSION, bOldStyle bool) (fail_reason int32, err error) {
    type Receipt struct {
        Receipt_data string `json:"receipt-data,omitempty"`
    }
    var (
        verifyUrl string
        receipt   Receipt
        rcpt      []byte
    )
    fail_reason = xyerror.IAP_SUCCESS
    xylog.Debug(uid, "Verify receipt via apple, receipt(%d): \n%s", len(receipt_data), receipt_data)

    if isProduction {
        verifyUrl = "https://buy.itunes.apple.com/verifyReceipt"
    } else {
        verifyUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
    }

    receipt.Receipt_data = receipt_data

    rcpt, err = json.Marshal(receipt)
    if err != xyerror.ErrOK {
        fail_reason = xyerror.IAP_FAILED_BEFORE_INVOKE_APPLE_VARIFY
        xylog.Error(uid, "[VerifyReceipt] Error marshal json msg: %s", err.Error())
        return
    }
    if len(rcpt) > 20 {
        xylog.Debug(uid, "[VerifyReceipt] app store verify receipt:%s...%s (%d)", string(rcpt[:10]), string(rcpt[len(rcpt)-10:]), len(rcpt))
    } else {
        xylog.Debug(uid, "[VerifyReceipt] app store verify receipt:%s (%d)", string(rcpt), len(rcpt))
    }

    var verifyResp *http.Response
    verifyResp, err = http.Post(verifyUrl, "application/json", bytes.NewReader(rcpt))

    xylog.Debug(uid, " [VerifyReceipt] get response from apple, err: %v", err)

    if err != xyerror.ErrOK {
        fail_reason = xyerror.IAP_FAILED_TO_INVOKE_APPLE_VERIFY
        xylog.Error(uid, "[VerifyReceipt] app store verify rst err: %v", err)
        return
    } else {
        // IAP_STEP_VERIFY_RECEIPT_POST
        defer verifyResp.Body.Close()
        var (
            r           interface{}
            res         map[string]interface{}
            isType      bool
            verify_data []byte
        )

        verify_data, err = ioutil.ReadAll(verifyResp.Body)
        decoder := json.NewDecoder(bytes.NewReader(verify_data))
        decoder.Decode(&r)

        res, isType = r.(map[string]interface{})

        if !isType {
            fail_reason = xyerror.IAP_APPLE_VARIFY_INVALID_RESPONSE
            xylog.Error(uid, "respone is not a valid json msg: %s", string(verify_data))
            err = errors.New("respone is not a valid json msg")
        } else {
            *receipt_str = fmt.Sprintf("%v", res)

            //对收据进行详细校验
            fail_reason, err = api.verifyReceipt(uid, setTids, res, iapTransactions, items, client_ver, bOldStyle)
        }
    }
    return
}

func (api *XYAPI) verifyReceipt(uid string, setTids Set, res map[string]interface{}, iapTransactions *[]battery.IapTransaction, items *[]*battery.TransactionItem, client_ver battery.CLIENT_VERSION, bOldStyle bool) (fail_reason int32, err error) {
    xylog.Debug(uid, "r : %v", res)

    //校验收据状态
    //var value interface{}
    if value, ok := res[RECEIPT_STATUS]; ok {
        xylog.Debug(uid, "[VerifyReceipt] status: %v", value)
        if value != float64(0) {
            fail_reason = xyerror.IAP_APPLE_VARIFY_FAIL
            s := fmt.Sprintf("apple returns: %0.0f", value.(float64))
            err = errors.New(s)
            xylog.Error(uid, "[VerifyReceipt] failed with status: %v", value)
            return
        }
    }
    if bOldStyle {
        //1.00.00 版本客户端请求中无系统版本，只能一个个试，从新到旧的做
        fail_reason, err = api.verifyIOS7(uid, setTids, res, iapTransactions, items, bOldStyle)
        if fail_reason != xyerror.IAP_SUCCESS || err != nil {
            fail_reason, err = api.verifyIOS6(uid, setTids, res, iapTransactions, items, bOldStyle)
        }
    } else {
        //不同的客户端版本，不同的解析
        switch client_ver {
        case battery.CLIENT_VERSION_IOS_6:
            fail_reason, err = api.verifyIOS6(uid, setTids, res, iapTransactions, items, bOldStyle)
        case battery.CLIENT_VERSION_IOS_7:
            fail_reason, err = api.verifyIOS7(uid, setTids, res, iapTransactions, items, bOldStyle)
        }
    }

    return
}

func (api *XYAPI) verifyIOS6(uid string, setTids Set, res map[string]interface{}, iapTransactions *[]battery.IapTransaction, items *[]*battery.TransactionItem, bOldStyle bool) (fail_reason int32, err error) {
    //校验app bundle信息
    if _, ok := res[RECEIPT]; ok {
        receipt, isType := res[RECEIPT].(map[string]interface{})
        if !isType {
            fail_reason = xyerror.IAP_INVALID_RECEIPT
            s := fmt.Sprintf("iap error res[RECEIPT] [%v] can't convert to map", res[RECEIPT])
            err = errors.New(s)
            xylog.Error(uid, "[VerifyReceipt] res[RECEIPT] [%v] can't convert to map", res[RECEIPT])
            return
        }

        //校验bundle_id
        if bundle_id, ok := receipt[RECEIPT_BID]; ok {
            xylog.Debug(uid, "[VerifyReceipt] bid : [%v]", bundle_id)
            //if bundle_id != SB_BUNDLE_ID {
            if _, found := DefConfigCache.Master().SetBundleId[bundle_id]; !found {
                fail_reason = xyerror.IAP_APPLE_VARIFY_BUNDLEID_FAIL
                s := fmt.Sprintf("iap error receipt[bid] : [%v]", bundle_id)
                err = errors.New(s)
                xylog.Error(uid, "[VerifyReceipt] failed with bundle id : [%v]", bundle_id)
                return
            }
        } else {
            fail_reason = xyerror.IAP_APPLE_VARIFY_BUNDLEID_FAIL
            s := fmt.Sprintf("iap error receipt[bid] not exist")
            err = errors.New(s)
            xylog.Error(uid, "[VerifyReceipt] iap error receipt[bid] not exist")
            return
        }

        iapTransaction := battery.IapTransaction{}

        var tid string
        if rtid, ok := receipt[RECEIPT_TRACSACTION]; ok {

            //查找receipt中的transactionid是否存在于客户端上报的transactionid列表中
            //只针对1.00.00后的客户端版本
            if rtid != "" && !bOldStyle {
                _, ok = setTids[rtid]
            }

            if rtid == "" || !ok {
                fail_reason = xyerror.IAP_APPLE_VARIFY_TRANSACTIONID_FAIL
                s := fmt.Sprintf("iap error receipt[transaction_id] [%v] not good, setTids : [%v]", rtid, setTids)
                err = errors.New(s)
                xylog.Debug(uid, "[VerifyReceipt] receipt[transaction_id] [%v] not good, setTids : [%v]", rtid, setTids)
                return
            } else { //设置transactionid
                if srtid, isType := rtid.(string); isType {
                    iapTransaction.TransactionId = proto.String(srtid)
                    tid = srtid
                }
            }

        } else {
            fail_reason = xyerror.IAP_APPLE_VARIFY_TRANSACTIONID_FAIL
            s := fmt.Sprintf("iap error receipt[transaction_id] not exist %v", receipt)
            err = errors.New(s)
            xylog.Error(uid, "[VerifyReceipt] receipt[transaction_id] not exist %v", receipt)
            return
        }

        //校验product_id信息
        if value, ok := receipt[RECEIPT_PRODUCT_ID]; ok {
            product_id, isType := value.(string)
            if isType {
                canBuy, mallItem, errStruct := api.isIapGoodCanBuy(uid, product_id)
                if !canBuy || errStruct.GetCode() != battery.ErrorCode_NoError || mallItem == nil {
                    api.appendTransactionItem(tid, xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL, items)
                    return
                } else { //设置一下本次transaction对应的product_id
                    iapTransaction.ItemId = proto.String(product_id)
                }

            } else {
                fail_reason = xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL
                s := fmt.Sprintf("iap error receipt[product_id] : %v no string", product_id)
                err = errors.New(s)
                xylog.Error(uid, "[VerifyReceipt] failed with product_id : %v no string", product_id)
                api.appendTransactionItem(tid, xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL, items)
                return
            }
        } else {
            fail_reason = xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL
            s := fmt.Sprintf("iap error receipt[product_id] not exist")
            err = errors.New(s)
            xylog.Error(uid, "[VerifyReceipt] iap error receipt[product_id] not exist")
            api.appendTransactionItem(tid, xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL, items)
            return
        }

        xylog.Debug(uid, "append iapTransaction %v to iapTransactions %v", iapTransaction, *iapTransactions)
        iapTransaction.Receipt = proto.String(fmt.Sprintf("%v", receipt))
        *iapTransactions = append(*iapTransactions, iapTransaction)
    }

    fail_reason = xyerror.IAP_SUCCESS
    return
}

func (api *XYAPI) verifyIOS7(uid string, setTids Set, res map[string]interface{}, iapTransactions *[]battery.IapTransaction, items *[]*battery.TransactionItem, bOldStyle bool) (fail_reason int32, err error) {

    //校验app bundle信息
    if _, ok := res[RECEIPT]; ok {
        receipt, isType := res[RECEIPT].(map[string]interface{})
        if !isType {
            fail_reason = xyerror.IAP_INVALID_RECEIPT
            s := fmt.Sprintf("iap error res[RECEIPT] [%v] can't convert to map", res[RECEIPT])
            err = errors.New(s)
            xylog.Error(uid, "[VerifyReceipt] res[RECEIPT] [%v] can't convert to map", res[RECEIPT])
            return
        }

        //校验bundle_id
        if bundle_id, ok := receipt[RECEIPT_BUNDLE_ID]; ok {
            xylog.Debug(uid, "[VerifyReceipt] bundle id : [%v]", bundle_id)
            if _, found := DefConfigCache.Master().SetBundleId[bundle_id]; !found {
                fail_reason = xyerror.IAP_APPLE_VARIFY_BUNDLEID_FAIL
                s := fmt.Sprintf("iap error receipt[bundle_id] : [%v]", bundle_id)
                err = errors.New(s)
                xylog.Error(uid, "[VerifyReceipt] failed with bundle id : [%v]", bundle_id)
                return
            }
        } else {
            fail_reason = xyerror.IAP_APPLE_VARIFY_BUNDLEID_FAIL
            s := fmt.Sprintf("iap error receipt[bundle_id] not exist")
            err = errors.New(s)
            xylog.Error(uid, "[VerifyReceipt] iap error receipt[bundle_id] not exist")
            return
        }

        if in_app, ok := receipt[RECEIPT_INAPP]; ok {
            transactions, isType := in_app.([]interface{}) //可能存在多个
            if !isType {
                fail_reason = xyerror.IAP_APPLE_VARIFY_FAIL
                s := fmt.Sprintf("iap receipt : %v can't convert to slice", in_app)
                err = errors.New(s)
                xylog.Error(uid, "[VerifyReceipt] res[in_app] : %v can't convert to slice", in_app)
                return
            } else {
                bGetDone := false
                iapTransaction := battery.IapTransaction{}
                for k, t := range transactions {
                    xylog.Debug(uid, "transactions[%d] : [%v]", k, t)
                    info, isType := t.(map[string]interface{})

                    if !isType {
                        xylog.Error(uid, "[VerifyReceipt] can't convert transaction %v to map", t)
                        continue
                    }

                    var tid string
                    if rtid, ok := info[RECEIPT_TRACSACTION]; ok {

                        //查找receipt中的transactionid是否存在于客户端上报的transactionid列表中
                        //只针对1.00.00后的客户端版本
                        if rtid != "" && !bOldStyle {
                            _, ok = setTids[rtid]
                        }

                        if rtid == "" || !ok {
                            fail_reason = xyerror.IAP_APPLE_VARIFY_TRANSACTIONID_FAIL
                            s := fmt.Sprintf("iap error receipt[transaction_id] [%v] not good, setTids : [%v]", rtid, setTids)
                            err = errors.New(s)
                            xylog.Debug(uid, "[VerifyReceipt] receipt[transaction_id] [%v] not good, setTids : [%v]", rtid, setTids)
                            continue
                        } else { //设置transactionid
                            if srtid, isType := rtid.(string); isType {
                                iapTransaction.TransactionId = proto.String(srtid)
                                tid = srtid
                            }
                        }

                    } else {
                        fail_reason = xyerror.IAP_APPLE_VARIFY_TRANSACTIONID_FAIL
                        s := fmt.Sprintf("iap error receipt[transaction_id] not exist %v", info)
                        err = errors.New(s)
                        xylog.Error(uid, "[VerifyReceipt] receipt[transaction_id] not exist %v", info)
                        continue
                    }

                    //校验product_id信息
                    if value, ok := info[RECEIPT_PRODUCT_ID]; ok {
                        product_id, isType := value.(string)
                        if isType {
                            canBuy, mallItem, errStruct := api.isIapGoodCanBuy(uid, product_id)
                            if !canBuy || errStruct.GetCode() != battery.ErrorCode_NoError || mallItem == nil {
                                api.appendTransactionItem(tid, xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL, items)
                                continue
                            } else { //设置一下本次transaction对应的product_id
                                iapTransaction.ItemId = proto.String(product_id)
                            }

                        } else {
                            fail_reason = xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL
                            s := fmt.Sprintf("iap error receipt[product_id] : %v no string", product_id)
                            err = errors.New(s)
                            xylog.Error(uid, "[VerifyReceipt] failed with product_id : %v no string", product_id)
                            api.appendTransactionItem(tid, xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL, items)
                            continue
                        }
                    } else {
                        fail_reason = xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL
                        s := fmt.Sprintf("iap error receipt[product_id] not exist")
                        err = errors.New(s)
                        xylog.Error(uid, "[VerifyReceipt] iap error receipt[product_id] not exist")
                        api.appendTransactionItem(tid, xyerror.IAP_APPLE_VARIFY_PRODUCTID_FAIL, items)
                        continue
                    }

                    xylog.Debug(uid, "append iapTransaction %v to iapTransactions %v", iapTransaction, *iapTransactions)
                    iapTransaction.Receipt = proto.String(fmt.Sprintf("%v", info))
                    *iapTransactions = append(*iapTransactions, iapTransaction)

                    bGetDone = true
                }

                //只要匹配到一个，就返回成功
                if bGetDone {
                    fail_reason = xyerror.IAP_SUCCESS
                    err = xyerror.ErrOK
                }
            }
        }
    }

    return
}

func (api *XYAPI) appendTransactionItem(tid string, fail_reason int32, items *[]*battery.TransactionItem) {
    item := &battery.TransactionItem{
        TransactionId: proto.String(tid),
        FailReason:    proto.Int32(fail_reason),
    }
    *items = append(*items, item)
}

//判断iapId对应的商品是否存在
// iapId string 商品的iapid
func (api *XYAPI) iapGoodExist(iapId string) bool {
    iapGood := xybusinesscache.DefGoodsCacheManager.IapGood(iapId)
    if iapGood != nil {
        return true
    } else {
        return false
    }
}

//判断商品是否可以购买
// uid string 玩家id
// iapGoodId string iap商品标识
func (api *XYAPI) isIapGoodCanBuy(uid, iapGoodId string) (canBuy bool, mallItem *battery.MallItem, errStruct battery.Error) {
    var (
        errStr string
        //expiredDate   int64
        //amountPerUser uint32
    )

    //初始化返回值
    canBuy = true
    errStruct.Code = xyerror.Resp_NoError.GetCode().Enum()

    //商品是否存在
    mallItem = xybusinesscache.DefGoodsCacheManager.IapGood(iapGoodId)
    if nil == mallItem {
        errStr = fmt.Sprintf("[%s] iapgood(%s) doesn't exists", uid, iapGoodId)
        errStruct.Code = xyerror.Resp_QueryGoodsError.GetCode().Enum()
        goto ErrHandle
    }

    //to delete:
    //对于iap商品，只要上报了校验信息就说明在appstore完成了付款，需要无条件地为其发放对应物品。
    //以下代码可删除。
    ////判断商品是否过期
    //expiredDate = mallItem.GetExpiretimestamp()
    //if expiredDate > 0 { //当expiredDate大于0时才有效
    //	curTime := xyutil.CurTimeSec()
    //	if expiredDate < curTime {
    //		errStr = fmt.Sprintf("[%s] iapgood(%d) expired, can't buy", uid, iapGoodId)
    //		errStruct.Code = xyerror.Resp_QueryGoodsError.GetCode().Enum()
    //		goto ErrHandle
    //	}
    //}

    ////是否超过玩家购买上限
    //amountPerUser = mallItem.GetAmountperuser()
    //if amountPerUser > 0 {
    //	curCount, _ := api.GetIapShoppingCount(uid, iapGoodId)
    //	if amountPerUser <= uint32(curCount) {
    //		errStr = fmt.Sprintf("[%s] iapGoodId(%d) shopping count (%d) is over amountPerUser(%d)", uid, iapGoodId, curCount, amountPerUser)
    //		errStruct.Code = xyerror.Resp_BuyGoodOverAmountPerUser.GetCode().Enum()
    //		goto ErrHandle
    //	}
    //}

    //没有错误，直接返回
    return

ErrHandle:
    xylog.Error(uid, errStr)
    canBuy = false
    return
}

//获取玩家购买iap商品的次数
// uid string 玩家标识
// iapGoodId string iap商品标识
func (api *XYAPI) GetIapShoppingCount(uid, iapGoodId string) (count int, err error) {
    count, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_IAPTRANSACTION).GetIapShoppingCount(uid, iapGoodId)
    return
}
