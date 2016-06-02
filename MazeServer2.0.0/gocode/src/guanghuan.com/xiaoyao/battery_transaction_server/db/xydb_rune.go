// xydb_rune
package batterydb

import (
	//"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) UpsertRune(userRune *battery.Rune) (err error) {
	condition := bson.M{"uid": userRune.GetUid(), "id": userRune.GetId()}
	return db.UpsertData(xybusiness.DB_TABLE_RUNE, condition, &userRune)
}

func (db *BatteryDB) IsRuneExisting(uid string, runeId uint64) (isExisting bool, err error) {
	condition := bson.M{"uid": uid, "id": runeId}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	var runeItem battery.Rune
	err = db.GetOneData(xybusiness.DB_TABLE_RUNE, condition, selector, &runeItem, mgo.Strong)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			isExisting = false
		}
	} else {
		isExisting = true
	}
	return
}

func (db *BatteryDB) GetRuneIdList(uid string, runelist *[]uint64) (err error) {
	var runeItemList []*battery.Rune
	condition := bson.M{"uid": uid}
	selector := bson.M{"id": 1} //查询结果只需要id就可以了
	err = db.GetAllData(xybusiness.DB_TABLE_RUNE, condition, selector, 0, &runeItemList, mgo.Strong)
	if err == xyerror.ErrOK {
		for _, runeItem := range runeItemList {
			*runelist = append(*runelist, runeItem.GetId())
		}
	}
	return
}

// 查询玩家拥有的符文信息
// uid string 玩家标识
// runeList []*battery.RuneUnit 符文信息列表
func (db *BatteryDB) GetRuneList(uid string) (runeList []*battery.Rune, err error) {
	condition := bson.M{"uid": uid}
	selector := bson.M{ /*"id": 1, "expiredlimitation": 1*/ }
	err = db.GetAllData(xybusiness.DB_TABLE_RUNE, condition, selector, 0, &runeList, mgo.Strong)
	return
}

// 查询玩家背包中的某一符文信息
// uid string 玩家标识
// id uint64 符文标识
func (db *BatteryDB) GetRune(uid string, id uint64) (userRune *battery.Rune, err error) {
	userRune = &battery.Rune{}
	condition := bson.M{"uid": uid, "id": id}
	err = db.GetOneData(xybusiness.DB_TABLE_RUNE, condition, bson.M{}, userRune, mgo.Strong)
	return
}
