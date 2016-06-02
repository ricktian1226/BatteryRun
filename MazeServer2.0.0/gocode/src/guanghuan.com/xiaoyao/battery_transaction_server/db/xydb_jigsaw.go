// xydb_jigsaw
package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

var G_MAXSubJigsawCount uint64 = 30

//玩家是否拥有拼图
// uid string 玩家id
// jigsawId uint64 父拼图id
//return:
// isExisting bool true 玩家拥有拼图，false 玩家未拥有拼图
func (db *BatteryDB) IsJigsawExisting(uid string, jigsawId uint64) (isExisting bool, err error) {
	condition := bson.M{"uid": uid, "id": jigsawId}
	isExisting, err = db.IsRecordExisting(xybusiness.DB_TABLE_JIGSAW, condition, mgo.Strong)
	return
}

//获取子拼图数目
// uid string 玩家id
// parentJigsawId uint64 父拼图id
//return:
// count int 子拼图数目
//func (db *BatteryDB) GetSubJigsawCount(uid string, parentJigsawId uint64) (count int) {
//	condition := bson.M{"uid": uid, "id": bson.M{"$gt": parentJigsawId, "$lt": parentJigsawId + G_MAXSubJigsawCount}}
//	count, err := db.GetRecordCount(DB_TABLE_JIGSAW, condition, mgo.Strong)
//	if err != nil {
//		count = 0
//	}
//	return
//}

//获取玩家的拼图列表
// uid string 玩家id
// jigsawIdList *[]uint64 拼图id列表
func (db *BatteryDB) GetJigsawList(uid string, jigsawIdList *[]uint64) (err error) {
	var jigsawList []*battery.Jigsaw
	condition := bson.M{"uid": uid}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetAllData(xybusiness.DB_TABLE_JIGSAW, condition, selector, 0, &jigsawList, mgo.Strong)
	if err == xyerror.ErrOK {
		for _, jigsaw := range jigsawList {
			*jigsawIdList = append(*jigsawIdList, jigsaw.GetId())
		}
	}
	return
}

//增加玩家拼图
// uid string 玩家id
// id uint64 拼图id
func (db *BatteryDB) AddJigsaw(uid string, id uint64) (err error) {
	condition := bson.M{"uid": uid, "id": id}
	jigsaw := &battery.Jigsaw{
		Uid: &uid,
		Id:  &id,
	}
	return db.UpsertData(xybusiness.DB_TABLE_JIGSAW, condition, jigsaw)
}
