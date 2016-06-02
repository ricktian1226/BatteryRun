// db
package xymoney

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	xydb "guanghuan.com/xiaoyao/common/db"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

var DefDB *xydb.XYDB

func Init(db *xydb.XYDB) {
	DefDB = db
}

func set(uid string, db *xydb.XYDB, money *battery.Money, mainType battery.MoneyType) (err error) {

	index := int(mainType)

	field := fmt.Sprintf("wallet.%d", index)
	err = db.UpdateData(DB_TABLE_ACCOUNT, bson.M{"uid": uid}, bson.M{"$set": bson.M{field: money}})
	return
}

func add(uid string, db *xydb.XYDB, amount uint32, mainType battery.MoneyType, subType battery.MoneySubType) (err error) {

	index := int(mainType)

	var subTypeStr string
	switch subType {
	case battery.MoneySubType_iap:
		subTypeStr = "iapamount"
	case battery.MoneySubType_oap:
		subTypeStr = "oapamount"
	case battery.MoneySubType_gain:
		subTypeStr = "gainamount"
	default:
		xylog.Error("[%s] Wrong MoneySubType : %d", uid, subType)
		err = xyerror.ErrBadInputData
	}
	field := fmt.Sprintf("wallet.%d.%s", index, subTypeStr)
	err = db.UpdateData(DB_TABLE_ACCOUNT, bson.M{"uid": uid}, bson.M{"$inc": bson.M{field: amount}})
	return
}

// 账户增加多种代币
// uid string 玩家id
// moneys []*battery.MoneyItem 代币列表
// factor uint32 系数(100代表100%)
func addMultiple(uid string, db *xydb.XYDB, moneys []*battery.MoneyItem, factor uint32) (err error) {

	fields := bson.M{}

	for _, money := range moneys {
		index := money.GetType()
		amount := money.GetAmount()
		if amount > 0 {
			name := fmt.Sprintf("wallet.%d.gainamount", index)
			if _, ok := fields[name]; ok {
				fields[name] = fields[name].(uint32) + amount*uint32(factor)/100
			} else {
				fields[name] = amount * factor / 100
			}
		}
	}

	err = db.UpdateData(DB_TABLE_ACCOUNT, bson.M{"uid": uid}, bson.M{"$inc": fields})

	return
}

func query(uid string, db *xydb.XYDB, mainType battery.MoneyType) (money *battery.Money, err error) {

	index := int(mainType)

	condition := bson.M{"uid": uid}

	selector := bson.M{"wallet": 1} //只要查钱包信息即可

	account := &battery.DBAccount{}

	err = db.GetOneData(DB_TABLE_ACCOUNT, condition, selector, account, mgo.Strong) //钱包信息，必须从主节点上读取，保证一致性

	if nil == err {
		wallet := account.GetWallet()
		xylog.Debug(uid, "wallet : %v", wallet)
		if index >= len(wallet) {
			xylog.Error(uid, "[%s] Wrong index %d to len(wallet) %d", uid, index, len(wallet))
			err = xyerror.ErrBadInputData
			return
		}

		money = wallet[index]
	}

	return
}

func QueryWalletFromDB(uid string, db *xydb.XYDB) (wallet []*battery.Money, err error) {

	condition := bson.M{"uid": uid}
	selector := bson.M{"wallet": 1} //只要查钱包信息即可
	account := &battery.DBAccount{}

	err = db.GetOneData(DB_TABLE_ACCOUNT, condition, selector, account, mgo.Strong) //钱包信息，必须从主节点上读取，保证一致性

	if nil != err {
		xylog.Error(uid, "[%s] Get wallet failed : ", uid, err)
		return
	}

	wallet = account.GetWallet()

	return
}
