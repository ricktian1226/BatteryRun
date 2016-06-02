package batteryapi

import (
	//proto "code.google.com/p/goprotobuf/proto"
	//"encoding/binary"
	//apns "github.com/timehop/apns"
	batterydb "guanghuan.com/xiaoyao/battery_apns_server/db"
	//"guanghuan.com/xiaoyao/common/apn"
	xyconf "guanghuan.com/xiaoyao/common/conf"
	"guanghuan.com/xiaoyao/common/db"
	//"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//"sync"
)

//apn客户端
//var apnClient *apns.Client

//func SetApnClient(c *apns.Client) {
//	apnClient = c
//}

type XYAPI struct {
	Config   *ConfigCache
	platform battery.PLATFORM_TYPE //请求的平台类型
}

var apiConfigUtil xyconf.ApiConfigUtil

func NewXYAPI() *XYAPI {
	return &XYAPI{
		Config:   DefConfigCache,
		platform: battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN,
	}
}

//根据终端平台类型设置业务数据库指针，在业务请求时调用
// platform battery.PLATFORM_TYPE 终端平台类型
func (api *XYAPI) SetDB(platform battery.PLATFORM_TYPE) {
	api.platform = platform
}

//获取业务数据库会话
// index int 业务索引
func (api *XYAPI) GetDB(index int) *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, api.platform)
	return dbInterface.(*batterydb.BatteryDB)
}
func (api *XYAPI) GetXYDB(index int) *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, api.platform)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}

func (api *XYAPI) GetCommonDB(index int) *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return dbInterface.(*batterydb.BatteryDB)
}

func (api *XYAPI) GetCommonXYDB(index int) *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}

func (api *XYAPI) GetLogDB() *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return dbInterface.(*batterydb.BatteryDB)
}

func (api *XYAPI) GetLogXYDB() *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}
