package batterydb

//import (
//	"code.google.com/p/goprotobuf/proto"
//	xylog "guanghuan.com/xiaoyao/common/log"
//	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
//	"labix.org/v2/mgo/bson"
//)

//import "time"

//const (
//	State_NotVerify   int32 = 0 // 未验证
//	State_SuccVecrify int32 = 1 // 成功
//	State_FailVecrify int32 = 2 // 失败
//)

//const (
//	Token_in  int32 = 0 // 产出
//	Token_out int32 = 1 // 消耗
//)

//func (db *BatteryDB) InsertNewOrder(order_id string, itemId string, uid string) error {

//	var iapOrder battery.IapOrder
//	iapOrder.Uid = proto.String(uid)
//	iapOrder.OrderId = proto.String(order_id)
//	iapOrder.ItemId = proto.String(itemId)
//	iapOrder.OrderRequestTimeSecond = proto.Int64(time.Now().Unix())
//	iapOrder.OrderFinishTimeSecond = proto.Int64(0)
//	iapOrder.ServerState = proto.Int32(State_NotVerify)
//	iapOrder.IapState = proto.Int32(State_NotVerify)

//	err := db.AddData(DB_TABLE_IAPORDER, iapOrder)
//	err = xyerror.DBError(err)
//	return err
//}

//func (db *BatteryDB) GetOrderByOrderId(order_id string) (battery.IapOrder, error) {

//	c := db.OpenTable(DB_TABLE_IAPORDER)
//	defer c.Close()

//	queryStr := bson.M{"orderid": order_id}
//	var query = c.Find(queryStr)
//	var orderData battery.IapOrder

//	err := query.One(&orderData)

//	err = xyerror.DBError(err)
//	return orderData, err
//}

//func (db *BatteryDB) UpdateOrder(iapOrder battery.IapOrder) error {
//	c := db.OpenTable(DB_TABLE_IAPORDER)
//	defer c.Close()
//	queryStr := bson.M{"orderid": iapOrder.GetOrderId()}
//	err := c.Update(queryStr, iapOrder)
//	if err != nil {
//		xylog.Error("[DB] update iapOrder failed: %v", err)
//	}

//	err = xyerror.DBError(err)
//	return err
//}
