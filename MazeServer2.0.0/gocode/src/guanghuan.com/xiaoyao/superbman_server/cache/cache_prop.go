// cache_prop
// 道具信息的缓存管理器定义

package xybusinesscache

import (
    "guanghuan.com/xiaoyao/common/cache"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    //"time"
)

//缓存中道具信息结构体
type PropStruct struct {
    Type         battery.PropType     //道具类型
    ResolveValue []*battery.MoneyItem //道具分解价值
    Items        []*battery.PropItem  //道具子项列表
    LottoValue   uint32               //道具抽奖权值
}

//保存道具信息的map propid <-> PropStruct
type MAPProp map[uint64]PropStruct

func (m *MAPProp) Print() {
    for id, propStruct := range *m {
        xylog.DebugNoId("--prop(%d)--", id)
        xylog.DebugNoId("%v", propStruct)
    }
}

type PropCache struct {
    props MAPProp //道具列表
}

func NewPropCache() *PropCache {
    return &PropCache{}
}

type PropCacheManager struct {
    cache [2]PropCache
    CacheDB
    xycache.CacheBase
}

//道具缓存管理器
var DefPropCacheManager = NewPropCacheManager()

func NewPropCacheManager() *PropCacheManager {
    return &PropCacheManager{}
}

func (pm *PropCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
    //初始化数据库操作指针
    pm.Init()
    //加载资源配置信息
    failReason, err = DefPropCacheManager.ReLoad()
    if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
        xylog.ErrorNoId("PropResLoad failed : %v ", err)
        //os.Exit(-1) //加载失败，进程直接退出
    }
    return

}

func (pm *PropCacheManager) Init() {
    //pm.index = 0
}

func (pm *PropCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {
    //begin := time.Now()
    //defer xylog.Debug("PropCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

    failReason, err = pm.Load()

    return
}

func (pm *PropCacheManager) Prop(pid uint64) *PropStruct {

    if p, ok := pm.cache[pm.Major()].props[pid]; ok {
        //找到
        return &p
    }

    return nil
}

func (pm *PropCacheManager) Props() *MAPProp {
    return &(pm.cache[pm.Major()].props)
}

func (pm *PropCacheManager) SecondaryProps() *MAPProp {
    return &(pm.cache[pm.Secondary()].props)
}

func (pm *PropCacheManager) Load() (failReason battery.ErrorCode, err error) {
    failReason, err = pm.loadProps()
    return
}

func (pm *PropCacheManager) loadProps() (failReason battery.ErrorCode, err error) {
    props := make([]battery.Prop, 0)
    err = DefCacheDB.LoadProps(&props)
    if err != nil || len(props) <= 0 {
        failReason = xyerror.Resp_QueryPropsFromDBError.GetCode()
        return
    }

    mapProps := pm.SecondaryProps()

    *mapProps = make(MAPProp, 0)

    var prop PropStruct
    for _, p := range props {
        prop.Items = make([]*battery.PropItem, 0)
        for _, item := range p.GetItems() {
            prop.Items = append(prop.Items, item)
        }
        prop.Type = p.GetType()
        prop.LottoValue = p.GetLottovalue()
        prop.ResolveValue = p.GetResolvevalue()
        (*mapProps)[p.GetId()] = prop
    }

    //mapProps.Print()

    pm.switchCache()

    return
}

func (pm *PropCacheManager) switchCache() (fail_reason int32, err error) {
    pm.Switch()
    //xylog.Debug("now prop cache switch to %d", pm.Major())
    return
}

//分发类型定义
type DispenseUnit struct {
    Items  []*battery.PropItem   //分发礼包列表
    Type   battery.DISPENSE_TYPE //分发类型
    MailId int32                 //关联的邮件id
}

func (d *DispenseUnit) Print() {
    xylog.DebugNoId("-------- type %v -------", d.Type)
    for _, item := range d.Items {
        xylog.DebugNoId("%v", item)
    }
}

//登录礼包source->propitems
type MAPNewAccountProp map[battery.ID_SOURCE][]DispenseUnit

//登录礼包缓存
type NewAccountPropCache struct {
    M MAPNewAccountProp
}

// 缓存清理函数
func (c *NewAccountPropCache) Clear() {
    c.M = make(MAPNewAccountProp, 0)
}

//登录礼包缓存管理器对象
var DefNewAccountPropManager = NewNewAccountPropManager()

//登录礼包缓存管理器定义
type NewAccountPropManager struct {
    cache [2]NewAccountPropCache
    xycache.CacheBase
}

func NewNewAccountPropManager() *NewAccountPropManager {
    return &NewAccountPropManager{}
}

//登录礼包缓存管理器启动时初始化
func (m *NewAccountPropManager) InitWhileStart() {
    m.Init()
    m.Reload()
}

//登录礼包缓存管理器初始化
func (m *NewAccountPropManager) Init() {

}

//登录礼包缓存管理器重载
func (m *NewAccountPropManager) Reload() (err error) {
    secondary := m.SecondaryCache()

    //从数据库加载登录礼包信息
    dbNewAccountProps := make([]*battery.DBNewAccountProp, 0)
    err = DefCacheDB.LoadNewAccountProps(&dbNewAccountProps)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("LoadNewAccountProps failed : %v", err)
        return
    }

    secondary.Clear()

    for _, dbNewAccountProp := range dbNewAccountProps {
        unit := DispenseUnit{
            Items:  dbNewAccountProp.GetItems(),
            Type:   dbNewAccountProp.GetDispenseType(),
            MailId: dbNewAccountProp.GetMailId(),
        }
        secondary.M[dbNewAccountProp.GetSource()] = append(secondary.M[dbNewAccountProp.GetSource()], unit)
    }

    m.Switch()

    m.Print()

    return
}

func (m *NewAccountPropManager) Print() {
    xylog.DebugNoId("------- NewAccountPropManager cache -------")
    for source, items := range m.MajorCache().M {
        xylog.DebugNoId("------- source %v -------\n%v", source, items)
        //for _, propItem := range propItems {
        //	xylog.DebugNoId("-------  %v -------", propItem)
        //}
    }
}

//获取主缓存
func (m *NewAccountPropManager) MajorCache() *NewAccountPropCache {
    return &(m.cache[m.Major()])
}

//获取备缓存
func (m *NewAccountPropManager) SecondaryCache() *NewAccountPropCache {
    return &(m.cache[m.Secondary()])
}

// 获取某个账户来源的登录礼包信息
// source battery.ID_SOURCE 账户来源
//returns:
// []DispenseUnit 登录礼包列表，如果查询失败，返回nil
func (m *NewAccountPropManager) Units(source battery.ID_SOURCE) []DispenseUnit {
    major := m.MajorCache()
    if units, ok := major.M[source]; ok { //找到就返回对应的礼包信息
        return units
    } else { //没找到返回nil
        return nil
    }
}

type ShareActivity struct {
    Id                             uint32
    ShareType                      battery.SHARE_TYPE
    DailyLimit, GoalValue, MailId  int32
    StartTime, EndTime             int64
    AwardsStartTime, AwardsEndTime int64
    Restart, Valid                 bool
    Items                          []*battery.ShareItem
    DispenseType                   battery.DISPENSE_TYPE
}

// 分享礼包
type MAPShareActivity map[uint32]*ShareActivity

// 分享礼包缓存
type SharedGiftPropCache struct {
    M MAPShareActivity
}

func (c *SharedGiftPropCache) Clear() {
    c.M = make(MAPShareActivity, 0)
}

var DefSharedActivityManager = NewShareActivityManager()

// 分享礼包缓存管理器定义
type ShareActivityManager struct {
    cache [2]SharedGiftPropCache
    xycache.CacheBase
}

func NewShareActivityManager() *ShareActivityManager {
    return &ShareActivityManager{}
}

func (m *ShareActivityManager) InitWhileStart() {
    m.Init()
    m.Reload()
}

func (m *ShareActivityManager) Init() {

}

func (m *ShareActivityManager) Reload() (err error) {
    secondary := m.SecondaryCache()

    dbShareAwards := make([]*battery.DBShareAward, 0)
    err = DefCacheDB.LoadSharedAwards(&dbShareAwards)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("LoadShareWards failed :%v", err)
        return
    }

    mapID2ShareItem := make(map[uint32][]*battery.ShareItem, 0)
    for _, shareAward := range dbShareAwards {
        shareItem := &battery.ShareItem{
            Counter: shareAward.Counter,
            Award:   shareAward.Items,
        }

        mapID2ShareItem[shareAward.GetId()] = append(mapID2ShareItem[shareAward.GetId()], shareItem)

    }

    // 从数据库加载分享礼包信息
    dbShareActivity := make([]*battery.DBShareActivity, 0)
    err = DefCacheDB.LoadSharedActivity(&dbShareActivity)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("LoadShareGiftProps failed :%v", err)
        return
    }

    secondary.Clear()
    for _, shareActivity := range dbShareActivity {
        activity := &ShareActivity{
            Id:              shareActivity.GetId(),
            ShareType:       shareActivity.GetShareType(),
            DailyLimit:      shareActivity.GetDailyLimit(),
            GoalValue:       shareActivity.GetGoalValue(),
            MailId:          shareActivity.GetMailID(),
            StartTime:       shareActivity.GetActivityStartTime(),
            EndTime:         shareActivity.GetActivityEndTime(),
            AwardsStartTime: shareActivity.GetAwardsBeginTime(),
            AwardsEndTime:   shareActivity.GetAwardsEndTime(),
            Restart:         shareActivity.GetRestart(),
            Valid:           shareActivity.GetValid(),
            Items:           mapID2ShareItem[shareActivity.GetId()],
            DispenseType:    shareActivity.GetDispenseType(),
        }
        secondary.M[activity.Id] = activity
    }
    m.Switch()
    m.Print()
    return
}

func (m *ShareActivityManager) Print() {
    xylog.DebugNoId("------SharedGiftPropManager cache -------")
    for source, items := range m.MajorCache().M {
        xylog.DebugNoId("----source %v ------\n %v", source, items)
    }
}

func (m *ShareActivityManager) MajorCache() *SharedGiftPropCache {
    return &(m.cache[m.Major()])
}

func (m *ShareActivityManager) SecondaryCache() *SharedGiftPropCache {
    return &(m.cache[m.Secondary()])
}

// 更具分享天数获取对应礼包信息
func (m *ShareActivityManager) GetShareActivityByID(id uint32) *ShareActivity {
    major := m.MajorCache()
    if units, ok := major.M[id]; ok {
        return units
    } else { //没找到返回nil
        return nil
    }
}

func (m *ShareActivityManager) GetAllShareActivity() (activitys []*ShareActivity) {
    major := m.MajorCache()
    activitys = make([]*ShareActivity, 0)
    for _, activity := range major.M {
        activitys = append(activitys, activity)
    }
    return
}
