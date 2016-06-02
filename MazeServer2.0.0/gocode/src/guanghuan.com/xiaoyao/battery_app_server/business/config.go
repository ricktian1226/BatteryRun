package batteryapi

import (
    "fmt"
    "guanghuan.com/xiaoyao/common/cache"
    xylog "guanghuan.com/xiaoyao/common/log"
    xyversion "guanghuan.com/xiaoyao/common/version"
    //"guanghuan.com/xiaoyao/superbman_server/cache"
    "time"
)

type ApiConfig struct {
    Name                        string            // 1 配置名称
    DefaultStamina              uint32            // 2 初始体力
    DefaultMaxRegStamina        int32             // 3 体力自动增长的最大值
    DefaultDiamond              uint32            // 4 初始的钻石数量
    DefaultStaminRegIntervalSec int64             // 5 体力自动恢复的时间
    DefaultGiftValidTimeSec     int64             // 6 体力请求的有效期
    DefaultGiftGiveCooldown     int64             // 7 向好友发起体力请求的时间间隔
    DefaultGiftAskCooldown      int64             // 8 向好友赠送体力的时间间隔
    DefaultGiftAcceptMaxCount   int32             // 9 最多接受的体力数量
    GiftMaxBatchSize            int               // 10 一次操作最多包含多少条体力请求
    DefaultSysAccount           string            // 11 系统账号
    EnableSecurity              bool              // 12 是否对通讯数据进行加密
    MailboxLimit                int               // 13 邮箱上限
    IsProduction                bool              // 14 是否生产环境
    MaxRequestSize              int64             // 15 每个请求数据的最大值(byte)
    LogLevel                    int32             // 16 日志级别
    MaxFriendsRequestCount      int               // 17 每次最多向多少个好友发起请求(包括查询好友数据，请求体力，等)
    MinClientVersion            xyversion.Version // 18 支持最低的客户端版本
    MaxClientVersion            xyversion.Version // 19 支持最高的客户端版本
    CurClientVersion            xyversion.Version // 20 当前最新的客户端版本
    TransactionNodeCount        int               // 21 事务服务节点数（路由选择：按照uid后三位 mode 该服务节点数选择）
    NatsTimeOut                 int               // 22 nats消息超市时间，单位:秒
    Panic                       bool              // 23 是否打开程序crash panic开关
    InvalidScore                uint64            // 非法分数边界值
}

func (config *ApiConfig) String() (str string) {
    str = fmt.Sprintf(`
	Name         = %s
	IsProduction = %t
	--
	Cur Client Version = %s
	Min Client Version = %s
	Max Client Version = %s
	--
	DefaultStamina              = %d
	DefaultMaxRegStamina        = %d
	DefaultStaminRegIntervalSec = %d
	--
	DefaultDiamond          = %d
	--
	DefaultGiftValidTimeSec = %d
	DefaultGiftAskCooldown  = %d
	DefaultGiftGiveCooldown = %d
	GiftMaxBatchSize        = %d
	MaxFriendsRequestCount  = %d
	--
	DefaultSysAccount = %s
	--
	EnableSecurity    = %t
	--
	MailBoxLimit      = %d
	--
	MaxRequestSize    = %d
	LogLevel          = %d
	--
	NatsTimeOut       = %d
	--
	TransactionNodeCount = %d
	--
	Panic             = %t
	InvalidScore      =%d
	`,
        config.Name,
        config.IsProduction,
        config.CurClientVersion.String(),
        config.MinClientVersion.String(),
        config.MaxClientVersion.String(),
        config.DefaultStamina,
        config.DefaultMaxRegStamina,
        config.DefaultStaminRegIntervalSec,
        config.DefaultDiamond,
        config.DefaultGiftValidTimeSec,
        config.DefaultGiftAskCooldown,
        config.DefaultGiftGiveCooldown,
        config.GiftMaxBatchSize,
        config.MaxFriendsRequestCount,
        config.DefaultSysAccount,
        config.EnableSecurity,
        config.MailboxLimit,
        config.MaxRequestSize,
        config.LogLevel,
        config.NatsTimeOut,
        config.TransactionNodeCount,
        config.Panic,
        config.InvalidScore)
    return
}

func (config *ApiConfig) Init(name string) {
    config.Name = name
    config.IsProduction = true
    config.CurClientVersion = DefVersion
    config.MinClientVersion = DefMinClientVersion
    config.MaxClientVersion = DefMaxClientVersion
    config.DefaultStamina = 5
    config.DefaultMaxRegStamina = 5
    config.DefaultDiamond = 100
    config.DefaultStaminRegIntervalSec = 600 // 10 min
    config.DefaultGiftValidTimeSec = int64(7 * 24 * time.Hour / time.Second)
    config.DefaultGiftAskCooldown = int64(24 * time.Hour / time.Second) // 24小时
    config.DefaultGiftGiveCooldown = int64(2 * time.Hour / time.Second) // 2小时
    config.DefaultGiftAcceptMaxCount = -1
    config.GiftMaxBatchSize = 100
    config.DefaultSysAccount = "sys"
    config.EnableSecurity = true
    config.MailboxLimit = 30
    config.MaxRequestSize = 20 * 1024 // 20kb
    if !config.IsProduction {
        config.LogLevel = int32(xylog.DebugLevel) // 测试环境，默认打开debug
    } else {
        config.LogLevel = int32(xylog.InfoLevel) // 生产环境，默认关闭debug
    }
    config.MaxFriendsRequestCount = 10
    config.TransactionNodeCount = 1 //下游服务节点数，默认一个
    config.NatsTimeOut = 10         //nats消息超时时间，单位：秒
    config.Panic = true
    config.InvalidScore = 500000
}

//配置结构体定义
type ConfigStruct struct {
    Configs ApiConfig //配置项
}

func (c *ConfigStruct) Clear() {

}

//配置缓存
type ConfigCache struct {
    caches [2]ConfigStruct
    xycache.CacheBase
}

func (c *ConfigCache) Init() {
    for i := 0; i < 2; i++ {
        c.caches[i].Clear()
    }
}

func (c *ConfigCache) Master() *ConfigStruct {
    return &(c.caches[c.Major()])
}

func (c *ConfigCache) Slave() *ConfigStruct {
    return &(c.caches[c.Secondary()])
}

func (c *ConfigCache) Configs() *ApiConfig {
    return &(c.Master().Configs)
}

func NewConfigCache() (c *ConfigCache) {
    c = &ConfigCache{}
    c.Init()
    return
}

var DefConfigCache = NewConfigCache()
