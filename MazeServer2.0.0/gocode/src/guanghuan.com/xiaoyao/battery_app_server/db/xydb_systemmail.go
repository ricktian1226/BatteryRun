package batterydb

import (
	//proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//"time"
)

//----------------系统邮箱信息配置相关------------

//获取所有配置信息
/*
func (db *BatteryDB) GetAllMailInfoConfig() (configList []*battery.SystemMailInfoConfig, err error) {
	condition := bson.M{"mailid": bson.M{"$gte": 0}}
//	err = db.GetAllData(DB_TABLE_MAILINFO_CONFIG, condition, 0, &configList)
	xylog.Debug("[app len %d]", len(configList))
	return configList, err
}*/

//删除所有老数据
func (db *BatteryDB) RemoveAllMaillnfoConfig() (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_MAILINFO_CONFIG, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"mailid": bson.M{"$gte": 0}}
	_, err = tbl.RemoveAll(condition)
	return err
}

//导入新的配置数据
func (db *BatteryDB) UpsertMailInfoConfig(config battery.SystemMailInfoConfig) (err error) {
	condition := bson.M{"mailid": config.GetMailID()}
	err = db.UpsertData(xybusiness.DB_TABLE_MAILINFO_CONFIG, condition, config)
	return err
}
