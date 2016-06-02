// xydb_consumable
//赛前道具
package batterydb

import (
	proto "code.google.com/p/goprotobuf/proto"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//判断玩家是否拥有赛前道具
// uid string 玩家id
// id uint64 道具id
// amount uint32 数目
//return:
// isExisting bool true 拥有足够数量 false 未拥有指定数量
func (db *BatteryDB) IsConsumableExisting(uid string, id uint64, amount uint32) (isExisting bool) {
	isExisting = true
	var err error
	consumable := &battery.Consumable{}
	condition := bson.M{"uid": uid, "id": id}
	selector := bson.M{"amount": 1} //只要查到amount字段就可以了
	err = db.GetOneData(xybusiness.DB_TABLE_CONSUMABLE, condition, selector, consumable, mgo.Strong)
	if err != xyerror.ErrOK || consumable.GetAmount() < amount { //数目不够
		isExisting = false
		return
	}

	return
}

//查询某一赛前道具信息
// uid string 玩家id
// id uint64 道具id
// consumable *battery.Consumable 赛前道具信息
func (db *BatteryDB) QueryConsumable(uid string, id uint64, consumable *battery.Consumable) (err error) {
	condition := bson.M{"uid": uid, "id": id}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_CONSUMABLE, condition, selector, consumable, mgo.Strong)
	return
}

//查询玩家拥有的随机赛前道具信息
// uid string 玩家id
func (db *BatteryDB) GetRandomConsumableID(uid string) (id uint64) {
	var consumable battery.Consumable
	condition := bson.M{"uid": uid, "random": true} //random true表示是随机道具
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	db.GetOneData(xybusiness.DB_TABLE_CONSUMABLE, condition, selector, &consumable, mgo.Strong)
	if consumable.GetAmount() == 1 {
		return consumable.GetId()
	}
	return 0
}

//查询玩家拥有的普通赛前道具列表
// uid string 玩家id
func (db *BatteryDB) GetNoRandomConsumables(uid string, consumables *[]*battery.Consumable) (err error) {
	condition := bson.M{"uid": uid, "random": false}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetAllData(xybusiness.DB_TABLE_CONSUMABLE, condition, selector, 0, consumables, mgo.Strong)
	return
}

//查询玩家拥有的所有赛前道具（普通+随机）列表
func (db *BatteryDB) GetConsumables(uid string, consumables *[]*battery.Consumable) (err error) {
	condition := bson.M{"uid": uid}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetAllData(xybusiness.DB_TABLE_CONSUMABLE, condition, selector, 0, consumables, mgo.Strong)
	return
}

//增加玩家普通赛前道具
// uid string 玩家id
// id uint64 道具id
// amount uint32 道具数目
func (db *BatteryDB) AddConsumable(uid string, id uint64, amount uint32) (err error) {

	consumable := &battery.Consumable{}

	err = db.QueryConsumable(uid, id, consumable)
	if err == xyerror.ErrOK { //如果找到，增加
		consumable.Amount = proto.Uint32(consumable.GetAmount() + amount)
	} else if err == xyerror.ErrNotFound { //如果没找到，设置
		consumable.Uid = proto.String(uid)
		consumable.Id = proto.Uint64(id)
		consumable.Amount = proto.Uint32(amount)
		consumable.Random = proto.Bool(false)
	}

	//upsert一下
	condition := bson.M{"uid": uid, "id": id}
	err = db.UpsertData(xybusiness.DB_TABLE_CONSUMABLE, condition, consumable)

	return
}

//增加赛前随机道具
// uid string 玩家id
// id uint64 道具id
// amount uint32 道具数目
func (db *BatteryDB) UpsertRandomConsumable(uid string, id uint64, amount uint32) (err error) {
	condition := bson.M{"uid": uid, "random": true}
	consumable := &battery.Consumable{
		Uid:    proto.String(uid),
		Id:     proto.Uint64(id),
		Amount: proto.Uint32(amount),
		Random: proto.Bool(true),
	}

	err = db.UpsertData(xybusiness.DB_TABLE_CONSUMABLE, condition, consumable)

	return
}

//使用赛前道具
// uid string 玩家id
// id uint64 道具id
// amount uint32 道具数目
func (db *BatteryDB) DecreaseConsumable(uid string, id uint64, amount uint32) (err error) {

	consumable := &battery.Consumable{}
	err = db.QueryConsumable(uid, id, consumable)
	if err == xyerror.ErrOK {
		if consumable.GetAmount() >= amount { //找到并且数目足够
			consumable.Amount = proto.Uint32(consumable.GetAmount() - amount)
			condition := bson.M{"uid": uid, "id": id}
			err = db.UpsertData(xybusiness.DB_TABLE_CONSUMABLE, condition, consumable)
		} else { //返回错误
			errStr := fmt.Sprintf("[%s] consumable %d amount %d no enough for %d", uid, id, consumable.GetAmount(), amount)
			err = errors.New(errStr)
		}
	}

	return
}
