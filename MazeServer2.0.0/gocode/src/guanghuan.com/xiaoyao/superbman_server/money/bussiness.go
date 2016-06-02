package xymoney

import (
	"code.google.com/p/goprotobuf/proto"
	//"fmt"
	"guanghuan.com/xiaoyao/common/db"
	xylog "guanghuan.com/xiaoyao/common/log"
	//xyutil "guanghuan.com/xiaoyao/common/util"
	//xyversion "guanghuan.com/xiaoyao/common/version"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

func Amount(uid string, mainType battery.MoneyType, wallet []*battery.Money) (amount uint32) {

	index := int(mainType)

	if nil == wallet || index < 0 || index > len(wallet)-1 {
		xylog.Error("[%s] invalid index %d len(wallet) %d", uid, index, len(wallet))
		return 0
	}

	amount = wallet[index].GetIapamount() + wallet[index].GetOapamount() + wallet[index].GetGainamount()
	return
}

func AmountDetail(uid string, mainType battery.MoneyType, wallet []*battery.Money) (detail *battery.Money) {

	index := int(mainType)

	if nil == wallet || index < 0 || index > len(wallet)-1 {
		xylog.Error("[%s] invalid index %d", uid, index)
		return nil
	}

	detail = wallet[index]
	return
}

func CheckWallet(uid string, moneys []*battery.MoneyItem, wallet []*battery.Money) bool {
	for _, money := range moneys {
		if !checkwallet(uid, money, wallet) {
			return false
		}
	}
	return true
}

func Check(uid string, spendAmount uint32, money *battery.Money) bool {
	if !check(uid, spendAmount, money) {
		return false
	}
	return true
}

//消费代币接口
// uid string 玩家id
// moneys []*battery.MoneyItem 消费代币列表
// bCheck bool 是否需要校验钱包
// delay bool 是否延迟刷新钱包
func Consum(uid string, db *xydb.XYDB, account *battery.DBAccount, moneys []*battery.MoneyItem, bCheck, delay bool) (err error) {

	wallet := account.Wallet

	if bCheck && !CheckWallet(uid, moneys, wallet) {
		err = xyerror.ErrNoEnoughMoney
		return
	}

	for _, money := range moneys {
		err = consum(uid, db, money, &wallet, delay)
	}

	//todo log

	return
}

func check(uid string, spendAmount uint32, money *battery.Money) bool {
	iapAmount := money.GetIapamount()
	oapAmount := money.GetOapamount()
	gainAmount := money.GetGainamount()

	totalAmount := iapAmount + oapAmount + gainAmount

	if spendAmount > totalAmount {
		xylog.Debug("[%s] no enough amount %d for %d", uid, totalAmount, spendAmount)
		return false
	}

	return true
}

func checkwallet(uid string, price *battery.MoneyItem, wallet []*battery.Money) bool {
	index := int(price.GetType())

	if index >= len(wallet) {
		xylog.Error("[%s] Wrong index %d of wallet", uid, index)
		return false
	}

	iapAmount := wallet[index].GetIapamount()
	oapAmount := wallet[index].GetOapamount()
	gainAmount := wallet[index].GetGainamount()
	totalAmount := (iapAmount + oapAmount + gainAmount)

	if price.GetAmount() > totalAmount {
		xylog.Debug("[%s] no enough amount %d for %d of index %v", uid, totalAmount, price.GetAmount(), price.GetType())
		return false
	}

	return true
}

//代币消费基础接口
// uid string 玩家id
// price *battery.MoneyItem 消费金额
// wallet *[]*battery.Money 玩家钱包指针
// delay bool 是否延迟刷新
func consum(uid string, db *xydb.XYDB, price *battery.MoneyItem, wallet *[]*battery.Money, delay bool) (err error) {

	index := int(price.GetType())

	amountSpend := price.GetAmount()

	amountIap := (*wallet)[index].GetIapamount()
	amountOap := (*wallet)[index].GetOapamount()
	amountGained := (*wallet)[index].GetGainamount()
	amountTotal := amountIap + amountOap + amountGained

	xylog.Info("[%s] Spend %v: needs [%d], has [%d]+[%d]+[%d]=[%d]",
		uid, price.GetType(), amountSpend, amountIap, amountOap, amountGained, amountTotal)
	var (
		amountIapSpend    uint32
		amountOapSpend    uint32
		amountGainedSpend uint32
		amountToSpend     uint32 = amountSpend
	)

	if amountToSpend > 0 {
		if amountIap >= amountToSpend {
			amountIapSpend = amountToSpend
		} else {
			amountIapSpend = amountIap
		}
		amountToSpend -= amountIapSpend
		amountIap -= amountIapSpend
	}

	if amountToSpend > 0 {
		if amountOap >= amountToSpend {
			amountOapSpend = amountToSpend
		} else {
			amountOapSpend = amountOap
		}
		amountToSpend -= amountOapSpend
		amountOap -= amountOapSpend
	}

	if amountToSpend > 0 {
		if amountGained >= amountToSpend {
			amountGainedSpend = amountToSpend
		} else {
			amountGainedSpend = amountGained
			xylog.Warning("[%s] something wrong here.", uid)
		}
		amountToSpend -= amountGainedSpend
		amountGained -= amountGainedSpend
	}

	//var newMoney battery.Money
	(*wallet)[index].Type = price.GetType().Enum()
	(*wallet)[index].Iapamount = proto.Uint32(amountIap)
	(*wallet)[index].Oapamount = proto.Uint32(amountOap)
	(*wallet)[index].Gainamount = proto.Uint32(amountGained)

	//var left, spend battery.Money
	if !delay {
		err = set(uid, db, (*wallet)[index], price.GetType())
	}
	//if err == xyerror.DBErrOK {
	//	spend.Iapamount = proto.Uint32(amountIapSpend)
	//	spend.Oapamount = proto.Uint32(amountOapSpend)
	//	spend.Gainamount = proto.Uint32(amountGainedSpend)

	//	left.Iapamount = proto.Uint32(amountIap)
	//	left.Oapamount = proto.Uint32(amountOap)
	//	left.Gainamount = proto.Uint32(amountGained)
	//}

	return
}

//代币增加接口
// uid string 玩家id
// amount uint32 代币数目
// mainType battery.MoneyType 代币主类型
// subType battery.MoneySubType 代币子类型
// money *battery.Money 代币结构体
// delay bool 是否延迟刷新
func Add(uid string, db *xydb.XYDB, amount uint32, mainType battery.MoneyType, subType battery.MoneySubType, account *battery.DBAccount, delay bool) (err error) {

	err = xyerror.ErrOK

	if delay {
		switch subType {
		case battery.MoneySubType_iap:
			account.Wallet[int(mainType)].Iapamount = proto.Uint32(account.Wallet[int(mainType)].GetIapamount() + amount)
		case battery.MoneySubType_oap:
			account.Wallet[int(mainType)].Oapamount = proto.Uint32(account.Wallet[int(mainType)].GetOapamount() + amount)
		case battery.MoneySubType_gain:
			account.Wallet[int(mainType)].Gainamount = proto.Uint32(account.Wallet[int(mainType)].GetGainamount() + amount)
		default: //donothing
		}
		return
	}

	//先把代币增加到账户上
	err = add(uid, db, amount, mainType, subType)
	if err != nil {
		xylog.Error("[%s] Add %d to index %d subType %d failed", uid, amount, mainType, subType)
	}

	xylog.Debug("[%s] Add %d to index %d subType %d succeed", uid, amount, mainType, subType)

	////重新查一遍代币账户，用于记录日志
	//money, err = query(uid, mainType)

	return
}

func AddMultiple(uid string, db *xydb.XYDB, moneys []*battery.MoneyItem, factor uint32) (err error) {
	return addMultiple(uid, db, moneys, factor)
}

func Set(uid string, db *xydb.XYDB, money *battery.Money, mainType battery.MoneyType) (err error) {
	//设置账户信息
	err = set(uid, db, money, mainType)
	if err != nil {
		xylog.Error("[%s] set %v to %v failed", uid, mainType, *money)
	}

	return
}

//查询玩家某种
func Query(uid string, db *xydb.XYDB, mainType battery.MoneyType) (money *battery.Money, err error) {
	money, err = query(uid, db, mainType)
	return
}

//查询玩家钱包信息
// uid string 玩家id
//return:
// wallet []*battery.MoneyItem 消息用钱包数据
func QueryWallet(uid string, db *xydb.XYDB) (wallet []*battery.MoneyItem, err error) {
	var dbWallet = []*battery.Money{}
	dbWallet, err = QueryWalletFromDB(uid, db)
	if err == xyerror.ErrOK {
		//wallet = make([]*battery.MoneyItem, 0)
		//for _, money := range dbWallet {
		//	item := &battery.MoneyItem{
		//		Type:   money.GetType().Enum(),
		//		Amount: proto.Uint32(money.GetIapamount() + money.GetOapamount() + money.GetGainamount()),
		//	}
		//	wallet = append(wallet, item)
		//}
		wallet = GetMoneyItemsFromMoneys(dbWallet)
	}

	return
}

//将数据库钱包转换为消息钱包
// wallet []*battery.Money 数据库钱包
// moneyItems []*battery.MoneyItem 消息用钱包
func GetMoneyItemsFromMoneys(wallet []*battery.Money) (moneyItems []*battery.MoneyItem) {
	for _, money := range wallet {
		moneyItem := &battery.MoneyItem{
			Type:   money.Type.Enum(),
			Amount: proto.Uint32(money.GetIapamount() + money.GetOapamount() + money.GetGainamount()),
		}
		moneyItems = append(moneyItems, moneyItem)
	}
	return
}

//查询玩家钱包信息（数据库）
// uid string 玩家id
func QueryDBWallet(uid string, db *xydb.XYDB) (wallet []*battery.Money, err error) {
	wallet, err = QueryWalletFromDB(uid, db)
	return
}
