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

//func (db *BatteryDB) GetAccountDirect(uid string, account *battery.DBAccount, consistency mgo.Mode) (err error) {
//	var condition = bson.M{"uid": uid}
//	err = db.GetOneData(xybusiness.DB_TABLE_ACCOUNT, condition, account, consistency)
//	if err != xyerror.ErrOK {
//		xylog.Error("[DB] GetAccountDirect failed : %v ", err)
//	}

//	return
//}

//func (db *BatteryDB) GetAccountsDirect(uids []string, accounts *[]*battery.DBAccount, consistency mgo.Mode) (err error) {
//	var condition = bson.M{"uid": bson.M{"$in": uids}}
//	err = db.GetAllData(xybusiness.DB_TABLE_ACCOUNT, condition, 0, accounts, consistency)
//	if err != xyerror.ErrOK {
//		xylog.Error("[DB] GetAccountsDirect failed : %v ", err)
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

func (db *BatteryDB) AddAccount(account *battery.DBAccount) error {

	xylog.Debug("BatteryDB Add account : %s", account.String())

	uid := account.GetUid()
	var err error
	// 只检查Uid是否为空
	if uid != "" {
		err = db.AddData(xybusiness.DB_TABLE_ACCOUNT, *account)
		if err != nil {
			xylog.Error(xylog.DefaultLogId, "[DB] fail to add new account: %v ", err)
		}
	} else {
		xylog.Error(xylog.DefaultLogId, "[DB] fail to add new account: uid is null ")
		err = xyerror.ErrBadInputData
	}

	err = xyerror.DBError(err)
	return err
}

func (db *BatteryDB) UpdateAccountDeviceId(uid string, device_id string) (err error) {
	condition := bson.M{"uid": uid}

	err = db.UpdateOneField(xybusiness.DB_TABLE_ACCOUNT, condition, "deviceid", device_id, false)
	return
}

func (db *BatteryDB) AddAccountLog(l battery.AccountLog) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_ACCOUNT_LOG, l)
	return
}

func (db *BatteryDB) IsUserExist(uid string) (isExist bool) {
	condition := bson.M{"uid": uid}
	isExist, _ = db.IsRecordExisting(xybusiness.DB_TABLE_ACCOUNT, condition, mgo.Strong)
	if !isExist {
		xylog.Error(uid, "user doesn't exist")
	}
	return
}

func (db *BatteryDB) IsIdentityExist(identity uint64) (isExist bool) {
	condition := bson.M{"identity": identity}
	isExist, _ = db.IsRecordExisting(xybusiness.DB_TABLE_ACCOUNT, condition, mgo.Strong)
	return
}
