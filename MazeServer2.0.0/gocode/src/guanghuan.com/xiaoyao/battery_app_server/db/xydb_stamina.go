package batterydb

//import (
//	//	"fmt"
//	xylog "guanghuan.com/xiaoyao/common/log"
//	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
//	"labix.org/v2/mgo/bson"
//)

//// 查询
//func (db *BatteryDB) GetStamina(uid string) (stamina battery.Stamina, err error) {
//	c := db.OpenTable(DB_TABLE_STAMINA)
//	defer c.Close()

//	// 只检查Uid是否为空
//	if uid != "" {
//		err = c.Find(bson.M{"uid": uid}).One(&stamina)
//		if err != nil {
//			xylog.Error("[DB] fail to add new stamina record: %v ", err)
//		}
//	} else {
//		err = xyerror.DBErrFailedDueToClientError
//	}

//	return stamina, xyerror.DBError(err)
//}

//// 更新
//func (db *BatteryDB) UpdateStamina(stamina battery.Stamina) (err error) {
//	uid := stamina.GetUid()
//	if uid == "" {
//		return xyerror.DBErrFailedDueToClientError
//	}
//	c := db.OpenTable(DB_TABLE_STAMINA)
//	defer c.Close()
//	err = c.Update(bson.M{"uid": uid}, stamina)

//	if err != nil {
//		xylog.Error("[DB] update stamina record failed: %v", err)
//	}
//	return xyerror.DBError(err)
//}

//// 添加
//func (db *BatteryDB) AddStamina(stamina battery.Stamina) (err error) {
//	c := db.OpenTable(DB_TABLE_STAMINA)
//	defer c.Close()
//	uid := stamina.GetUid()
//	// 只检查Uid是否为空
//	if uid != "" {
//		err = c.Insert(stamina)
//		if err != nil {
//			xylog.Error("[DB] fail to add new stamina record: %v ", err)
//		}
//	} else {
//		err = xyerror.DBErrFailedDueToClientError
//	}

//	return xyerror.DBError(err)
//}
