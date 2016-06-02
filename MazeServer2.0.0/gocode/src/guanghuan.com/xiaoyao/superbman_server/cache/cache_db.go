// cache_db
package xybusinesscache

import (
	//"guanghuan.com/xiaoyao/common/db"
	//"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//默认的数据库操作接口
var DefCacheDB CacheDB

type CacheDB struct {
	//db *xydb.XYDB
	dbCommon *xybusiness.XYBusinessDB
	dbIos    *xybusiness.XYBusinessDB
}

//func (ldb *CacheDB) SetDB(db *xydb.XYDB) {
func (ldb *CacheDB) SetDB(dbCommon *xybusiness.XYBusinessDB,
	dbIos *xybusiness.XYBusinessDB) {
	ldb.dbCommon = dbCommon
	ldb.dbIos = dbIos
	//xylog.Debug("db : %v", ldb.db)
}

func (ldb *CacheDB) DB() *xybusiness.XYBusinessDB {
	return ldb.dbCommon
}

func (ldb *CacheDB) DBIos() *xybusiness.XYBusinessDB {
	return ldb.dbIos
}
