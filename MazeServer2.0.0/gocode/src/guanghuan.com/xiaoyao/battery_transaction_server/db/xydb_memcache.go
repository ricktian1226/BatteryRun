package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//to delete
func (db *BatteryDB) GetMemCacheValueOld(uid string, key battery.MemCacheEnum, consistency mgo.Mode) (value string, err error) {
	condition := bson.M{"uid": uid, "key": key}
	selector := bson.M{"value": 1}
	cache := &battery.DBMemCache{}
	err = db.GetOneData(xybusiness.DB_TABLE_MEMCACHE, condition, selector, cache, consistency)
	if err != xyerror.ErrOK {
		return
	}

	value = cache.GetValue()

	return
}

func (db *BatteryDB) RemoveMemCacheValueOld(uid string, key battery.MemCacheEnum) (err error) {
	condition := bson.M{"uid": uid, "key": key}
	err = db.RemoveData(xybusiness.DB_TABLE_MEMCACHE, condition)
	if err != xyerror.ErrOK {
		//xylog.Error(uid, "[DB] RemoveMemCacheValueOld(%v) failed : %v ", condition, err)
		return
	}
	return
}

func (db *BatteryDB) GetMemCacheValue(uid string, key battery.MemCacheEnum, platform battery.PLATFORM_TYPE, consistency mgo.Mode) (value string, err error) {
	condition := bson.M{"uid": uid, "key": key, "platform": platform}
	selector := bson.M{"value": 1}
	cache := &battery.DBMemCache{}
	err = db.GetOneData(xybusiness.DB_TABLE_MEMCACHE, condition, selector, cache, consistency)
	if err != xyerror.ErrOK {
		return
	}

	value = cache.GetValue()

	return
}

func (db *BatteryDB) GetMemCacheValues(uid string, keys []battery.MemCacheEnum, platform battery.PLATFORM_TYPE, consistency mgo.Mode) (caches []*battery.DBMemCache, err error) {
	condition := bson.M{"uid": uid, "platform": platform, "key": bson.M{"$in": keys}}
	selector := bson.M{"key": 1, "value": 1}
	err = db.GetAllData(xybusiness.DB_TABLE_MEMCACHE, condition, selector, 0, &caches, consistency)
	if err != xyerror.ErrOK {
		return
	}

	return
}

func (db *BatteryDB) SetMemCacheValue(uid, value string, key battery.MemCacheEnum, platform battery.PLATFORM_TYPE, consistency mgo.Mode) (err error) {
	condition := bson.M{"uid": uid, "key": key, "platform": platform}
	cache := &battery.DBMemCache{
		Uid:      &uid,
		Key:      key.Enum(),
		Value:    &value,
		Platform: platform.Enum(),
	}
	err = db.UpsertData(xybusiness.DB_TABLE_MEMCACHE, condition, cache)
	if err != xyerror.ErrOK {
		return
	}

	return
}
