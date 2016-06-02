package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

func (db *BatteryDB) GetGidBySid(sid string, source battery.ID_SOURCE) (gid string, err error) {
	idmap := &battery.IDMap{}
	selector := bson.M{"gid": 1}

	err = db.QueryTpidBySid(selector, sid, source, idmap, mgo.Strong)
	if err != xyerror.ErrOK {
		xylog.Debug("DB Get IDMap for sid(%s) failed: %v", sid, err)
	} else {
		gid = idmap.GetGid()
	}
	return
}
