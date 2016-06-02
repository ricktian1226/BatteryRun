package main

import (
    "crypto/md5"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "sort"
    "strconv"
    "strings"
    "time"

    proto "code.google.com/p/goprotobuf/proto"
    martini "github.com/codegangsta/martini"
    nats "github.com/nats-io/nats"

    batteryapi "guanghuan.com/xiaoyao/battery_maintenance_server/bussiness"
    xyencoder "guanghuan.com/xiaoyao/common/encoding"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
    "guanghuan.com/xiaoyao/superbman_server/error"
)

func HttpSDKCallBack(w http.ResponseWriter, r *http.Request, params martini.Params) (status int, resp string) {
    var (
        uri string = r.RequestURI

        respData []byte
        err      error
    )
    xylog.DebugNoId("req url :%v", r.RequestURI)

    status = http.StatusOK

    // uri, token = ProcessUri(uri)
    // xylog.DebugNoId("uri=%s, token=%s, user agent=%s", uri, token, r.UserAgent())
    // respData, err = ProcessHttpMsg(token, r)

    respData, err = ProcessCallBackMsg(uri, r)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[%s] failed: %s", uri, err.Error())
        status = http.StatusInternalServerError //处理失败，返回服务端错误
    } else {
        resp = getCallBackResp(respData)
        xylog.DebugNoId("response.content : %s", resp)
    }
    return
}

func ProcessCallBackMsg(uri string, r *http.Request) (resp []byte, err error) {
    var (
        req   proto.Message
        route *HttpPostToNatsRoute
        subj  string
        reply *nats.Msg
        data  []byte
    )

    req, err = constructCallBackMsg(r)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("ConstructPbMsg failed : %v", err)
        return
    }
    xylog.DebugNoId("PbMsg : %v", req)

    //进行编码
    data, err = xyencoder.PbEncode(req)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("xyencoder.PbEncode failed : %v", err)
        return
    }

    //进行加密
    data, err = crypto.Encrypt(data)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("crypto.Encrypt failed : %v", err)
        return
    }

    route = DefHttpPostTable.GetRoutePath(uri)
    xylog.DebugNoId("route:%v,item:%v", route, uri)
    if route == nil {
        err = errors.New("No route for item :" + uri)
        return
    }
    subj = route.GetNatsSubject()
    if subj == "" {
        err = errors.New("No subject for uri:" + uri)
        return
    }
    xylog.DebugNoId("forward request to %s", subj)

    reply, err = nats_service.Request(subj, data, time.Duration(DefConfig.NatsTimeout)*time.Second)
    if err != nil {
        xylog.ErrorNoId("<%s> Error: %s", subj, err.Error())
        return
    } else {
        if reply != nil {
            resp = reply.Data
        } else {
            err = errors.New("no reply data")
        }
    }
    return
}

const (
    CallParameter_UUid       = "uuid"             // sdk唯一id
    CallParameter_OrderId    = "order_id"         // sdk 订单号
    CallParameter_AppOrderId = "app_order_id"     // 游戏订单号
    CallParameter_EXT        = "app_callback_ext" //扩展参数,保存商品id
    CallParameter_UserUID    = "app_player_id"    // 游戏uid
    CallParameter_Amount     = "pay_amount"       // 充值数量
    CallParameter_PayTime    = "pay_time"
    CallParameter_Sandbox    = "sandbox" // 是否测试

    CallParameter_Sign   = "sign" // 签名校验值
    CallParameter_ZoneId = "app_zone_id"
    CallParameter_Time   = "time"
    CallParameter_UserID = "app_user_id" // 账号id
)

func constructCallBackMsg(r *http.Request) (message proto.Message, err error) {
    var (
        goodsId            uint64
        sandbox, payAmount int
        parameterSlice     ParameterSlice = make([]ParameterStruct, 0)
    )
    r.ParseForm()
    parameterSlice.Add(CallParameter_Sign, r.PostFormValue(CallParameter_Sign))
    parameterSlice.Add(CallParameter_ZoneId, r.PostFormValue(CallParameter_ZoneId))
    parameterSlice.Add(CallParameter_Time, r.PostFormValue(CallParameter_Time))
    parameterSlice.Add(CallParameter_UserID, r.PostFormValue(CallParameter_UserID))
    parameterSlice.Add(CallParameter_UUid, r.PostFormValue(CallParameter_UUid))
    parameterSlice.Add(CallParameter_AppOrderId, r.PostFormValue(CallParameter_AppOrderId))

    req := &battery.SDKAddOrderRequest{}

    v := r.PostFormValue(CallParameter_UserUID)
    if v == "" {
        xylog.ErrorNoId("get uid from parameter failed")
        return nil, xyerror.ErrBadInputData
    }
    req.Uid = proto.String(v)
    parameterSlice.Add(CallParameter_UserUID, v)

    v = r.PostFormValue(CallParameter_OrderId)
    if v == "" {
        xylog.ErrorNoId("get orderid from parameter fail")
        return nil, xyerror.ErrBadInputData
    }
    req.OrderId = proto.String(v)
    parameterSlice.Add(CallParameter_OrderId, v)

    v = r.PostFormValue(CallParameter_EXT)
    if v == "" {
        xylog.ErrorNoId("get goodsid from parameter fail")
        return nil, xyerror.ErrBadInputData
    }
    parameterSlice.Add(CallParameter_EXT, v)
    goodsId, err = strconv.ParseUint(v, 10, 64)
    if err != xyerror.ErrOK {
        xylog.WarningNoId("strconv.ParseUint for uid failed : %v", err)
        err = xyerror.ErrOK
    } else {
        req.GoodsId = proto.Uint64(goodsId)
    }

    v = r.PostFormValue(CallParameter_Sandbox)
    if v == "" {
        xylog.ErrorNoId("get sandbox from parameterfail")
        return nil, xyerror.ErrBadInputData
    }
    parameterSlice.Add(CallParameter_Sandbox, v)
    sandbox, err = strconv.Atoi(v)
    if err != xyerror.ErrOK {
        xylog.WarningNoId("strconv.Atoi for uid failed : %v", err)
        err = xyerror.ErrOK
    } else {
        req.Sandbox = proto.Int32(int32(sandbox))
    }

    v = r.PostFormValue(CallParameter_Amount)
    if v == "" {
        xylog.DebugNoId("get paytime from parameter fail")
    } else {
        parameterSlice.Add(CallParameter_Amount, v)
        payAmount, err = strconv.Atoi(v)
        if err != xyerror.ErrOK {
            xylog.WarningNoId("strconv.Atoi for uid failed : %v", err)
            err = xyerror.ErrOK
        } else {
            req.PayAmount = proto.Int32(int32(payAmount))
        }
    }

    v = r.PostFormValue(CallParameter_PayTime)
    if v == "" {
        xylog.WarningNoId("get paytime from parameter fail")
    } else {
        parameterSlice.Add(CallParameter_PayTime, v)
        req.PayTime = proto.String(v)
    }

    // 数据验证
    var key string
    key = batteryapi.DefConfigCache.Configs().AppSecretkey
    if !SDKVerification(parameterSlice, key) {
        xylog.ErrorNoId("verify error,invalid sig")
        err = xyerror.ErrBadInputData
        return
    }

    message = req
    return
}

func getCallBackResp(respData []byte) (resp string) {

    respData, err := crypto.Decrypt(respData)
    if err != nil {
        resp = "decrypt false"
        return
    }

    respone := &battery.SDKAddOrderResponse{}
    err = proto.Unmarshal(respData, respone)
    if err != nil {
        resp = "unmarshal false"
        return
    }

    if respone.Error.GetCode() != battery.ErrorCode_NoError {
        resp = "add order fail"
        return
    }
    resp = "ok"

    return
}

// 被坑，get请求使用
// func GetCallBackMsg(parameterSlice ParameterSlice) (message proto.Message, err error) {
//     var (
//         goodsId            uint64
//         sandbox, payAmount int
//     )
//     req := &battery.SDKAddOrderRequest{}

//     v, result := parameterSlice.Get(CallParameter_UserUID)
//     if !result {
//         xylog.ErrorNoId("get uid from parameter failed")
//         return nil, xyerror.ErrBadInputData
//     }
//     req.Uid = proto.String(v)

//     v, result = parameterSlice.Get(CallParameter_OrderId)
//     if !result {
//         xylog.ErrorNoId("get orderid from parameter fail")
//         return nil, xyerror.ErrBadInputData
//     }
//     req.OrderId = proto.String(v)

//     v, result = parameterSlice.Get(CallParameter_EXT)
//     if !result {
//         xylog.ErrorNoId("get goodsid from parameter fail")
//         return nil, xyerror.ErrBadInputData
//     }
//     goodsId, err = strconv.ParseUint(v, 10, 64)
//     if err != xyerror.ErrOK {
//         xylog.WarningNoId("strconv.ParseUint for uid failed : %v", err)
//         err = xyerror.ErrOK
//     } else {
//         req.GoodsId = proto.Uint64(goodsId)
//     }

//     v, result = parameterSlice.Get(CallParameter_Sandbox)
//     if !result {
//         xylog.ErrorNoId("get sandbox from parameterfail")
//         return nil, xyerror.ErrBadInputData
//     }
//     sandbox, err = strconv.Atoi(v)
//     if err != xyerror.ErrOK {
//         xylog.WarningNoId("strconv.Atoi for uid failed : %v", err)
//         err = xyerror.ErrOK
//     } else {
//         req.Sandbox = proto.Int32(int32(sandbox))
//     }

//     v, result = parameterSlice.Get(CallParameter_Amount)
//     if !result {
//         xylog.DebugNoId("get paytime from parameter fail")
//     } else {
//         payAmount, err = strconv.Atoi(v)
//         if err != xyerror.ErrOK {
//             xylog.WarningNoId("strconv.Atoi for uid failed : %v", err)
//             err = xyerror.ErrOK
//         } else {
//             req.PayAmount = proto.Int32(int32(payAmount))
//         }
//     }

//     v, result = parameterSlice.Get(CallParameter_PayTime)
//     if !result {
//         xylog.WarningNoId("get paytime from parameter fail")

//     } else {
//         req.PayTime = proto.String(v)
//     }

//     message = req
//     return
// }

// SDK数据验证
func SDKVerification(parameterslice ParameterSlice, key string) bool {
    var (
        uriStr    string
        verifyStr string
    )
    keys := make([]string, 0, len(parameterslice))
    for k := range parameterslice {
        if parameterslice[k].key != "sign" {
            keys = append(keys, parameterslice[k].key)
        }

    }
    sort.Strings(keys) // 参数升序排序
    for index, arg := range keys {
        v, result := parameterslice.Get(arg)
        if result {
            uriStr = fmt.Sprintf("%s%s=%s", uriStr, arg, v)
            if index != len(keys)-1 {
                uriStr = fmt.Sprintf("%s&", uriStr)
            }
        }

    }
    xylog.InfoNoId("uriStr:%s", uriStr)
    // 参数升序url编码
    verifyStr = url.QueryEscape(uriStr)
    verifyStr = fmt.Sprintf("%s&%s", verifyStr, key)
    // MD5加密
    h := md5.New()
    io.WriteString(h, verifyStr)
    verifyStr = fmt.Sprintf("%x", h.Sum(nil))
    verifyStr = strings.ToLower(verifyStr)
    xylog.InfoNoId("verifyStr:%s", verifyStr)
    sig, result := parameterslice.Get("sign")
    if (!result) || (sig != verifyStr) {
        return false
    }

    xylog.InfoNoId("verify succeed")
    return true

}
