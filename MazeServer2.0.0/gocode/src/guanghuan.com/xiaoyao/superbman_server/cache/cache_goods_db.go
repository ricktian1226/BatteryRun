// cache_goods_db
package xybusinesscache

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//加载所有的商品信息
func (cdb *CacheDB) LoadGoods(dbGoodsList *[]*battery.DBMallItem) (err error) {
    condition := bson.M{"valid": true}
    selector := bson.M{"_id": 0, "valid": 0, "xxx_unrecognized": 0}
    return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_GOODS, condition, selector, 0, dbGoodsList, mgo.Strong)
}

//查询某一类型的商品信息
func (cdb *CacheDB) LoadSpecificTypeGoods(mallType battery.MallType, dbGoodsList *[]*battery.DBMallItem) (err error) {
    condition := bson.M{"malltype": mallType, "valid": true}
    selector := bson.M{"_id": 0, "valid": 0, "xxx_unrecognized": 0}
    return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_GOODS, condition, selector, 0, dbGoodsList, mgo.Strong)
}

func (cdb *CacheDB) LoadCheckPointUnlockGoodsConfig(goodsList *[]*battery.DBCheckPointUnlockGoodsConfig) (err error) {
    condition := bson.M{}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
    return cdb.dbCommon.GetAllData(xybusiness.DB_TABLE_CHECKPOINTUNLOCK_GOODS_CONFIG, condition, selector, 0, goodsList, mgo.Strong)
}
