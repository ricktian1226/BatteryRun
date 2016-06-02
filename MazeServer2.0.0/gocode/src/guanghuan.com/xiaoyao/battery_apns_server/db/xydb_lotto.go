// 抽奖相关的推送业务
package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"guanghuan.com/xiaoyao/common/apn"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// QueryUidsSysLottoFreeCountNoNull 查询freecount大于0的uid列表
func (db *BatteryDB) QueryUidsSysLottoFreeCountNoNull() (lottoInfos []*battery.SysLottoInfo, err error) {
	lottoInfos = make([]*battery.SysLottoInfo, 0)
	selector := bson.M{"uid": 1}
	condition := bson.M{"freecount": bson.M{"$gt": 0}}
	err = db.GetAllData(xybusiness.DB_TABLE_SYS_LOTTO_INFO, condition, selector, 0, &lottoInfos, mgo.Monotonic)
	return
}

// QueryDeviceTokenByUids 根据uid列表查询devicetoken
func (db *BatteryDB) QueryDeviceTokenByUids(uids []string) (accounts []*battery.DBAccount, err error) {
	accounts = make([]*battery.DBAccount, 0)
	selector := bson.M{"deviceid": 1}
	condition := bson.M{"uid": bson.M{"$in": uids}}
	err = db.GetAllData(xybusiness.DB_TABLE_ACCOUNT, condition, selector, 0, &accounts, mgo.Monotonic)
	return
}
