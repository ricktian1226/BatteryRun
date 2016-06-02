package batteryapi

import (
	batterydb "guanghuan.com/xiaoyao/battery_transaction_server/db"
	xyconf "guanghuan.com/xiaoyao/common/conf"
	"guanghuan.com/xiaoyao/common/db"
	//"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

const (
	LOGTRACE_DEF                    = iota
	LOGTRACE_NEWLOGIN               //1
	LOGTRACE_ADDACCOUNT             //2
	LOGTRACE_ADDSYSLOTTOINFO        //3
	LOGTRACE_ADDROLEINFO            //4
	LOGTRACE_ADDTPID                //5
	LOGTRACE_OLDLOGIN               //6
	LOGTRACE_GETDBACCOUNT           //7
	LOGTRACE_UPDATEDBACCOUNT        //8
	LOGTRACE_QUERYUSERMISSION       //9
	LOGTRACE_CONFIRMUSERMISSION     //10
	LOGTRACE_QUERYUSERMISSIONDETAIL //11
	LOGTRACE_UPDATEUSERMISSION      //12
	LOGTRACE_UPDATEUSERMISSIONSTAT  //13
	LOGTRACE_NEWGAME                //14
	LOGTRACE_GAMERESULT             //15
	LOGTRACE_UPDATEGAME             //16
	LOGTRACE_AFTERGAMEREWARDS       //17
	LOGTRACE_GETCURRENTSTAMINA      //18
	LOGTRACE_UPDATEUSERCHECKPOINT   //19
	LOGTRACE_RUNEEXIST              //20
)

type empty struct{}
type Set map[interface{}]empty

type XYAPI struct {
	Config   *ConfigCache
	platform battery.PLATFORM_TYPE //请求的平台类型
}

var apiConfigUtil xyconf.ApiConfigUtil

//创建业务实例对象
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
	//xybusiness.DefBusinessDBSessionManager.Print()
	//if dbInterface == nil {
	//	xylog.ErrorNoId("DefBusinessDBSessionManager.Get index(%d) platform(%v) is nil", index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	//	return nil
	//} else {
	return dbInterface.(*batterydb.BatteryDB)
	//}
}

func (api *XYAPI) GetCommonXYDB(index int) *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}

func (api *XYAPI) GetCommonXYBusinessDB(index int) *xybusiness.XYBusinessDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYBusinessDB)
}

func (api *XYAPI) GetLogDB() *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return dbInterface.(*batterydb.BatteryDB)
}

func (api *XYAPI) GetLogXYDB() *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}
