package xybusinesscache

import (
    "guanghuan.com/xiaoyao/common/cache"
    xylog "guanghuan.com/xiaoyao/common/log"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

type BlacklistCache struct {
    list []BlackList
}
type BlackList struct {
    Uid string
}

func (b *BlacklistCache) Print() {
    xylog.DebugNoId("------------ranklist blacklist----------")
    for _, uid := range b.list {
        xylog.DebugNoId("blacklist :%s", uid)
    }
}

func (b *BlacklistCache) Clear() {
    b.list = make([]BlackList, 0)
}

// 排行榜黑名单
type BlacklistManager struct {
    cache [2]BlacklistCache
    xycache.CacheBase
}

var DefRankBlacklistManager = NewBlacklistCahceManager()

func NewBlacklistCahceManager() *BlacklistManager {
    return &BlacklistManager{}

}
func (b *BlacklistManager) InitWhileStart() {
    b.Init()
    b.Reload()
}

func (b *BlacklistManager) Init() {

}

func (b *BlacklistManager) Reload() (err error) {
    secondary := b.SecondaryCache()
    list := make([]BlackList, 0)
    err = DefCacheDB.LoadBlacklist(&list)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("load blacklist fail :%v", err)
    }
    secondary.Clear()
    secondary.list = list
    b.Switch()
    b.MajorCache().Print()
    return
}

func (b *BlacklistManager) MajorCache() *BlacklistCache {
    return &(b.cache[b.Major()])
}

func (b *BlacklistManager) SecondaryCache() *BlacklistCache {
    return &(b.cache[b.Secondary()])
}

func (b *BlacklistManager) List() (list []BlackList) {
    major := b.MajorCache()
    list = major.list
    return
}
