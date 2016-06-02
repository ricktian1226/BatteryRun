package batterydb

import (
	//proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//"time"
)

//获取玩家的好友体力赠送日志
// uid string 玩家id
// friendUids []string 好友uid列表
// items *[]*battery.StaminaGiveLogItem 日志信息列表
func (db *BatteryDB) GetStaminaGiveLogItems(uid string, friendUids []string, items *[]*battery.StaminaGiveLogItem) (err error) {
	condition := bson.M{"uid": uid, "frienduid": bson.M{"$in": friendUids}}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetAllData(xybusiness.DB_TABLE_STAMINA_GIVEAPPLY_LOG, condition, selector, 0, items, mgo.Monotonic) //读写分离，从secondary节点上读取
	return
}
