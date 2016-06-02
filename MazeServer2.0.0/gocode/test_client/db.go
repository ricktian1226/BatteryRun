// db
package main

import (
	xydb "guanghuan.com/xiaoyao/common/db"
	xylog "guanghuan.com/xiaoyao/common/log"
	//"labix.org/v2/mgo"
	proto "code.google.com/p/goprotobuf/proto"
	//"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//apnsdb "guanghuan.com/xiaoyao/battery_apns_server/db"
	batterydb "guanghuan.com/xiaoyao/battery_transaction_server/db"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//xymoney "guanghuan.com/xiaoyao/superbman_server/money"
	//"guanghuan.com/xiaoyao/superbman_server/server"
	"math/rand"
)

const (
	RoleType_A = iota
	RoleType_B
)

const (
	MoneyType_Coin = iota
	MoneyType_Diamond
)

type Role struct {
	Type int
}

type Money struct {
	Type   int
	Amount int
}

type UserInfo struct {
	Uid   string
	Money []Money
	Role  []Role
}

type DBConfig struct {
	DBUrl  string
	DBName string
}

var DefDB *xydb.XYDB

func DB_Init() {
	DefDB = xydb.NewXYDB(DefConfig.DBUrl, DefConfig.DBName)
	err := DefDB.OpenDB()
	if nil != err {
		xylog.ErrorNoId("Open DB failed : %v", err)
	}
	xylog.DebugNoId("Open DB Succeed")

	//batterydb.NewTestBatteryDB(DefConfig.DBUrl, DefConfig.DBName)
}

func test_db_insert_user() (err error) {

	userInfo := &UserInfo{
		Uid: "139752879435073500001623",
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			m := &Money{
				Type:   i,
				Amount: j + 1,
			}
			userInfo.Money = append(userInfo.Money, *m)
		}
	}

	err = DefDB.AddData("user", userInfo)
	return
}

func test_db_query_user() (err error) {
	dbsession := DefDB.CopySession()
	defer dbsession.Close()

	collection := dbsession.DB(DefConfig.DBName).C("user")

	match := bson.M{"$match": bson.M{"uid": "139752879435073500001622"}}
	project := bson.M{"$project": bson.M{"_id": 0, "money": 1}}
	unwind := bson.M{"$unwind": "$money"}
	//sort := bson.M{"$sort": bson.M{"money.type": -1}}
	match1 := bson.M{"$match": bson.M{"money.type": 1}}
	pipe := collection.Pipe([]bson.M{match, project, unwind /*, sort*/, match1})
	iter := pipe.Iter()

	var moneys []bson.M
	iter.All(&moneys)
	defer iter.Close()

	xylog.DebugNoId("user : %v", moneys)
	return
}

func test_db_update_user() (err error) {
	dbsession := DefDB.CopySession()
	defer dbsession.Close()

	collection := dbsession.DB(DefConfig.DBName).C("user")

	selector := bson.M{"uid": "139752879435073500001623", "money.type": 1}

	updater := bson.M{"$inc": bson.M{"money.0.amount": 1, "money.1.amount": 2}}

	err = collection.Update(selector, updater)

	xylog.DebugNoId("err : %v", err)
	return
}

func test_db_update_user_money() (err error) {
	//money1 := &battery.MoneyItem{
	//	Type:   battery.MoneyType_chip.Enum(),
	//	Amount: proto.Uint32(10),
	//}
	//money2 := &battery.MoneyItem{
	//	Type:   battery.MoneyType_diamond.Enum(),
	//	Amount: proto.Uint32(10),
	//}

	//moneys := []*battery.MoneyItem{money1, money2}

	//xymoney.AddMultiple("1411365678693098337130", moneys, 100)

	return
}

func test_db_insert_user_checkpoint() (err error) {
	for j := 0; j < 10; j++ {
		checkPoint := battery.DBUserCheckPoint{
			Uid:          proto.String("1413430690598372090525"),
			CheckPointId: proto.Uint32(uint32(j)),
			Score:        proto.Uint64(uint64(rand.Int63())),
			Charge:       proto.Uint64(uint64(rand.Int63())),
		}

		DefDB.AddData("usercheckpoint", checkPoint)
	}
	return
}

func test_db_query_user_checkpoints() (err error) {
	dbUserCheckPoints := []*battery.DBUserCheckPoint{}
	begin := uint32(0)
	end := uint32(10)
	err = batterydb.DefBatteryDB.QueryUserCheckPoints("1411365678693098300000", begin, end, &dbUserCheckPoints)
	xylog.DebugNoId("QueryUserCheckPoints [%d, %d] result : %v , Error : %v", begin, end, dbUserCheckPoints, err)
	return
}

func test_db_query_user_checkpoint_detail() (err error) {
	dbUserCheckPoint := &battery.DBUserCheckPoint{}
	err = batterydb.DefBatteryDB.QueryUserCheckPointDetail("1411365678693098300000", uint32(2), dbUserCheckPoint)
	xylog.DebugNoId("QueryUserCheckPointDetail result : %v , Error : %v", dbUserCheckPoint, err)
	return
}

func test_db_query_friends_uid() (err error) {
	tpids := make([]*battery.IDMap, 0)
	sids := []string{"123456", "123457"}
	//err = xybusiness.QueryTpidsBySids(sids, battery.ID_SOURCE_SRC_SINA_WEIBO, &tpids)
	xylog.DebugNoId("QueryFriendsUid for (%v) result : %v , Error : %v", sids, tpids, err)
	return
}

func test_db_upsert_user_checkpoint() (err error) {
	dbDetail := &battery.DBUserCheckPoint{
		Uid:          proto.String("1411365678693098300000"),
		CheckPointId: proto.Uint32(0),
		Score:        proto.Uint64(101),
		Charge:       proto.Uint64(101),
	}

	err = batterydb.DefBatteryDB.UpsertUserCheckPoint(dbDetail)
	xylog.DebugNoId("UpsertUserCheckPoint (%v), Error : %v", dbDetail, err)
	return
}

func test_db_query_user_checkpoints_sum() (err error) {
	checkPointsSum := &battery.DBUserCheckPoint{}
	err = batterydb.DefBatteryDB.QueryUserCheckPointsSum("1411365678693098337130", checkPointsSum)
	return
}

func test_db_increase_user_lotto_ticket() (err error) {
	err = batterydb.DefBatteryDB.IncreaseUserLottoTicket("1411365678693098337130", 3)
	return
}

func test_db_query_user_account_somefields() (err error) {
	uid := "1425352623264100099847"
	account := &battery.DBAccount{}
	condition := bson.M{"uid": uid}
	selector := bson.M{"uid": 1, "name": 1}
	err = batterydb.DefBatteryDB.GetOneData(xybusiness.DB_TABLE_ACCOUNT, condition, selector, account, mgo.Strong)
	xylog.DebugNoId("account[%s] : %v, Error : %v", uid, account, err)
	return
}

func test_db_query_user_account_distinct() (err error) {
	accounts := make([]*battery.DBAccount, 0)
	condition := bson.M{}
	selector := bson.M{"deviceid": 1}
	//err = DefDB.GetAllDataDistinct(xybusiness.DB_TABLE_ACCOUNT, condition, selector, "deviceid", 0, &accounts, mgo.Strong)
	err = DefDB.GetAllData(xybusiness.DB_TABLE_ACCOUNT, condition, selector, 0, &accounts, mgo.Strong)
	xylog.DebugNoId("accounts: %v, Error : %v", accounts, err)
	return
}
