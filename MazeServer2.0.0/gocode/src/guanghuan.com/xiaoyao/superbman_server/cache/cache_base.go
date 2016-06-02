// cache_base
package xybusinesscache

import (
	//xydb "guanghuan.com/xiaoyao/common/db"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//"sync/atomic"
)

//初始化函数：
// 设置数据库操作指针
//func Init(db *xydb.XYDB) {
func Init(dbCommon *xybusiness.XYBusinessDB,
	dbIos *xybusiness.XYBusinessDB) {
	DefCacheDB.SetDB(dbCommon, dbIos)
}

////缓存下标管理器
//type CacheBase struct {
//	index int32
//}

////获取主缓存下标
//func (c *CacheBase) Major() int32 {
//	return atomic.LoadInt32(&c.index)
//}

////获取备缓存下标
//func (c *CacheBase) Secondary() int32 {
//	return (atomic.LoadInt32(&c.index) + 1) % 2
//}

////切换缓存
//func (c *CacheBase) Switch() {
//	atomic.StoreInt32(&c.index, (atomic.LoadInt32(&c.index)+1)%2)
//}
