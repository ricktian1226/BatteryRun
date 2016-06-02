// cache_advertisement
// 广告信息的缓存管理器定义
package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	//"time"
	"math/rand"
	"sort"
)

// MAPAdvertisement 广告信息列表 ，key 广告标识
type MAPAdvertisement map[uint32]*battery.Advertisement

// Print 打印广告缓存信息
func (m *MAPAdvertisement) Print() {
	xylog.DebugNoId("================= Advertisement begin =================")
	for _, v := range *m {
		xylog.DebugNoId("%v", v)
	}
	xylog.DebugNoId("================= Advertisement end =================")
}

// Clear 重置函数
func (m *MAPAdvertisement) Clear() {
	*m = make(MAPAdvertisement, 0)
}

// MAPWeight2AdvertisementId 广告位缓存信息
type MAPWeight2AdvertisementId map[uint32]uint32

// Print 打印广告缓存信息
func (m *MAPWeight2AdvertisementId) Print() {
	for k, v := range *m {
		xylog.DebugNoId("Weight(%d) AdvertisementId(%d)", k, v)
	}
}

// Clear 重置函数
func (m *MAPWeight2AdvertisementId) Clear() {
	*m = make(MAPWeight2AdvertisementId, 0)
}

// AdvertisementPool 广告池信息
type AdvertisementPool struct {
	enable      bool                         //广告位是否播放广告
	totalWeight uint32                       //权重总值
	M           MAPWeight2AdvertisementId    //权重值->广告Id
	S           sort.IntSlice                //排序的权重key值
	flag        []*battery.AdvertisementFlag //播放标志
}

func NewAdvertisementPool() *AdvertisementPool {
	return &AdvertisementPool{
		enable:      false,
		totalWeight: 0,
		M:           make(MAPWeight2AdvertisementId, 0),
		flag:        make([]*battery.AdvertisementFlag, 0),
	}
}

//非法的广告id
const InvalidAdvertisementId uint32 = 0

// elect 获取一个可播放广告
//returns:
// bool   选取结果。 true 选取到广告；false 未选取到广告
// uint32 广告id。选取到广告，则返回广告id，未选取到广告，则返回一个非法广告id
func (p *AdvertisementPool) elect() (bool, uint32) {
	if len(p.M) > 0 {
		if len(p.M) == 1 { //如果只有一个广告，直接返回该广告Id
			for _, v := range p.M {
				return true, v
			}
		} else { //如果广告多余一个，随机选取一个
			//获取随机数,根据随机数获取广告Id
			ranNum := rand.Intn(int(p.totalWeight))
			for _, w := range p.S {
				if ranNum <= w {
					xylog.DebugNoId("AdvertisementPool.Elect ranNum %d totalWeight %d : goods %d", ranNum, p.totalWeight, p.M[uint32(w)])
					return true, p.M[uint32(w)]
				}
			}
		}
	}

	return false, InvalidAdvertisementId
}

// parse 解析数据库查出的广告位信息
func (p *AdvertisementPool) parse(dbAdvertisementSpace *battery.AdvertisementSpace) {
	p.enable = dbAdvertisementSpace.GetEnable()
	p.totalWeight = 0
	for _, item := range dbAdvertisementSpace.GetItems() {
		p.totalWeight += item.GetWeight()
		p.M[p.totalWeight] = item.GetId()
		p.S = append(p.S, int(p.totalWeight))
		p.flag = dbAdvertisementSpace.GetFlags()
	}

	sort.Sort(p.S)
}

// Print 打印广告池信息
func (m *AdvertisementPool) Print() {
	m.M.Print()
}

// MAPAdvertisementSpace 广告位信息列表
type MAPAdvertisementSpace map[uint32]*AdvertisementPool

// Clear 重置函数
func (m *MAPAdvertisementSpace) Clear() {
	*m = make(MAPAdvertisementSpace, 0)
}

// Print 重置函数
func (m *MAPAdvertisementSpace) Print() {
	xylog.DebugNoId("================= AdvertisementSpace begin =================")
	for k, v := range *m {
		xylog.DebugNoId("[%d] %v", k, v)
	}
	xylog.DebugNoId("================= AdvertisementSpace end =================")
}

// AdvertisementCache 广告缓存结构体定义
type AdvertisementCache struct {
	MSpace MAPAdvertisementSpace //广告位信息列表
	M      MAPAdvertisement      //广告信息列表
}

// elect 选取一个可播放广告
// spaceId uint32 广告位标识
//returns:
// ret bool 选取结果
// advertisement *battery.Advertisement 广告信息指针
func (ac *AdvertisementCache) Elect(spaceId uint32) (ret bool, enable bool, flags []*battery.AdvertisementFlag, advertisement *battery.Advertisement) {

	ret = false

	if pool, ok := ac.MSpace[spaceId]; ok { //查找是否存在该广告位
		enable, flags = pool.enable, pool.flag

		if !enable { //如果广告位不播放广告，直接返回
			return
		}

		var id uint32
		ret, id = pool.elect()
		if ret && id != InvalidAdvertisementId {
			if advertisement, ok = ac.M[id]; ok { //查找对应的广告信息
				xylog.DebugNoId("elect AdvertisementId result(%v)", pool)
				ret = true
				return
			}
		} else {
			xylog.ErrorNoId("elect AdvertisementId from AdvertisementPool(%v) failed", pool)
			return
		}

	} else {
		xylog.ErrorNoId("Get AdvertisementSpaceInfo(%d) from ac.MSpace(%v) failed", spaceId, ac.MSpace)
		return
	}

	return
}

// AdvertisementManager 广告信息管理器
type AdvertisementManager struct {
	cache [2]AdvertisementCache
	xycache.CacheBase
}

//广告缓存管理器
var DefAdvertisementManager = NewAdvertisementManager()

func NewAdvertisementManager() *AdvertisementManager {
	return &AdvertisementManager{}
}

// InitWhileStart 进程启动时初始化函数
func (am *AdvertisementManager) InitWhileStart() (err error) {
	//加载资源配置信息
	err = DefAdvertisementManager.Reload()
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("Advertisement ResLoad failed : %v ", err)
		return
	}

	return
}

// Elect 选取一个可播放广告
// spaceId uint32 广告位标识
//returns:
// ret bool 选取结果
// advertisement *battery.Advertisement 广告信息指针
func (am *AdvertisementManager) Elect(spaceId uint32) (ret bool, enable bool, flags []*battery.AdvertisementFlag, advertisement *battery.Advertisement) {
	return am.cache[am.Major()].Elect(spaceId)
}

// Reload 重载资源
func (am *AdvertisementManager) Reload() (err error) {

	//加载广告配置信息
	err = am.reloadAdvertisement()
	if err != xyerror.ErrOK {
		return
	}

	//加载广告位配置信息
	err = am.reloadAdvertisementSpace()
	if err != xyerror.ErrOK {
		return
	}

	//切换缓存
	am.switchCache()

	return
}

// reloadAdvertisement 重载广告配置信息
func (am *AdvertisementManager) reloadAdvertisement() (err error) {
	//从数据库读取
	advertisements := make([]*battery.Advertisement, 0)
	err = DefCacheDB.LoadAdvertisements(&advertisements)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		xylog.ErrorNoId("LoadAdvertisements from db failed : %v", err)
		return
	}

	xylog.DebugNoId("Advertisements : %v", advertisements)

	secondaryAdvertisement := am.SecondaryAdvertisements()
	secondaryAdvertisement.Clear()
	for _, advertisement := range advertisements {
		(*secondaryAdvertisement)[advertisement.GetId()] = advertisement
	}

	secondaryAdvertisement.Print()

	return xyerror.ErrOK
}

// reloadAdvertisementSpace 重载广告位配置信息
func (am *AdvertisementManager) reloadAdvertisementSpace() (err error) {
	//从数据库读取
	advertisementSpaces := make([]*battery.AdvertisementSpace, 0)
	err = DefCacheDB.LoadAdvertisementSpaces(&advertisementSpaces)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		xylog.ErrorNoId("LoadAdvertisementSpaces from db failed : %v", err)
		return
	}

	xylog.DebugNoId("AdvertisementSpaces : %v", advertisementSpaces)

	secondaryAdvertisementSpace := am.SecondaryAdvertisementSpaces()
	secondaryAdvertisementSpace.Clear()

	for _, advertisementSpace := range advertisementSpaces {

		advertisementPool := NewAdvertisementPool()
		advertisementPool.parse(advertisementSpace)

		(*secondaryAdvertisementSpace)[advertisementSpace.GetId()] = advertisementPool
	}

	secondaryAdvertisementSpace.Print()

	return xyerror.ErrOK
}

func (am *AdvertisementManager) Advertisements() *MAPAdvertisement {
	return &(am.cache[am.Major()].M)
}

func (am *AdvertisementManager) AdvertisementSpaces() *MAPAdvertisementSpace {
	return &(am.cache[am.Major()].MSpace)
}

func (am *AdvertisementManager) SecondaryAdvertisements() *MAPAdvertisement {
	return &(am.cache[am.Secondary()].M)
}

func (am *AdvertisementManager) SecondaryAdvertisementSpaces() *MAPAdvertisementSpace {
	return &(am.cache[am.Secondary()].MSpace)
}

func (am *AdvertisementManager) switchCache() (fail_reason int32, err error) {
	am.Switch()
	xylog.DebugNoId("now AdvertisementManager cache switch to %d", am.Major())
	return
}
