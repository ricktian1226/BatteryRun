package batteryapi

import (
	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/panic"
	xyservice "guanghuan.com/xiaoyao/common/service"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

type BatteryService struct {
	Defsvc      *xyservice.DefaultService
	configName  string
	ConfigCache *ConfigCache
}

func NewBatteryService(name string, config string) (svc *BatteryService) {
	svc = &BatteryService{
		Defsvc:     xyservice.NewDefaultService(name),
		configName: config,
	}

	return
}

func (svc *BatteryService) Init() (err error) {
	return
}

func (svc *BatteryService) Start() (err error) {
	return svc.Defsvc.Start()
}
func (svc *BatteryService) Stop() (err error) {
	return svc.Defsvc.Stop()
}
func (svc *BatteryService) Name() string {
	return svc.Defsvc.Name()
}
func (svc *BatteryService) IsRunning() bool {
	return svc.Defsvc.IsRunning()
}

func (svc *BatteryService) LoadConfig() bool {

	DefConfigCache.Slave().Clear()

	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_APICONFIG,
		battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	if nil == dbInterface {
		xylog.ErrorNoId("Get dbInterface for %v, %v failed. damn~~~",
			xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_APICONFIG,
			battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)

		return false
	}

	if apiConfigUtil.Load(dbInterface, svc.configName, &(DefConfigCache.Slave().Configs)) {
		xypanic.Panic_Switch = DefConfigCache.Slave().Configs.Panic //设置一下panic开关
		DefConfigCache.Switch()
		//xylog.DebugNoId("ApiConfig : %s", DefConfigCache.Configs().String())
	} else {
		xylog.ErrorNoId("LoadConfig failed.")
		return false
	}

	return true
}

//配置项刷新后置处理
func (svc *BatteryService) ApplyConfig() {

}

func (svc *BatteryService) PrintConfig() {
	xylog.InfoNoId("Api Settings:%v", DefConfigCache.Configs())
}
