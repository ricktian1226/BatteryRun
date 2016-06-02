package batterydb

import (
    "gopkg.in/mgo.v2/bson"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//func (db *BatteryDB) GetGoodsByMallTypeAndCategory(mallType battery.MallType, category battery.PropType /*showInvalid bool, */, maxCount int32) (goodslist []*battery.MallItem, err error) {
//	var condition interface{} = nil

//	if category > 0 {
//		condition = bson.M{"malltype": mallType, "category": category, "valid": true}
//	} else {
//		condition = bson.M{"malltype": mallType, "valid": true}
//	}

//	err = db.GetAllData(DB_TABLE_GOODS, condition, 0, &goodslist)
//	return
//}

//增加商品配置信息
func (db *BatteryDB) AddMallItem(mallitem *battery.DBMallItem) error {
    return db.AddData(xybusiness.DB_TABLE_GOODS, mallitem)
}

//增加商品配置信息
func (db *BatteryDB) UpsertMallItem(mallitem *battery.DBMallItem) (err error) {
    condition := bson.M{"id": mallitem.GetId()}
    return db.UpsertData(xybusiness.DB_TABLE_GOODS, condition, mallitem)
}

//删除所有的商品信息
func (db *BatteryDB) RemoveAllMallItems() error {
    return db.RemoveAllData(xybusiness.DB_TABLE_GOODS, bson.M{})
}

func (db *BatteryDB) AddCheckPointUnlockGoodsConfig(goods *battery.DBCheckPointUnlockGoodsConfig) error {
    return db.AddData(xybusiness.DB_TABLE_CHECKPOINTUNLOCK_GOODS_CONFIG, goods)
}

func (db *BatteryDB) RemoveAllUnlockGoodsConfig() error {
    return db.RemoveAllData(xybusiness.DB_TABLE_CHECKPOINTUNLOCK_GOODS_CONFIG, bson.M{})
}
