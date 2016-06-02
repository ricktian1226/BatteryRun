package main

import (
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "sync"
    //"reflect"
    "time"

    proto "code.google.com/p/goprotobuf/proto"
    martini "github.com/codegangsta/martini"
    nats "github.com/nats-io/nats"
    xylog "guanghuan.com/xiaoyao/common/log"
    xypanic "guanghuan.com/xiaoyao/common/panic"
    xyperf "guanghuan.com/xiaoyao/common/performance"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xycache "guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/error"
)

//func HttpGetWorker(w http.ResponseWriter, r *http.Request, params martini.Params) (status int, resp_data string) {
//	status = http.StatusOK
//	return
//}

//重载配置信息统一接口
// resType string 加载的资源类型字符串
func HttpConfigReload(resType string) (status int, resp string) {
    err := nats_service.Publish(resType, []byte(""))
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("Publish %s message failed: %s", resType, err.Error())
        resp = err.Error()
    } else {
        resp = resType + " ConfigReload OK"
    }
    status = http.StatusOK
    return
}

//广播加载商品配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpGoodsConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_GOODS)
    return
}

//广播加载商品配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpBeforeGameConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_BEFOREGAME)
    return
}

//广播加载抽奖配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpLottoConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_LOTTO)
    return
}

//广播加载道具配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpPropConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_PROP)
    return
}

//广播加载符文配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpRuneConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_RUNE)
    return
}

//广播加载任务配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpMissionConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_MISSION)
    return
}

//广播加载商品配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpPickUpConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_PICKUP)
    return
}

//广播加载角色配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpRoleInfoConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_ROLE_INFO)
    return
}

//广播加载角色加成配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpRoleLevelBonusConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_ROLE_LEVEL_BONUS)
    return
}

//广播加载公告配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpAnnouncementConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_ANNOUNCEMENT)
    return
}

//广播加载广告配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpAdvertisementConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_ADVERTISEMENT)
    return
}

//广播加载tip配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpTipConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_TIP)
    return
}

//广播加载tip配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpNewAccountPropConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_NEWACCOUNTPROP)
    return
}

//广播加载所有资源配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpAllConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_ALL)
    return
}

//广播加载签到配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpSigninConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_SIGNIN)
    return
}

//广播加载系统配置信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpApiConfigReload() (status int, resp string) {
    status, resp = HttpConfigReload(xycache.RES_RELOAD_API)
    return
}

//广播加载debuguser信息
//return:
// status int 操作状态
// resp string 返回前端处理结果
func HttpDebugUsersReload() (status int, resp string) {
    status, resp = HttpConfigReload(xylog.DEBUGUSER_SUBJECT)
    return
}

func HttpConfigProfile(params martini.Params) (status int, resp string) {
    var isValid = true
    var op string = params["op"]
    xylog.InfoNoId("Profiling, op=%s", op)
    switch op {
    case "start":
        pm.Start()
        resp = fmt.Sprintf("Profiler started : %s", pm.StartTime.String())
    case "stop":
        pm.Stop()
        resp = fmt.Sprintf("Profiler Stopped : %s", pm.StopTime.String())
    case "reset":
        pm.Reset()
        resp = fmt.Sprintf("Profiler Reset : %s", time.Now().String())
    default:
        xylog.InfoNoId("Profiler Result \n" + pm.String())
        resp = pm.String()
    }
    if isValid {
        nats_service.Publish("profile", []byte(op))
    }

    status = http.StatusOK
    return
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

func HttpPostWorkerNoToken(w http.ResponseWriter, r *http.Request, params martini.Params) (status int, resp string) {
    var (
        uri string = r.RequestURI
        err error
        //resp_data []byte
        respData = bytePool.Get().(*[]byte)
    )

    defer bytePool.Put(respData)

    status = http.StatusOK
    xylog.DebugNoId("uri=%s, token=n/a, user agent=%s", uri, r.UserAgent())
    err = ProcessHttpMsg(uri, "", r, respData)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[%s] failed: %s", uri, err.Error())
        status = http.StatusInternalServerError //处理失败，返回服务端错误
    } else {
        resp = string(*respData)
    }
    return
}

var bytePool = &sync.Pool{
    New: func() interface{} {
        s := make([]byte, 0)
        return &s
    },
}

func HttpPostWorker(w http.ResponseWriter, r *http.Request, params martini.Params) (status int, resp string) {
    var (
        uri   string = r.RequestURI
        token string
        err   error
        //resp_data []byte
        respData = bytePool.Get().(*[]byte)
    )

    defer bytePool.Put(respData)

    status = http.StatusOK
    uri, token = ProcessUri(uri)
    xylog.DebugNoId("uri=%s, token=%s, user agent=%s", uri, token, r.UserAgent())
    err = ProcessHttpMsg(uri, token, r, respData)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[%s] failed: %s", uri, err.Error())
        status = http.StatusInternalServerError //处理失败，返回服务端错误
    } else {
        resp = string(*respData)
    }
    return
}

func ProcessHttpMsg(uri string, token string, r *http.Request, respData *[]byte) (err error) {

    begin := time.Now()
    defer xyperf.Trace(xyperf.DefLogId, &begin)
    defer xypanic.Crash()

    var (
        reqData = bytePool.Get().(*[]byte)
        subj    string
        reply   *nats.Msg
        route   *HttpPostToNatsRoute
    )
    if token != "" {
        route = DefHttpPostTable.GetRoutePath(uri)
    } else {
        route = DefHttpPostNoTokenTable.GetRoutePath(uri)
    }

    if route == nil {
        err = errors.New("No route for uri:" + uri)
        return
    }
    subj = route.GetNatsSubject()
    if subj == "" {
        err = errors.New("No subject for uri:" + uri)
        return
    }

    // 登录操作记录ip
    // 协议定义缺欠，，需要解序列pb

    err = ReadRequestData(r, reqData)
    if err != nil {
        return
    }

    if subj == "login" {
        loginIP := r.RemoteAddr
        req_type := &battery.LoginRequest{}

        err = Before(*reqData, req_type)
        if err != nil {
            xylog.ErrorNoId("error pre-Processing :%s", err.Error())
            return
        }
        req_type.LoginIp = proto.String(loginIP)
        *reqData, err = After(req_type)
        if err != nil {
            xylog.ErrorNoId("Error after_Processing :%S", err.Error())
            return
        }

    }
    xylog.DebugNoId("forward request to %s", subj)

    var (
        start   time.Time
        dur     time.Duration
        job_ok  bool
        in_len  int
        out_len int
    )

    start = time.Now()
    reply, err = nats_service.Request(subj, *reqData, time.Duration(DefConfig.NatsTimeout)*time.Second)

    dur = time.Since(start)
    in_len = len(*reqData)

    if err != nil {
        xylog.ErrorNoId("<%s> Error: %s", subj, err.Error())
        goto ErrorHandle
    } else {
        if reply != nil {
            *respData = reply.Data
            job_ok = true
            out_len = len(*respData)
        } else {
            err = errors.New("no reply data")
        }
    }
ErrorHandle:
    pm.AddJobResult(uri, job_ok, dur, int64(out_len), int64(in_len))
    return
}

//func UriToSubject(uri string) (subj string) {
//	return "echo"
//}

func ReadRequestData(r *http.Request, data *[]byte) (err error) {
    if r.ContentLength > DefConfig.MaxRequestSize {
        xylog.ErrorNoId("Request is too large: %d > %d", r.ContentLength, DefConfig.MaxRequestSize)
        err = errors.New("Request is too large")
    } else if r.ContentLength > 0 {
        *data, err = ioutil.ReadAll(r.Body)
        defer r.Body.Close()
        xylog.DebugNoId("expect: %d, read: %d", r.ContentLength, len(*data))
        if err != nil {
            xylog.ErrorNoId("Error Reading request data: %s", err.Error())
        }
    } else {
        xylog.ErrorNoId("no conent length in request")
        err = errors.New("Needs Content length")
    }
    return
}
