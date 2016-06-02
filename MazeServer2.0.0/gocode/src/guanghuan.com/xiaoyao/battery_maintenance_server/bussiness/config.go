package bussiness

import (
    "fmt"
    "guanghuan.com/xiaoyao/common/cache"
)

type ApiConfig struct {
    Name string // 1 配置名称

    Panic        bool   // 23 是否打开程序crash panic开关
    Appkey       string // 加密串
    AppSecretkey string // 充值加密串
    Appid        string // 游戏后台id
}

func (config *ApiConfig) String() (str string) {
    str = fmt.Sprintf(`
	Name         = %s
	
	Panic           = %t
	Appkey    =%s
	AppSecretkey =%s
	Appid     =%s
	`,
        config.Name,
        config.Panic,
        config.Appkey,
        config.AppSecretkey,
        config.Appid)
    return
}

func (config *ApiConfig) Init(name string) {

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
