package batteryapi

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "sort"
    "strings"
    "time"

    "code.google.com/p/goprotobuf/proto"

    "guanghuan.com/xiaoyao/common/log"

    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

// RetryIapStatistic 重试未上报成功的内购信息
func (api *XYAPI) RetryIapStatistics() (errCode battery.ErrorCode) {
    //从数据库查询未完成的内购信息
    statistics, err := api.GetDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_IAPSTATISTIC).QueryUndoneIapStatistics()
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("QueryUndoneIapStatistics failed : %v", err)
        errCode = battery.ErrorCode_IapQueryUndoneIapStatisticsError
        return
    }

    now := time.Now().Unix()
    //遍历所有未完成的内购信息，进行重发处理
    for _, statistic := range statistics {
        api.TrySingleDBIapStatistic(statistic, now)
    }

    return
}

//重试次数和重试时间间隔的映射关系：
// 已重试次数	下次重试时间
// [0,3)	10s
// [3,6)	60s
// [6,12)	600s
// [12,)	不重试，放弃
//var RetryCount2Time = map[int]int64{
//	3:  10,
//	6:  60,
//	12: 600,
//}

type RetryCount2TimeManager struct {
    M   map[int]int64
    S   sort.IntSlice
}

func (m *RetryCount2TimeManager) Insert(k int, v int64) {
    m.M[k] = v
    m.S = append(m.S, k)
}

func (m *RetryCount2TimeManager) Value(retryCount int) int64 {
    for _, k := range m.S {
        if retryCount < k {
            return m.M[k]
        }
    }
    return 0
}

func NewRetryCount2TimeManager() (manager *RetryCount2TimeManager) {
    manager = &RetryCount2TimeManager{
        M:  make(map[int]int64, 0),
        S:  make(sort.IntSlice, 0),
    }

    //这部分配置信息可以修改为数据库配置，这样方便修改
    manager.Insert(3, 10)
    manager.Insert(6, 60)
    manager.Insert(12, 600)
    sort.Sort(manager.S)

    return
}

var DefRetryCount2TimeManager = NewRetryCount2TimeManager()

var DefProductIdentity2ProductId = MAPProductIdentity2ProductId{
    "Buy_Goods_0704":       "76",
    "Buy_Goods_0705":       "77",
    "Buy_Goods_0702":       "78",
    "Buy_Goods_0703":       "79",
    "Buy_Goods_0110010001": "80",
    "Buy_Goods_0110040000": "81",
    "Buy_Goods_0110030000": "82",
    "Buy_Goods_0110020000": "83",
    "Buy_Goods_0706":       "84",
    "Buy_Goods_00701":      "85",
    "Buy_Goods_0901":       "114",
    "Buy_Goods_0902":       "115",
    "Buy_Goods_0903":       "116",
    "Buy_Goods_0904":       "117",
    "Buy_Goods_0905":       "118",
    "Buy_Goods_0906":       "119",
    "Buy_Goods_0907":       "120",
    "Buy_Goods_0707":       "122",
    "Buy_Goods_0708":       "123",
    "Buy_Goods_0710":       "124",
    "Buy_Goods_0712":       "125",
    "Buy_Goods_0908":       "127",
}

type MAPProductIdentity2ProductId map[string]string

// ProductId 根据productIdentity获取productId
func (m *MAPProductIdentity2ProductId) ProductId(productIdentity string) (productId string) {
    if v, ok := (*m)[productIdentity]; ok {
        productId = v
    }
    return
}

// TrySingleIapStatistic
// statistic *battery.IapStatistic 请求的IapStatistic结构体
// now int64 当前的时间戳
//returns:
//  errCode battery.ErrorCode 操作错误码
func (api *XYAPI) TrySingleIapStatistic(statistic *battery.IapStatistic, now int64) (errCode battery.ErrorCode) {
    dbIapStatistic := api.getDBIapStatistic(statistic)
    return api.TrySingleDBIapStatistic(dbIapStatistic, now)
}

// getDBIapStatistic 根据请求的IapStatistic生成DBIapStatistic信息结构体
// statistic *battery.IapStatistic 请求的IapStatistic结构体
// now int64 当前的时间戳
//return:
//  *battery.DBIapStatistic
func (api *XYAPI) getDBIapStatistic(statistic *battery.IapStatistic) *battery.DBIapStatistic {
    statistic.AppId = proto.String(DefConfigCache.Configs().AppId)
    statistic.Bis = proto.String(Default_IapStatisticValue_BusinessId)
    statistic.Ac = proto.String(Default_IapStatisticValue_Action)

    return &battery.DBIapStatistic{
        TransactionId: statistic.Oid,
        RetryCount:    proto.Int32(0),    //初始化重试次数为0
        Done:          proto.Bool(false), //初始化发送未完成
        Detail:        statistic,
    }
}

// TrySingleIapStatistic 处理单个内购统计信息的
// statistic *battery.DBIapStatistic 内购统计信息
func (api *XYAPI) TrySingleDBIapStatistic(statistic *battery.DBIapStatistic, now int64) (errCode battery.ErrorCode) {

    retryCount, retryTimestamp := statistic.GetRetryCount(), statistic.GetRetryTimestamp()

    if api.needRetry(retryCount, now, retryTimestamp) {
        //发送到数据中心
        errCode = api.send2DataCenter(statistic.GetDetail())
        if errCode == battery.ErrorCode_IapStatisticServerError { //如果是数据中心服务错误，则需要对该内购信息进行重发
            statistic.Done = proto.Bool(false)
            xylog.ErrorNoId("send2DataCenter failed : %d", errCode)
        } else {
            statistic.Done = proto.Bool(true)
        }
        statistic.RetryCount = proto.Int32(retryCount + 1) //发送次数加1
        statistic.RetryTimestamp = &now                    //上次重试时间设置为当前时间戳
        statistic.Result = errCode.Enum()                  //设置错误码

        //更新该内购统计信息
        xylog.DebugNoId("UpsertIapStatistic : %v", statistic)
        err := api.GetDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_IAPSTATISTIC).UpsertIapStatistic(statistic)
        if err != xyerror.ErrOK {
            errCode = battery.ErrorCode_IapStatisticDBUpsertError
            xylog.ErrorNoId("UpsertIapStatistic failed : %v", err)
            return
        }
    }

    return
}

//内功统计上报字段名称
const (
    IapStatisticKey_AppId          = "app_id"           //产品id
    IapStatisticKey_BusinessId     = "bis"              //业务标识
    IapStatisticKey_Action         = "ac"               //动作标识
    IapStatisticKey_ProductId      = "pid"              //产品标识
    IapStatisticKey_Quantity       = "quantity"         //产品数目
    IapStatisticKey_Amount         = "amount"           //金额数目，单位：分
    IapStatisticKey_RoleId         = "rid"              //角色标识
    IapStatisticKey_Platform       = "platform"         //平台类型
    IapStatisticKey_ClientIp       = "client_ip"        //客户端ip
    IapStatisticKey_ChannelId      = "channel_id"       //渠道标识
    IapStatisticKey_Uid            = "uid"              //玩家标识
    IapStatisticKey_UserName       = "username"         //玩家名称
    IapStatisticKey_RoleName       = "rname"            //角色名称
    IapStatisticKey_Level          = "level"            //角色等级
    IapStatisticKey_OrderId        = "oid"              //订单号
    IapStatisticKey_SandBox        = "sandbox"          //是否是沙盒
    IapStatisticKey_PlatformUserId = "platform_user_id" //第三方平台玩家标识
    IapStatisticKey_DeviceId       = "device_id"        //玩家的设备token
    IapStatisticKey_Memo           = "memo"             //备注

    Default_IapStatisticValue_BusinessId = "order"    //业务id，默认用order
    Default_IapStatisticValue_Action     = "report"   //action，默认用report
    Default_IapStatisticValue_Quantity   = "quantity" //quantity，默认用1
)

// send2DataCenter 将统计信息上报到数据中心
//detail *battery.IapStatistic 统计信息的明细
//returns:
// errCode battery.ErrorCode 操作的错误码
//       battery.ErrorCode_IapStatisticResponseUnmarshalError 数据中心返回信息解码失败
//       battery.ErrorCode_IapStatisticResponseError
func (api *XYAPI) send2DataCenter(detail *battery.IapStatistic) (errCode battery.ErrorCode) {

    errCode = battery.ErrorCode_NoError

    keys, values := make([]string, 0), make(map[string]string, 0)

    values[IapStatisticKey_AppId] = DefConfigCache.Configs().AppId
    keys = append(keys, IapStatisticKey_AppId)

    values[IapStatisticKey_BusinessId] = Default_IapStatisticValue_BusinessId
    keys = append(keys, IapStatisticKey_BusinessId)

    values[IapStatisticKey_Action] = Default_IapStatisticValue_Action
    keys = append(keys, IapStatisticKey_Action)

    // 角色名使用唯一uid标示
    values[IapStatisticKey_RoleId] = detail.GetUid()
    keys = append(keys, IapStatisticKey_RoleId)

    if nil != detail.Quantity {
        values[IapStatisticKey_Quantity] = fmt.Sprintf("%d", detail.GetQuantity())

    } else { //默认是1
        values[IapStatisticKey_Quantity] = Default_IapStatisticValue_Quantity
    }
    keys = append(keys, IapStatisticKey_Quantity)

    productId := DefProductIdentity2ProductId.ProductId(detail.GetPid())
    if productId == "" {
        xylog.ErrorNoId("GetProductId for ProductIdentity(%s) failed", detail.GetPid())
        errCode = battery.ErrorCode_IapStatisticGetProductIdError
        return
    }
    values[IapStatisticKey_ProductId] = productId
    keys = append(keys, IapStatisticKey_ProductId)

    if nil != detail.Username {
        values[IapStatisticKey_UserName] = detail.GetUsername()
        keys = append(keys, IapStatisticKey_UserName)
    }

    if nil != detail.Rname {
        values[IapStatisticKey_RoleName] = detail.GetRname()
        keys = append(keys, IapStatisticKey_RoleName)
    }

    if nil != detail.Level {
        values[IapStatisticKey_Level] = fmt.Sprintf("%d", detail.GetLevel())
        keys = append(keys, IapStatisticKey_Level)
    }

    values[IapStatisticKey_SandBox] = fmt.Sprintf("%d", detail.GetSandbox())
    keys = append(keys, IapStatisticKey_SandBox)

    values[IapStatisticKey_OrderId] = detail.GetOid()
    keys = append(keys, IapStatisticKey_OrderId)

    values[IapStatisticKey_DeviceId] = detail.GetDeviceId()
    keys = append(keys, IapStatisticKey_DeviceId)

    sort.Strings(keys)
    xylog.DebugNoId("iapstatstics key :%v", keys)

    var encodeSrc string
    for i, k := range keys {
        if i == 0 {
            encodeSrc += k + "=" + values[k]
        } else {
            encodeSrc += "&" + k + "=" + values[k]
        }
    }

    urlEncode := url.QueryEscape(encodeSrc)
    urlEncode += fmt.Sprintf("&%s", DefConfigCache.Configs().AppSecretKey)
    //xylog.DebugNoId("\nencodeSrc : %v\nurlencode : %v\n", encodeSrc, urlEncode)

    md5Ctx := md5.New()
    md5Ctx.Write([]byte(urlEncode))
    cipherStr := md5Ctx.Sum(nil)

    sign := strings.ToLower(hex.EncodeToString(cipherStr))
    encodeSrc += "&sign=" + sign

    type Response struct {
        Code int    `json:"code"`
        Desc string `json:"desc"`
        Data string `json:"data"`
    }

    reqUrl := fmt.Sprintf("%s?%v", DefConfigCache.Configs().DataCenterUrl, encodeSrc)
    resp, err := http.Get(reqUrl)
    if err != xyerror.ErrOK { //返回信息
        errCode = battery.ErrorCode_IapStatisticServerError
        xylog.ErrorNoId("Get with req(%s) failed : %v", reqUrl, err)
        return
    }
    xylog.DebugNoId("resp : %v", resp)
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    xylog.DebugNoId("resp body : %v", string(body))

    respData := &Response{}
    err = json.Unmarshal(body, respData)
    if err != nil {
        xylog.DebugNoId("failed. %v", err)
        errCode = battery.ErrorCode_IapStatisticResponseUnmarshalError
        return
    } else {
        if respData.Code != 0 {
            xylog.DebugNoId("failed.")
            errCode = battery.ErrorCode_IapStatisticResponseError
            return
        } else {
            xylog.DebugNoId("succeed.")
        }
    }

    return

}

// needRetry 判断是否需要发送
// retryCount int32 已经重试的次数
// now int64 当前的时间戳
// retryTimestamp int64 上次发送的时间戳
//return:
// bool 是否需要重发
func (api *XYAPI) needRetry(retryCount int32, now, retryTimestamp int64) bool {

    xylog.DebugNoId("retryCount(%d)  now(%d) retryTimestamp(%d)", retryCount, now, retryTimestamp)

    if retryCount == 0 { //如果是0表示是第一次发送，直接返回true
        xylog.DebugNoId("retryCount(0) first try time")
        return true
    }

    interval := DefRetryCount2TimeManager.Value(int(retryCount))
    if interval > 0 && retryTimestamp+interval <= now { //到了重发时间
        xylog.DebugNoId("retryCount(%d) interval (%d), ", retryCount, interval)
        return true
    }

    return false
}
