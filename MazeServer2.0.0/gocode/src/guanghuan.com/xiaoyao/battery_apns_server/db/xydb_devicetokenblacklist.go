package batterydb

//import (
//	"gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	"guanghuan.com/xiaoyao/common/apn"
//	"guanghuan.com/xiaoyao/superbman_server/server"
//)

////查询device token黑名单
//func (db *BatteryDB) QueryDeviceTokenBlackList() (dts []*xyapn.BlackDeviceToken, err error) {
//	dts = make([]*xyapn.BlackDeviceToken, 0)
//	selector := bson.M{"token": 1}
//	err = db.GetAllData(xybusiness.DB_TABLE_DEVICETOKEN_BLICKLIST, bson.M{}, selector, 0, &dts, mgo.Strong)
//	return
//}

////把device token从黑名单中删除
//func (db *BatteryDB) RemoveDeviceTokenFromBlackList(dt string) (err error) {
//	return db.RemoveAllData(xybusiness.DB_TABLE_DEVICETOKEN_BLICKLIST, bson.M{"token": dt})
//}

////device token从黑名单中删除
//func (db *BatteryDB) UpsertDeviceTokenToBlackList(dt string, timestamp int64) (err error) {
//	bdt := &xyapn.BlackDeviceToken{
//		Token:     &dt,
//		Timestamp: &timestamp,
//	}
//	return db.UpsertData(xybusiness.DB_TABLE_DEVICETOKEN_BLICKLIST, bson.M{"token": dt}, bdt)
//}
