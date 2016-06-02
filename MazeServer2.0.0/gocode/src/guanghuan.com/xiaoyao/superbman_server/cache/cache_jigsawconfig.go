// cache_jigsawconfig
// 签到活动信息的缓存管理器定义

package xybusinesscache

import (
	//"code.google.com/p/goprotobuf/proto"
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	//"sync/atomic"
	//"time"
)

type MAPJigsawConfig map[uint64]*battery.JigsawConfig

type JigsawConfigCache struct {
	JigsawConfigs MAPJigsawConfig
}

func NewJigsawConfigCache() *JigsawConfigCache {
	return &JigsawConfigCache{
		JigsawConfigs: make(MAPJigsawConfig, 0),
	}
}

type JigsawConfigCacheManager struct {
	cache [2]*JigsawConfigCache
	xycache.CacheBase
}

//系统道具配置缓存管理器
var DefJigsawConfigCacheManager = NewJigsawConfigCacheManager()

func NewJigsawConfigCacheManager() (manager *JigsawConfigCacheManager) {

	manager = &JigsawConfigCacheManager{}

	for i := 0; i < 2; i++ {
		manager.cache[i] = NewJigsawConfigCache()
	}

	return
}

//sm?! damn!~
func (sm *JigsawConfigCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	sm.Init()
	//加载资源配置信息
	failReason, err = DefJigsawConfigCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("JigsawConigResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return
}

func (sm *JigsawConfigCacheManager) Init() {
	//sm.index = 0
	for i := 0; i < 2; i++ {
		sm.cache[i] = NewJigsawConfigCache()
	}
}

func (sm *JigsawConfigCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {

	//begin := time.Now()
	//defer xylog.Debug("JigsawConfigCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = sm.Load()

	return
}

func (sm *JigsawConfigCacheManager) JigsawConfig(id uint64) *battery.JigsawConfig {

	if p, ok := sm.cache[sm.Major()].JigsawConfigs[id]; ok {
		//找到
		return p
	}

	return nil
}

func (sm *JigsawConfigCacheManager) JigsawConfigs() *MAPJigsawConfig {
	return &(sm.cache[sm.Major()].JigsawConfigs)
}

func (sm *JigsawConfigCacheManager) SecondaryJigsawConfigs() *MAPJigsawConfig {
	return &(sm.cache[sm.Secondary()].JigsawConfigs)
}

func (sm *JigsawConfigCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = sm.loadJigsawConfigs()
	return
}

func (sm *JigsawConfigCacheManager) loadJigsawConfigs() (failReason battery.ErrorCode, err error) {
	mapJigsawConfigs := sm.SecondaryJigsawConfigs()
	*mapJigsawConfigs = make(MAPJigsawConfig, 0)
	/*
		//加载SignInItem信息
		items := make([]*battery.DBSignInItem, 0)
		err = DefCacheDB.LoadSignInItems(&items)
		if err != nil || len(items) <= 0 {
			failReason = xyerror.Resp_QuerySignInActivitysFromDBError.GetCode()
			return
		}
		mapId2Items := make(map[uint64][]*battery.SignInItem, 0)

		for _, item := range items {
			var itemTmp battery.SignInItem
			itemTmp.Value = item.Value
			itemTmp.Award = item.Award
			mapId2Items[item.GetId()] = append(mapId2Items[item.GetId()], &itemTmp)
		}

		//将Items按照Value排下序，方便随机存取
		for k, items := range mapId2Items {
			for _, item := range items {
				mapId2Items[k][int(item.GetValue())] = item
			}
		}
	*/
	//加载系统道具配置信息
	dbJigsawConfigs := make([]*battery.JigsawConfig, 0)
	err = DefCacheDB.LoadJigsawConfigs(&dbJigsawConfigs)
	if err != nil || len(dbJigsawConfigs) <= 0 {
		failReason = xyerror.Resp_QueryJigsawConfigsFromDBError.GetCode()
		return
	}

	for _, dbJigsawConfig := range dbJigsawConfigs {
		(*mapJigsawConfigs)[dbJigsawConfig.GetJigsawid()] = dbJigsawConfig
	}

	sm.switchCache()

	//xylog.Debug("JigsawConfigs : %v ", *mapJigsawConfigs)

	return
}

func (sm *JigsawConfigCacheManager) switchCache() (fail_reason int32, err error) {
	sm.Switch()
	//xylog.Debug("now JigsawConfigs cache switch to %d", sm.Major())
	return
}
