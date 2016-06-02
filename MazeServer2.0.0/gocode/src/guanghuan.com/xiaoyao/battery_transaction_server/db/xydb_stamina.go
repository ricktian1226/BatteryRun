package batterydb

import (
	"gopkg.in/mgo.v2/bson"

	"guanghuan.com/xiaoyao/common/util"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// UpdateStaminaLastUpdateTimeAndLimitAddtional 刷新玩家的上次体力刷新时间和体力上限加成
// uid string 玩家标识
// addtional int32 加成值
func (db *BatteryDB) UpdateStaminaLastUpdateTimeAndLimitAddtional(uid string, addtional int32) (err error) {

	// 先试一下刷新玩家体力刷新时间为-1的情况
	condition := bson.M{"uid": uid, "staminalastupdatetime": -1} //只刷新之前已经停止刷新体力的记录
	setter := bson.M{"$set": bson.M{"staminalastupdatetime": xyutil.CurTimeSec(), "staminalimitaddtional": addtional}}
	err = db.UpdateData(xybusiness.DB_TABLE_ACCOUNT, condition, setter)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			//如果没找到体力刷新时间为-1的记录，则说明玩家当前未达到体力上限，值刷新体力加成即可
			condition = bson.M{"uid": uid} //只刷新之前已经停止刷新体力的记录
			setter = bson.M{"$set": bson.M{"staminalimitaddtional": addtional}}
			err = db.UpdateData(xybusiness.DB_TABLE_ACCOUNT, condition, setter)
		}
	}

	return
}

// UpdateCoinAddtional 更新玩家的金币加成信息
// uid string 玩家标识
// addtional int32 加成值
func (db *BatteryDB) UpdateCoinAddtional(uid string, addtional int32) error {
	condition := bson.M{"uid": uid}
	setter := bson.M{"$set": bson.M{"coinaddtional": addtional}}
	return db.UpdateData(xybusiness.DB_TABLE_ACCOUNT, condition, setter)
}

// UpdateResolveAddtional 更新玩家的分解加成信息
// uid string 玩家标识
// addtional int32 加成值
func (db *BatteryDB) UpdateResolveAddtional(uid string, addtional int32) error {
	condition := bson.M{"uid": uid}
	setter := bson.M{"$set": bson.M{"resolveaddtional": addtional}}
	return db.UpdateData(xybusiness.DB_TABLE_ACCOUNT, condition, setter)
}
