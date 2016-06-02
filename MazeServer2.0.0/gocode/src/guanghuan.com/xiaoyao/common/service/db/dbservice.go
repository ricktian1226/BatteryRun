package xydbservice

import (
	xydb "guanghuan.com/xiaoyao/common/db"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyservice "guanghuan.com/xiaoyao/common/service"
)

type DBService struct {
	xyservice.DefaultService
	dburl   string
	dbname  string
	db_inst xydb.DBInterface
}

func NewDBService(svc_name string, dburl string, dbname string, db xydb.DBInterface) *DBService {
	dbsvc := &DBService{
		DefaultService: *xyservice.NewDefaultService(svc_name),
		dburl:          dburl,
		dbname:         dbname,
		db_inst:        db,
	}
	return dbsvc
}

func NewXYDBService(svc_name string, dburl string, dbname string) *DBService {
	dbsvc := NewDBService(svc_name, dburl, dbname, xydb.NewXYDB(dburl, dbname))
	return dbsvc
}

func (svc *DBService) GetDB() xydb.DBInterface {
	return svc.db_inst
}

func (svc *DBService) Init() (err error) {
	xylog.InfoNoId("Connecting to dbserver(%s) on db(%s)", svc.dburl, svc.dbname)
	svc.DefaultService.Init()
	err = svc.db_inst.OpenDB()
	if err == nil {
		svc.DefaultService.Start()
		xylog.InfoNoId("Connected")
	} else {
		xylog.ErrorNoId("Failed to connect to DB (%s), ERROR: %s", svc.dburl, err.Error())
	}
	return
}
func (svc *DBService) Start() (err error) {
	if !svc.IsRunning() {
		return svc.DefaultService.Start()
	}
	return
}

func (svc *DBService) Stop() (err error) {
	if svc.IsRunning() {
		svc.db_inst.CloseDB()
		svc.DefaultService.Stop()
		xylog.InfoNoId("Disconnect to DB server")
	}
	return
}
