package db

import (
	xydbservice "guanghuan.com/xiaoyao/common/service/db"
)

func NewBatteryDBService(svc_name string, dburl string, dbname string) *xydbservice.DBService {
	db := NewBatteryDB(dburl, dbname)
	dbsvc := xydbservice.NewDBService(svc_name, dburl, dbname, db)
	return dbsvc
}
