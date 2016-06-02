package xybusinesscache

import (
    "time"

    "guanghuan.com/xiaoyao/common/cache"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

// 排行榜玩家列表
type GlobalRankList []string

type GlobalRankListManager struct {
    cache [2]GlobalRankList
    xycache.CacheBase
}

type RankListConfig struct {
    GlobalRankReLoadSecs int64 // 重载时间（单位秒）
    GlobalRankSize       int   // 全局排行榜列表长度
}

var DefRankListConfig *RankListConfig

func (g *GlobalRankListManager) Print() {

}

func (g *GlobalRankList) Clear() {
    *g = make(GlobalRankList, 0)
}

var DefGlobalRankListManager GlobalRankListManager

func (g *GlobalRankListManager) InitWhileStart(config *RankListConfig) (failReason battery.ErrorCode, err error) {
    DefRankListConfig = config
    xylog.DebugNoId("ranklist config :%v", DefRankListConfig)
    g.Init()
    return
}

func (g *GlobalRankListManager) Init() {
    go g.ReloadTick()
}

func (g *GlobalRankListManager) ReloadTick() {
    g.Reload()
    timer := time.NewTicker(time.Second * time.Duration(DefRankListConfig.GlobalRankReLoadSecs))
    for {
        select {
        case <-timer.C:
            g.Reload()
        }
    }
}

func (g *GlobalRankListManager) Reload() (err error) {
    err = g.load(&g.cache[g.Secondary()])
    if err == xyerror.ErrOK {
        g.Switch()
    }
    return
}

func (g *GlobalRankListManager) load(ranklist *GlobalRankList) (err error) {
    ranklist.Clear()
    dbranklists := make([]*battery.RankLisk, 0)
    err = DefCacheDB.LoadRankList(DefRankListConfig.GlobalRankSize, &dbranklists)

    if err == xyerror.ErrOK && len(dbranklists) > 0 {
        for _, dbranklist := range dbranklists {
            uid := dbranklist.GetUid()
            if !isUserInBlackList(uid) {
                *ranklist = append(*ranklist, uid)
            }
        }
    }
    return
}

func (g *GlobalRankListManager) GlobalRank() GlobalRankList {
    return g.cache[g.Major()]
}
