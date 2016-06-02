// cache_mailconfig
// 签到活动信息的缓存管理器定义

package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	//"code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	//"sync/atomic"
	//"time"
)

type MAPMailConfig map[int32]*battery.DBMailInfoConfig

type MailConfigCache struct {
	MailConfigs MAPMailConfig
}

func NewMailConfigCache() *MailConfigCache {
	return &MailConfigCache{}
}

type MailConfigCacheManager struct {
	cache [2]MailConfigCache
	xycache.CacheBase
}

//邮件配置缓存管理器
var DefMailConfigCacheManager = NewMailConfigCacheManager()

func NewMailConfigCacheManager() *MailConfigCacheManager {
	return &MailConfigCacheManager{}
}

//sm?! damn!~
func (sm *MailConfigCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	sm.Init()
	//加载资源配置信息
	failReason, err = DefMailConfigCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("MailConigResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return
}

func (sm *MailConfigCacheManager) Init() {
	//sm.index = 0
}

func (sm *MailConfigCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {

	//begin := time.Now()
	//defer xylog.Debug("MailConfigCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = sm.Load()

	return
}

func (sm *MailConfigCacheManager) MailConfig(id int32) *battery.DBMailInfoConfig {

	if p, ok := sm.cache[int(sm.Major())].MailConfigs[id]; ok {
		//找到
		return p
	}

	return nil
}

func (sm *MailConfigCacheManager) MailConfigs() *MAPMailConfig {
	return &(sm.cache[sm.Major()].MailConfigs)
}

func (sm *MailConfigCacheManager) SecondaryMailConfigs() *MAPMailConfig {
	return &(sm.cache[sm.Secondary()].MailConfigs)
}

func (sm *MailConfigCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = sm.loadMailConfigs()
	return
}

func (sm *MailConfigCacheManager) loadMailConfigs() (failReason battery.ErrorCode, err error) {
	mapMailConfigs := sm.SecondaryMailConfigs()
	*mapMailConfigs = make(MAPMailConfig, 0)

	//加载系统道具配置信息
	dbMailConfigs := make([]*battery.DBMailInfoConfig, 0)
	err = DefCacheDB.LoadMailConfigs(&dbMailConfigs)
	if err != nil || len(dbMailConfigs) <= 0 {
		failReason = xyerror.Resp_QueryMailConfigsFromDBError.GetCode()
		return
	}

	for _, dbMailConfig := range dbMailConfigs {
		(*mapMailConfigs)[dbMailConfig.GetMailID()] = dbMailConfig
	}

	sm.switchCache()

	//xylog.Debug("MailConfigs : %v ", *mapMailConfigs)

	return
}

func (sm *MailConfigCacheManager) switchCache() (fail_reason int32, err error) {
	sm.Switch()
	//xylog.Debug("now MailConfigs cache switch to %d", sm.Major())
	return
}
