// cache_goods
// 商城信息的缓存管理器定义
package xybusinesscache

import (
    proto "code.google.com/p/goprotobuf/proto"
    "guanghuan.com/xiaoyao/common/cache"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    //"time"
)

type MAPGoods map[uint64]*battery.MallItem

func (m *MAPGoods) Print() {
    for id, item := range *m {
        xylog.DebugNoId("--good(%d)--", id)
        xylog.DebugNoId("%v", *item)
    }
}

func (m *MAPGoods) Clear() {
    *m = make(MAPGoods, 0)
}

type MAPMallSubType map[battery.MallSubType][]*battery.MallItem
type MAPMallType map[battery.MallType]MAPMallSubType

func (m *MAPMallType) Print() {
    for mallType, mapSubMallType := range *m {
        xylog.DebugNoId("--mallType(%d)--", mallType)
        for mallSubType, items := range mapSubMallType {
            xylog.DebugNoId("--mallSubType(%d)--", mallSubType)
            for _, item := range items {
                xylog.DebugNoId("%v", item)
            }
        }
    }
}

func (m *MAPMallType) Clear() {
    *m = make(MAPMallType, 0)
}

type MAPIapGoods map[string]*battery.MallItem

func (m *MAPIapGoods) Print() {
    for iapid, item := range *m {
        xylog.DebugNoId("iapgood(%s, %v)", iapid, item)
    }
}

func (m *MAPIapGoods) Clear() {
    *m = make(MAPIapGoods, 0)
}

type GoodsCache struct {
    goods         MAPGoods    //按照商品id索引的map
    goodsSpecific MAPMallType //按照mallType+mallSubType索引的map
    iapGoods      MAPIapGoods //按照iapid索引的map
}

func NewGoodsCache() *GoodsCache {
    return &GoodsCache{
        goods: make(MAPGoods, 0),
    }
}

type GoodsCacheManager struct {
    cache [2]GoodsCache
    //CacheDB
    xycache.CacheBase
}

//道具缓存管理器
var DefGoodsCacheManager = NewGoodsCacheManager()

func NewGoodsCacheManager() *GoodsCacheManager {
    return &GoodsCacheManager{}
}

func (gm *GoodsCacheManager) InitWhileStart() (failReason battery.ErrorCode, err error) {
    //初始化数据库操作指针
    gm.Init()
    //加载资源配置信息
    failReason, err = DefGoodsCacheManager.ReLoad()
    if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
        xylog.ErrorNoId("GoodsResLoad failed : %v ", err)
        //os.Exit(-1) //加载失败，进程直接退出
    }
    return
}

func (gm *GoodsCacheManager) Init() {
    //gm.index = 0
}

func (gm *GoodsCacheManager) ReLoad() (failReason battery.ErrorCode, err error) {
    failReason, err = gm.Load()
    return
}

//获取某个商品的详细信息
// id uint64 商品id
func (gm *GoodsCacheManager) Good(id uint64) *battery.MallItem {

    //gm.cache[gm.Major()].goods.Print()

    if g, ok := gm.cache[gm.Major()].goods[id]; ok {
        //找到
        return g
    }
    return nil
}

//获取特定商店下的商品列表信息
// mallType battery.MallType 商城类型
// mallSubType battery.MallSubType 商品类型
func (gm *GoodsCacheManager) SpecificMall(mallType battery.MallType, mallSubType battery.MallSubType, platform battery.PLATFORM_TYPE) (goodsList []*battery.MallItem) {

    goodsList = make([]*battery.MallItem, 0)

    if mallType == battery.MallType_Mall_Unkown { //如果未指定商场主类型，则返回所有商场类型的商品
        for _, mall := range gm.cache[gm.Major()].goodsSpecific {
            for _, goodsListTmp := range mall {

                goodsList = append(goodsList, goodsListTmp...)
            }
        }
    } else {
        if mall, ok := gm.cache[gm.Major()].goodsSpecific[mallType]; ok {
            //如果未指定mallSubType，则返回所有的mallType下的所有商品
            if mallSubType == battery.MallSubType_MallSubType_Unkown {
                for subtype, goodsListTmp := range mall {

                    if subtype == battery.MallSubType_MallSubType_Android_Diamond {
                        switch platform {
                        case battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID:
                            // 归一返回子类为钻石子类
                            for _, goods := range goodsListTmp {
                                goods.MallSubType = battery.MallSubType_MallSubType_Android_Diamond.Enum()
                                goodsList = append(goodsList, goods)

                            }
                            // ios 平台不返回安卓钻石商品
                        case battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS:
                        }
                    } else {
                        goodsList = append(goodsList, goodsListTmp...)
                    }
                }

            } else if subMall, ok := mall[mallSubType]; ok {
                goodsList = subMall
            }
        }
    }

    return
}

func (gm *GoodsCacheManager) IapGood(iapId string) *battery.MallItem {
    if g, ok := gm.cache[gm.Major()].iapGoods[iapId]; ok {
        //找到
        return g
    }
    return nil
}

func (gm *GoodsCacheManager) Goods() *MAPGoods {
    return &(gm.cache[gm.Major()].goods)
}

func (gm *GoodsCacheManager) SecondaryGoods() *MAPGoods {
    return &(gm.cache[gm.Secondary()].goods)
}

func (gm *GoodsCacheManager) GoodsSpecific() *MAPMallType {
    return &(gm.cache[gm.Major()].goodsSpecific)
}

func (gm *GoodsCacheManager) SecondaryGoodsSpecific() *MAPMallType {
    return &(gm.cache[gm.Secondary()].goodsSpecific)
}

func (gm *GoodsCacheManager) IapGoods() *MAPIapGoods {
    return &(gm.cache[gm.Major()].iapGoods)
}

func (gm *GoodsCacheManager) SecondaryIapGoods() *MAPIapGoods {
    return &(gm.cache[gm.Secondary()].iapGoods)
}

func (gm *GoodsCacheManager) Load() (failReason battery.ErrorCode, err error) {
    failReason, err = gm.loadGoods()
    return
}

func (gm *GoodsCacheManager) loadGoods() (failReason battery.ErrorCode, err error) {
    dbgoods := make([]*battery.DBMallItem, 0)
    err = DefCacheDB.LoadGoods(&dbgoods)
    if err != nil || len(dbgoods) <= 0 {
        failReason = xyerror.Resp_QueryGoodsError.GetCode()
        return
    }

    mapGoods, mapMallType, mapIapGoods := gm.SecondaryGoods(), gm.SecondaryGoodsSpecific(), gm.SecondaryIapGoods()

    mapGoods.Clear()
    mapMallType.Clear()
    mapIapGoods.Clear()

    //xylog.DebugNoId("dbgoods : %v", dbgoods)

    for _, g := range dbgoods {
        //id<->goods
        good := gm.getMallItemFromDBMallItem(g)
        //xylog.DebugNoId("good : %v", good)
        (*mapGoods)[g.GetId()] = good

        //mallType -> mallSubType -> goods
        mallType, mallSubType := good.GetMallType(), good.GetMallSubType()
        if _, ok := (*mapMallType)[mallType]; !ok {
            (*mapMallType)[mallType] = make(MAPMallSubType, 0)
        }
        if _, ok := (*mapMallType)[mallType][mallSubType]; !ok {
            (*mapMallType)[mallType][mallSubType] = make([]*battery.MallItem, 0)
        }
        (*mapMallType)[mallType][mallSubType] = append((*mapMallType)[mallType][mallSubType], good)

        //iapid -> goods
        iapId := g.GetIapid()
        if iapId != "" {
            (*mapIapGoods)[iapId] = good
        }
    }

    //mapGoods.Print()
    //mapMallType.Print()
    //mapIapGoods.Print()

    gm.switchCache()

    return
}

func (gm *GoodsCacheManager) getMallItemFromDBMallItem(mallItem *battery.DBMallItem) *battery.MallItem {
    return &battery.MallItem{
        Id:              proto.Uint64(mallItem.GetId()),
        MallType:        mallItem.GetMallType().Enum(),
        MallSubType:     mallItem.GetMallSubType().Enum(),
        PosIndex:        proto.Uint32(mallItem.GetPosIndex()),
        Discount:        proto.Uint32(mallItem.GetDiscount()),
        Price:           mallItem.GetPrice(),
        Items:           mallItem.GetItems(),
        Amountperuser:   proto.Uint32(mallItem.GetAmountperuser()),
        Amountpergame:   proto.Uint32(mallItem.GetAmountpergame()),
        Amountperday:    proto.Uint32(mallItem.GetAmountperday()),
        Bestdeal:        proto.Bool(mallItem.GetBestdeal()),
        Tesell:          proto.Bool(mallItem.GetTesell()),
        Expiretimestamp: proto.Int64(mallItem.GetExpiretimestamp()),
        Iapid:           proto.String(mallItem.GetIapid()),
        Icon:            proto.Int32(mallItem.GetIcon()),
        Name:            proto.String(mallItem.GetName()),
        Description:     proto.String(mallItem.GetDescription()),
        Label:           proto.Int32(mallItem.GetLabel()),
        ItemFlags:       mallItem.GetItemFlags(),
        Multiple:        mallItem.Multiple,
    }
}

func (gm *GoodsCacheManager) switchCache() (fail_reason int32, err error) {
    gm.Switch()
    //xylog.DebugNoId("now Goods cache switch to %d", gm.Major())
    return
}

// 解锁商品 checkpoint-> goodsid
type MAPCheckPointGoods map[uint32]uint64

type CheckPointUnlockGoodsCache struct {
    M MAPCheckPointGoods
}

func (c *CheckPointUnlockGoodsCache) Clear() {
    c.M = make(MAPCheckPointGoods, 0)
}

var DefCheckPointUnlockGoodsManager = NewCheckPointUnlockGoodsConfig()

type CheckPointUnlockGoodsManager struct {
    cache [2]CheckPointUnlockGoodsCache
    xycache.CacheBase
}

func NewCheckPointUnlockGoodsConfig() *CheckPointUnlockGoodsManager {
    return &CheckPointUnlockGoodsManager{}
}

func (m *CheckPointUnlockGoodsManager) InitWhileStart() {
    m.Init()
    m.Reload()
}

func (m *CheckPointUnlockGoodsManager) Init() {

}

func (m *CheckPointUnlockGoodsManager) Reload() (err error) {
    secondary := m.SecondaryCache()

    dbcheckpointUnlockGoodsConfig := make([]*battery.DBCheckPointUnlockGoodsConfig, 0)
    err = DefCacheDB.LoadCheckPointUnlockGoodsConfig(&dbcheckpointUnlockGoodsConfig)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("LoadCheckPointUnlockGoodsConfig failed %v", err)
        return
    }
    secondary.Clear()
    for _, v := range dbcheckpointUnlockGoodsConfig {
        secondary.M[v.GetCheckPointId()] = v.GetGoodsId()
    }
    m.Switch()
    m.Print()
    return
}

func (m *CheckPointUnlockGoodsManager) Print() {
    xylog.DebugNoId("------CheckPointUnlockGoodsConfig cache -------")
    for source, items := range m.MajorCache().M {
        xylog.DebugNoId("----source %v ------\n %v", source, items)
    }
}

func (m *CheckPointUnlockGoodsManager) MajorCache() *CheckPointUnlockGoodsCache {
    return &(m.cache[m.Major()])
}

func (m *CheckPointUnlockGoodsManager) SecondaryCache() *CheckPointUnlockGoodsCache {
    return &(m.cache[m.Secondary()])
}

// 更具需要解锁的关卡获取需要购买的商品id
func (m *CheckPointUnlockGoodsManager) GoodsId(checkpointId uint32) (goodId uint64) {
    major := m.MajorCache()
    goodId = major.M[checkpointId]
    return
}
