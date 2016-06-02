package batteryapi

import (
    "code.google.com/p/goprotobuf/proto"
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "guanghuan.com/xiaoyao/common/cache"
    "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/util"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    //"guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

const RELOAD_BANNEDUSERS_SUBJECT = "reloadbanneduser"

//运营操作（道具）
func (api *XYAPI) OperationMaintenanceProp(req *battery.MaintenancePropRequest, resp *battery.MaintenancePropResponse) (err error) {
    var (
        identity = req.GetIdentity()
    )

    //获取请求的终端平台类型
    api.SetDB(req.GetPlatformType())

    var (
        uid    string
        optype = req.GetMaintenanceType()

        timestamp = req.GetTimestamp()
        message   string
    )

    resp.Error = &battery.Error{}
    *(resp.Error) = *(xyerror.Resp_NoError)

    //查找玩家信息
    if optype == battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_CDKEY_EXCHANGE {
        uid = req.GetUid()
    } else {
        uid, err = api.GetUid(identity)
        if err != xyerror.ErrOK {
            resp.Error.Code = battery.ErrorCode_GetAccountByUidError.Enum()
            resp.Error.Desc = proto.String(api.GetRespJson("", fmt.Sprintf("get uid for identity(%d) failed : %v", identity, err)))
            goto ErrHandle
        }
    }

    xylog.Debug(uid, "req:%s", req)
    //查找tpid信息
    {
        selector, tpid := bson.M{"sid": 1, "note": 1}, &battery.IDMap{}
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).XYBusinessDB.QueryTpidByUid(selector, uid, tpid, mgo.Strong)
        if err != xyerror.ErrOK {
            resp.Error.Code = battery.ErrorCode_GetAccountByUidError.Enum()
            resp.Error.Desc = proto.String(api.GetRespJson("", fmt.Sprintf("get tpid for identity(%s) failed : %v", identity, err)))
            goto ErrHandle
        }

        message += fmt.Sprintf("uid\":\"%s\",\"sid\":\"%v\",\"name\":\"%v\"", uid, tpid.GetSid(), tpid.GetNote())
    }

    xylog.Debug(uid, "optype : %v, message : %s", optype, message)

    switch optype {
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_PROP_ADD, battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_CDKEY_EXCHANGE: //分发物品
        req.Uid = &uid
        var failReason battery.ErrorCode
        failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_MaintenanceProp, req, resp)
        if failReason != battery.ErrorCode_NoError {
            resp.Error = xyerror.ConstructError(failReason)
        }
        resp.Error.Desc = proto.String(api.GetRespJson(message, resp.Error.GetDesc()))
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_QUERY_USER_INFO: //如果只是查询玩家信息，直接返回，不记录日志
        resp.Error.Code = battery.ErrorCode_NoError.Enum()
        resp.Error.Desc = proto.String(api.GetRespJson(message, resp.Error.GetDesc()))
        return
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_BAN_USER,
        battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_UNBAN_USER:
        api.BanUser(uid, timestamp, optype, resp.Error)
        resp.Error.Desc = proto.String(api.GetRespJson(message, resp.Error.GetDesc()))
        return
    }

ErrHandle:

    l := &battery.MaintenanceLog{
        Identity:  &identity,
        Error:     resp.Error,
        Timestamp: proto.Int64(xyutil.CurTimeSec()),
        Detail:    proto.String(req.String()),
    }
    go api.AddMaintenanceLog(l)

    return

}

//对玩家进行封号或者解封处理
// uid string 玩家标识
// timestamp int64 解封时间戳
// optype battery.MAINTENANCE_TYPE 操作类型
// errStruct *battery.Error 错误信息
func (api *XYAPI) BanUser(uid string, timestamp int64, optype battery.MAINTENANCE_TYPE, errStruct *battery.Error) {
    xylog.WarningNoId("BanUser(%d) optype(%v)", uid, optype)
    switch optype {
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_BAN_USER: //封号
        api.banUser(uid, -1, errStruct) //默认指定时间戳为-1，表示永久封号。后续如果要打开定制解封时间，需要修改这个地方的时间戳参数。
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_UNBAN_USER: //解封
        api.unbanUser(uid, errStruct)
    default:
        xylog.ErrorNoId("BanUser unkown optype(%v)", optype)
    }

    xylog.WarningNoId("BanUser(%s) end", uid, optype)
}

//封号函数
// uid string 玩家标识
// timestamp int64 解封时间戳
// errStruct *battery.Error 错误返回信息
func (api *XYAPI) banUser(uid string, timestamp int64, errStruct *battery.Error) {

    //将玩家封号信息插入数据库中
    err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_BANNEDUSER).UpsertBannedUser(uid, timestamp)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("UpsertBannedUser failed : %v", err)
        errStruct.Code = battery.ErrorCode_UpsertBannedUserDBError.Enum()
        errStruct.Desc = proto.String(fmt.Sprintf("UpsertBannedUserDBError uid(%s) timestamp(%d) failed : %v", uid, timestamp, err))
        return
    } else {
        xylog.DebugNoId("UpsertBannedUser(%s) succeed", uid)
    }

    //广播通知所有节点更新黑名单信息
    api.banUserNatsPublish(uid, battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_BAN_USER)
}

//解封函数
func (api *XYAPI) unbanUser(uid string, errStruct *battery.Error) {

    //玩家封号信息从数据库中删除
    err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_BANNEDUSER).RemoveBannedUser(uid)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_RemoveBannedUsersDBError.Enum()
        errStruct.Desc = proto.String(fmt.Sprintf("RemoveBannedUserDBError uid(%s) failed : %v", uid, err))
        return
    }

    //广播通知所有节点更新黑名单信息
    api.banUserNatsPublish(uid, battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_UNBAN_USER)
}

func (api *XYAPI) banUserNatsPublish(uid string, maintenanceType battery.MAINTENANCE_TYPE) {
    banUnit := &battery.MaintenanceBanUnit{
        Uid:             &uid,
        MaintenanceType: maintenanceType.Enum(),
    }

    data, err := proto.Marshal(banUnit)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("proto.Marshal for banUnit(%s) failed : %v", uid, err)
    } else {
        xylog.DebugNoId("NatsPublish for banUnit : %v", banUnit)
        api.NatsPublish(RELOAD_BANNEDUSERS_SUBJECT, data)
    }
}

//查询封号玩家信息
//return:
// uids []string 封号玩家uid列表
// err error 操作错误
func (api *XYAPI) QueryBannedUsers() (uids []string, err error) {
    uids, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_BANNEDUSER).QueryBannedUsers(xyutil.CurTimeSec())
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("QueryBannedUsers failed : %v", err)
        return
    }

    return
}

//拼装返回字符串
func (api *XYAPI) GetRespJson(message, desc string) (result string) {
    head := true
    result = "{"

    if message != "" {
        if !head { //如果不是第一个元素，需要加上","
            result += ","
        } else {
            head = false
        }

        result += message

    }

    if desc != "" {
        if !head { //如果不是第一个元素，需要加上","
            result += ","
        } else {
            head = false
        }

        result += "\"desc\":\"" + desc + "\""
    }

    result += "}"

    return
}

func (api *XYAPI) AddMaintenanceLog(l *battery.MaintenanceLog) error {
    return api.GetLogDB().AddMaintenanceLog(l)
}

//黑名单玩家管理器定义
var DefBannedUserManager = NewBannedUserManager()

type MapBannedUser map[string]empty

type BannedUserCache struct {
    uids MapBannedUser
}

func (b *BannedUserCache) Print() {
    xylog.DebugNoId("==================== Banned Users begin ==================== ")
    for k, _ := range b.uids {
        xylog.DebugNoId(k)
    }
    xylog.DebugNoId("==================== Banned Users end ==================== ")
}

func NewBannedUserCache() *BannedUserCache {
    return &BannedUserCache{
        uids: make(MapBannedUser, 0),
    }
}

type BannedUserManager struct {
    caches [2]*BannedUserCache
    xycache.CacheBase
}

//加载封号玩家信息
func (manager *BannedUserManager) Load() (err error) {
    var uids []string
    uids, err = NewXYAPI().QueryBannedUsers()

    xylog.DebugNoId("BannedUsers : %v", uids)

    if err == xyerror.ErrOK {
        secondary := manager.caches[int(manager.Secondary())]
        secondary.uids = make(MapBannedUser, 0)
        for _, uid := range uids {
            secondary.uids[uid] = empty{}
        }
        //secondary.Print()

        manager.Switch()

        manager.caches[int(manager.Major())].Print()
        manager.caches[int(manager.Secondary())].Print()
    }

    //manager.Print()

    return
}

//玩家是否被封号
// uid string 玩家标识
//returns:
// true  是
// false 否
func (manager *BannedUserManager) IsUidBanned(uid string) bool {
    major := manager.caches[int(manager.Major())]

    if _, ok := major.uids[uid]; ok {
        return true
    }
    return false
}

func (manager *BannedUserManager) ProcessBanUnit(banUnit *battery.MaintenanceBanUnit) {
    maintenanceType := banUnit.GetMaintenanceType()
    switch maintenanceType {
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_BAN_USER:
        manager.BanUser(banUnit.GetUid())
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_UNBAN_USER:
        manager.UnbanUser(banUnit.GetUid())
    default:
        xylog.ErrorNoId("error maintenance type (%v)", maintenanceType)
    }
}

//对玩家进行封号
// uid string 玩家标识
func (manager *BannedUserManager) BanUser(uid string) {
    major, secondary := manager.caches[int(manager.Major())], manager.caches[int(manager.Secondary())]
    secondary.uids = major.uids
    secondary.uids[uid] = empty{} //把uid加入到封号列表中
    manager.Switch()
}

//对玩家进行解封
// uid string 玩家标识
func (manager *BannedUserManager) UnbanUser(uid string) {
    major, secondary := manager.caches[int(manager.Major())], manager.caches[int(manager.Secondary())]
    secondary.uids = major.uids
    delete(secondary.uids, uid) //把uid从封号列表中删除
    manager.Switch()
}

func (manager *BannedUserManager) Print() {
    manager.caches[manager.Major()].Print()
}

func NewBannedUserManager() (manager *BannedUserManager) {
    manager = new(BannedUserManager)

    for i := 0; i < 2; i++ {
        manager.caches[i] = NewBannedUserCache()
    }

    return
}
