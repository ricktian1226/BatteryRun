// cache_announcement
// 公告信息的缓存管理器定义
package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"time"
)

type MAPAnnouncements map[uint64]battery.Announcement

func (m *MAPAnnouncements) Print() {
	for id, item := range *m {
		xylog.DebugNoId("--announcement(%d)--", id)
		xylog.DebugNoId("%v", &item)
	}
}

func (m *MAPAnnouncements) Clear() {
	*m = make(MAPAnnouncements, 0)
}

type AnnouncementsCache struct {
	announcements MAPAnnouncements //按照公告id索引的map
}

func NewAnnouncementsCache() *AnnouncementsCache {
	return &AnnouncementsCache{
		announcements: make(MAPAnnouncements, 0),
	}
}

type AnnouncementsCacheManager struct {
	cache               [2]AnnouncementsCache
	lastUpdateTimestamp int64 //上次刷新时间
	//CacheDB
	xycache.CacheBase
}

//道具缓存管理器
var DefAnnouncementsCacheManager = NewAnnoucenmentsCacheManager()

func NewAnnoucenmentsCacheManager() *AnnouncementsCacheManager {
	return &AnnouncementsCacheManager{}
}

func (am *AnnouncementsCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	//初始化数据库操作指针
	am.Init()
	//加载资源配置信息
	failReason, err = DefAnnouncementsCacheManager.ReLoad(true)
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("Announcements ResLoad failed : %v ", err)
		return
	}

	//启动一个单独的gorountine来定时加载当天的公告信息资源
	go am.CycleReload(time.Duration(10))

	return
}

//公告信息跨天的时候才需要加载，除非是修改配置强制加载。通过定时任务，10分钟一次进行加载
func (am *AnnouncementsCacheManager) CycleReload(interval time.Duration) {
	timer := time.NewTicker(time.Minute * interval)
	for {
		select {
		case <-timer.C:
			DefAnnouncementsCacheManager.ReLoad(false)
		}
	}

	return
}

func (am *AnnouncementsCacheManager) Init() {
	//am.index = 0
	am.lastUpdateTimestamp = 0
}

func (am *AnnouncementsCacheManager) ReLoad(force bool) (failReason battery.ErrorCode, err error) {

	reload, now := false, time.Now().Unix()

	if force {
		reload = true
	} else {
		if xyutil.DayDiff(am.lastUpdateTimestamp, now) > 0 {
			reload = true
		}
	}

	if reload {
		failReason, err = am.Load()
		if failReason == battery.ErrorCode_NoError && err == nil {
			am.lastUpdateTimestamp = now
		}
	}

	return
}

//获取某个公告的详细信息
// id uint64 公告id
func (am *AnnouncementsCacheManager) Announcement(id uint64) *battery.Announcement {
	if a, ok := am.cache[am.Major()].announcements[id]; ok {
		//找到
		return &a
	}
	return nil
}

func (am *AnnouncementsCacheManager) Announcements() *MAPAnnouncements {
	return &(am.cache[am.Major()].announcements)
}

//获取当前的公告列表
func (am *AnnouncementsCacheManager) AnnouncementsSlice() (announcements []battery.Announcement) {
	for _, announcement := range am.cache[am.Major()].announcements {
		announcements = append(announcements, announcement)
	}
	return
}

func (am *AnnouncementsCacheManager) SecondaryAnnouncements() *MAPAnnouncements {
	return &(am.cache[am.Secondary()].announcements)
}

func (am *AnnouncementsCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = am.loadAnnouncements()
	return
}

//从数据库加载当天的公告信息，转换为消息用公告信息
func (am *AnnouncementsCacheManager) loadAnnouncements() (failReason battery.ErrorCode, err error) {
	dbAnnouncements := make([]*battery.DBAnnouncementConfig, 0)
	err = DefCacheDB.LoadAnnouncements(&dbAnnouncements)
	if err != nil || len(dbAnnouncements) <= 0 {
		failReason = battery.ErrorCode_QueryAnnouncementConfigFromDBError
		return
	}

	mapAnnouncements := am.SecondaryAnnouncements()
	mapAnnouncements.Clear()

	for _, a := range dbAnnouncements {
		(*mapAnnouncements)[a.GetId()] = am.getAnnouncementFromDBAnnouncementConfig(a)
	}

	am.switchCache()

	return
}

//将数据库用公告信息转换为消息用公告信息
func (am *AnnouncementsCacheManager) getAnnouncementFromDBAnnouncementConfig(dbAnnouncementConfig *battery.DBAnnouncementConfig) (announcement battery.Announcement) {
	announcement.Id = dbAnnouncementConfig.Id
	announcement.Title = dbAnnouncementConfig.Title
	announcement.Content = dbAnnouncementConfig.Description
	announcement.BeginTime = dbAnnouncementConfig.BeginTime
	announcement.EndTime = dbAnnouncementConfig.EndTime
	announcement.State = battery.ANNOUNCEMENT_STATE_ANNOUNCEMENT_STATE_VALID.Enum()
	return
}

func (am *AnnouncementsCacheManager) switchCache() (fail_reason int32, err error) {
	am.Switch()
	xylog.DebugNoId("now Annoucements cache switch to %d", am.Major())
	return
}
