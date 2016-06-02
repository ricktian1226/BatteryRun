package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xycache "guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//upsert玩家角色信息
// userRoleInfo *battery.UserRoleInfo 玩家角色信息指针
func (db *BatteryDB) UpsertUserRoleInfo(userRoleInfo *battery.UserRoleInfo) (err error) {
	condition := bson.M{"uid": userRoleInfo.GetUid()}
	err = db.UpsertData(xybusiness.DB_TABLE_ROLEINFO, condition, userRoleInfo)
	return
}

//刷新玩家当前使用的角色
// uid string 玩家id
// selectRoleID uint64 当前使用角色
func (db *BatteryDB) UpdateSelectRoleID(uid string, selectRoleID uint64) (err error) {
	// 查询用户是否存在
	condition := bson.M{"uid": uid}
	err = db.UpdateOneField(xybusiness.DB_TABLE_ROLEINFO, condition, "selectroleid", selectRoleID, false)
	return
}

//判断玩家是否拥有角色
// uid string 玩家id
// roleId uint64 角色id
//return:
// isExist bool true 玩家已拥有对应角色，false 玩家未拥有对应角色
func (db *BatteryDB) IsRoleExisting(uid string, roleId uint64) (isExist bool) {

	var (
		userRoleInfo *battery.UserRoleInfo
		err          error
	)

	userRoleInfo, err = db.GetUserRoleInfo(uid)
	if err != xyerror.ErrOK {
		//xylog.Error("[%s] GetUserRoleInfo Error : %s", uid, err)
	} else {
		for _, roleInfoItem := range userRoleInfo.GetRoleInfoItemList() {
			if roleInfoItem != nil && roleInfoItem.GetRoleId() == roleId && roleInfoItem.GetCurLevel() > 0 && !roleInfoItem.GetIsLock() {
				return true
			}
		}
	}
	return false
}

//查询玩家角色等级
// uid string 玩家id
// roleId uint64 角色id
//return:
// level int32 等级
func (db *BatteryDB) GetRoleLevel(uid string, roleId uint64) (level int32) {

	var (
		err          error
		userRoleInfo *battery.UserRoleInfo
	)
	userRoleInfo, err = db.GetUserRoleInfo(uid)
	if err == xyerror.ErrOK {
		for _, roleInfoItem := range userRoleInfo.GetRoleInfoItemList() {
			if roleInfoItem != nil && roleInfoItem.GetRoleId() == roleId {
				return roleInfoItem.GetCurLevel()
			}
		}
	}

	return xycache.INVALID_LEVEL
}

//查询玩家选中的角色
// uid string 玩家id
//return:
// roleId uint64 选中的角色id
func (db *BatteryDB) GetSelectRoleID(uid string) (roleId uint64) {
	roleId = xycache.INVALID_ROLEID //不存在该职业
	var (
		err          error
		userRoleInfo *battery.UserRoleInfo
	)

	userRoleInfo, err = db.GetUserRoleInfo(uid)
	if err == xyerror.ErrOK {
		roleId = userRoleInfo.GetSelectRoleId()
	}

	return
}

//解锁玩家角色
// uid string 玩家id
// roleId uint64 角色id
func (db *BatteryDB) UnlockRole(uid string, roleId uint64) (err error) {
	var userRoleInfo *battery.UserRoleInfo
	userRoleInfo, err = db.GetUserRoleInfo(uid)
	if err == xyerror.ErrOK {
		for _, roleInfoItem := range userRoleInfo.GetRoleInfoItemList() {
			if roleInfoItem != nil && roleInfoItem.GetRoleId() == roleId && roleInfoItem.GetIsLock() && roleInfoItem.GetCurLevel() == 0 {
				*roleInfoItem.IsLock = false
				err = db.UpsertUserRoleInfo(userRoleInfo)
				break
			}
		}
	}
	return
}

//升级玩家角色等级
// uid string 玩家id
// roleId uint64 角色id
// level int32 升级后的等级
func (db *BatteryDB) UpgradeRole(uid string, roleId uint64, level int32) (err error) {
	var userRoleInfo *battery.UserRoleInfo
	userRoleInfo, err = db.GetUserRoleInfo(uid)
	if err == xyerror.ErrOK {
		for _, roleInfoItem := range userRoleInfo.GetRoleInfoItemList() {
			if roleInfoItem != nil && roleInfoItem.GetRoleId() == roleId && !roleInfoItem.GetIsLock() && roleInfoItem.GetCurLevel() == level-1 && db.GetRoleMaxLevel(roleId) >= level {
				*roleInfoItem.CurLevel = level
				err = db.UpsertUserRoleInfo(userRoleInfo)
				break
			}
		}
	}
	return
}

//查询玩家的角色信息
// uid string 玩家id
// userRoleInfo *battery.UserRoleInfo 角色信息指针
func (db *BatteryDB) GetUserRoleInfo(uid string) (userRoleInfo *battery.UserRoleInfo, err error) {
	userRoleInfo = &battery.UserRoleInfo{}
	condition := bson.M{"uid": uid}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetOneData(xybusiness.DB_TABLE_ROLEINFO, condition, selector, userRoleInfo, mgo.Strong)
	return
}

//查询好友的角色信息
// uids []string 好友uid列表
// maxCount int 返回的最大条数
//return:
// usersRoleInfo []*battery.UserRoleInfo 角色信息列表
func (db *BatteryDB) GetUserRoleInfos(uids []string, maxCount int) (userRoleInfos []*battery.UserRoleInfo, err error) {
	userRoleInfos = make([]*battery.UserRoleInfo, 0)
	condition := bson.M{"uid": bson.M{"$in": uids}}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	err = db.GetAllData(xybusiness.DB_TABLE_ROLEINFO, condition, selector, maxCount, &userRoleInfos, mgo.Monotonic)
	return
}

//判断某个角色配置信息有效，这个有效值将决定角色列表，角色等级配置表，和 角色等级商品表的数据哪些是有效的
//func (db *BatteryDB) IsRoleConfigValid(roleID uint64) (isValid bool) {
//	condition := bson.M{"roleid": roleID, "isvalid": true}
//	isValid, _ = db.IsRecordExisting(DB_TABLE_ROLEINFO_CONFIG, condition, mgo.Strong)
//	return isValid
//}

//获取角色的最大等级
func (db *BatteryDB) GetRoleMaxLevel(roleId uint64) (level int32) {
	level = xycache.INVALID_LEVEL
	roleInfo := xycache.DefRoleInfoCacheManager.Info(roleId)
	if nil != roleInfo {
		level = roleInfo.MaxLevel
	}
	return
}
