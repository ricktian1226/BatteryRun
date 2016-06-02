package main

import (
    "crypto/md5"
    "encoding/json"
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
    xypanic "guanghuan.com/xiaoyao/common/panic"
    xyperf "guanghuan.com/xiaoyao/common/performance"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
    "guanghuan.com/xiaoyao/superbman_server/error"
)

const (
    MaintenanceParameter_Source    = "source"
    MaintenanceParameter_Optype    = "optype"
    MaintenanceParameter_Platform  = "platform"
    MaintenanceParameter_Uid       = "identity"
    MaintenanceParameter_Userid    = "userid" //  兑换玩家uid
    MaintenanceParameter_Amount    = "amount"
    MaintenanceParameter_PROP_Id   = "propid"
    MaintenanceParameter_PROP_Type = "proptype"
    MaintenanceParameter_Sig       = "sig"
    MaintenanceParameter_TIMESTAMP = "ts"
    MaintenanceParameter_ZONEID    = "zoneid"
)

type ParameterStruct struct {
    key, value string
}

type ParameterSlice []ParameterStruct

func (s *ParameterSlice) Add(k, v string) {
    *s = append(*s, ParameterStruct{key: k, value: v})
}

//为什么用slice不用map，因为golang的map实际上是hashmap，在元素未上千级时，查询的效率不如slice遍历
func (s *ParameterSlice) Get(k string) (v string, result bool) {
    result = false
    for _, e := range *s {
        if e.key == k {
            result = true
            v = e.value
        }
    }
    return
}

func (s *ParameterSlice) Parse(values []string) {
    for _, tmp := range values {
        tmps := strings.Split(tmp, "=")
        if len(tmps) != 2 {
            xylog.ErrorNoId("invalid parameter : %v", tmp)
            return
        }
        s.Add(tmps[0], tmps[1])
    }
}

//处理pprof消息
func HttpPProf(params martini.Params) (status int, resp string) {

    var (
        appName = params["appname"]
        dcId    = params["dcid"]
        nodeId  = params["nodeid"]
        op      = params["op"]
        subject string
    )

    xylog.InfoNoId("HttpPProf : (appname:%s, dcid:%s, nodeid:%s, op:%s)", appName, dcId, nodeId, op)

    if appName == "all" { //表示所有应用
        subject = "pprof"
    } else if dcId == "all" || nodeId == "all" { //标识所有节点
        subject = fmt.Sprintf("pprof_%s", appName)
    } else { //某一特定的节点
        subject = fmt.Sprintf("pprof_%s_%s_%s", appName, dcId, nodeId)
    }

    nats_service.Publish(subject, []byte(op))

    status = http.StatusOK

    return
}

func ProcessUri(uri string) (u string, token string) {
    u = uri
    idx := strings.LastIndex(uri, "/")
    if idx > 0 {
        u = u[:idx]
        token = uri[idx+1:]
    }
    return
}

const (
    STATUS_SUCCEED = 1
    STATUS_FAILED  = 2
)

func HttpPostWorkerNoToken(w http.ResponseWriter, r *http.Request, params martini.Params) (status int, resp string) {
    var (
        uri       string = r.RequestURI
        err       error
        resp_data []byte
    )
    status = http.StatusOK
    xylog.DebugNoId("uri=%s, token=n/a, user agent=%s", uri, r.UserAgent())

    resp_data, err = ProcessHttpMsg(uri /*, ""*/, r)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[%s] failed: %s", uri, err.Error())
        status = http.StatusInternalServerError //处理失败，返回服务端错误
    } else {
        resp = getResponseContent(resp_data)
    }
    return
}

type RespInfo struct {
    Desc string
}
type RespJsonData struct {
    Status int      `json:"status"`
    Info   RespInfo `json:"info"`
    Text   string   `json:"text"`
}

func HttpGetWorker(w http.ResponseWriter, r *http.Request, params martini.Params) (status int, resp string) {
    var (
        uri       string = r.RequestURI
        token     string
        err       error
        resp_data []byte
    )
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    xylog.DebugNoId("request : %v", r.RequestURI)

    status = http.StatusOK
    uri, token = ProcessUri(uri)
    xylog.DebugNoId("uri=%s, token=%s, user agent=%s", uri, token, r.UserAgent())
    resp_data, err = ProcessHttpMsg(token, r)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[%s] failed: %s", uri, err.Error())
        status = http.StatusInternalServerError //处理失败，返回服务端错误
    } else {
        resp = getResponseContent(resp_data)
        xylog.DebugNoId("response.content : %s", resp)
    }
    return
}

const (
    RESP_STR_STATUS = "\"status\""
    RESP_STR_INFO   = "\"info\""
    RESP_STR_TEXT   = "\"text\""
    RESP_STR_DESC   = "\"desc\""
)

func getResponseContent(respData []byte) (content string) {
    //resp = string(resp_data)
    var (
        statusTmp = STATUS_FAILED
        err       error
        jsondata  = RespJsonData{}
    )

    respData, err = crypto.Decrypt(respData)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("Error decrypt: %v", err)
        // content = fmt.Sprintf("{%s:\"%d\",%s: \"crypto.Decrypt failed\",%s:{ %s:\"Error decrypt %v\"}}",
        // 	RESP_STR_STATUS,
        // 	statusTmp,
        // 	RESP_STR_TEXT,
        // 	RESP_STR_INFO,
        // 	RESP_STR_DESC,
        // 	err)
        jsondata = RespJsonData{statusTmp, RespInfo{"crypto.Decrypt failed"}, "Error decrypt"}
        data, _ := json.Marshal(jsondata)
        content = string(data)
        return
    } else {
        response := &battery.MaintenancePropResponse{}
        err = proto.Unmarshal(respData, response)
        if err != xyerror.ErrOK {
            xylog.ErrorNoId("Error unmarshal: %v", err)
            // content = fmt.Sprintf("{%s:\"%d\",%s: \"proto.Unmarshal failed\",%s:{ %s:\"Error unmarshal %v\"}}",
            // 	RESP_STR_STATUS,
            // 	statusTmp,
            // 	RESP_STR_TEXT,
            // 	RESP_STR_INFO,
            // 	RESP_STR_DESC,
            // 	err)
            jsondata = RespJsonData{statusTmp, RespInfo{"proto.Unmarshal failed"}, "Error unmarshal"}
            data, _ := json.Marshal(jsondata)
            content = string(data)
        } else {
            xylog.DebugNoId("response: %v", response)
            if response.Error.GetCode() == battery.ErrorCode_NoError {
                statusTmp = STATUS_SUCCEED
            }
            //content = fmt.Sprintf("{\"status\":\"%d\",\"info\":%s}", statusTmp)
            // content = fmt.Sprintf("{%s:%d,%s: \"operation succeed\",%s:\"%s\"}",
            // 	RESP_STR_STATUS,
            // 	statusTmp,
            // 	RESP_STR_TEXT,
            // 	RESP_STR_INFO,
            // 	response.Error.GetDesc())
            jsondata = RespJsonData{statusTmp, RespInfo{response.Error.GetDesc()}, "operation succeed"}
            data, _ := json.Marshal(jsondata)
            content = string(data)
        }
    }

    //content = fmt.Sprintf("{\"status\":1, \"info\":{\"name\":\"kkkk\",\"desc\":\"kkkk\"}}")
    //content = fmt.Sprintf("{\"status\": \"1\",\"info\":{\"source\": \"SRC_SINA_WEIBO\",\"identity\":96214777591758850,\"sid\": \"2844846670\",\"name\":\"ricktian\",\"desc\": \"NoError\"}}")

    return
}

// 兑换礼包数据验证
func Verification(parameterslice ParameterSlice, key string) bool {
    var (
        // parameter = []string{
        // 	MaintenanceParameter_Amount,
        // 	MaintenanceParameter_Optype,
        // 	MaintenanceParameter_Platform,
        // 	MaintenanceParameter_PROP_Id,
        // 	MaintenanceParameter_Source,
        // 	MaintenanceParameter_TIMESTAMP,
        // 	MaintenanceParameter_Userid,
        // 	MaintenanceParameter_ZONEID,
        // }

        uriStr    string
        verifyStr string
    )
    keys := make([]string, 0, len(parameterslice))
    for k := range parameterslice {
        if parameterslice[k].key != "sig" {
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
    sig, result := parameterslice.Get("sig")
    if (!result) || (sig != verifyStr) {
        return false
    }

    xylog.InfoNoId("verify succeed")
    return true

}
func ProcessHttpMsg(uri string, r *http.Request) (respData []byte, err error) {

    begin := time.Now()
    defer xyperf.Trace(xyperf.DefLogId, &begin)

    defer xypanic.Crash()

    var (
        reqData        []byte
        subj           string
        reply          *nats.Msg
        route          *HttpPostToNatsRoute
        item           string
        parameterSlice ParameterSlice
    )
    xylog.DebugNoId("maintenanceurl:%s", uri)
    item, parameterSlice, err = ParseUri(uri)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("ParseUri failed : %v", err)
        return
    }

    // sdk使用secret秘钥，其余使用app秘钥
    var key string  
    if item== SDKCallBackSubItem_CallBack{
        key = batteryapi.DefConfigCache.Configs().AppSecretkey
    }else {
     key = batteryapi.DefConfigCache.Configs().Appkey
    }
    // 解密验证
    if !Verification(parameterSlice, key) {
        xylog.ErrorNoId("verify error,invalid sig")
        err = xyerror.ErrBadInputData
        return
    }

    route = DefHttpPostTable.GetRoutePath(item)
    xylog.DebugNoId("route:%v,item:%v", route, item)
    if route == nil {
        err = errors.New("No route for item :" + item)
        return
    }

    subj = route.GetNatsSubject()
    if subj == "" {
        err = errors.New("No subject for uri:" + uri)
        return
    }

    reqData, err = ReadRequestData(item, parameterSlice)
    if err != xyerror.ErrOK {
        return
    }

    xylog.DebugNoId("forward request to %s", subj)

    reply, err = nats_service.Request(subj, reqData, time.Duration(DefConfig.NatsTimeout)*time.Second)

    if err != nil {
        xylog.ErrorNoId("<%s> Error: %s", subj, err.Error())
        goto ErrorHandle
    } else {
        if reply != nil {
            respData = reply.Data
        } else {
            err = errors.New("no reply data")
        }
    }
ErrorHandle:
    return
}

//根据请求uri获取消息体
func ReadRequestData(item string, parameterSlice ParameterSlice) (data []byte, err error) {

    //构造消息结构
    var req proto.Message
    req, err = ConstructPbMsg(item, parameterSlice)
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

    return
}

//根据uri构造pb消息
// uri string 请求的uri
//return:
// proto.Message 返回的pb消息
// error 处理结果
func ConstructPbMsg(item string, parameterSlice ParameterSlice) (proto.Message, error) {
    return GetPbMsg(item, parameterSlice)
}

//解析请求uri
// uri string 请求uri
//return:
// item string 运营请求item
// parameterSlice ParameterSlice 请求参数列表
// err error 返回错误
func ParseUri(uri string) (item string, parameterSlice ParameterSlice, err error) {
    //分解出subject和参数
    contents := strings.Split(uri, "?")

    if len(contents) < 2 {
        xylog.ErrorNoId("error contents : %v ", contents)
        err = xyerror.ErrBadInputData
        return
    }

    var parameters []string
    item, parameters = contents[0], strings.Split(contents[1], "&")
    parameterSlice.Parse(parameters)
    if len(parameterSlice) <= 0 {
        err = xyerror.ErrBadInputData
        return
    }

    return
}

//获取pb消息内容
func GetPbMsg(item string, parameterSlice ParameterSlice) (proto.Message, error) {
    switch item {
    case MaintenanceSubItem_Prop:
        return GetPropPbMsg(parameterSlice)
    case MaintenanceSubItem_CDkey:
        return GetCDkeyPropPbMsg(parameterSlice)
    }
    return nil, xyerror.ErrBadInputData
}

// 获取兑换码礼包信息
func GetCDkeyPropPbMsg(parameterSlice ParameterSlice) (message proto.Message, err error) {
    var (
        propId                   uint64
        amount, platform, optype int
    )
    req := &battery.MaintenancePropRequest{}
    // userid ,requested
    v, result := parameterSlice.Get(MaintenanceParameter_Userid)
    if !result {
        xylog.ErrorNoId("get uid from parameterSlice failed")
        return nil, xyerror.ErrBadInputData
    }
    req.Uid = proto.String(v)

    // platform ,requested
    v, result = parameterSlice.Get(MaintenanceParameter_Platform)
    if !result {
        xylog.ErrorNoId("get platform from parameterSlice failed")
        return nil, xyerror.ErrBadInputData
    }
    platform, err = strconv.Atoi(v)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("strconv.Atoi for platform failed : %v", err)
        return nil, err
    }
    req.PlatformType = battery.PLATFORM_TYPE(platform).Enum()

    // optype, requested
    v, result = parameterSlice.Get(MaintenanceParameter_Optype)
    if !result {
        xylog.ErrorNoId("get optype from parameterSlice failed")
        return nil, xyerror.ErrBadInputData
    }
    optype, err = strconv.Atoi(v)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("strconv.Atoi for optype failed : %v", err)
        return nil, err
    }
    req.MaintenanceType = battery.MAINTENANCE_TYPE(optype).Enum()

    // propId, optional
    v, result = parameterSlice.Get(MaintenanceParameter_PROP_Id)
    if !result {
        xylog.DebugNoId("get prop id from parameterSlice failed")
    } else {
        propId, err = strconv.ParseUint(v, 10, 64)
        if err != xyerror.ErrOK {
            xylog.WarningNoId("strconv.ParseUint for uid failed : %v", err)
            err = xyerror.ErrOK
        } else {
            req.PropId = proto.Uint64(propId)
        }
    }

    //amount, optional
    v, result = parameterSlice.Get(MaintenanceParameter_Amount)
    if !result {
        xylog.DebugNoId("get prop amount from parameterSlice failed")
    } else {
        amount, err = strconv.Atoi(v)
        if err != xyerror.ErrOK {
            xylog.WarningNoId("strconv.Atoi for uid failed : %v", err)
            err = xyerror.ErrOK
        } else {
            req.Amount = proto.Int32(int32(amount))
        }
    }

    message = req
    return

}

// GetPropPbMsg 获取道具类pb消息内容
func GetPropPbMsg(parameterSlice ParameterSlice) (message proto.Message, err error) {
    var (
        propId                                     uint64
        propType, amount, platform, optype, source int
    )

    // identity ,requested
    req := &battery.MaintenancePropRequest{}
    v, result := parameterSlice.Get(MaintenanceParameter_Uid)
    if !result {
        xylog.ErrorNoId("get uid from parameterSlice failed")
        return nil, xyerror.ErrBadInputData
    }

    req.Identity = proto.String(v)

    // platform ,requested
    v, result = parameterSlice.Get(MaintenanceParameter_Platform)
    if !result {
        xylog.ErrorNoId("get platform from parameterSlice failed")
        return nil, xyerror.ErrBadInputData
    }

    platform, err = strconv.Atoi(v)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("strconv.Atoi for platform failed : %v", err)
        return nil, err
    }

    req.PlatformType = battery.PLATFORM_TYPE(platform).Enum()

    // platform ,requested
    v, result = parameterSlice.Get(MaintenanceParameter_Source)
    if !result {
        xylog.ErrorNoId("get source from parameterSlice failed")
        return nil, xyerror.ErrBadInputData
    }

    source, err = strconv.Atoi(v)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("strconv.Atoi for source failed : %v", err)
        return nil, err
    }

    req.Source = battery.ID_SOURCE(source).Enum()

    // optype, requested
    v, result = parameterSlice.Get(MaintenanceParameter_Optype)
    if !result {
        xylog.ErrorNoId("get optype from parameterSlice failed")
        return nil, xyerror.ErrBadInputData
    }

    optype, err = strconv.Atoi(v)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("strconv.Atoi for optype failed : %v", err)
        return nil, err
    }

    req.MaintenanceType = battery.MAINTENANCE_TYPE(optype).Enum()

    // prop type, optional
    v, result = parameterSlice.Get(MaintenanceParameter_PROP_Type)
    if !result {
        xylog.DebugNoId("get prop type from parameterSlice failed")
    } else {
        propType, err = strconv.Atoi(v)
        if err != xyerror.ErrOK {
            xylog.ErrorNoId("strconv.Atoi for uid failed : %v", err)
            return nil, err
        }

        req.PropType = battery.PropType(propType).Enum()

    }

    // propId, optional
    v, result = parameterSlice.Get(MaintenanceParameter_PROP_Id)
    if !result {
        xylog.DebugNoId("get prop id from parameterSlice failed")
    } else {
        propId, err = strconv.ParseUint(v, 10, 64)
        if err != xyerror.ErrOK {
            xylog.WarningNoId("strconv.ParseUint for uid failed : %v", err)
            err = xyerror.ErrOK
        } else {
            req.PropId = proto.Uint64(propId)
        }
    }

    //amount, optional
    v, result = parameterSlice.Get(MaintenanceParameter_Amount)
    if !result {
        xylog.DebugNoId("get prop amount from parameterSlice failed")
    } else {
        amount, err = strconv.Atoi(v)
        if err != xyerror.ErrOK {
            xylog.WarningNoId("strconv.Atoi for uid failed : %v", err)
            err = xyerror.ErrOK
        } else {
            req.Amount = proto.Int32(int32(amount))
        }
    }

    message = req

    return
}
