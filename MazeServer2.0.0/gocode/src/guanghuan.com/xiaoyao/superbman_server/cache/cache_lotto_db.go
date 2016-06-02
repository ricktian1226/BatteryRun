// cache_lotto_db
package xybusinesscache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (ldb *CacheDB) LoadSlotItems(slotitems *[]battery.LottoSlotItem) (err error) {
	c := ldb.dbCommon.OpenTable(xybusiness.DB_TABLE_LOTTO_SLOTITEMS, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"valid": true}
	query := c.Find(queryStr)
	err = query.All(slotitems)
	err = xyerror.DBError(err)

	return
}

func (ldb *CacheDB) LoadWeights(weights *[]battery.LottoWeight) (err error) {
	c := ldb.dbCommon.OpenTable(xybusiness.DB_TABLE_LOTTO_WEIGHT, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"valid": true}
	query := c.Find(queryStr)
	err = query.All(weights)
	err = xyerror.DBError(err)

	return
}

func (ldb *CacheDB) LoadSerialNumSlots(weights *[]battery.LottoSerialNumSlot) (err error) {
	c := ldb.dbCommon.OpenTable(xybusiness.DB_TABLE_LOTTO_SERIALNUM_SLOT, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"valid": true}
	query := c.Find(queryStr)
	err = query.All(weights)
	err = xyerror.DBError(err)

	return
}

//func (ldb *CacheDB) LoadAfterGameWeights(weights *[]battery.LottoAfterGameWeight) (err error) {
//	c := ldb.db.OpenTable(xybusiness.DB_TABLE_LOTTO_AFTERGAME_WEIGHT, mgo.Strong)
//	defer c.Close()

//	//mne := bson.M{"$ne": false}
//	queryStr := bson.M{"valid": true}
//	query := c.Find(queryStr)
//	err = query.All(weights)
//	err = xyerror.DBError(err)

//	return
//}

//func (ldb *CacheDB) LoadStages(stages *[]battery.LottoStageItem) (err error) {
//	c := ldb.db.OpenTable(xybusiness.DB_TABLE_LOTTO_STAGE, mgo.Strong)
//	defer c.Close()

//	//mne := bson.M{"$ne": false}
//	queryStr := bson.M{"valid": true}
//	query := c.Find(queryStr)
//	err = query.All(stages)
//	err = xyerror.DBError(err)

//	return
//}
