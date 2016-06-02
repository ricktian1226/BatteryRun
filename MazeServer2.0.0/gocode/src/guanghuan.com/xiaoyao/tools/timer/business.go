// business
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"flag"
	"gopkg.in/mgo.v2/bson"
	xydb "guanghuan.com/xiaoyao/common/db"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"time"
)

type Config struct {
	DBUrl        string
	DBName       string
	InfoInterval uint
	FreeInterval uint
}

var DefConfig = &Config{}
var DefDB *xydb.XYDB

const DB_TABLE_LOTTO_CONFIG = "lottoconfig"

func ProcessCmd() {
	flag.StringVar(&DefConfig.DBUrl, "dburl", "mongodb://localhost:27017/brdb01", "mongodb url")
	flag.StringVar(&DefConfig.DBName, "dbname", "brdb01", "mongodb db name")
	flag.UintVar(&DefConfig.InfoInterval, "infointerval", 0, "lotto sysinfo timer interval")
	flag.UintVar(&DefConfig.FreeInterval, "freeinterval", 0, "lotto free rest timer interval")
}

func Init() {
	ProcessCmd()
	xylog.ProcessCmdAndApply()
	//flag.Usage()

	DefDB = xydb.NewXYDB(DefConfig.DBUrl, DefConfig.DBName)
	err := DefDB.OpenDB()
	if nil != err {
		xylog.Error("Open DB failed : %v", err)
	}
	xylog.Debug("Open DB Succeed")
}

func UpdateSysLottoInfo() (err error) {
	c := DefDB.OpenTable(DB_TABLE_LOTTO_CONFIG)
	defer c.Close()

	fields := bson.M{"refreshtimestamp": time.Now().Unix()}

	err = DefDB.UpdateMultipleFields(DB_TABLE_LOTTO_CONFIG, bson.M{}, fields, true)
	if nil != err {
		xylog.Error("refresh lottosysinfo failed : %v", err)
	}
	xylog.Debug("refresh lottosysinfo Succeed")
	return
}

func UpsertSysLottoInfo() (err error) {
	c := DefDB.OpenTable(DB_TABLE_LOTTO_CONFIG)
	defer c.Close()

	t := time.Now().Unix()
	info := &battery.LottoConfigItem{
		Name:  proto.String("syslottoinfotimestamp"),
		Value: proto.Int64(t),
	}

	condition := bson.M{"name": "syslottoinfotimestamp"}
	_, err = c.Upsert(condition, info)
	if nil != err {
		xylog.Error("upsert lottosysinfo.syslottoinfotimestamp failed : %v", err)
	}
	xylog.Debug("upsert lottosysinfo.syslottoinfotimestamp Succeed")
	return
}

func UpsertSysFreeRest() (err error) {
	c := DefDB.OpenTable(DB_TABLE_LOTTO_CONFIG)
	defer c.Close()

	t := time.Now().Unix()
	info := &battery.LottoConfigItem{
		Name:  proto.String("syslottofreetimestamp"),
		Value: proto.Int64(t),
	}

	condition := bson.M{"name": "syslottofreetimestamp"}
	_, err = c.Upsert(condition, info)
	if nil != err {
		xylog.Error("upsert lottosysinfo.syslottofreetimestamp failed : %v", err)
	}
	xylog.Debug("upsert lottosysinfo.syslottofreetimestamp Succeed")
	return
}
