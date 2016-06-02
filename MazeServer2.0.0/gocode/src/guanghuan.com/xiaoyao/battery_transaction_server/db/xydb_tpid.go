package batterydb

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//查询玩家是否已经注册
// tpid *battery.TPID 玩家的第三方账户信息 tpid.Source 账户来源，tpid.Id 第三方账户标识
//return:
// isExisting bool 是否存在
func (db *BatteryDB) IsTpidRegistered(tpid *battery.TPID) (isExisting bool, err error) {
    condition := bson.M{"source": tpid.GetSource(), "sid": tpid.GetId()}
    isExisting, err = db.IsRecordExisting(xybusiness.DB_TABLE_TPID_MAP, condition, mgo.Strong)
    return
}

//查询玩家第三方账户信息
// tpid *battery.TPID 玩家的第三方账户信息 tpid.Source 账户来源，tpid.Id 第三方账户标识
//return:
// isExisting bool 是否存在
func (db *BatteryDB) GetIdMapByTpid(tpid *battery.TPID) (idMap *battery.IDMap, err error) {
    idMap = new(battery.IDMap)
    condition := bson.M{"source": tpid.GetSource(), "sid": tpid.GetId()}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0, "createdate": 0}
    err = db.GetOneData(xybusiness.DB_TABLE_TPID_MAP, condition, selector, idMap, mgo.Strong)
    return
}

//查询玩家uid
// tpid *battery.TPID 玩家的第三方账户信息 tpid.Source 账户来源，tpid.Id 第三方账户标识
//return:
// gid string 玩家uid
func (db *BatteryDB) GetGidByTpid(tpid *battery.TPID) (gid string, err error) {
    var idMap *battery.IDMap
    idMap, err = db.GetIdMapByTpid(tpid)
    if err != xyerror.DBErrOK {
        return
    } else {
        gid = idMap.GetGid()
    }
    return
}

//根据玩家uid查询第三方账户信息
// gid string 玩家uid
// source battery.ID_SOURCE 玩家第三方账户来源
//return:
// idmap battery.IDMap 第三方账户结构体
func (db *BatteryDB) GetIdMapByGid(gid string) (idmap battery.IDMap, err error) {
    condition := bson.M{"gid": gid}
    selector := bson.M{"_id": 0, "xxx_unrecognized": 0, "createdate": 0}
    err = db.GetOneData(xybusiness.DB_TABLE_TPID_MAP, condition, selector, &idmap, mgo.Strong)
    return
}

//根据玩家uid查询玩家sid
// gid string 玩家uid
// source battery.ID_SOURCE 玩家第三方账户来源
//return:
// sid string 玩家sid
func (db *BatteryDB) GetSidByGid(gid string, source battery.ID_SOURCE) (sid string, err error) {
    var idmap battery.IDMap
    condition := bson.M{"gid": gid, "source": source}
    selector := bson.M{"sid": 1} //查询结果只返回sid字段就可以了
    err = db.GetOneData(xybusiness.DB_TABLE_TPID_MAP, selector, condition, &idmap, mgo.Strong)
    if err != xyerror.ErrOK {
        return
    } else {
        sid = idmap.GetSid()
    }
    return
}

// GetSidByGid 根据玩家uid查询玩家昵称
// gid string 玩家uid
//returns:
// note string 玩家sid
// err error 操作错误
func (db *BatteryDB) GetNoteByGid(gid string) (note string, err error) {
    var idmap battery.IDMap
    condition := bson.M{"gid": gid}
    selector := bson.M{"note": 1} //查询结果只返回note字段就可以了
    err = db.GetOneData(xybusiness.DB_TABLE_TPID_MAP, condition, selector, &idmap, mgo.Strong)
    if err != xyerror.ErrOK {
        return
    } else {
        note = idmap.GetNote()
    }
    return
}

//根据玩家的sid获取玩家的uid
// sid string 玩家sid
// source battery.ID_SOURCE 玩家第三方账户来源
//return:
// gid string 玩家uid
func (db *BatteryDB) GetGidBySid(sid string, source battery.ID_SOURCE) (gid string, err error) {
    var idmap battery.IDMap
    condition := bson.M{"sid": sid, "source": source}
    selector := bson.M{"gid": 1} //查询结果只返回gid字段就可以了
    err = db.GetOneData(xybusiness.DB_TABLE_TPID_MAP, condition, selector, &idmap, mgo.Strong)
    if err != xyerror.ErrOK {
        return
    } else {
        gid = idmap.GetGid()
    }
    return
}

//增加玩家第三方账户信息
// tpid *battery.TPID 玩家第三方账户信息
//return:
// gid string 玩家uid
//func (db *BatteryDB) AddTPID(tpid *battery.TPID) (gid string, err error) {
func (db *BatteryDB) AddTPID(idMap *battery.IDMap) (gid string, err error) {
    err = db.AddData(xybusiness.DB_TABLE_TPID_MAP, idMap)
    if err != xyerror.ErrOK {
        return
    } else {
        gid = idMap.GetGid()
    }
    return
}

//刷新玩家的名称和图标url信息
// tpid *battery.TPID 玩家第三方账户信息
func (db *BatteryDB) UpsertTPID(idMap *battery.IDMap) (err error) {
    condition := bson.M{"sid": idMap.GetSid(), "source": idMap.GetSource()}
    return db.UpsertData(xybusiness.DB_TABLE_TPID_MAP, condition, idMap)
}

//删除玩家tpid信息
// tpid *battery.TPID 玩家第三方账户信息
func (db *BatteryDB) RemoveTPID(gid string, source battery.ID_SOURCE) (err error) {
    condition := bson.M{"gid": gid, "source": source}
    return db.RemoveData(xybusiness.DB_TABLE_TPID_MAP, condition)
}

func (db *BatteryDB) UpsertTPIDName(idmap *battery.IDMap) (err error) {
    condition := bson.M{"gid": idmap.GetGid()}
    return db.UpsertData(xybusiness.DB_TABLE_TPID_MAP, condition, idmap)
}
