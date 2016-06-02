// batterydb
package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (db *BatteryDB) IsPropExist(id uint64) (isExist bool) {
	condition := bson.M{"id": id}
	isExist, _ = db.IsRecordExisting(xybusiness.DB_TABLE_PROP, condition, mgo.Monotonic)
	return
}
