// cache_signin
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

type MAPSignInActivity map[uint64]*battery.SignInActivity

type SignInCache struct {
	Activitys MAPSignInActivity
}

func NewSignInCache() *SignInCache {
	return &SignInCache{}
}

type SignInCacheManager struct {
	cache [2]SignInCache
	xycache.CacheBase
}

//道具缓存管理器
var DefSignInCacheManager = NewSignInCacheManager()

func NewSignInCacheManager() *SignInCacheManager {
	return &SignInCacheManager{}
}

//sm?! damn!~
func (sm *SignInCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	sm.Init()
	//加载资源配置信息
	failReason, err = DefSignInCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("SignInResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return
}

func (sm *SignInCacheManager) Init() {
	//sm.index = 0
}

func (sm *SignInCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {

	//begin := time.Now()
	//defer xylog.Debug("SignInCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = sm.Load()

	return
}

func (sm *SignInCacheManager) Activity(id uint64) *battery.SignInActivity {

	if p, ok := sm.cache[sm.Major()].Activitys[id]; ok {
		//找到
		return p
	}

	return nil
}

func (sm *SignInCacheManager) Activitys() *MAPSignInActivity {
	return &(sm.cache[sm.Major()].Activitys)
}

func (sm *SignInCacheManager) SecondaryActivitys() *MAPSignInActivity {
	return &(sm.cache[sm.Secondary()].Activitys)
}

func (sm *SignInCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = sm.loadActivitys()
	return
}

func (sm *SignInCacheManager) loadActivitys() (failReason battery.ErrorCode, err error) {
	mapActivitys := sm.SecondaryActivitys()
	*mapActivitys = make(MAPSignInActivity, 0)
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
		mapId2Items[k] = make([]*battery.SignInItem, len(items)) //这里必须重新分配一块内存，否则mapId2Items[k]和items引用的是同一块内存。会出现覆盖的bug
		for _, item := range items {
			mapId2Items[k][int(item.GetValue())] = item
		}
	}

	//加载活动项
	dbActivitys := make([]*battery.DBSignInActivity, 0)
	err = DefCacheDB.LoadSignInActivitys(&dbActivitys)
	if err != nil || len(dbActivitys) <= 0 {
		failReason = xyerror.Resp_QuerySignInActivitysFromDBError.GetCode()
		return
	}

	for _, dbActivity := range dbActivitys {
		activityTmp := &battery.SignInActivity{
			Id:        proto.Uint64(dbActivity.GetId()),
			Type:      dbActivity.GetType().Enum(),
			GoalValue: proto.Uint32(dbActivity.GetGoalValue()),
			BeginTime: proto.Int64(dbActivity.GetBeginTime()),
			EndTime:   proto.Int64(dbActivity.GetEndTime()),
		}
		if items, ok := mapId2Items[dbActivity.GetId()]; ok {
			activityTmp.Items = items
		} else {
			activityTmp.Items = make([]*battery.SignInItem, 0)
		}

		(*mapActivitys)[dbActivity.GetId()] = activityTmp
	}

	sm.switchCache()

	//xylog.Debug("Activitys : %v ", *mapActivitys)

	return
}

func (sm *SignInCacheManager) switchCache() (fail_reason int32, err error) {
	sm.Switch()
	//xylog.Debug("now signinactivity cache switch to %d", sm.Major())
	return
}
