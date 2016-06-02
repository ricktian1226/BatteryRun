// cache_bgpropconfig
// 签到活动信息的缓存管理器定义

package xybusinesscache

import (
	//"code.google.com/p/goprotobuf/proto"
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"math/rand"
	"sort"
	"time"
	//"sync/atomic"
)

type Empty struct{}
type Set map[interface{}]Empty

//type MAPBGPropConfig map[uint64]*battery.BeforeGamePropConfig
type MAPBeforeGameWeight2RandomGood map[uint32]uint64
type BeforeGameWeightCache struct {
	M              MAPBeforeGameWeight2RandomGood //weight->goodid
	S              sort.IntSlice
	goodsSet       Set    //保存所有游戏前商品id的set
	randomGoodsSet Set    //保存游戏前随机商品id的set
	totalWeight    uint32 //所有商品的权重总和
}

func (cache *BeforeGameWeightCache) Clear() {
	cache.M = make(MAPBeforeGameWeight2RandomGood, 0)
	cache.S = make(sort.IntSlice, 0)
	cache.goodsSet = make(Set, 0)
	cache.randomGoodsSet = make(Set, 0)
	cache.totalWeight = 0
}

func (cache *BeforeGameWeightCache) Print() {
	for weight, id := range cache.M {
		xylog.DebugNoId("BeforeGameWeightCache.M[%d] : %d", weight, id)
	}

	xylog.DebugNoId("BeforeGameWeightCache.goodsSet :")

	for id, _ := range cache.goodsSet {
		xylog.DebugNoId("%d", id)
	}

	xylog.DebugNoId("BeforeGameWeightCache.totalWeight : %d", cache.totalWeight)
}

func NewBeforeGameWeightCache() *BeforeGameWeightCache {
	return &BeforeGameWeightCache{}
}

type BeforeGameWeightCacheManager struct {
	cache [2]BeforeGameWeightCache
	xycache.CacheBase
}

//赛前道具配置缓存管理器
var DefBeforeGameWeightCacheManager = NewBeforeGameWeightCacheManager()

func NewBeforeGameWeightCacheManager() *BeforeGameWeightCacheManager {
	return &BeforeGameWeightCacheManager{}
}

//sm?! damn!~
func (sm *BeforeGameWeightCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	sm.Init()
	//加载资源配置信息
	failReason, err = DefBeforeGameWeightCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != xyerror.ErrOK {
		xylog.ErrorNoId("BeforeGameWeightResLoad failed : %d, %v ", failReason, err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return
}

func (sm *BeforeGameWeightCacheManager) Init() {
	//sm.index = 0
}

func (sm *BeforeGameWeightCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {
	begin := time.Now()
	defer xylog.DebugNoId("BeforeGameWeightCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = sm.Load()
	return
}

func (sm *BeforeGameWeightCacheManager) MajorCache() *BeforeGameWeightCache {
	return &(sm.cache[sm.Major()])
}

func (sm *BeforeGameWeightCacheManager) SecondaryCache() *BeforeGameWeightCache {
	return &(sm.cache[sm.Secondary()])
}

func (sm *BeforeGameWeightCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = sm.loadBeforeGameWeight()
	if failReason != battery.ErrorCode_NoError || err != xyerror.ErrOK {
		return
	}
	return
}

func (sm *BeforeGameWeightCacheManager) loadBeforeGameWeight() (failReason battery.ErrorCode, err error) {
	beforeGameWeightCache := sm.SecondaryCache()
	beforeGameWeightCache.Clear()

	//加载系统道具配置信息
	dbBeforeGameRandomGoodWeights := make([]*battery.DBBeforeGameRandomGoodWeight, 0)
	err = DefCacheDB.LoadBeforeGameRandomWeights(&dbBeforeGameRandomGoodWeights)
	if err != xyerror.ErrOK || len(dbBeforeGameRandomGoodWeights) <= 0 {
		failReason = xyerror.Resp_QueryBGPropConfigsFromDBError.GetCode()
		return
	}

	totalWeight := uint32(0)

	for _, dbBeforeGameRandomGoodWeight := range dbBeforeGameRandomGoodWeights {

		id := dbBeforeGameRandomGoodWeight.GetGoodId()
		weight := dbBeforeGameRandomGoodWeight.GetWeight()
		if weight > 0 {
			totalWeight += weight
			beforeGameWeightCache.M[totalWeight] = id
			beforeGameWeightCache.S = append(beforeGameWeightCache.S, int(totalWeight))
			beforeGameWeightCache.randomGoodsSet[id] = Empty{}
		}
	}

	beforeGameWeightCache.totalWeight = totalWeight

	sort.Sort(beforeGameWeightCache.S) //排一下序

	//加载游戏前商品信息
	err = sm.loadBeforeGameGoods(&(beforeGameWeightCache.goodsSet))
	if err != xyerror.ErrOK {
		return
	}

	//beforeGameWeightCache.Print()

	sm.switchCache()

	return
}

//加载游戏前道具信息
func (sm *BeforeGameWeightCacheManager) loadBeforeGameGoods(goodsSet *Set) (err error) {
	goods := make([]*battery.DBMallItem, 0)
	err = DefCacheDB.LoadSpecificTypeGoods(battery.MallType_Mall_BeforeGame, &goods)
	if err != xyerror.ErrOK {
		xylog.DebugNoId("LoadSpecificTypeGoods failed : %v", err)
		return
	}

	//xylog.Debug("BeforeGame goods : %v", goods)

	for _, good := range goods {
		(*goodsSet)[good.GetId()] = Empty{}
	}

	return
}

//非法的随机赛前道具编码
const (
	InvalidRandomGood = uint64(0)
)

//获取赛前随机道具
func (sm *BeforeGameWeightCacheManager) RandomGood() uint64 {
	beforeGameWeightCache := sm.cache[sm.Major()]
	if len(beforeGameWeightCache.M) > 0 {
		totalWeight := beforeGameWeightCache.totalWeight
		ranNum := rand.Intn(int(totalWeight))
		for _, w := range beforeGameWeightCache.S {
			if ranNum <= w {
				xylog.DebugNoId("BeforeGameWeightCacheManager.RandomGood ranNum %d totalWeight %d : goods %d", ranNum, totalWeight, beforeGameWeightCache.M[uint32(w)])
				return beforeGameWeightCache.M[uint32(w)]
			}
		}
	}

	xylog.ErrorNoId("BeforeGameWeightCacheManager.RandomGood Error no found")

	return InvalidRandomGood
}

func (sm *BeforeGameWeightCacheManager) GoodsSet() *Set {
	return &(sm.MajorCache().goodsSet)
}

func (sm *BeforeGameWeightCacheManager) RandomGoodsSet() *Set {
	return &(sm.MajorCache().randomGoodsSet)
}

func (sm *BeforeGameWeightCacheManager) switchCache() (fail_reason int32, err error) {
	sm.Switch()
	xylog.DebugNoId("now RuneConfigs cache switch to %d", sm.Major())
	return
}
