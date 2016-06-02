package batterydb

import (
	xydbservice "guanghuan.com/xiaoyao/common/service/db"
)

func NewBatteryDBService(serviceName string, dburl string, dbname string) *xydbservice.DBService {
	db := NewBatteryDB(dburl, dbname)
	dbsvc := xydbservice.NewDBService(serviceName, dburl, dbname, db)
	return dbsvc
}
