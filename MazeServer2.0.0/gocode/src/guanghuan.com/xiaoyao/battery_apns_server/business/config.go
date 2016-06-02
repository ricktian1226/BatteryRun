package batteryapi

import (
	"fmt"
	"guanghuan.com/xiaoyao/common/cache"
)

type ApiConfig struct {
	Name                   string // 配置名称
	IsApnsProduction       bool   //是否是生产环境
	DefaultGiftNotifyCount int32  //玩家好友邮件的推送阈值
	Panic                  bool   //是否打开程序crash panic开关
	AppKey                 string
	AppMasterSecret        string
	ApnNotifyDevicePerReq  int //apn没请求推送设备数目控制
}

func (config *ApiConfig) String() (str string) {
	str = fmt.Sprintf(`
	Name         = %s
	IsApnsProduction = %t
	DefaultGiftNotifyCount = %d
	Panic             = %t
	AppKey  =   %s
	AppMasterSecret = %s
	ApnNotifyDevicePerReq = %d
	`,
		config.Name,
		config.IsApnsProduction,
		config.DefaultGiftNotifyCount,
		config.Panic,
		config.AppKey,
		config.AppMasterSecret,
		config.ApnNotifyDevicePerReq)
	return
}

func (config *ApiConfig) Init(name string) {
	config.Name = name
	config.IsApnsProduction = true
	config.DefaultGiftNotifyCount = 5
	config.Panic = true
	config.ApnNotifyDevicePerReq = 100
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
