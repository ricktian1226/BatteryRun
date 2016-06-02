package batteryapi

import (
	"fmt"
	"guanghuan.com/xiaoyao/common/cache"
)

type ApiConfig struct {
	Name          string // 配置名称
	Panic         bool   //是否打开程序crash panic开关
	AppKey        string
	AppId         string
	AppSecretKey  string //与数据中心通信密钥
	DataCenterUrl string //数据中心链接
}

func (config *ApiConfig) String() (str string) {
	str = fmt.Sprintf(`
	Name         = %s
	Panic             = %t
	AppKey  =   %s
	AppId =   %s
	AppSecretKey = %s
	DataCenterUrl = %s
	`,
		config.Name,
		config.Panic,
		config.AppKey,
		config.AppId,
		config.AppSecretKey,
		config.DataCenterUrl)
	return
}

func (config *ApiConfig) Init(name string) {
	config.Name = name
	config.Panic = true
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
