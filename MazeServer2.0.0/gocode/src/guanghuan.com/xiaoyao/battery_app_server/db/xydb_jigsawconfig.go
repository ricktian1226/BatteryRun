// xydb_jigsaw
package batterydb

import (
	//proto "code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//----------------拼图信息配置相关------------
//将老的配置数据全部设置为无效
func (db *BatteryDB) RemoveAllJigsawConfig() (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_JIGSAW_CONFIG, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"jigsawid": bson.M{"$gte": 0}}
	_, err = tbl.RemoveAll(condition)
	return err
}

//导入新的配置数据
func (db *BatteryDB) UpsertJigsawConfig(config battery.JigsawConfig) (err error) {
	condition := bson.M{"jigsawid": config.GetJigsawid()}
	err = db.UpsertData(xybusiness.DB_TABLE_JIGSAW_CONFIG, condition, config)
	return err
}

//----------------系统道具信息配置相关------------
//将老的系统道具配置数据全部设置为无效
func (db *BatteryDB) RemoveAllRuneConfig() (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_RUNE_CONFIG, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"propid": bson.M{"$gte": 0}}
	_, err = tbl.RemoveAll(condition)
	return err
}

//将老的赛前道具配置数据全部设置为无效
func (db *BatteryDB) RemoveAllBGPropConfig() (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_BEFOREGAME_RANDOM_WEIGHT, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"goodid": bson.M{"$gte": 0}}
	_, err = tbl.RemoveAll(condition)
	return err
}

//upsert 系统道具配置信息
func (db *BatteryDB) UpsertRuneConfig(config *battery.RuneConfig) (err error) {
	condition := bson.M{"propid": config.GetPropid(), "value": config.GetValue()}
	return db.UpsertData(xybusiness.DB_TABLE_RUNE_CONFIG, condition, config)
}

//upsert 赛前道具配置信息
func (db *BatteryDB) UpsertBeforeGameRandomGoodsConfig(config *battery.DBBeforeGameRandomGoodWeight) (err error) {
	//condition := bson.M{"goodid": config.GetGoodId(), "weight": config.GetWeight()}
	condition := bson.M{"goodid": config.GetGoodId()}
	return db.UpsertData(xybusiness.DB_TABLE_BEFOREGAME_RANDOM_WEIGHT, condition, config)
}
