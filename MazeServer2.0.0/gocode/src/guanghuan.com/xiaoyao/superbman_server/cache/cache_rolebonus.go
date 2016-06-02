// cache_rolebonus
// 角色加成信息缓存
package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

type RoleLevelBonus struct {
	GoldBonus  int32
	ScoreBonus int32
}

//保存角色加成信息的map id <-> RoleLevelBonusStruct

type MAPRoleLevelBonus map[uint64]*RoleLevelBonus

func (m MAPRoleLevelBonus) Print() {
	for id, bouns := range m {
		xylog.DebugNoId("id(%d), bonus(gold %d, score %d)", id, bouns.GoldBonus, bouns.ScoreBonus)
	}
}

type RoleLevelBonusStruct struct {
	M MAPRoleLevelBonus
}

func (p *RoleLevelBonusStruct) Print() {
	p.M.Print()
}

func (p *RoleLevelBonusStruct) Clear() {
	p.M = make(MAPRoleLevelBonus, 0)
}

func NewRoleLevelBonusStruct() *RoleLevelBonusStruct {
	return &RoleLevelBonusStruct{
		M: make(MAPRoleLevelBonus, 0),
	}
}

type RoleLevelBonusCache struct {
	roleLevelBonuses RoleLevelBonusStruct //收集物信息列表
}

func (p RoleLevelBonusCache) Print() {
	p.roleLevelBonuses.Print()
}

func NewRoleLevelBonusCache() *RoleLevelBonusCache {
	return &RoleLevelBonusCache{}
}

type RoleLevelBonusCacheManager struct {
	cache [2]RoleLevelBonusCache
	CacheDB
	xycache.CacheBase
}

//收集物缓存管理器
var DefRoleLevelBonusCacheManager = NewRoleLevelBonusCacheManager()

func NewRoleLevelBonusCacheManager() *RoleLevelBonusCacheManager {
	return &RoleLevelBonusCacheManager{}
}

func (pm *RoleLevelBonusCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	//初始化数据库操作指针
	pm.Init()
	//加载资源配置信息
	failReason, err = DefRoleLevelBonusCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("RoleLevelBonusResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return

}

func (pm *RoleLevelBonusCacheManager) Init() {
	//pm.index = 0
}

func (pm *RoleLevelBonusCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {
	//begin := time.Now()
	//defer xylog.Debug("RoleLevelBonusCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = pm.Load()

	return
}

//获取收集物id
// checkPointId uint32 记忆点id
// propType battery.PropType 道具类型
//return:
// uint64 道具id，没找到，返回INVALID_PROPID
func (pm *RoleLevelBonusCacheManager) Bonus(id uint64) *RoleLevelBonus {

	if bonus, ok := pm.cache[pm.Major()].roleLevelBonuses.M[id]; ok {
		return bonus
	}

	return nil
}

func (pm *RoleLevelBonusCacheManager) Bonuses() *RoleLevelBonusStruct {
	return &(pm.cache[pm.Major()].roleLevelBonuses)
}

func (pm *RoleLevelBonusCacheManager) SecondaryBonuses() *RoleLevelBonusStruct {
	return &(pm.cache[pm.Secondary()].roleLevelBonuses)
}

func (pm *RoleLevelBonusCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = pm.loadRoleLevelBonus()
	return
}

func (pm *RoleLevelBonusCacheManager) loadRoleLevelBonus() (failReason battery.ErrorCode, err error) {
	dbItems := make([]*battery.DBRoleLevelBonusItem, 0)
	err = DefCacheDB.LoadRoleLevelBonus(&dbItems)
	if err != xyerror.ErrOK || len(dbItems) <= 0 {
		failReason = battery.ErrorCode_QueryRoleLevelBonusFromDBError
		return
	}

	roleLevelBonuses := pm.SecondaryBonuses()

	roleLevelBonuses.Clear()

	for _, dbItem := range dbItems {
		id := dbItem.GetId()
		bonus := &RoleLevelBonus{
			GoldBonus:  dbItem.GetGoldBonus(),
			ScoreBonus: dbItem.GetScoreBonus(),
		}
		roleLevelBonuses.M[id] = bonus
	}

	pm.switchCache()

	//pm.cache[pm.Major()].Print()

	return
}

func (pm *RoleLevelBonusCacheManager) switchCache() (fail_reason int32, err error) {
	pm.Switch()
	//xylog.Debug("now rolelevelbonus cache switch to %d", pm.Major())
	return
}
