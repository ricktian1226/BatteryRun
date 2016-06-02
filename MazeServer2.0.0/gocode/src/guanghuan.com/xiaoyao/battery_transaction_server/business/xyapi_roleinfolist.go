package batteryapi

// xyapi_rolelist

//购买角色(返回购买确认信息（成功则扣除相关消耗：碎片、金币、宝石），失败则提示错误信息)

//返回游戏结算结果(在游戏结算中应该有了加进去就行)

import (
	"code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xycache "guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	xymoney "guanghuan.com/xiaoyao/superbman_server/money"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"time"
)

const (
	UNLOCK_ROLE_NEED_ADD_JIGSAW    = true  //解锁角色增加拼图
	UNLOCK_ROLE_NO_NEED_ADD_JIGSAW = false //解锁角色不增加拼图
)

//查询玩家角色列表消息接口
func (api *XYAPI) OperationRoleInfoList(req *battery.RoleInfoListRequest, resp *battery.RoleInfoListResponse) (err error) {

	var (
		uid       string = req.GetUid()
		cmd              = req.GetCmd()
		errStruct        = xyerror.DefaultError()
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	errStruct.Code = battery.ErrorCode_NoError.Enum()

	//xylog.Debug(uid, "[%s] [RoleInfoListRequest] cmd %d start", uid, cmd)
	//defer xylog.Debug(uid, "[%s] [RoleInfoListRequest] cmd %d done", uid, cmd)

	//初始化返回信息
	resp.Uid = req.Uid
	resp.Cmd = req.Cmd
	resp.RoleId = req.RoleId

	switch cmd {
	//查询玩家角色列表
	case battery.RoleInfoListCmd_RoleInfoListCmd_RoleList:
		var userRoleInfo *battery.UserRoleInfo
		userRoleInfo, err = api.getUserRoleInfo(uid)
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_QueryUserRoleInfoError.Enum()
			goto ErrHandle
		}

		userRoleInfo.Uid = nil //外面已经有uid了，减少消息内容，这个uid置为nil
		resp.RoleInfos = append(resp.RoleInfos, userRoleInfo)

	//查询好友角色列表
	case battery.RoleInfoListCmd_RoleInfoListCmd_FriendRoleList:
		resp.RoleInfos, err = api.getFriendRoleInfos(uid, req.GetFriendSids(), req.GetSource())
		if err == xyerror.ErrNotFound || len(resp.RoleInfos) <= 0 {
			errStruct.Code = battery.ErrorCode_NotUsefulSidError.Enum()
			goto ErrHandle
		} else if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_QueryUserRoleInfoError.Enum()
			goto ErrHandle
		}

	//设置玩家选中角色
	case battery.RoleInfoListCmd_RoleInfoListCmd_Select:
		err = api.setUserSelectRole(uid, req.GetRoleId())
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_SetUserSelectRoleError.Enum()
			goto ErrHandle
		}

	//升级/购买玩家角色
	case battery.RoleInfoListCmd_RoleInfoListCmd_Upgrade:
		api.upgradeUserRole(uid, req.GetRoleId(), &(resp.Wallet), errStruct)
		if errStruct.GetCode() != battery.ErrorCode_NoError {
			goto ErrHandle
		}
	}

ErrHandle:

	resp.Error = errStruct
	resp.Error.Desc = nil
	//xylog.Debug("[%s]  resp : %v", uid, resp)
	return
}

//设置玩家当前选中角色
func (api *XYAPI) setUserSelectRole(uid string, roleId uint64) (err error) {
	if nil == xycache.DefRoleInfoCacheManager.Info(roleId) {
		xylog.Error(uid, "invalid roleId(%d)", roleId)
		err = xyerror.ErrQueryRoleInfoConfigFromCacheError
		return
	}

	err = api.UpdateSelectRoleID(uid, roleId)

	return
}

//升级玩家角色
// uid string 玩家id
// roleLevelId uint64 升级后的角色id
// wallet *[]*battery.MoneyItem 钱包列表指针
func (api *XYAPI) upgradeUserRole(uid string, roleLevelId uint64, wallet *[]*battery.MoneyItem, errStruct *battery.Error) {
	var (
		roleId    = roleLevelId / 10000 * 10000
		nextLevel = int32(roleLevelId % 10000)
		curLevel  = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).GetRoleLevel(uid, roleId) //获取角色的当前等级
	)

	xylog.Debug(uid, "userRole(%d), upgrade from %d to %d", roleId, curLevel, nextLevel)
	//errStruct = new(battery.Error)
	//*errStruct = *(xyerror.Resp_NoError)

	if curLevel >= 0 && curLevel == (nextLevel-1) { //升级后的级别必须是当且级别+1
		xylog.Debug(uid, " userRole(%d) upgrade from(%v) to (%d) ready", roleId, curLevel, nextLevel)
		api.BuyGoods(uid, "", roleLevelId, errStruct)         //商品id就是角色等级id，bad idea!
		if errStruct.GetCode() == battery.ErrorCode_NoError { //返回钱包信息
			*wallet, _ = xymoney.QueryWallet(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT))
		}
	} else {
		errStruct.Code = battery.ErrorCode_UpgradeUserRoleError.Enum() // 返回升级错误  change at 2015/9/9
		xylog.Warning(uid, "level is not right,pls check.")
	}

	return
}

//玩家角色道具发放
// uid string 玩家id
// roleID uint64 角色id
// accountWithFlag *AccountWithFlag 玩家账户信息
// needAddJigsaw bool 是否需要增加拼图
func (api *XYAPI) updateUserDataWithRole(uid string, accountWithFlag *AccountWithFlag, id uint64, delay bool, needAddJigsaw bool) (err error) {
	xylog.Debug(uid, "updateUserDataWithRole %d", id)
	//非已经拥有则，升级或者解锁
	var (
		roleLevel    = int32(id % 10000)
		roleId       = id / 10000 * 10000
		curRoleLevel int32
		isLock       bool = false
		userRoleInfo *battery.UserRoleInfo
		item         *battery.RoleInfoItem //临时变量，指向待更新的roleInfoItem
	)

	//查询玩家对应角色的信息
	userRoleInfo, err = api.getUserRoleInfo(uid)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "GetUserRoleInfo failed : %v", err)
		return
	}
	isFound := false
	for _, roleInfoItem := range userRoleInfo.GetRoleInfoItemList() {
		if roleInfoItem.GetRoleId() == roleId {
			item = roleInfoItem //找到，赋值
			curRoleLevel = roleInfoItem.GetCurLevel()
			isLock = roleInfoItem.GetIsLock()
			isFound = true
		}
	}
	if !isFound { //没找到对应的角色信息，返回错误
		err = xyerror.ErrQueryUserRoleInfoError
		xylog.Error(uid, "no found roleId(%d) in userRoleInfo : %v", roleId, userRoleInfo)
		return
	}

	if (roleLevel == 0 && curRoleLevel == 0) || //购买角色或者碎片合成角色
		(roleLevel == 1 && curRoleLevel == 0) { //购买level 1 角色
		if isLock {
			err = api.unlockRole(uid, accountWithFlag, roleId, item, needAddJigsaw, delay)
		}
		if err != xyerror.ErrOK {
			return
		}
		xylog.Debug(uid, "set role(%d) to level 1 ", roleId)
		item.CurLevel = proto.Int32(1)

		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).UpsertUserRoleInfo(userRoleInfo)
	} else if roleLevel >= 2 && curRoleLevel == roleLevel-1 && !isLock { //升级
		xylog.Debug(uid, "upgrade role(%d) from %d to %d", roleId, curRoleLevel, roleLevel)
		item.CurLevel = proto.Int32(roleLevel)
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).UpsertUserRoleInfo(userRoleInfo)
	} else if roleLevel <= curRoleLevel && curRoleLevel >= 1 { //重复角色分解
		xylog.Debug(uid, "role(%d) curLevel")
		var propItem battery.PropItem
		propItem.Id = &roleId
		propItem.Amount = proto.Uint32(1)
		err = api.ResolveProp(uid, accountWithFlag, &propItem, delay)
	}
	return
}

//解锁玩家角色
// roleId uint64 角色id
// item *battery.RoleInfoItem 角色信息指针
func (api *XYAPI) unlockRole(uid string, accountWithFlag *AccountWithFlag, roleId uint64, item *battery.RoleInfoItem, needAddJigsaw, delay bool) (err error) {
	xylog.Debug(uid, "unlock role(%d)", roleId)
	item.IsLock = proto.Bool(false) //角色解锁
	if needAddJigsaw {              //需要增加拼图
		roleInfo := xycache.DefRoleInfoCacheManager.Info(roleId)
		if nil != roleInfo && roleInfo.JigsawId > 0 { //发放拼图
			xylog.Debug(uid, "unlockRole addJigsaw(%d)", roleInfo.JigsawId)
			err = api.AddJigsaw(uid, accountWithFlag, roleInfo.JigsawId, JIGSAW_NO_NEED_UNLOCK_PROP, delay)
		}
	} else {
		xylog.Debug(uid, "unlockRole no need to addJigsaw")
	}
	return
}

//增加玩家角色信息
// uid string 玩家id
func (api *XYAPI) addDefaultUserRoleInfo(uid string) (userRoleInfo *battery.UserRoleInfo, err error) {

	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_ADDROLEINFO, &begin)

	userRoleInfo = api.defaultUserRoleInfo(uid)
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).UpsertUserRoleInfo(userRoleInfo)
	return
}

//获取玩家角色信息
//如果玩家是第一次登录则往数据库里创建新的动态数据
// uid string 玩家id
// roleInfoItems *[]*battery.RoleInfoItem 角色等级信息指针
// selectId *uint64 玩家当前使用的角色id指针
func (api *XYAPI) getUserRoleInfo(uid string) (userRoleInfo *battery.UserRoleInfo, err error) {
	userRoleInfo, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).GetUserRoleInfo(uid)
	if err == xyerror.ErrNotFound { //如果没找到，插入一个。 用于修，实际上线后不会有该流程
		xylog.Warning(uid, "[%s] UserRoleInfo not found ,add one", uid)
		userRoleInfo, err = api.addDefaultUserRoleInfo(uid)
	}
	return
}

//生成默认的角色信息，在玩家第一次登录的时候，就会为其生成默认的角色信息
// uid string 玩家id
//return:
// userRoleInfo *battery.UserRoleInfo 角色信息指针
func (api *XYAPI) defaultUserRoleInfo(uid string) (userRoleInfo *battery.UserRoleInfo) {

	userRoleInfo = &battery.UserRoleInfo{
		Uid: &uid,
		//SelectRoleId: &defaultRoleID,
	}

	defaultRoleID := xycache.DefRoleInfoCacheManager.DefaultOwn()
	if defaultRoleID <= 0 {
		xylog.Warning(uid, "[%s] GetDeaultOwnRoleID Error : %d ", uid, defaultRoleID)
	} else {
		userRoleInfo.SelectRoleId = &defaultRoleID
	}

	mapRoleInfos := xycache.DefRoleInfoCacheManager.Infos()
	for id, roleInfo := range *mapRoleInfos {
		if roleInfo == nil {
			continue
		}
		roleInfoItem := &battery.RoleInfoItem{
			RoleId: proto.Uint64(id),
		}
		if defaultRoleID == id { //默认角色，设置等级为1，解锁
			roleInfoItem.CurLevel = proto.Int32(1)
			roleInfoItem.IsLock = proto.Bool(false)
		} else { //非默认角色，，设置等级为0，未解锁
			roleInfoItem.CurLevel = proto.Int32(0)
			roleInfoItem.IsLock = proto.Bool(true)
		}
		userRoleInfo.RoleInfoItemList = append(userRoleInfo.RoleInfoItemList, roleInfoItem)
	}

	return
}

//查询好友的角色信息
// sid string 玩家的第三方id
// source battery.ID_SOURCE 第三方账户来源
// roleInfoItems *[]*battery.RoleInfoItem 角色等级信息指针
// selectId *uint64 玩家当前使用的角色id指针
func (api *XYAPI) getFriendRoleInfos(uid string, sids []string, source battery.ID_SOURCE) (userRoleInfos []*battery.UserRoleInfo, err error) {

	//查询玩家好友的uid列表
	var (
		uids       []string
		mapUid2Sid map[string]string
	)
	uids, mapUid2Sid, err = api.GetUidsFromSids(uid, sids, source)
	if err != xyerror.ErrOK {
		return
	}
	xylog.Debug(uid, "uids : %v", uids)
	userRoleInfos, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).GetUserRoleInfos(uids, DefConfigCache.Configs().MaxFriendsRequestCount)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "GetUserRoleInfos of %v failed : %v", uids, err)
		return
	}

	//替换uid为sid，前台需要的sid
	for _, userRoleInfo := range userRoleInfos {
		uid := userRoleInfo.GetUid()
		if sid, ok := mapUid2Sid[uid]; ok {
			userRoleInfo.Uid = proto.String(sid)
		}
	}

	return
}

//判断玩家是否拥有角色
func (api *XYAPI) IsRoleExisting(uid string, roleId uint64) bool {
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).IsRoleExisting(uid, roleId)
}

//刷新玩家当前使用的角色id
// uid string 玩家id
// roleID uint64 角色id
func (api *XYAPI) UpdateSelectRoleID(uid string, roleID uint64) (err error) {
	var userRoleInfo *battery.UserRoleInfo
	userRoleInfo, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).GetUserRoleInfo(uid)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "[DB] fail to GetRoleInfoTableItem: %v ", err)
		return
	} else {
		for _, roleInfoItem := range userRoleInfo.RoleInfoItemList {
			if roleInfoItem != nil && roleInfoItem.GetRoleId() == roleID && roleInfoItem.GetCurLevel() > 0 && !roleInfoItem.GetIsLock() {
				err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).UpdateSelectRoleID(uid, roleID)
				return
			}
		}
	}
	return
}
