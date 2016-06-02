// xyconf
//后台配置公共类，可以为各个服务调用
//该文件的配置类用于应用内业务配置信息，可以在不停止应用服务的情况在reload

package xyconf

import (
	"gopkg.in/mgo.v2"
	xydb "guanghuan.com/xiaoyao/common/db"
	xylog "guanghuan.com/xiaoyao/common/log"
)

const (
	DB_TABLE_API_CONFIG = "apiconfig" // api 配置参数
)

type ApiConfigInterface interface {
	Init(name string)
}

type ApiConfigUtil struct {
}

func (config *ApiConfigUtil) Load(db xydb.DBInterface, config_name string, aci ApiConfigInterface) (isSuccess bool) {
	isSuccess, _ = config.ReadFromDB(db, config_name, aci)
	if isSuccess {
		xylog.InfoNoId("Reading config(%s) from DB success", config_name)
	}
	return
}

func (config *ApiConfigUtil) ReadFromDB(db xydb.DBInterface, config_name string, aci ApiConfigInterface) (isSuccess bool, err error) {
	var (
		isExisting bool
	)

	isExisting, err = db.IsRecordExistingWithField(DB_TABLE_API_CONFIG, "name", config_name, mgo.Strong)
	if isExisting {
		err = db.GetOneDataWithField(DB_TABLE_API_CONFIG, "name", config_name, nil, aci, mgo.Strong)
		if err == nil {
			isSuccess = true
		}
	}

	xylog.DebugNoId("isSuccess : %v, config : %v", isSuccess, aci)

	return
}

func (config *ApiConfigUtil) WriteToDB(db xydb.DBInterface, aci ApiConfigInterface) (err error) {
	err = db.AddData(DB_TABLE_API_CONFIG, aci)
	return
}
func (config *ApiConfigUtil) Reload(config_name string, aci ApiConfigInterface) (db xydb.DBInterface, isSuccess bool, err error) {
	isSuccess, err = config.ReadFromDB(db, config_name, aci)
	return
}
func (config *ApiConfigUtil) Write(db xydb.DBInterface, aci ApiConfigInterface) (err error) {
	err = config.WriteToDB(db, aci)
	return
}
