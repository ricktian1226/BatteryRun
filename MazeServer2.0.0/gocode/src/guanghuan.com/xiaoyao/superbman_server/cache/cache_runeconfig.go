// cache_runeconfig
// 签到活动信息的缓存管理器定义

package xybusinesscache

import (
	"code.google.com/p/goprotobuf/proto"
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	//"sync/atomic"
	//"time"
)

type MAPRuneConfig map[uint64]*battery.RuneConfig

type RuneConfigCache struct {
	RuneConfigs MAPRuneConfig
}

func NewRuneConfigCache() *RuneConfigCache {
	return &RuneConfigCache{}
}

type RuneConfigCacheManager struct {
	cache [2]RuneConfigCache
	xycache.CacheBase
}

//系统道具配置缓存管理器
var DefRuneConfigCacheManager = NewRuneConfigCacheManager()

func NewRuneConfigCacheManager() *RuneConfigCacheManager {
	return &RuneConfigCacheManager{}
}

//sm?! damn!~
func (sm *RuneConfigCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	sm.Init()
	//加载资源配置信息
	failReason, err = DefRuneConfigCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("RuneConigResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return
}

func (sm *RuneConfigCacheManager) Init() {
	//sm.index = 0
}

func (sm *RuneConfigCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {

	//begin := time.Now()
	//defer xylog.Debug("RuneConfigCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = sm.Load()

	return
}

func (sm *RuneConfigCacheManager) RuneConfig(id uint64) *battery.RuneConfig {

	if p, ok := sm.cache[sm.Major()].RuneConfigs[id]; ok {
		//找到
		return p
	}

	return nil
}

func (sm *RuneConfigCacheManager) RuneConfigs() *MAPRuneConfig {
	return &(sm.cache[sm.Major()].RuneConfigs)
}

func (sm *RuneConfigCacheManager) SecondaryRuneConfigs() *MAPRuneConfig {
	return &(sm.cache[sm.Secondary()].RuneConfigs)
}

func (sm *RuneConfigCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = sm.loadRuneConfigs()
	return
}

func (sm *RuneConfigCacheManager) loadRuneConfigs() (failReason battery.ErrorCode, err error) {
	mapRuneConfigs := sm.SecondaryRuneConfigs()
	*mapRuneConfigs = make(MAPRuneConfig, 0)

	//加载系统道具配置信息
	dbRuneConfigs := make([]*battery.RuneConfig, 0)
	err = DefCacheDB.LoadRuneConfigs(&dbRuneConfigs)
	if err != nil || len(dbRuneConfigs) <= 0 {
		failReason = xyerror.Resp_QueryRuneConfigsFromDBError.GetCode()
		return
	}

	for _, dbRuneConfig := range dbRuneConfigs {
		runeConfig := &battery.RuneConfig{
			Propid: proto.Uint64(dbRuneConfig.GetPropid()),
			Value:  proto.Int32(dbRuneConfig.GetValue()),
			//ExpiredLimitation: proto.Int64(dbRuneConfig.GetExpiredLimitation()),
		}

		(*mapRuneConfigs)[dbRuneConfig.GetPropid()] = runeConfig
	}

	sm.switchCache()

	//xylog.Debug("RuneConfigs : %v ", *mapRuneConfigs)

	return
}

func (sm *RuneConfigCacheManager) switchCache() (fail_reason int32, err error) {
	sm.Switch()
	//xylog.Debug("now RuneConfigs cache switch to %d", sm.Major())
	return
}
