package xydb

import (
	//	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	//xyutil "guanghuan.com/xiaoyao/common/util"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	//	"math/rand"
	//	"time"
	//"reflect"
)

type DBInterface interface {
	OpenDB() (err error)
	CloseDB()
	IsRecordExistingWithField(tablename string, field string, value interface{}, consistency mgo.Mode) (isExisting bool, err error)
	GetOneDataWithField(tablename, field string, value, selector, obj_ret interface{}, consistency mgo.Mode) (err error)
	AddData(tablename string, obj interface{}) (err error)
}

type XYDB struct {
	session *mgo.Session
	dburl   string
	dbname  string
}

func NewXYDB(url string, dbname string) *XYDB {
	db := &XYDB{
		dburl:  url,
		dbname: dbname,
	}
	return db
}

// Connect to database
func (db *XYDB) OpenDB() error {
	var err error
	db.session, err = mgo.Dial(db.dburl)
	//xylog.Debug("db.session %p", db.session)
	if err != nil {
		xylog.ErrorNoId("!!! DB ERROR: %s", err.Error())
	}

	err = xyerror.DBError(err)
	return err
}

// Close connections
func (db *XYDB) CloseDB() {
	if db.session != nil {
		db.session.Close()
	}
}

func (db *XYDB) CopySession() *mgo.Session {
	return db.session.Copy()
}

type XYTable struct {
	mgo.Collection
	session  *mgo.Session
	database *mgo.Database
	//	table   *mgo.Collection
}

func (db *XYDB) OpenTable(tablename string, consistency mgo.Mode) (table *XYTable) {
	s := db.CopySession()
	s.SetMode(consistency, false)
	d := s.DB(db.dbname)
	c := d.C(tablename)

	table = &XYTable{*c, s, d}
	//	table.session = s

	return
}
func (db *XYDB) CloseTable(table *XYTable) {
	table.Close()
}
func (tbl *XYTable) Close() {
	if tbl.session != nil {
		tbl.session.Close()
	}
}
func (tbl *XYTable) OpenATable(tablename string) (table *XYTable) {
	c := tbl.database.C(tablename)
	table = &XYTable{*c, tbl.session, tbl.database}

	return
}

//func (db *XYDB) NewId() string {
//	return xyutil.NewId()
//}

// add a record to table
func (db *XYDB) AddData(tablename string, obj interface{}) (err error) {

	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	err = tbl.Insert(obj)

	err = xyerror.DBError(err)
	return
}

// add several datas to table
func (db *XYDB) AddDatas(tablename string, objs []interface{}) (err error) {

	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	err = tbl.Insert(objs...)

	err = xyerror.DBError(err)
	return
}

func (db *XYDB) UpsertData(tablename string, condition interface{}, obj interface{}) (err error) {
	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	_, err = tbl.Upsert(condition, obj)

	err = xyerror.DBError(err)
	return
}

// update a record with condition
func (db *XYDB) UpdateData(tablename string, condition interface{}, obj interface{}) (err error) {
	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	err = tbl.Update(condition, obj)

	err = xyerror.DBError(err)
	return
}

// get first record from table with condition
//func (db *XYDB) GetOneData(tablename string, condition interface{}, obj interface{}, consistency mgo.Mode) (err error) {
//	dbsession := db.session.Copy()
//	defer dbsession.Close()

//	dbsession.SetMode(consistency, false)

//	tbl := dbsession.DB(db.dbname).C(tablename)

//	err = tbl.Find(condition).One(obj)

//	err = xyerror.DBError(err)

//	return
//}

//查询单条记录
// tablename string 集合名称
// condition interface{} 查询条件
// selector  interface{} 返回字段
// obj interface{} 返回结果保存
// consistency mgo.Mode 查询模式
func (db *XYDB) GetOneData(tablename string, condition, selector, obj interface{}, consistency mgo.Mode) (err error) {
	dbsession := db.session.Copy()
	defer dbsession.Close()

	dbsession.SetMode(consistency, false)

	tbl := dbsession.DB(db.dbname).C(tablename)
	if nil != selector {
		err = tbl.Find(condition).Select(selector).One(obj)
	} else {
		err = tbl.Find(condition).One(obj)
	}

	err = xyerror.DBError(err)

	return
}

// get all records from table with condition
// tablename string 集合名称
// condition interface{} 查询条件
// selector  interface{} 返回字段
// obj interface{} 返回结果保存
// consistency mgo.Mode 查询模式
func (db *XYDB) GetAllData(tablename string, condition, selector interface{}, max_count int, obj interface{}, consistency mgo.Mode) (err error) {
	tbl := db.OpenTable(tablename, consistency)
	defer tbl.Close()

	query := tbl.Find(condition)
	if nil != selector {
		query = query.Select(selector)
	}

	if max_count > 0 {
		query = query.Limit(max_count)
	}

	err = query.All(obj)

	err = xyerror.DBError(err)
	return
}

// GetAllDataDistinct get all records from table with condition
// tablename string 集合名称
// condition interface{} 查询条件
// selector  interface{} 返回字段
// obj interface{} 返回结果保存
// consistency mgo.Mode 查询模式
func (db *XYDB) GetAllDataDistinct(tablename string, condition, selector interface{}, distinct string, max_count int, obj interface{}, consistency mgo.Mode) (err error) {
	tbl := db.OpenTable(tablename, consistency)
	defer tbl.Close()

	//query := tbl.Find(condition)
	//if nil != selector {
	//	query = query.Select(selector)
	//}

	//if max_count > 0 {
	//	query = query.Limit(max_count)
	//}

	//err = query.Distinct(distinct, obj)
	err = tbl.Find(condition).Distinct(distinct, obj)

	//err = query.All(obj)

	err = xyerror.DBError(err)
	return
}

// get count of records
func (db *XYDB) GetRecordCount(tablename string, condition interface{}, consistency mgo.Mode) (count int, err error) {
	tbl := db.OpenTable(tablename, consistency)
	defer tbl.Close()

	//只是获取记录条数的话，我们只要获取_id字段就可以了
	selector := bson.M{"_id": 1}
	count, err = tbl.Find(condition).Select(selector).Count()
	err = xyerror.DBError(err)
	return
}

// check if certain record exists by checking record count
func (db *XYDB) IsRecordExisting(tablename string, condition interface{}, consistency mgo.Mode) (isExisting bool, err error) {
	var count int
	count, err = db.GetRecordCount(tablename, condition, consistency)
	isExisting = (count > 0)
	return
}

// check if certain record exist by verifying one field
func (db *XYDB) IsRecordExistingWithField(tablename string, field string, value interface{}, consistency mgo.Mode) (isExisting bool, err error) {
	condition := bson.M{field: value}
	return db.IsRecordExisting(tablename, condition, consistency)
}

// update first record where its certain filed has desired value
func (db *XYDB) UpdateDataWithField(tablename string, field string, value interface{}, obj_in interface{}) (err error) {
	condition := bson.M{field: value}
	err = db.UpdateData(tablename, condition, obj_in)

	err = xyerror.DBError(err)
	return
}

// get first record where its certain filed has desired value
//查询返回一个文档，查询条件只有一个字段
// tablename string 集合名称
// field string 字段名称
// value interface{} 字段值
// selector interface{} 返回字段定义
// obj_ret interface{} 返回文档信息
// consistency mgo.Mode 查询模式
func (db *XYDB) GetOneDataWithField(tablename, field string, value, selector, obj_ret interface{}, consistency mgo.Mode) (err error) {
	condition := bson.M{field: value}
	err = db.GetOneData(tablename, condition, selector, obj_ret, consistency)

	err = xyerror.DBError(err)
	return
}

//查询返回第一个文档，查询带sort排序
// tablename string 集合名称
// field string 字段名称
// value interface{} 字段值
// selector interface{} 返回字段定义
// sort_field string 排序字段
// obj_ret interface{} 返回文档信息
// consistency mgo.Mode 查询模式
func (db *XYDB) GetFirstData(tablename string, condition, selector interface{}, sort_field string, obj_ret interface{}, consistency mgo.Mode) (err error) {
	tbl := db.OpenTable(tablename, consistency)
	defer tbl.Close()
	err = tbl.Find(condition).Select(selector).Sort(sort_field).One(obj_ret)
	err = xyerror.DBError(err)

	return
}

// update one field of all records meet the condition
func (db *XYDB) UpdateOneField(tablename string, condition interface{}, field_name string, field_value interface{}, updateAll bool) (err error) {
	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	update_str := bson.M{"$set": bson.M{field_name: field_value}}
	if updateAll {
		_, err = tbl.UpdateAll(condition, update_str)
	} else {
		err = tbl.Update(condition, update_str)
	}

	err = xyerror.DBError(err)
	return
}

func (db *XYDB) UpdateMultipleFields(tablename string, condition interface{}, fields bson.M, updateAll bool) (err error) {

	if 0 >= len(fields) {
		return
	}

	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	update_str := bson.M{"$set": fields}

	if updateAll {
		_, err = tbl.UpdateAll(condition, update_str)
	} else {
		err = tbl.Update(condition, update_str)
	}

	err = xyerror.DBError(err)
	return
}

func (db *XYDB) UpdateMultipleFieldsWithPush(tablename string, condition interface{}, seter bson.M, pusher bson.M, updateAll bool) (err error) {

	if 0 >= len(seter) && 0 >= len(pusher) {
		return
	}

	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	update_str := bson.M{"$set": seter, "$push": pusher}

	if updateAll {
		_, err = tbl.UpdateAll(condition, update_str)
	} else {
		err = tbl.Update(condition, update_str)
	}

	err = xyerror.DBError(err)
	return
}

// remove one records meet the condition
func (db *XYDB) RemoveData(tablename string, condition interface{}) (err error) {
	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	err = tbl.Remove(condition)

	err = xyerror.DBError(err)
	return
}

// remove all records meet the condition
func (db *XYDB) RemoveAllData(tablename string, condition interface{}) (err error) {
	tbl := db.OpenTable(tablename, mgo.Strong)
	defer tbl.Close()

	_, err = tbl.RemoveAll(condition)

	err = xyerror.DBError(err)
	return
}

func (db *XYDB) SetField(name string, value interface{}, fields bson.M) {

	switch t := value.(type) {
	case string:
		if "" != t {
			fields[name] = t
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if 0 != t {
			fields[name] = t
		}
	case []interface{}:
		if nil != t {
			fields[name] = t
		}
	}
}

func (db *XYDB) GetByPipe(collectionName string, pipeMap []bson.M, consistency mgo.Mode) (result []interface{}, err error) {
	table := db.OpenTable(collectionName, consistency)
	defer table.Close()

	result = make([]interface{}, 0)

	pipe := table.Pipe(pipeMap)
	iter := pipe.Iter()

	err = iter.All(&result)
	defer iter.Close()

	//xylog.Debug("result : %v", result)

	return
}
