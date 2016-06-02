// cache_pickup
// 收集物信息的缓存管理器定义
package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"math/rand"
	"sort"
	//"time"
)

const (
	INVALID_PROPID = uint64(0)
)

//保存道具信息的map propid <-> PropStruct
type MAPPickUpWeight2PropId map[uint32]uint64

func (m MAPPickUpWeight2PropId) Print() {
	for weight, id := range m {
		xylog.DebugNoId("weight(%d), propid(%d)", weight, id)
	}
}

type PickUpWeightStruct struct {
	M           MAPPickUpWeight2PropId
	S           sort.IntSlice
	totalWeight int
}

func (p *PickUpWeightStruct) Print() {
	xylog.DebugNoId("sortSlice : %v", p.S)
	p.M.Print()
}

func (p *PickUpWeightStruct) Clear() {
	p.M = make(MAPPickUpWeight2PropId, 0)
	p.S = make(sort.IntSlice, 0)
}

func (p *PickUpWeightStruct) Sort() {
	sort.Sort(p.S)
}

func NewPickUpWeightStruct() *PickUpWeightStruct {
	return &PickUpWeightStruct{
		M:           make(MAPPickUpWeight2PropId, 0),
		S:           make(sort.IntSlice, 0),
		totalWeight: 0,
	}
}

type MAPPickUpPropType map[battery.PropType]*PickUpWeightStruct

func (m *MAPPickUpPropType) Print() {
	for propType, subMap := range *m {
		xylog.DebugNoId("--type(%v)--", propType)
		subMap.Print()
	}
}

type MAPPickUpCheckPointId map[uint32]*MAPPickUpPropType

func (m *MAPPickUpCheckPointId) Print() {
	for checkPointId, subMap := range *m {
		xylog.DebugNoId("--checkpoint(%d)--", checkPointId)
		subMap.Print()
	}
}

func (m *MAPPickUpCheckPointId) Clear() {
	*m = make(MAPPickUpCheckPointId, 0)
}

func (m *MAPPickUpCheckPointId) Sort() {
	for _, mapPropType := range *m {
		for _, pickUpStruct := range *mapPropType {
			pickUpStruct.Sort()
		}
	}
}

type PickUpCache struct {
	pickUps MAPPickUpCheckPointId //收集物信息列表
}

func (p PickUpCache) Print() {
	p.pickUps.Print()
}

func NewPickUpCache() *PickUpCache {
	return &PickUpCache{}
}

type PickUpCacheManager struct {
	cache [2]PickUpCache
	CacheDB
	xycache.CacheBase
}

//收集物缓存管理器
var DefPickUpCacheManager = NewPickUpCacheManager()

func NewPickUpCacheManager() *PickUpCacheManager {
	return &PickUpCacheManager{}
}

func (pm *PickUpCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	//初始化数据库操作指针
	pm.Init()
	//加载资源配置信息
	failReason, err = DefPickUpCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("PickUpResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return

}

func (pm *PickUpCacheManager) Init() {
	//pm.index = 0
}

func (pm *PickUpCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {
	//begin := time.Now()
	//defer xylog.DebugNoId("PickUpCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = pm.Load()

	return
}

//获取收集物id
// checkPointId uint32 记忆点id
// propType battery.PropType 道具类型
//return:
// uint64 道具id，没找到，返回INVALID_PROPID
func (pm *PickUpCacheManager) PickUp(checkPointId uint32, propType battery.PropType) uint64 {

	if types, ok := pm.cache[pm.Major()].pickUps[checkPointId]; ok {
		if weightStruct, ok := (*types)[propType]; ok {
			weight := rand.Intn(weightStruct.totalWeight)
			if weight < 0 {
				return INVALID_PROPID
			}
			for _, w := range weightStruct.S {
				if w > weight {
					if id, ok := weightStruct.M[uint32(w)]; ok {
						return id
					}
				}
			}
		}
	}

	return INVALID_PROPID
}

func (pm *PickUpCacheManager) PickUps() *MAPPickUpCheckPointId {
	return &(pm.cache[pm.Major()].pickUps)
}

func (pm *PickUpCacheManager) SecondaryPickUps() *MAPPickUpCheckPointId {
	return &(pm.cache[pm.Secondary()].pickUps)
}

func (pm *PickUpCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = pm.loadPickUps()
	return
}

func (pm *PickUpCacheManager) loadPickUps() (failReason battery.ErrorCode, err error) {
	dbPickUpItems := make([]*battery.DBPickUpItem, 0)
	err = DefCacheDB.LoadPickUps(&dbPickUpItems)
	if err != xyerror.ErrOK || len(dbPickUpItems) <= 0 {
		failReason = battery.ErrorCode_QueryPickUpConfigsFromDBError
		return
	}

	mapPickUps := pm.SecondaryPickUps()

	mapPickUps.Clear()

	for _, dbPickUpItem := range dbPickUpItems {
		checkPointId := dbPickUpItem.GetCheckPointId()
		propType := dbPickUpItem.GetPropType()
		propId := dbPickUpItem.GetPropId()
		weight := dbPickUpItem.GetWeight()

		if _, ok := (*mapPickUps)[checkPointId]; !ok {
			mapTmp := make(MAPPickUpPropType, 0)
			(*mapPickUps)[checkPointId] = &mapTmp
		}

		if _, ok := (*((*mapPickUps)[checkPointId]))[propType]; !ok {
			(*((*mapPickUps)[checkPointId]))[propType] = NewPickUpWeightStruct()
		}

		(*((*mapPickUps)[checkPointId]))[propType].totalWeight += int(weight)
		totalWeight := (*((*mapPickUps)[checkPointId]))[propType].totalWeight
		(*((*mapPickUps)[checkPointId]))[propType].M[uint32(totalWeight)] = propId
		(*((*mapPickUps)[checkPointId]))[propType].S = append((*((*mapPickUps)[checkPointId]))[propType].S, totalWeight)
	}

	mapPickUps.Sort()

	pm.switchCache()

	//pm.cache[pm.Major()].Print()

	return
}

func (pm *PickUpCacheManager) switchCache() (fail_reason int32, err error) {
	pm.Switch()
	//xylog.Debug("now pickup cache switch to %d", pm.Major())
	return
}
