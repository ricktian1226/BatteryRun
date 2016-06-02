// cache_mission
// 任务信息的缓存管理器定义
package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

//保存任务信息的map missionType <-> []*battery.MissionItem
type MAPMissionType2Slice map[battery.MissionType][]*battery.MissionItem

func (m MAPMissionType2Slice) Print() {
	for missionType, missionItems := range m {
		xylog.DebugNoId("--missionType(%v)--", missionType)
		for _, missionItem := range missionItems {
			xylog.DebugNoId("%v", missionItem)
		}
	}
}

func (m *MAPMissionType2Slice) Clear() {
	*m = make(MAPMissionType2Slice, 0)
}

type MAPMissions map[uint64]battery.MissionItem

func (m *MAPMissions) Print() {
	for _, mission := range *m {
		xylog.DebugNoId("mission[%d] %v", mission.GetId(), mission)
	}
}

func (m *MAPMissions) Clear() {
	*m = make(MAPMissions, 0)
}

type MissionCache struct {
	typeMissions MAPMissionType2Slice
	missions     MAPMissions
}

func (p MissionCache) Print() {
	p.missions.Print()
	//p.typeMissions.Print()
}

func NewMissionCache() *MissionCache {
	return &MissionCache{}
}

type MissionCacheManager struct {
	cache [2]MissionCache
	CacheDB
	xycache.CacheBase
}

//任务缓存管理器
var DefMissionCacheManager = NewMissionCacheManager()

func NewMissionCacheManager() *MissionCacheManager {
	return &MissionCacheManager{}
}

func (pm *MissionCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	//初始化数据库操作指针
	pm.Init()
	//加载资源配置信息
	failReason, err = DefMissionCacheManager.ReLoad()
	if failReason != battery.ErrorCode_NoError || err != xyerror.ErrOK {
		xylog.ErrorNoId("MissionResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return

}

func (pm *MissionCacheManager) Init() {
	//pm.index = 0
}

func (pm *MissionCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {
	//begin := time.Now()
	//defer xylog.Debug("MissionCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = pm.Load()

	return
}

//获取指定任务类型的任务配置信息
// missionTypes []battery.MissionType 任务类型列表
// now int64 当前时间戳
//return:
// missions []*battery.MissionItem 任务配置信息列表
func (pm *MissionCacheManager) TypeMissionItems(missionTypes []battery.MissionType, now int64) (missions []*battery.MissionItem) {

	for _, missionType := range missionTypes {
		if missionItems, ok := pm.cache[pm.Major()].typeMissions[missionType]; ok {
			for _, missionItem := range missionItems {
				missionItemTmp := missionItem
				if now >= missionItem.GetBegintime() && now <= missionItem.GetEndtime() {
					missions = append(missions, missionItemTmp)
				}
			}
		}
	}

	return
}

//获取指定id的任务配置信息
// uint64 id 任务id
//return:
// *battery.MissionItem 任务配置信息指针
func (pm *MissionCacheManager) Mission(id uint64) *battery.MissionItem {
	if mission, ok := pm.cache[pm.Major()].missions[id]; ok {
		return &mission
	} else {
		return nil
	}
}

func (pm *MissionCacheManager) TypeMissions() *MAPMissionType2Slice {
	return &(pm.cache[pm.Major()].typeMissions)
}

func (pm *MissionCacheManager) SecondaryTypeMissions() *MAPMissionType2Slice {
	return &(pm.cache[pm.Secondary()].typeMissions)
}

func (pm *MissionCacheManager) Missions() *MAPMissions {
	return &(pm.cache[pm.Major()].missions)
}

func (pm *MissionCacheManager) SecondaryMissions() *MAPMissions {
	return &(pm.cache[pm.Secondary()].missions)
}

func (pm *MissionCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = pm.loadMissions()
	return
}

func (pm *MissionCacheManager) loadMissions() (failReason battery.ErrorCode, err error) {
	dbMissionItems := make([]*battery.MissionItem, 0)
	err = DefCacheDB.LoadMissionItems(&dbMissionItems)
	if err != xyerror.ErrOK || len(dbMissionItems) <= 0 {
		failReason = battery.ErrorCode_QueryMissionError
		return
	}

	typeMissions := pm.SecondaryTypeMissions()
	missions := pm.SecondaryMissions()

	typeMissions.Clear()
	missions.Clear()

	for _, dbMissionItem := range dbMissionItems {

		missionType := dbMissionItem.GetType()
		dbMissionItemTmp := dbMissionItem
		if _, ok := (*typeMissions)[missionType]; !ok {
			(*typeMissions)[missionType] = make([]*battery.MissionItem, 0)
		}

		(*typeMissions)[missionType] = append((*typeMissions)[missionType], dbMissionItemTmp)
		(*missions)[dbMissionItemTmp.GetId()] = *dbMissionItem
	}

	pm.switchCache()

	//pm.cache[pm.Major()].Print()

	return
}

func (pm *MissionCacheManager) switchCache() (fail_reason int32, err error) {
	pm.Switch()
	//xylog.Debug("now mission cache switch to %d", pm.Major())
	return
}
