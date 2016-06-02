package main

import (
	//proto "code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	xydb "guanghuan.com/xiaoyao/common/db"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"labix.org/v2/mgo/bson"
	"os"
)

type Config struct {
	DBUrl  string
	DBName string
	Limit  int
}

var DefConfig = &Config{}
var DefDB *xydb.XYDB

const DB_TABLE_LOTTO_LOG = "lottolog"

const FILE_NAME = "lottolog_"

var DefFile *os.File

func ProcessCmd() {
	flag.StringVar(&DefConfig.DBUrl, "dburl", "mongodb://localhost:27017/brdb01", "mongodb url")
	flag.StringVar(&DefConfig.DBName, "dbname", "brdb01", "mongodb db name")
	flag.IntVar(&DefConfig.Limit, "limit", 10000, "query limit")
}

func Init() (err error) {
	ProcessCmd()
	xylog.ProcessCmdAndApply()

	DefDB = xydb.NewXYDB(DefConfig.DBUrl, DefConfig.DBName)
	err = DefDB.OpenDB()
	if nil != err {
		xylog.Error("Open DB failed : %v", err)
	}
	xylog.Debug("Open DB Succeed")

	fileName := string(FILE_NAME) + xyutil.CurTimeStr() + ".csv"
	DefFile, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 666)
	if err != nil {
		xylog.Error("Open file %s failed : %v", fileName, err)
	}
	xylog.Debug("Open file %s Succeed", fileName)
	return
}

func QueryLottoLog() (err error) {
	DefFile.WriteString("uid,slots,selected,value\n")

	condition := bson.M{"cmd": 3}

	logs := make([]*battery.LottoLog, 0)

	tbl := DefDB.OpenTable(DB_TABLE_LOTTO_LOG)
	defer tbl.Close()
	query := tbl.Find(condition)
	query = query.Limit(DefConfig.Limit).Sort("timestamp")
	err = query.All(&logs)
	err = xyerror.DBError(err)

	if nil != err {
		xylog.Error("query lotto log failed : %v", err)
	}

	xylog.Debug("query lotto log succeed")

	for _, log := range logs {
		Export2File(log)
	}
	return
}

func Export2File(log *battery.LottoLog) {
	uid := log.GetUid()
	slots := log.GetSlots()
	selected := log.GetSelected()
	value := log.GetValue()
	decuct := log.GetDeduct()
	opdate := log.GetOpdate()

	var strSlots string
	for i, slot := range slots {
		if i != 0 {
			strSlots += ":"
		}
		item := slot.GetItems()
		strSlots += fmt.Sprintf("%d", item.GetId())
	}
	record := fmt.Sprintf("%s,%s,%d,%d,%d,%s\n", uid, strSlots, selected, value, decuct, opdate)
	DefFile.WriteString(record)
	//xylog.Debug("record : %s", record)
}

func main() {
	err := Init()
	if err != nil {
		os.Exit(0)
	}

	err = QueryLottoLog()
	if err != nil {
		os.Exit(0)
	}
}
