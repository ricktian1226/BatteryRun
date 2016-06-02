package batterydb

import (
//"code.google.com/p/goprotobuf/proto"
//"guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

//func (db *BatteryDB) InsertNewGood(goodId string, diamondCount int32, price float32) error {
//	var good battery_run_net.IapGood
//	good.GoodId = proto.String(goodId)
//	good.Count = proto.Int32(diamondCount)
//	good.Price = proto.Float32(price)

//	c := db.OpenTable(DB_TABLE_IAPGOOD)
//	defer c.Close()
//	err := c.Insert(good)

//	err = xyerror.DBError(err)
//	return err
//}

//func (db *BatteryDB) GetGoodById(goodId string) (battery_run_net.IapGood, error) {

//	c := db.OpenTable(DB_TABLE_IAPGOOD)
//	defer c.Close()

//	var good battery_run_net.IapGood
//	queryStr := bson.M{"goodid": goodId}
//	query := c.Find(queryStr)
//	err := query.One(&good)

//	err = xyerror.DBError(err)
//	return good, err
//}

//func (db *BatteryDB) IsGoodExistBuy(goodId string) (bool, error) {
//	return db.IsRecordExistingWithField(DB_TABLE_IAPGOOD, "goodid", goodId)
//}
