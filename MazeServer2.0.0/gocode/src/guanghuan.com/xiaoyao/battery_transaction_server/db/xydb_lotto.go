// xydb_lotto
package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//增加玩家系统抽奖信息
func (db *BatteryDB) AddLottoInfo(i *battery.SysLottoInfo) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_SYS_LOTTO_INFO, *i)
	return
}

//获取玩家的系统抽奖信息
// uid string 玩家id
// i *battery.SysLottoInfo 系统抽奖信息
func (db *BatteryDB) QuerySysLottoInfo(uid string, i *battery.SysLottoInfo) (err error) {
	c := db.OpenTable(xybusiness.DB_TABLE_SYS_LOTTO_INFO, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"uid": uid}
	query := c.Find(queryStr)
	err = query.One(i)
	err = xyerror.DBError(err)

	return
}

//刷新玩家的系统抽奖信息
// i *battery.SysLottoInfo
// slots 奖池信息
// timestamp 奖池刷新时间戳信息
// value 玩家内部价值信息
// freecount 系统发放免费抽奖次数
// gainfreecount 获取的免费抽奖次数（购买、抽奖获取等）
func (db *BatteryDB) UpdateSysLottoInfo(uid string, i *battery.SysLottoInfo) (err error) {
	condition := bson.M{"uid": uid}
	fields := bson.M{}

	if nil != i.GetSlots() {
		fields["slots"] = i.GetSlots()
	}

	if nil != i.Timestamp {
		fields["timestamp"] = i.GetTimestamp()
	}

	if nil != i.Value {
		fields["value"] = i.GetValue()
	}

	if nil != i.FreeCount {
		fields["freecount"] = i.GetFreeCount()
	}

	if nil != i.GainFreeCount {
		fields["gainfreecount"] = i.GetGainFreeCount()
	}

	if nil != i.FreecountRefreshTimestamp {
		fields["freecountrefreshtimestamp"] = i.GetFreecountRefreshTimestamp()
	}

	if nil != i.FreecountLimitation {
		fields["freecountlimitation"] = i.GetFreecountLimitation()
	}

	if nil != i.FreecountLimitationExpiredTimestamp {
		fields["freecountlimitationexpiredtimestamp"] = i.GetFreecountLimitationExpiredTimestamp()
	}

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_SYS_LOTTO_INFO, condition, fields, true)

	return
}

//增加玩家抽奖券
// uid string 玩家id
// amount uint32 抽奖券数目
func (db *BatteryDB) IncreaseUserLottoTicket(uid string, amount uint32) (err error) {
	fields := bson.M{"$inc": bson.M{"gainfreecount": amount}}
	condition := bson.M{"uid": uid}
	err = db.UpdateData(xybusiness.DB_TABLE_SYS_LOTTO_INFO, condition, fields)
	return
}

// 增加抽奖事务记录
//t *battery.LottoTransaction 抽奖事务指针
func (db *BatteryDB) AddLottoTransaction(t *battery.LottoTransaction) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_LOTTO_TRANSACTION, *t)
	return
}

// 获取抽奖事务信息
// uid string 玩家标识
// lottoid uint64 抽奖事务id
// parentlottoid int64 抽奖父id事务
// t *battery.LottoTransaction 抽奖事务指针
func (db *BatteryDB) GetLottoTransaction(uid string, lottoid, parentlottoid uint64, t *battery.LottoTransaction) (err error) {
	c := db.OpenTable(xybusiness.DB_TABLE_LOTTO_TRANSACTION, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"uid": uid, "lottoid": lottoid, "parentlottoid": parentlottoid}
	query := c.Find(queryStr)
	err = query.One(t)
	err = xyerror.DBError(err)

	return
}

func (db *BatteryDB) QueryLottoTransactionsByParentLottoId(uid string, parentlottoid uint64, lottoTransactions *[]*battery.LottoTransaction) (err error) {
	c := db.OpenTable(xybusiness.DB_TABLE_LOTTO_TRANSACTION, mgo.Strong)
	defer c.Close()

	queryStr := bson.M{"uid": uid, "parentlottoid": parentlottoid}
	query := c.Find(queryStr)
	err = query.All(lottoTransactions)
	err = xyerror.DBError(err)

	return
}

func (db *BatteryDB) UpdateLottoTransaction(t *battery.LottoTransaction) (err error) {
	c := db.OpenTable(xybusiness.DB_TABLE_LOTTO_TRANSACTION, mgo.Strong)
	defer c.Close()

	condition := bson.M{"uid": t.GetUid(), "lottoid": t.GetLottoid(), "parentlottoid": t.GetParentlottoid()}
	fields := bson.M{}

	if nil != t.GetStates() {
		fields["states"] = t.GetStates()
	}

	err = db.UpdateMultipleFields(xybusiness.DB_TABLE_LOTTO_TRANSACTION, condition, fields, true)

	return
}

func (db *BatteryDB) PushLottoTransactionState(t *battery.LottoTransaction, state *battery.LottoStateEntry) (err error) {
	condition := bson.M{"uid": t.GetUid(), "lottoid": t.GetLottoid(), "parentlottoid": t.GetParentlottoid()}
	fields := bson.M{"$push": bson.M{"states": state}, "$set": bson.M{"state": state.GetState()}}
	err = db.UpdateData(xybusiness.DB_TABLE_LOTTO_TRANSACTION, condition, fields)
	return
}

func (db *BatteryDB) AddLottoLog(l *battery.LottoLog) (err error) {
	err = db.AddData(xybusiness.DB_TABLE_LOTTO_LOG, l)
	return
}
