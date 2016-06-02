// cache_checkpoint
// 记忆点信息的缓存管理器定义
package xybusinesscache

import (
    proto "code.google.com/p/goprotobuf/proto"
    "guanghuan.com/xiaoyao/common/cache"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    //"guanghuan.com/xiaoyao/superbman_server/server"
    "gopkg.in/mgo.v2/bson"
    "time"
)

//记忆点排行榜
type SliceCheckPoint []*battery.UserCheckPointDetail

func (s *SliceCheckPoint) Print(id uint32) {
    xylog.DebugNoId("--------- checkpoint(%d) ---------", id)
    for _, checkPoint := range *s {
        xylog.DebugNoId("(%v)", checkPoint)
    }
}

func (s *SliceCheckPoint) Clear() {
    *s = make(SliceCheckPoint, 0)
}

type CheckPointCache struct {
    id           uint32
    androidCache [2]SliceCheckPoint // 安卓排行版缓存
    iosCache     [2]SliceCheckPoint // ios 排行版缓存
    xycache.CacheBase
    //CacheDB
}

//重载对应记忆点的全局排行榜信息
// config *CheckPointConfig 记忆点配置信息
func (cpc *CheckPointCache) Reload() {
    err := cpc.load(cpc.id, &(cpc.androidCache[cpc.Secondary()]), battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID)
    if xyerror.ErrOK == err { //加载成功的情况下才做切换
        cpc.Switch()
    }
    // ios 没需要先不加载
    // err = cpc.load(cpc.id, &(cpc.iosCache[cpc.Secondary()]))
    // if xyerror.ErrOK == err { //加载成功的情况下才做切换
    // 	cpc.Switch()
    // }
}

func isUserInBlackList(uid string) bool {
    for _, list := range DefRankBlacklistManager.List() {
        if list.Uid == uid {
            xylog.WarningNoId("uid :%v in blacklist", uid)
            return true
        }
    }
    return false
}

//加载最新的全局排行榜信息
// checkPointId uint32 记忆点id
// checkPoints *SliceCheckPoint 查询到的记忆点排行榜列表
func (cpc *CheckPointCache) load(checkPointId uint32, checkPoints *SliceCheckPoint, platform battery.PLATFORM_TYPE) (err error) {
    //将原来的信息清除
    checkPoints.Clear()
    dbCheckPoints := make([]*battery.DBUserCheckPoint, 0)
    tempdbCheckPoints := make([]*battery.DBUserCheckPoint, 0)
    err = DefCacheDB.LoadCheckPoint(DefCheckPointConfig.GlobalRankSize, checkPointId, &tempdbCheckPoints, platform)
    for _, dbCheckPoint := range tempdbCheckPoints {
        if !isUserInBlackList(dbCheckPoint.GetUid()) {
            dbCheckPoints = append(dbCheckPoints, dbCheckPoint)
        }
    }
    if xyerror.ErrOK == err && len(dbCheckPoints) > 0 {

        uids := make([]string, 0)
        for _, dbCheckPoint := range dbCheckPoints {

            uids = append(uids, dbCheckPoint.GetUid())
        }

        //查询排行榜所有玩家的tpid信息
        tpids := make([]*battery.IDMap, 0)
        selector := bson.M{"sid": 1, "gid": 1, "note": 1, "iconurl": 1}
        err = DefCacheDB.DB().QueryTpidsByUids(selector, uids, &tpids)
        if err != xyerror.ErrOK {
            xylog.ErrorNoId("QueryTpidsByUids for %v failed : %v", uids, err)
            return
        } else if len(tpids) <= 0 {
            //	xylog.WarningNoId("QueryTpidsByUids for %v len(tpids) (%d) <= 0", uids, len(tpids))
            return
        }

        mapUid2Sid := make(map[string]string, 0)
        mapUid2IconUrl := make(map[string]string, 0)
        mapUid2Note := make(map[string]string, 0)
        for _, tpid := range tpids {
            mapUid2Sid[tpid.GetGid()] = tpid.GetSid()
            mapUid2IconUrl[tpid.GetGid()] = tpid.GetIconUrl()
            mapUid2Note[tpid.GetGid()] = tpid.GetNote()
        }

        //查询排行榜所有玩家的account信息
        //dbAccounts := make([]*battery.DBAccount, 0)
        //err = xybusiness.QueryAccountsByUids(DefCacheDB.db, uids, &dbAccounts)
        //if err != xyerror.ErrOK || len(dbAccounts) <= 0 {
        //	xylog.Error("QueryAccountsByUids for %v failed : %v", uids, err)
        //	return
        //}
        //xylog.Debug("QueryAccountsByUids(%v) length : %d", uids, len(dbAccounts))

        for _, dbCheckPoint := range dbCheckPoints {
            checkpoint := GetCheckPointFromDBUserCheckPoint(dbCheckPoint)
            checkpoint.CheckPointId = nil //不需要保存记忆点id

            uidTmp := checkpoint.GetSid()
            //将玩家的uid指定为玩家的sid
            if sid, ok := mapUid2Sid[uidTmp]; ok {
                checkpoint.Sid = proto.String(sid)
            }
            //设置iconurl
            if iconUrl, ok := mapUid2IconUrl[uidTmp]; ok {
                checkpoint.IconUrl = proto.String(iconUrl)
            }
            if note, ok := mapUid2Note[uidTmp]; ok {
                checkpoint.Name = proto.String(note)
            }
            *checkPoints = append(*checkPoints, checkpoint)
        }
    }
    //checkPoints.Print(checkPointId)
    return
}

//根据数据库的记忆点信息获取缓存记忆点信息
// dbUserCheckPoint *battery.DBUserCheckPoint 数据库中的玩家checkpoint信息
// needUid bool 是否需要uid，在一些请求中，不需要返回uid，这样可以减小通信包的大小
//return:
// detail *battery.UserCheckPointDetail 返回的记忆点详细信息
func GetCheckPointFromDBUserCheckPoint(dbUserCheckPoint *battery.DBUserCheckPoint) (detail *battery.UserCheckPointDetail) {

    detail = &battery.UserCheckPointDetail{
        Sid:          proto.String(dbUserCheckPoint.GetUid()),
        CheckPointId: proto.Uint32(dbUserCheckPoint.GetCheckPointId()),
        Score:        proto.Uint64(dbUserCheckPoint.GetScore()),
        Charge:       proto.Uint64(dbUserCheckPoint.GetCharge()),
        Coin:         proto.Uint32(dbUserCheckPoint.GetCoin()),
        RoleId:       proto.Uint64(getRoleByRoleID(dbUserCheckPoint.GetRoleId())),
        Grade:        proto.Uint32(dbUserCheckPoint.GetGrade()),
        Collections:  dbUserCheckPoint.Collections,
    }

    return
}

func getRoleByRoleID(id uint64) (roleid uint64) {
    return id / 10000 * 10000
}

//创建新的记忆点缓存
func NewCheckPointCache(id uint32) (cpc *CheckPointCache) {
    cpc = &CheckPointCache{
        id:           id,                                                                     //设置缓存对应的记忆点id
        androidCache: [2]SliceCheckPoint{make(SliceCheckPoint, 0), make(SliceCheckPoint, 0)}, //初始化缓存为空
    }
    //cpc.index = 0 //默认主缓存为0
    return
}

type MAPCheckPointRank map[uint32]*CheckPointCache

func (m *MAPCheckPointRank) Print() {
    for id, checkPoint := range *m {
        for i, cache := range checkPoint.androidCache {
            xylog.DebugNoId("cache(%d)", i)
            cache.Print(id)
        }
    }
}

func (m *MAPCheckPointRank) Clear() {
    *m = make(MAPCheckPointRank, 0)
}

type CheckPointManager struct {
    checkPoints MAPCheckPointRank
}

var DefCheckPointGlobalRankManager CheckPointManager

type CheckPointConfig struct {
    IdNum                uint32 //记忆点个数
    GlobalRankReLoadSecs int64  //重载时间（单位秒）
    GlobalRankSize       int    //全局排行榜列表长度
}

var DefCheckPointConfig *CheckPointConfig

//进程启动时调用的初始化函数
// db *xydb.XYDB 数据库操作类指针
// config *CheckPointConfig 配置信息
func (cpm *CheckPointManager) InitWhileStart(config *CheckPointConfig) (failReason battery.ErrorCode, err error) {
    DefCheckPointConfig = config
    xylog.DebugNoId("DefCheckPointConfig : %v", DefCheckPointConfig)
    cpm.init()
    return
}

//初始化记忆点管理器：
//初始化记忆点缓存列表
// db *xydb.XYDB 数据库操作指针
func (cpm *CheckPointManager) init() {
    //保存记忆点对应的缓存信息
    var i uint32
    cpm.checkPoints = make(MAPCheckPointRank, 0)
    for i = 0; i < DefCheckPointConfig.IdNum; i++ {
        cpm.checkPoints[i] = NewCheckPointCache(i) //每个记忆点创建一个记忆点缓存
    }

    go cpm.ReloadAllTick() //启动一个goruntine加载所有的记忆点全局排行榜
}

//循环加载记忆点全局排行榜信息的函数，由go runtine调用
//func (cpm *CheckPointManager) ReloadFunc(checkPointCache *CheckPointCache) {

//checkPointCache.Reload()

//timer := time.NewTicker(time.Second * time.Duration(DefCheckPointConfig.GlobalRankReLoadSecs))
//for {
//	select {
//	case <-timer.C:
//		checkPointCache.Reload()
//	}
//}
//}

func (cpm *CheckPointManager) ReloadAllTick() {

    cpm.ReloadAll()

    timer := time.NewTicker(time.Second * time.Duration(DefCheckPointConfig.GlobalRankReLoadSecs))

    for {
        select {
        case <-timer.C: //用定时任务进行所有记忆点排行榜的刷新
            cpm.ReloadAll()
        }
    }
}

func (cpm *CheckPointManager) ReloadAll() {
    for i := uint32(0); i < DefCheckPointConfig.IdNum; i++ {
        cpm.checkPoints[i].Reload()
    }
}

//根据记忆点id获取其对应的全局排行榜信息
// checkPointId uint32 记忆点id
//return:
// *SliceCheckPoint 找到返回对应的排行榜列表指针，没找到则返回nil
func (cpm *CheckPointManager) GlobalRank(checkPointId uint32, platform battery.PLATFORM_TYPE) *SliceCheckPoint {
    switch platform {
    case battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID:
        if checkPointCache, ok := cpm.checkPoints[checkPointId]; ok {
            return &(checkPointCache.androidCache[checkPointCache.Major()])
        }
    case battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS:
        if checkPointCache, ok := cpm.checkPoints[checkPointId]; ok {
            return &(checkPointCache.iosCache[checkPointCache.Major()])
        }

    }

    return nil
}
