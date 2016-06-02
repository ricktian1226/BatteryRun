package batteryapi

import (
	"fmt"
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyversion "guanghuan.com/xiaoyao/common/version"
	"time"
)

//var SetBundleId Set = make(Set, 2) //bundle_id列表

type ApiConfig struct {
	Name                        string            // 1 配置名称
	DefaultStamina              uint32            // 2 初始体力
	DefaultMaxRegStamina        int32             // 3 体力自动增长的最大值
	DefaultStaminRegIntervalSec int64             // 4 体力自动恢复的时间
	DefaultDiamond              uint32            // 3 初始的钻石数量
	DefaultChip                 uint32            // 3 初始的碎片数量
	DefaultCoin                 uint32            // 3 初始的金币数量
	DefaultGiftValidTimeSec     int64             // 6 体力请求的有效期
	DefaultGiftGiveCooldown     int64             // 7 向好友发起体力请求的时间间隔
	DefaultGiftAskCooldown      int64             // 8 向好友赠送体力的时间间隔
	DefaultGiftNotifyCount      int               //   体力请求发送apn推送的阈值
	DefaultGiftAcceptMaxCount   int32             // 9 最多接受的体力数量
	GiftMaxBatchSize            int               // 10 一次操作最多包含多少条体力请求
	DefaultSysAccount           string            // 11 系统账号
	EnableSecurity              bool              // 12 是否对通讯数据进行加密
	MailboxLimit                int               // 13 邮箱上限
	MaxFriendsRequestCount      int               // 17 每次最多向多少个好友发起请求(包括查询好友数据，请求体力，等)
	LogLevel                    int32             // 16 日志级别
	IsProduction                bool              // 2 是否生产环境
	MaxRequestSize              int64             // 15 每个请求数据的最大值(byte)
	MinClientVersion            xyversion.Version // 18 支持最低的客户端版本
	MaxClientVersion            xyversion.Version // 19 支持最高的客户端版本
	CurClientVersion            xyversion.Version // 20 当前最新的客户端版本
	//TransactionNodeCount           uint32            // 4 事务服务节点数（路由选择：按照uid后三位 mode 该服务节点数选择）
	ChannelCount                   uint32  // 5 用户channel节点数（路由选择：按照uid后三位 mode 该服务节点数选择）
	ChannelMaxMsg                  uint32  // 6 channel 可缓存的最大消息数
	BundleId                       string  // 7 bundle_id
	LottoSlotCount                 uint32  //抽奖格子数
	LottoInitUserValue             int32   //默认权值
	LottoCostPerTime               int32   //抽奖消耗
	LottoDeduct                    int32   //抽水比率
	SysLottoFreeCount              int32   //免费次数初值
	SysLottoRefreshTime            int64   //系统抽奖刷新时间（单位:秒）
	AfterGameLottoDeleteSlotLimit  int     //游戏后抽奖删除格子次数上限
	AfterGameLottoMallItems        string  //游戏后抽奖商品列表
	MissionCountLimit              int     //任务上限（教学任务+主线任务）
	DailyMissionCountLimit         int     //日常任务上限
	DailyMissionRefreshHour        int     //日常任务刷新时间，第几个小时
	CheckPointIdNum                uint32  //记忆点个数
	CheckPointGlobalRankReLoadSecs int64   //记忆点全局排行榜重载时间（单位秒）
	CheckPointGlobalRankSize       int     //记忆点全局排行榜列表长度
	GlobalRankReLoadSecs           int64   // 全局排行榜重载时间
	GlobalRankSize                 int     // 全局排行榜列表长度
	AfterGameAwardFactor           uint64  //游戏后奖励系数（万分之几）
	Panic                          bool    //程序panic开关
	QuotasNeedFinish               []int32 //需要记忆点完成的任务指标

	SDKLoginUrl string // sdk登录验证url
	Appkey      string // 验证秘钥
}

func (config *ApiConfig) String() (str string) {
	str = fmt.Sprintf(`
	Name           = %s
	--
	DefaultStamina              = %d
	DefaultMaxRegStamina        = %d
	DefaultStaminRegIntervalSec = %d
	--
	DefaultDiamond = %d
	DefaultChip = %d
	DefaultCoin = %d
	--
	DefaultGiftValidTimeSec   = %d
	DefaultGiftGiveCooldown   = %d
	DefaultGiftAskCooldown    = %d
	DefaultGiftNotifyCount    = %d
	DefaultGiftAcceptMaxCount = %d
	GiftMaxBatchSize          = %d
	--
	DefaultSysAccount = %s
	--
	EnableSecurity = %t
	--
	MailboxLimit = %d
	--
	MaxFriendsRequestCount = %d
	--
	MinClientVersion = %d
	MaxClientVersion = %d
	CurClientVersion = %d
	--
	LogLevel         = %d
	--
	IsProduction = %t
	--
	MaxRequestSize = %d bytes
	--
	ChannelCount         = %d
	ChannelMaxMsg        = %d
	--
	BundleId = %s
	--
	LottoSlotCount = %d
	LottoInitUserValue = %d
	LottoCostPerTime = %d
	LottoDeduct   = %d
	SysLottoFreeCount = %d
	SysLottoRefreshTime  = %d
	AfterGameLottoDeleteSlotLimit = %d
	--
	MissionCountLimit = %d
	DailyMissionCountLimit = %d
	DailyMissionRefreshHour = %d
	--
	CheckPointIdNum   = %d
	CheckPointGlobalRankReLoadSecs = %d
	CheckPointGlobalRankSize = %d
	--
	AfterGameAwardFactor = %d/10000
	--
	Panic = %t
	--
	QuotasNeedFinish = %v
	SDKLoginUrl = %v
	`,
		config.Name,
		config.DefaultStamina,
		config.DefaultMaxRegStamina,
		config.DefaultStaminRegIntervalSec,
		config.DefaultDiamond,
		config.DefaultChip,
		config.DefaultCoin,
		config.DefaultGiftValidTimeSec,
		config.DefaultGiftGiveCooldown,
		config.DefaultGiftAskCooldown,
		config.DefaultGiftNotifyCount,
		config.DefaultGiftAcceptMaxCount,
		config.GiftMaxBatchSize,
		config.DefaultSysAccount,
		config.EnableSecurity,
		config.MailboxLimit,
		config.MaxFriendsRequestCount,
		config.MinClientVersion,
		config.MaxClientVersion,
		config.CurClientVersion,
		config.LogLevel,
		config.IsProduction,
		config.MaxRequestSize,
		config.ChannelCount,
		config.ChannelMaxMsg,
		config.BundleId,
		config.LottoSlotCount,
		config.LottoInitUserValue,
		config.LottoCostPerTime,
		config.LottoDeduct,
		config.SysLottoFreeCount,
		config.SysLottoRefreshTime,
		config.AfterGameLottoDeleteSlotLimit,
		config.MissionCountLimit,
		config.DailyMissionCountLimit,
		config.DailyMissionRefreshHour,
		config.CheckPointIdNum,
		config.CheckPointGlobalRankReLoadSecs,
		config.CheckPointGlobalRankSize,
		config.AfterGameAwardFactor,
		config.Panic,
		config.QuotasNeedFinish,
		config.SDKLoginUrl,
	)
	return
}

func (config *ApiConfig) Init(name string) {
	config.Name = name
	config.IsProduction = true
	config.CurClientVersion = DefVersion
	config.MinClientVersion = DefMinClientVersion
	config.MaxClientVersion = DefMaxClientVersion
	config.MaxFriendsRequestCount = 10
	config.DefaultStamina = 5
	config.DefaultMaxRegStamina = 5
	config.DefaultDiamond = 100
	config.DefaultChip = 0
	config.DefaultCoin = 0
	config.DefaultStaminRegIntervalSec = 600 // 10 min
	config.DefaultGiftValidTimeSec = int64(7 * 24 * time.Hour / time.Second)
	config.DefaultGiftAskCooldown = int64(24 * time.Hour / time.Second) // 24小时
	config.DefaultGiftGiveCooldown = int64(2 * time.Hour / time.Second) // 2小时
	config.DefaultGiftNotifyCount = 10
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
	//config.TransactionNodeCount = 1
	config.ChannelCount = 1
	config.ChannelMaxMsg = 100
	config.BundleId = "com.guanghuan.SuperBMan,com.737.batteryrun"
	config.LottoSlotCount = 8 //默认8个抽奖格子
	config.LottoInitUserValue = 100
	config.LottoCostPerTime = 80
	config.LottoDeduct = 20
	config.SysLottoFreeCount = 3
	config.SysLottoRefreshTime = 180
	config.AfterGameLottoDeleteSlotLimit = 3    //默认只能删除3次抽奖格子
	config.MissionCountLimit = 3                //默认限制3个激活任务
	config.DailyMissionCountLimit = 1           //日常任务默认限制1个激活任务
	config.DailyMissionRefreshHour = 8          //日常任务默认每天8点刷新
	config.CheckPointIdNum = 100                //默认100个记忆点
	config.CheckPointGlobalRankSize = 20        //默认全局排行榜前20名
	config.CheckPointGlobalRankReLoadSecs = 600 //默认10分钟加载一次
	config.AfterGameAwardFactor = 100           //默认万分之100
	config.Panic = true                         //默认打开程序的panic开关
}

//配置结构体定义
type ConfigStruct struct {
	Configs                            ApiConfig //配置项
	SetBundleId                        Set       //内购的bundleid列表
	AfterGameLottoDeleteCount2MallItem map[int]uint64
}

func (c *ConfigStruct) Clear() {
	c.SetBundleId = make(Set, 0)
	c.AfterGameLottoDeleteCount2MallItem = make(map[int]uint64, 0)
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

//返回
func (c *ConfigCache) Configs() *ApiConfig {
	return &(c.Master().Configs)
}

func (c *ConfigCache) BuldleIds() Set {
	return c.Master().SetBundleId
}

func (c *ConfigCache) AfterGameLottoDeleteMallItems() map[int]uint64 {
	return c.Master().AfterGameLottoDeleteCount2MallItem
}

func NewConfigCache() (c *ConfigCache) {
	c = &ConfigCache{}
	c.Init()
	return
}

var DefConfigCache = NewConfigCache()
