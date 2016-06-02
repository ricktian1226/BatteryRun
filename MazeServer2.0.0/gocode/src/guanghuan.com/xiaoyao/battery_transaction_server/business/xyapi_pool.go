// xyapi_pool
/*
   业务临时变量的缓存池实现，基于sync.Pool
*/
package batteryapi

import (
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"sync"
)

//AccountWithFlag
func FuncPoolAccountWithFlag() interface{} {
	return &AccountWithFlag{}
}

var PoolAccountWithFlag = sync.Pool{
	New: FuncPoolAccountWithFlag,
}

func GetPoolAccountWithFlag() (accountWithFlag *AccountWithFlag) {
	ok := false
	if accountWithFlag, ok = PoolAccountWithFlag.Get().(*AccountWithFlag); ok {
		return
	} else {
		accountWithFlag = &AccountWithFlag{}
	}
	return
}

func PutPoolAccountWithFlag(accountWithFlag *AccountWithFlag) {
	PoolAccountWithFlag.Put(accountWithFlag)
	return
}

//DBAccount
func FuncPoolDBAccount() interface{} {
	return &battery.DBAccount{}
}

var PoolDBAccount = sync.Pool{
	New: FuncPoolDBAccount,
}

func GetPoolDBAccount() (dbAccount *battery.DBAccount) {
	ok := false
	if dbAccount, ok = PoolDBAccount.Get().(*battery.DBAccount); ok {
		return
	} else {
		dbAccount = &battery.DBAccount{}
	}
	return
}

func PutPoolDBAccount(dbAccount *battery.DBAccount) {
	PoolDBAccount.Put(dbAccount)
	return
}

//TPID
func FuncPoolTPID() interface{} {
	return &battery.TPID{}
}

var PoolTPID = sync.Pool{
	New: FuncPoolTPID,
}

func GetPoolTPID() (tpid *battery.TPID) {
	ok := false
	if tpid, ok = PoolTPID.Get().(*battery.TPID); ok {
		return
	} else {
		tpid = &battery.TPID{}
	}
	return
}

func PutPoolTPID(tpid *battery.TPID) {
	PoolTPID.Put(tpid)
	return
}

//Error
func FuncPoolError() interface{} {
	return &battery.Error{}
}

var PoolError = sync.Pool{
	New: FuncPoolError,
}

func GetPoolError() (errStruct *battery.Error) {
	ok := false
	if errStruct, ok = PoolError.Get().(*battery.Error); ok {
		return
	} else {
		errStruct = &battery.Error{}
	}
	return
}

func PutPoolError(errStruct *battery.Error) {
	PoolError.Put(errStruct)
	return
}

//Account
func FuncPoolAccount() interface{} {
	return &battery.Account{}
}

var PoolAccount = sync.Pool{
	New: FuncPoolAccount,
}

func GetPoolAccount() (account *battery.Account) {
	ok := false
	if account, ok = PoolAccount.Get().(*battery.Account); ok {
		return
	} else {
		account = &battery.Account{}
	}
	return
}
func PutPoolAccount(account *battery.Account) {
	PoolAccount.Put(account)
	return
}
