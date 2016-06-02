package batterydb

//import (
//"code.google.com/p/goprotobuf/proto"
//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//xyerror "guanghuan.com/xiaoyao/superbman_server/error"
//)

//func (db *BatteryDB) GetGoodByIapID(iapid string) (good battery.DBMallItem, err error) {

//	c := db.OpenTable(DB_TABLE_GOODS, mgo.Monotonic)
//	defer c.Close()

//	queryStr := bson.M{"iapid": iapid}
//	query := c.Find(queryStr)
//	err = query.One(&good)

//	err = xyerror.DBError(err)
//	return good, err
//}

//func (db *BatteryDB) IsGoodExistByIapID(iapid string) (bool, error) {
//	return db.IsRecordExistingWithField(DB_TABLE_GOODS, "iapid", iapid, mgo.Monotonic)
//}
