// cache_useridentity_db
// 玩家identity管理器
package xybusinesscache

import (
	"fmt"

	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// LoadUserIdentityCounter 加载当前节点对应的useridentitycounter信息
// prefix string 当前节点前缀
//returns:
// userIdentity *battery.UserIdentity 当前节点对应的useridentitycounter信息结构指针
// err error 数据库操作错误
func (cdb *CacheDB) LoadUserIdentityCounter(prefix string, userIdentityCounter *battery.DBUserIdentityCounter) (err error) {
	condition := bson.M{"prefix": prefix}
	selector := bson.M{ /*"counter": 1*/ }
	err = cdb.dbCommon.GetOneData(xybusiness.DB_TABLE_USERIDENTITYCOUNTER, condition, selector, userIdentityCounter, mgo.Strong)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound { //不存在的话，就插入一条
			userIdentityCounter = cdb.defaultUserIdentityCounter(prefix)
			//userIdentityCounter.Counter[index] = int32(1)
			err = cdb.dbCommon.UpsertData(xybusiness.DB_TABLE_USERIDENTITYCOUNTER, condition, userIdentityCounter)
		}
	}

	return
}

// UpsertUserIdentityCounter 更新当前节点的useridentitycounter信息
// prefix string 当前节点前缀
// userIdentity *battery.UserIdentity 当前节点对应的useridentitycounter信息结构指针
//return:
// err error 数据库操作错误
func (cdb *CacheDB) UpsertUserIdentityCounter(prefix string, userIdentity *battery.DBUserIdentityCounter) (err error) {
	condition := bson.M{"prefix": prefix}
	return cdb.dbCommon.UpsertData(xybusiness.DB_TABLE_USERIDENTITYCOUNTER, condition, userIdentity)
}

// IncreaseUserIdentityCounter 增加userIdentityCounter计数
// prefix string 服务节点前缀
// index int 计数器下标
//returns:
// err error 操作错误
func (cdb *CacheDB) IncreaseUserIdentityCounter(prefix string, index int) (err error) {
	condition := bson.M{"prefix": prefix}
	setter := bson.M{"$inc": bson.M{fmt.Sprintf("counter.%d", index): 1}}
	err = cdb.dbCommon.UpdateData(xybusiness.DB_TABLE_USERIDENTITYCOUNTER, condition, setter)
	return
}

//// 有几个平台类型， 参见 proto中PLATFORM_TYPE的定义
//const (
//	PLATFORM_COUNT = 3
//)

// defaultUserIdentityCounter 生成默认的useridentitycounter对象
// prefix string 服务节点前缀
// returns:
// userIdentity *battery.DBUserIdentityCounter 计数器管理器对象
func (cdb *CacheDB) defaultUserIdentityCounter(prefix string) (userIdentity *battery.DBUserIdentityCounter) {
	counter := make([]int32, len(battery.PLATFORM_TYPE_value)) //根据协议中平台类型
	return &battery.DBUserIdentityCounter{
		Prefix:  proto.String(prefix),
		Counter: counter,
	}
}
