package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xycache "guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//const INVALID_ROLEID = 150010000 //非法的角色id

//获取玩家当前拥有的角色信息
// uid string 玩家id
//return:
// userRoleInfo *battery.RoleInfoTableItem 玩家角色信息结构体
func (db *BatteryDB) GetUserRoleInfo(uid string) (userRoleInfo *battery.UserRoleInfo, err error) {
	userRoleInfo = &battery.UserRoleInfo{}
	condition := bson.M{"uid": uid}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_ROLEINFO, condition, selector, &userRoleInfo, mgo.Monotonic)
	return
}

//查询玩家当前选中的角色
// uid string 玩家id
//return:
// roleID uint64 当前使用的角色id
func (db *BatteryDB) GetSelectRoleID(uid string) (roleID uint64) {
	roleID = xycache.INVALID_ROLEID //不存在该职业
	var err error
	var userRoleInfo *battery.UserRoleInfo
	userRoleInfo, err = db.GetUserRoleInfo(uid)
	if err != xyerror.ErrOK {
		xylog.Error(xylog.DefaultLogId, "GetUserRoleInfo failed : %v", err)
	} else {
		roleID = userRoleInfo.GetSelectRoleId()
	}
	return
}

//to delete
//判断某个角色配置信息有效，这个有效值将决定角色列表，角色等级配置表，和 角色等级商品表的数据哪些是有效的
//func (db *BatteryDB) IsRoleConfigValid(roleID uint64) (isValid bool) {
//	condition := bson.M{"roleid": roleID}
//	isValid, _ = db.IsRecordExisting(DB_TABLE_ROLEINFO_CONFIG, condition, mgo.Strong)
//	return isValid
//}

//----------------角色信息配置相关------------
//将老的配置数据全部设置为无效
func (db *BatteryDB) ResetAllRoleInfoConfigDataInvalid() (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_ROLEINFO_CONFIG, mgo.Strong)
	defer tbl.Close()
	condition := bson.M{"id": bson.M{"$gte": 0}}
	update_str := bson.M{"$set": bson.M{"isvalid": false}}
	_, err = tbl.UpdateAll(condition, update_str)
	return
}

//删除所有的roleinfoconfig数据
func (db *BatteryDB) RemoveAllRoleInfoConfig() (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_ROLEINFO_CONFIG, mgo.Strong)
	defer tbl.Close()
	//condition := bson.M{"id": bson.M{"$gte": 0}}
	_, err = tbl.RemoveAll(bson.M{})
	return
}

//导入角色配置信息
func (db *BatteryDB) UpsertRoleConfig(config battery.DBRoleInfoConfig) (err error) {
	condition := bson.M{"id": config.GetId()}
	err = db.UpsertData(xybusiness.DB_TABLE_ROLEINFO_CONFIG, condition, config)
	return
}

//导入角色等级加成信息
func (db *BatteryDB) UpsertRoleLevelBonus(config *battery.DBRoleLevelBonusItem) (err error) {
	condition := bson.M{"id": config.GetId()}
	err = db.UpsertData(xybusiness.DB_TABLE_ROLE_LEVEL_BONUS, condition, config)
	return err
}

//删除所有的roleinfoconfig数据
func (db *BatteryDB) RemoveAllRoleLevelBonus() (err error) {
	tbl := db.OpenTable(xybusiness.DB_TABLE_ROLE_LEVEL_BONUS, mgo.Strong)
	defer tbl.Close()
	//condition := bson.M{"id": bson.M{"$gte": 0}}
	_, err = tbl.RemoveAll(bson.M{})
	return
}

//to delete
//func (db *BatteryDB) RemoveAllRoleInfoInPropTB() (err error) {
//	//先删除老数据
//	tbl := db.OpenTable(DB_TABLE_PROP, mgo.Strong)
//	defer tbl.Close()
//	condition := (bson.M{"$and": []bson.M{bson.M{"id": bson.M{"$gte": 150000000}}, bson.M{"id": bson.M{"$lt": 160000000}}}})
//	_, err = tbl.RemoveAll(condition)
//	return err
//}
