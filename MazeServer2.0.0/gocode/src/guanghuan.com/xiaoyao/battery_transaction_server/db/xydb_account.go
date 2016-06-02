package batterydb

import (
    //"fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

////根据uid获取玩家账户信息
//// uid string 玩家id
//// account *battery.DBAccount 返回的账户信息
//func (db *BatteryDB) GetAccountDirect(uid string, account *battery.DBAccount, consistency mgo.Mode) (err error) {
//	var condition = bson.M{"uid": uid}
//	err = db.GetOneData(xybusiness.DB_TABLE_ACCOUNT, condition, account, consistency)
//	if err != nil {
//		xylog.Error("[DB] GetAccountDirect failed : %v ", err)
//	}

//	return
//}
func (db *BatteryDB) GetUsersAccoplishment(uids []string, accomplishments *[]*battery.DBUserAccomplishment, consistency mgo.Mode) (err error) {
    condition := bson.M{"uid": bson.M{"$in": uids}}
    selector := bson.M{"uid": 1, "accomplishment": 1}
    err = db.GetAllData(xybusiness.DB_TABLE_USER_ACCOMPLISHMENT, condition, selector, 0, accomplishments, consistency)
    if err != xyerror.ErrOK {
        xylog.Error(xylog.DefaultLogId, "[DB] GetUsersAccoplishment failed : %v ", err)
    }

    return
}

//增加玩家账户信息
// account *battery.DBAccount 玩家账户信息
func (db *BatteryDB) AddAccount(account *battery.DBAccount) error {

    uid := account.GetUid()
    xylog.Debug(uid, "BatteryDB Add account : %s", account.String())
    var err error
    // 只检查Uid是否为空
    if uid != "" {
        err = db.AddData(xybusiness.DB_TABLE_ACCOUNT, *account)
        if err != nil {
            xylog.Error(uid, "[DB] fail to add new account: %v ", err)
        }
    } else {
        xylog.ErrorNoId("[DB] fail to add new account: uid is null ")
        err = xyerror.ErrBadInputData
    }

    err = xyerror.DBError(err)
    return err
}

//刷新玩家账户信息
// account *battery.DBAccount 玩家账户信息
func (db *BatteryDB) UpdateAccount(account *battery.DBAccount) error {
    uid := account.GetUid()
    if uid == "" {
        xylog.ErrorNoId("[DB] fail to UpdateAccount : uid is null ")
        return xyerror.ErrBadInputData
    }

    err := db.UpdateData(xybusiness.DB_TABLE_ACCOUNT, bson.M{"uid": uid}, account)
    xylog.Debug(uid, "updating user(%s):\n%s", uid, account.String())
    if err != nil {
        xylog.ErrorNoId("[DB] update account failed: %v", err)
    }

    err = xyerror.DBError(err)
    return err
}

//增加玩家登录日志
// l battery.AccountLog 玩家登录日志结构体
func (db *BatteryDB) AddAccountLog(l battery.AccountLog) (err error) {
    err = db.AddData(xybusiness.DB_TABLE_ACCOUNT_LOG, l)
    return
}

//根据uid获取玩家成就信息
// uid string 玩家id
// userAccomplishment *battery.DBUserAccomplishment 玩家成就信息
// consistency mgo.Mode
func (db *BatteryDB) GetUserAccomplishment(uid string, userAccomplishment *battery.DBUserAccomplishment, consistency mgo.Mode) (err error) {
    condition := bson.M{"uid": uid}
    selector := bson.M{"accomplishment": 1}
    err = db.GetOneData(xybusiness.DB_TABLE_USER_ACCOMPLISHMENT, condition, selector, userAccomplishment, consistency)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "[DB] GetUserAccomplishment failed : %v ", err)
    }

    return
}

//增加玩家成就信息
// userAccomplishment *battery.DBUserAccomplishment 玩家成就信息
func (db *BatteryDB) AddUserAccomplishment(userAccomplishment *battery.DBUserAccomplishment) error {

    //xylog.Debug("BatteryDB Add userAccomplishment : %s", account.String())

    uid := userAccomplishment.GetUid()
    var err error
    // 只检查Uid是否为空
    if uid != "" {
        err = db.AddData(xybusiness.DB_TABLE_USER_ACCOMPLISHMENT, *userAccomplishment)
        if err != xyerror.ErrOK {
            xylog.Error(uid, "[DB] fail to add new userAccomplishment: %v ", err)
        }
    } else {
        xylog.Error(uid, "[DB] fail to add new userAccomplishment: uid is null ")
        err = xyerror.ErrBadInputData
    }

    //err = xyerror.DBError(err)
    return err
}

//刷新玩家成就信息
// userAccomplishment *battery.DBUserAccomplishment 玩家成就信息
func (db *BatteryDB) UpdateUserAccomplishment(userAccomplishment *battery.DBUserAccomplishment) error {
    uid := userAccomplishment.GetUid()
    if uid == "" {
        xylog.Error(uid, "[DB] fail to UpdateUserAccomplishment : uid is null ")
        return xyerror.ErrBadInputData
    }

    err := db.UpdateData(xybusiness.DB_TABLE_USER_ACCOMPLISHMENT, bson.M{"uid": uid}, userAccomplishment)
    //xylog.Debug("updating user(%s):\n%s", uid, account.String())
    if err != xyerror.ErrOK {
        xylog.Error(uid, "[DB] update UserAccomplishment failed: %v", err)
    }

    err = xyerror.DBError(err)
    return err
}

////刷新玩家的加成信息
//func (db *BatteryDB) UpdateGoldAddtional(uid string, addtional int32) error {
//	condition := bson.M{"uid": uid}
//	setter := bson.M{"$set": bson.M{"goldaddtional": addtional}}
//	return db.UpdateData(xybusiness.DB_TABLE_ACCOUNT, condition, setter)
//}
