// cache_role
package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

const INVALID_ROLEID = 0
const INVALID_LEVEL int32 = -1

type RoleInfoItem struct {
	//Id           uint64
	IsDefaultOwn bool
	MaxLevel     int32
	JigsawId     uint64
}

//保存角色加成信息的map id <-> RoleInfoItem

type MAPRoleInfoItems map[uint64]*RoleInfoItem

func (m MAPRoleInfoItems) Print() {
	for id, item := range m {
		xylog.DebugNoId("id(%d), info(IsDefaultOwn %t, MaxLevel %d, JigsawId %d)", id, item.IsDefaultOwn, item.MaxLevel, item.JigsawId)
	}
}

func (m *MAPRoleInfoItems) Clear() {
	*m = make(MAPRoleInfoItems, 0)
}

type RoleInfoCache struct {
	roleInfos  MAPRoleInfoItems //角色信息列表
	defaultOwn uint64           //默认拥有的角色
}

func (p *RoleInfoCache) Print() {
	p.roleInfos.Print()
	xylog.DebugNoId("defaultOwn %d", p.defaultOwn)
}

func (p *RoleInfoCache) Clear() {
	p.roleInfos.Clear()
	p.defaultOwn = INVALID_ROLEID
}

func NewRoleInfoCache() *RoleInfoCache {
	return &RoleInfoCache{}
}

type RoleInfoCacheManager struct {
	cache [2]RoleInfoCache
	CacheDB
	xycache.CacheBase
}

//收集物缓存管理器
var DefRoleInfoCacheManager = NewRoleInfoCacheManager()

func NewRoleInfoCacheManager() *RoleInfoCacheManager {
	return &RoleInfoCacheManager{}
}

func (pm *RoleInfoCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
	//初始化数据库操作指针
	pm.Init()
	//加载资源配置信息
	failReason, err = DefRoleInfoCacheManager.ReLoad()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("RoleInfoResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return

}

func (pm *RoleInfoCacheManager) Init() {
	//pm.index = 0
}

func (pm *RoleInfoCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {
	//begin := time.Now()
	//defer xylog.Debug("RoleLevelBonusCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	failReason, err = pm.Load()

	return
}

//获取角色信息
// id uint64 角色id
//return:
// *RoleInfoItem 角色信息指针
func (pm *RoleInfoCacheManager) Info(id uint64) *RoleInfoItem {

	if info, ok := pm.cache[pm.Major()].roleInfos[id]; ok {
		return info
	}

	return nil
}

//返回角色配置信息列表
func (pm *RoleInfoCacheManager) Infos() *MAPRoleInfoItems {
	return &(pm.cache[pm.Major()].roleInfos)
}

//获取默认的角色id
func (pm *RoleInfoCacheManager) DefaultOwn() uint64 {
	return pm.cache[pm.Major()].defaultOwn
}

func (pm *RoleInfoCacheManager) MajorCache() *RoleInfoCache {
	return &(pm.cache[pm.Major()])
}

func (pm *RoleInfoCacheManager) SecondaryCache() *RoleInfoCache {
	return &(pm.cache[pm.Secondary()])
}

func (pm *RoleInfoCacheManager) Load() (failReason battery.ErrorCode, err error) {
	failReason, err = pm.loadRoleInfo()
	return
}

func (pm *RoleInfoCacheManager) loadRoleInfo() (failReason battery.ErrorCode, err error) {
	dbItems := make([]*battery.DBRoleInfoConfig, 0)
	err = DefCacheDB.LoadRoleConfig(&dbItems)
	if err != xyerror.ErrOK || len(dbItems) <= 0 {
		failReason = battery.ErrorCode_QueryRoleInfoConfigFromDBError
		return
	}

	roleInfoCache := pm.SecondaryCache()

	roleInfoCache.Clear()

	defaultOwnFound := false

	for _, dbItem := range dbItems {
		id := dbItem.GetId()
		info := &RoleInfoItem{
			IsDefaultOwn: dbItem.GetIsDefaultOwn(),
			MaxLevel:     dbItem.GetMaxLevel(),
			JigsawId:     dbItem.GetJigsawId(),
		}
		roleInfoCache.roleInfos[id] = info
		if dbItem.GetIsDefaultOwn() { //玩家默认拥有的角色id
			roleInfoCache.defaultOwn = id
			defaultOwnFound = true
		}
	}

	//如果没有默认拥有的角色，返回退出
	if !defaultOwnFound {
		xylog.WarningNoId("RoleInfoConfig defaultOwn NotFound,pls check")
		//err = xyerror.ErrBadInputData
		//return
	}

	pm.switchCache()

	//pm.cache[pm.Major()].Print()

	return
}

func (pm *RoleInfoCacheManager) switchCache() (fail_reason int32, err error) {
	pm.Switch()
	//xylog.Debug("now roleinfos cache switch to %d", pm.Major())
	return
}
