package batterydb

import (
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) AddIapLog(l battery.IapLog) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_IAP_LOG, l)
	return
}
