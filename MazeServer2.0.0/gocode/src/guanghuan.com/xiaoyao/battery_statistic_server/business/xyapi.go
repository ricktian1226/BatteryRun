package batteryapi

import (
	batterydb "guanghuan.com/xiaoyao/battery_statistic_server/db"
	xyconf "guanghuan.com/xiaoyao/common/conf"
	"guanghuan.com/xiaoyao/common/db"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

type XYAPI struct {
	Config *ConfigCache
}

var apiConfigUtil xyconf.ApiConfigUtil

func NewXYAPI() *XYAPI {
	return &XYAPI{
		Config: DefConfigCache,
	}
}

//获取数据库会话指针
// index int 数据表索引
func (api *XYAPI) GetDB(index int) *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return dbInterface.(*batterydb.BatteryDB)
}

func (api *XYAPI) GetXYDB(index int) *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}
