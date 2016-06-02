// xybusiness
package xybusiness

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    beegoconf "github.com/astaxie/beego/config"

    "guanghuan.com/xiaoyao/common/db"
    "guanghuan.com/xiaoyao/common/log"

    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

const (
    DB_IOS       = "briosdb"
    DB_ANDROID   = "brandroiddb"   //android 玩家数据库
    DB_COMMON    = "brcommondb"    //互通数据库
    DB_LOG       = "brlogdb"       //业务日志数据库
    DB_STATISTIC = "brstatisticdb" //业务统计数据库
)

//数据库对象索引
const (
    DB_INDEX_IOS = iota
    DB_INDEX_ANDROID
    DB_INDEX_COMMON
    DB_INDEX_LOG
    DB_INDEX_STATISTIC
    DB_INDEX_MAX
)

//数据库名
var DB_NAME = []string{
    DB_IOS,       //ios 玩家数据库
    DB_ANDROID,   //android 玩家数据库
    DB_COMMON,    //common 玩家公共数据数据库
    DB_LOG,       //日志数据库
    DB_STATISTIC, //统计数据库
}

var INI_CONFIG_ITEM_DB = []string{
    "DBBriosdb",       //ios 玩家数据库主配置项
    "DBBrandroiddb",   //android 玩家数据库主配置项
    "DBBrcommondb",    //common 玩家数据库主配置项
    "DBBrlogdb",       //日志数据库主配置项
    "DBBrstatisticdb", //统计数据库主配置项
}

//common库的数据库表名
var INI_CONFIG_ITEM_COMMON_COLLECTION = []string{
    "apiconfig",
    "announcementconfig",
    "beforegamerandomweightconfig",
    "goodsconfig",
    "jigsawconfig",
    "lottoserialnumslotconfig",
    "lottoslotitemconfig",
    "lottoweightconfig",
    "mailinfoconfig",
    "missionconfig",
    "pickupweightconfig",
    "propconfig",
    "roleinfoconfig",
    "rolelevelbonusconfig",
    "runeconfig",
    "signawardconfig",
    "signinactivityconfig",
    "advertisementconfig",
    "advertisementspaceconfig",
    "tipconfig",
    "newaccountpropconfig",
    "shareactivity",
    "shareawards",
    "checkpointunlockgoodsconfig",
    "ranklist",
    "blacklist",
    "tpidmap",
    "useraccomplishment",
    "usercheckpoint",
    "devicetokenblacklist",
    "memcache",
    "banneduser",
    "debuguser",
    "useridentitycounter",
}

const ( //注意index和上面定义的collection的下标是一致的
    BUSINESS_COMMON_COLLECTION_INDEX_APICONFIG = iota
    BUSINESS_COMMON_COLLECTION_INDEX_ANNOUNCEMENT_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_BEFOREGAMERANDOM_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_GOODS_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_JIGSAW_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSERIALNUMSLOT_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSLOTITEM_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_LOTTOWEIGHT_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_MAILINFO_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_MISSION_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_PICKUPWEIGHT_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_PROP_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_ROLEINFO_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_ROLELEVELBONUS_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_RUNE_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_SIGNAWARD_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_SIGNINACTIVITY_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_ADVERTISEMENT_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_ADVERTISEMENT_SPACE_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_TIP_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_NEWACCOUNTPROP_CONFIG
    BUSINESS_COMMON_COLLECTION_INDEX_SHARE_ACTIVITY
    BUSINESS_COMMON_COLLECTION_INDEX_SHARE_AWARDS
    Business_COMMON_COLLECTION_INDEX_CHECKPOINTUNLOCK_GOODS_CONFIG
    Business_COMMON_COLLECTION_INDEX_RANKLIST
    Business_COMMON_COLLECTION_INDEX_BLACKLIST

    BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP
    BUSINESS_COMMON_COLLECTION_INDEX_USERACCOMPLISHMENT
    BUSINESS_COMMON_COLLECTION_INDEX_USERCHECKPOINT
    BUSINESS_COMMON_COLLECTION_INDEX_DEVICETOKENBLACKLIST
    BUSINESS_COMMON_COLLECTION_INDEX_MEMCACHE
    BUSINESS_COMMON_COLLECTION_INDEX_BANNEDUSER
    BUSINESS_COMMON_COLLECTION_INDEX_DEBUGUSER

    BUSINESS_COMMON_COLLECTION_INDEX_USERIDENTITYCOUNTER

    BUSINESS_COMMON_COLLECTION_INDEX_IAPSTATISTIC

    BUSINESS_COMMON_COLLECTION_INDEX_LOG

    BUSINESS_COMMON_COLLECTION_INDEX_MAX
)

//业务数据库名称
var INI_CONFIG_ITEM_COLLECTION = []string{
    "account",                  //玩家账户信息表
    "consumable",               //玩家消耗品表
    "friendmail",               //玩家好友邮件表
    "game",                     //玩家游戏信息表
    "jigsaw",                   //玩家拼图信息表
    "lottosysinfo",             //玩家抽奖信息表
    "lottotransaction",         //玩家抽奖事务表
    "memcache",                 //
    "receipt",                  //购买票据信息表
    "roleinfo",                 //玩家角色信息表
    "rune",                     //玩家符文信息表
    "shoppingtransaction",      //购买事务信息表
    "staminagiveapplylog",      //体力相关时间戳信息表
    "systemmaillist",           //系统邮件信息表
    "userdonecollectedmission", //玩家已完成任务信息表
    "usermission",              //玩家任务信息表
    "usersigninactivity",       //玩家签到信息表
    "iaptransaction",           //玩家iap信息表
    "friendmailcount",          //玩家好友邮件总数表
    "usershareinfo",            //　玩家分享信息
    "sdkorder",                 // sdk订单信息
}

const ( //注意index和上面定义的collection的下标是一致的
    BUSINESS_COLLECTION_INDEX_ACCOUNT = iota
    BUSINESS_COLLECTION_INDEX_CONSUMABLE
    BUSINESS_COLLECTION_INDEX_FRIENDMAIL
    BUSINESS_COLLECTION_INDEX_GAME
    BUSINESS_COLLECTION_INDEX_JIGSAW
    BUSINESS_COLLECTION_INDEX_LOTTOSYSINFO
    BUSINESS_COLLECTION_INDEX_LOTTOTRANSACTION
    BUSINESS_COLLECTION_INDEX_MEMCACHE
    BUSINESS_COLLECTION_INDEX_RECEIPT
    BUSINESS_COLLECTION_INDEX_ROLEINFO
    BUSINESS_COLLECTION_INDEX_RUNE
    BUSINESS_COLLECTION_INDEX_SHOPPINGTRANSACTION
    BUSINESS_COLLECTION_INDEX_STAMINAGIVEAPPLYLOG
    BUSINESS_COLLECTION_INDEX_SYSTEMMAILLIST
    BUSINESS_COLLECTION_INDEX_USERDONECOLLECTIONMISSION
    BUSINESS_COLLECTION_INDEX_USERMISSION
    BUSINESS_COLLECTION_INDEX_USERSIGNINACTIVITY
    BUSINESS_COLLECTION_INDEX_IAPTRANSACTION
    BUSINESS_COLLECTION_INDEX_FRIENDMAILCOUNT
    BUSINESS_COLLECTION_INDEX_USERSHARE
    BUSINESS_COLLECTION_INDEX_SDKORDER
    BUSINESS_COLLECTION_INDEX_MAX
)

////内购上报数据库名称
//var INI_CONFIG_ITEM_IAPSTATISTIC_COLLECTION = []string{
//	"iapstatistic", //内购上报信息表
//}

//-------------------------业务数据表信息管理器-------------------------

const (
    INI_CONFIG_ITEM_SERVER_CONFIG               = "Server::config"
    INI_CONFIG_ITEM_SERVER_NATSURL              = "Server::natsurl"
    INI_CONFIG_ITEM_SERVER_APNNATSURL           = "Server::apnnatsurl"
    INI_CONFIG_ITEM_SERVER_ALERT_NATSURL        = "Server::alertnatsurl"
    INI_CONFIG_ITEM_SERVER_IAPSTATISTIC_NATSURL = "Server::iapstatisticnatsurl"
    INI_CONFIG_ITEM_SERVER_TEST                 = "Server::testenv"
    INI_CONFIG_ITEM_SERVER_PID                  = "Server::pid"

    INI_CONFIG_ITEM_LOG_DB = "DBBrlogdb::logdb"

    INI_CONFIG_ITEM_IAPSTATISTIC_DB = "DBBrstatisticdb::iapstatistic"
)

//const (
//	BUSINESS_LOG_COLLECTION_INDEX_LOG = iota
//	BUSINESS_LOG_COLLECTION_INDEX_MAX
//)

type BusinessCollectinConfig struct {
    item, subItem string
    platform      battery.PLATFORM_TYPE
}

//保存业务模块对应的数据库信息
type BusinessCollection struct {
    index         int                   //同一平台下的数据库信息索引
    dburl, dbname string                //数据库url,数据库名称
    platform      battery.PLATFORM_TYPE //平台类型
}

func (c *BusinessCollection) Detail() (string, string, battery.PLATFORM_TYPE, int) {
    return c.dburl, c.dbname, c.platform, c.index
}

type BusinessCollections struct {
    collections []BusinessCollection
}

//插入业务模块对应的数据库信息
// index int BUSINESS_COLLECTION_INDEX
// platform battery.PLATFORM_TYPE 平台类型
// dbUrl string 数据库url链接
// dbName string 数据库名称
func (m *BusinessCollections) Insert(index int, platform battery.PLATFORM_TYPE, dbUrl, dbName string) {
    if dbUrl != "" && dbName != "" {
        m.collections = append(m.collections, BusinessCollection{
            index:    index,
            dburl:    dbUrl,
            dbname:   dbName,
            platform: platform,
        })
    }
}

func (m *BusinessCollections) Collections() *([]BusinessCollection) {
    return &(m.collections)
}

func NewBusinessCollections() *BusinessCollections {
    return &BusinessCollections{
        collections: make([]BusinessCollection, 0),
    }
}

type MAPDBSesssion map[int]xydb.DBInterface

type MAPPlatfrom2DBSesssions map[battery.PLATFORM_TYPE]MAPDBSesssion

type BusinessDBSessionManager struct {
    sessions MAPPlatfrom2DBSesssions
}

//全局变量，业务数据库session管理器
var DefBusinessDBSessionManager = NewBusinessDBSessionManager()

func (m *BusinessDBSessionManager) Insert(index int, platform battery.PLATFORM_TYPE, dbInterface xydb.DBInterface) {
    if _, ok := m.sessions[platform]; !ok {
        m.sessions[platform] = make(MAPDBSesssion, 0)
    }

    m.sessions[platform][index] = dbInterface
}

func (m *BusinessDBSessionManager) Get(index int, platform battery.PLATFORM_TYPE) xydb.DBInterface {
    if m.sessions == nil {
        return nil
    }

    if v, ok := m.sessions[platform]; !ok {
        return nil
    } else if dbInterface, ok := v[index]; !ok {
        return nil
    } else {
        return dbInterface
    }
}

func (m *BusinessDBSessionManager) Print() {
    xylog.DebugNoId("-----------BusinessDBSessionManager begin ----------")
    for platform, s := range m.sessions {
        xylog.DebugNoId("-----platform(%v)-----", platform)
        for index, ss := range s {
            xylog.DebugNoId("%v : %v", index, ss)
        }
    }
    xylog.DebugNoId("-----------BusinessDBSessionManager end ----------")
}

func NewBusinessDBSessionManager() *BusinessDBSessionManager {
    return &BusinessDBSessionManager{
        sessions: make(MAPPlatfrom2DBSesssions, 0),
    }
}

//从配置项中解析业务表对应的数据库连接信息
// config beegoconf.ConfigContainer 配置项信息容器
// businessCollections *BusinessCollections 保存业务表配置项信息容器
func ParseBusinessCollections(config beegoconf.ConfigContainer, businessCollections *BusinessCollections) {
    //common
    for i, v := range INI_CONFIG_ITEM_COMMON_COLLECTION {
        parseBusinessCollection(config, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN, INI_CONFIG_ITEM_DB[DB_INDEX_COMMON]+"::"+v, DB_NAME[DB_INDEX_COMMON], i, businessCollections)
    }

    //ios/android
    platforms := []battery.PLATFORM_TYPE{battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS, battery.PLATFORM_TYPE_PLATFORM_TYPE_ANDROID}
    dbs := []string{DB_IOS, DB_ANDROID}
    itemDbs := []string{"DBBriosdb", "DBBrandroiddb"}
    for n, platform := range platforms {
        for i, v := range INI_CONFIG_ITEM_COLLECTION {
            parseBusinessCollection(config, platform, itemDbs[n]+"::"+v, dbs[n], i, businessCollections)
        }
    }

    //log
    parseBusinessCollection(config, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN, INI_CONFIG_ITEM_LOG_DB, DB_LOG, BUSINESS_COMMON_COLLECTION_INDEX_LOG, businessCollections)

    //iapstatistic
    parseBusinessCollection(config, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN, INI_CONFIG_ITEM_IAPSTATISTIC_DB, DB_STATISTIC, BUSINESS_COMMON_COLLECTION_INDEX_IAPSTATISTIC, businessCollections)

}

//从文件配置项中解析业务数据库信息
// configs beegoconf.ConfigContainer 解析配置文件获取的业务数据库配置信息
// platform battery.PLATFORM_TYPE 平台类型
// item string 配置项名称
// dbname string 数据库名
// index int 业务数据库表索引枚举值
// businessCollections *BusinessCollections 业务数据库表信息列表
func parseBusinessCollection(configs beegoconf.ConfigContainer, platform battery.PLATFORM_TYPE, item, dbname string, index int, businessCollections *BusinessCollections) {
    tmp := configs.String(item)
    if tmp != "" {
        tmp = tmp + "/" + dbname
        businessCollections.Insert(index, platform, tmp, dbname)
    }
}

// 数据库表名
const (
    //动态业务表
    DB_TABLE_ACCOUNT                    = "account"                  // 玩家账户
    DB_TABLE_USER_ACCOMPLISHMENT        = "useraccomplishment"       // 玩家成就信息
    DB_TABLE_DEVICETOKEN_BLICKLIST      = "devicetokenblacklist"     // devicetoken黑名单
    DB_TABLE_TPID_MAP                   = "tpidmap"                  // 第三方账户信息表
    DB_TABLE_GAMEDATA                   = "game"                     // 游戏数据信息表
    DB_TABLE_JIGSAW                     = "jigsaw"                   // 玩家拼图信息表
    DB_TABLE_RUNE                       = "rune"                     // 玩家符文信息
    DB_TABLE_CONSUMABLE                 = "consumable"               // 玩家消耗类道具信息
    DB_TABLE_SIGN_IN_RECORD             = "usersigninactivity"       // 玩家签到活动信息
    DB_TABLE_USER_MISSION               = "usermission"              // 用户任务
    DB_TABLE_USER_DONECOLLECTED_MISSION = "userdonecollectedmission" // 用户donecollectedmission
    DB_TABLE_GIFT                       = "gift"                     // 体力请求（无用？）
    DB_TABLE_RECEIPT                    = "receipt"                  // 购买凭据信息
    DB_TABLE_IAPTRANSACTION             = "iaptransaction"           // 内购交易记录
    DB_TABLE_PUSH_Natice                = "pushnotification"         // 推送命令记录表
    DB_TABLE_PUSH_RECORD                = "pushrecord"               // 推送记录
    DB_TABLE_ANNOUNCEMENT               = "announcement"             // 玩家公告信息
    DB_TABLE_SYS_LOTTO_INFO             = "lottosysinfo"             // 用户系统抽奖信息
    DB_TABLE_LOTTO_TRANSACTION          = "lottotransaction"         // 抽奖事务信息
    DB_TABLE_SHOPPING_TRANSACTION       = "shoppingtransaction"      // 购买交易信息
    DB_TABLE_USER_CHECK_POINT           = "usercheckpoint"           // 玩家记忆点信息
    DB_TABLE_ROLEINFO                   = "roleinfo"                 // 玩家角色信息
    DB_TABLE_SYSTEMMAIL                 = "systemmaillist"           // 玩家系统邮件信息
    DB_TABLE_STAMINA_GIVEAPPLY_LOG      = "staminagiveapplylog"      // 玩家体力赠送/请求信息
    DB_TABLE_FRIENDMAIL                 = "friendmail"               // 玩家好友邮件信息
    DB_TABLE_MEMCACHE                   = "memcache"                 // 玩家缓存信息
    DB_TABLE_INVALID_DEVICE_TOKEN       = "invalid_device_token"     // 无效的设备id
    DB_TABLE_BANNED_USER                = "banneduser"               // 封号的玩家信息
    DB_TABLE_DEBUG_USER                 = "debuguser"                // 打开调试开关的玩家信息
    DB_TABLE_FRIENDMAILCOUNT            = "friendmailcount"          // 玩家好友邮件数目信息
    DB_TABLE_USERIDENTITYCOUNTER        = "useridentitycounter"      //玩家计数器
    DB_TABLE_USERSHARED_RECORD          = "usershareinfo"            // 玩家分享信息
    DB_TABLE_SDKORDER                   = "sdkorder"                 // sdk订单信息
    DB_TABLE_RANKLIST                   = "ranklist"                 // 玩家排行榜

    //静态配置表
    DB_TABLE_JIGSAW_CONFIG                 = "jigsawconfig"                 // 拼图配置信息
    DB_TABLE_SIGN_IN_ACTIVITY              = "signinactivityconfig"         // 活动配置信息
    DB_TABLE_SIGN_IN_ITEM                  = "signawardconfig"              // 活动奖励配置信息
    DB_TABLE_MISSION                       = "missionconfig"                // 任务配置信息
    DB_TABLE_GOODS                         = "goodsconfig"                  // 商品配置信息
    DB_TABLE_PROP                          = "propconfig"                   // 道具配置信息
    DB_TABLE_LOTTO_SLOTITEMS               = "lottoslotitemconfig"          // 抽奖格子子项配置信息
    DB_TABLE_LOTTO_WEIGHT                  = "lottoweightconfig"            // 抽奖权重配置信息
    DB_TABLE_LOTTO_SERIALNUM_SLOT          = "lottoserialnumslotconfig"     // 抽奖序号礼包配置信息
    DB_TABLE_RUNE_CONFIG                   = "runeconfig"                   // 符文配置信息
    DB_TABLE_ROLEINFO_CONFIG               = "roleinfoconfig"               // 角色配置信息
    DB_TABLE_MAILINFO_CONFIG               = "mailinfoconfig"               // 邮件配置信息
    DB_TABLE_BEFOREGAME_RANDOM_WEIGHT      = "beforegamerandomweightconfig" // 赛前随机道具权重配置信息
    DB_TABLE_PICKUP_WEIGHT                 = "pickupweightconfig"           // 游戏中收集物权重配置信息
    DB_TABLE_ROLE_LEVEL_BONUS              = "rolelevelbonusconfig"         // 角色加成配置信息
    DB_TABLE_ANNOUNCEMENT_CONFIG           = "announcementconfig"           // 公告配置信息
    DB_TABLE_ADVERTISEMENT_CONFIG          = "advertisementconfig"          // 广告配置信息
    DB_TABLE_ADVERTISEMENTSPACE_CONFIG     = "advertisementspaceconfig"     // 广告位配置信息
    DB_TABLE_TIP_CONFIG                    = "tipconfig"                    // 提示配置信息
    DB_TABLE_NEWACCOUNTPROP_CONFIG         = "newaccountpropconfig"         // 登录礼包配置信息
    DB_TABLE_SHARE_ACTIVITY                = "shareactivity"                // 分享活动配置
    DB_TABLE_SHARE_AWARDS                  = "sharewards"                   // 分享奖励配置
    DB_TABLE_CHECKPOINTUNLOCK_GOODS_CONFIG = "checkpointunlockgoodsconfig"  // 解锁关卡商品配置
    DB_TABLE_BLACKLIST                     = "blacklist"                    // 玩家黑名单
    //操作日志表
    DB_TABLE_GIFT_LOG        = "giftlog"        // 体力请求日志
    DB_TABLE_GAME_LOG        = "gamelog"        // 游戏日志
    DB_TABLE_SHOPPING_LOG    = "shoppinglog"    // 购买日志
    DB_TABLE_ACCOUNT_LOG     = "accountlog"     // 账户日志
    DB_TABLE_IAP_LOG         = "iaplog"         //内购日志
    DB_TABLE_LOTTO_LOG       = "lottolog"       // 抽奖日志
    DB_TABLE_CHECKPOINT_LOG  = "checkpointlog"  // 记忆点日志
    DB_TABLE_PUSH_LOG        = "pushlog"        // 推送日志
    DB_TABLE_MAINTENANCE_LOG = "maintenancelog" // 运营操作日志

    //内购信息
    DB_TABLE_IAPSTATISTIC = "iapstatistic" // 内购上报信息表

)

type XYBusinessDB struct {
    xydb.XYDB
}

func NewXYBusinessDB(dburl, dbname string) *XYBusinessDB {
    return &XYBusinessDB{
        XYDB: *xydb.NewXYDB(dburl, dbname),
    }
}

//根据sid查询tpid
// sids []string 玩家sid列表
// idSource battery.ID_SOURCE id来源
//return:
// tpids *[]*battery.IDMap
func (db *XYBusinessDB) QueryTpidBySid(selector interface{}, sid string, idSource battery.ID_SOURCE, tpid *battery.IDMap, consistency mgo.Mode) (err error) {
    condition := bson.M{"source": idSource, "sid": sid}
    err = db.GetOneData(DB_TABLE_TPID_MAP, condition, selector, tpid, consistency)
    return
}

//根据sid查询tpid列表
// sids []string 玩家sid列表
// idSource battery.ID_SOURCE id来源
//return:
// tpids *[]*battery.IDMap
func (db *XYBusinessDB) QueryTpidsBySids(selector interface{}, sids []string, idSource battery.ID_SOURCE, tpids *[]*battery.IDMap) (err error) {
    condition := bson.M{"source": idSource, "sid": bson.M{"$in": sids}}
    err = db.GetAllData(DB_TABLE_TPID_MAP, condition, selector, 0, tpids, mgo.Monotonic) //读写分离，允许从备节点上读取
    return
}

//根据uid查询tpid
// uid string 玩家uid
// idSource battery.ID_SOURCE id来源
//return:
// tpid *battery.IDMap
func (db *XYBusinessDB) QueryTpidByUid(selector interface{}, uid string, tpid *battery.IDMap, consistency mgo.Mode) (err error) {
    condition := bson.M{"gid": uid}
    err = db.GetOneData(DB_TABLE_TPID_MAP, condition, selector, tpid, consistency) //读写分离，允许从备节点上读取
    return
}

//根据uid查询tpid列表
// uids []string 玩家uid列表
// idSource battery.ID_SOURCE id来源
//return:
// tpids *[]*battery.IDMap
func (db *XYBusinessDB) QueryTpidsByUids(selector interface{}, uids []string, tpids *[]*battery.IDMap) (err error) {
    condition := bson.M{"gid": bson.M{"$in": uids}}
    err = db.GetAllData(DB_TABLE_TPID_MAP, condition, selector, 0, tpids, mgo.Monotonic) //读写分离，允许从备节点上读取
    return
}

//根据uid查询account列表
// uids []string 玩家uid列表
//return:
// dbAccounts *[]*battery.DBAccount 玩家的account信息列表
func (db *XYBusinessDB) QueryAccountsByUids(selector interface{}, uids []string, dbAccounts *[]*battery.DBAccount) (err error) {
    condition := bson.M{"uid": bson.M{"$in": uids}}
    err = db.GetAllData(DB_TABLE_ACCOUNT, condition, selector, 0, dbAccounts, mgo.Monotonic) //读写分离，允许从备节点上读取
    return
}

//根据uid查询玩家账户信息列表
// uid string 玩家标识
// account *battery.DBAccount 账户信息
// consistency mgo.Mode 查询模式
func (db *XYBusinessDB) GetAccountDirect(uid string, selector interface{}, account *battery.DBAccount, consistency mgo.Mode) (err error) {
    var condition = bson.M{"uid": uid}
    err = db.GetOneData(DB_TABLE_ACCOUNT, condition, selector, account, consistency)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[DB] GetAccountDirect failed : %v ", err)
    }

    return
}

//根据uids列表查询多个玩家账户信息列表
// uid string 玩家标识
// account *battery.DBAccount 账户信息
// consistency mgo.Mode 查询模式
func (db *XYBusinessDB) GetAccountsDirect(uids []string, selector interface{}, accounts *[]*battery.DBAccount, consistency mgo.Mode) (err error) {
    var condition = bson.M{"uid": bson.M{"$in": uids}}
    err = db.GetAllData(DB_TABLE_ACCOUNT, condition, selector, 0, accounts, consistency)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[DB] GetAccountsDirect failed : %v ", err)
    }

    return
}

//根据identity查询玩家账户信息列表
// uid string 玩家标识
// account *battery.DBAccount 账户信息
// consistency mgo.Mode 查询模式
func (db *XYBusinessDB) GetAccountDirectByIdentity(identity string, selector interface{}, account *battery.DBAccount, consistency mgo.Mode) (err error) {
    var condition = bson.M{"identitystring": identity}
    err = db.GetOneData(DB_TABLE_ACCOUNT, condition, selector, account, consistency)
    if err != xyerror.ErrOK {
        xylog.ErrorNoId("[DB] GetAccountDirectByIdentity failed : %v ", err)
    }

    return
}

//查询打开调试开关的玩家标识列表
func (db *XYBusinessDB) QueryDebugUsers() (ids []interface{}, err error) {
    items := make([]*battery.DBDebugUserItem, 0)
    err = db.GetAllData(DB_TABLE_DEBUG_USER, bson.M{}, bson.M{}, 0, &items, mgo.Strong)
    if err == nil {
        ids = make([]interface{}, 0)
        for _, item := range items {
            ids = append(ids, item.GetId())
        }
    }

    return
}
