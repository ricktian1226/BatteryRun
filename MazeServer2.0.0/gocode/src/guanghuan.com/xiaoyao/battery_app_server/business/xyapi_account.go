package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	xyversion "guanghuan.com/xiaoyao/common/version"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	//xymoney "guanghuan.com/xiaoyao/superbman_server/money"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//玩家登录消息
func (api *XYAPI) OperationLogin(req *battery.LoginRequest, resp *battery.LoginResponse, addr string) (err error) {
	var (
		loginId       battery.TPID
		clientVersion xyversion.Version = DefVersion
		uid, errStr   string
		failReason    battery.ErrorCode
	)

	loginId = *req.GetLoginId()

	if req.Version != nil {
		clientVersion = xyversion.Version(req.GetVersion())
	}

	if clientVersion.LowerThan(api.Config.Configs().MinClientVersion) || clientVersion.LargerThan(api.Config.Configs().MaxClientVersion) {
		errStr = fmt.Sprintf("[%s] Client version is not support: got %s, expect [%s,%s]", loginId.GetId(), clientVersion.String(),
			api.Config.Configs().MinClientVersion.String(), api.Config.Configs().MaxClientVersion.String())
		xylog.ErrorNoId(errStr)
		resp.Error = xyerror.Resp_ClientNotSupport

		go api.GetLogDB().AddAccountLog(battery.AccountLog{
			Sid: loginId.Id,
			Error: &battery.Error{
				Code: battery.ErrorCode_ClientVersionNotSupport.Enum(),
				Desc: &errStr,
			},
			Timestamp:     proto.Int64(xyutil.CurTimeSec()),
			Source:        loginId.Source,
			Name:          loginId.Name,
			ClientVersion: req.Version,
			OpDateStr:     proto.String(xyutil.CurTimeStr()),
		})

		goto ErrorHandle
	}

	//查询玩家的
	//根据loginId查询玩家的tpid信息
	uid, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetGidBySid(loginId.GetId(), loginId.GetSource())
	if err != xyerror.ErrOK {
		xylog.WarningNoId("GetGidBySid(%s) failed : %v", loginId.GetId(), err)
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_Login, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrorHandle:
	return
}

//查询玩家钱包信息消息
func (api *XYAPI) OperationQueryWallet(req *battery.QueryWalletRequest, resp *battery.QueryWalletResponse) (err error) {

	var (
		uid        = req.GetUid()
		failReason battery.ErrorCode
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if !api.isUidValid(uid) {
		errStr := fmt.Sprintf("[%s] uid invalid", uid)
		xylog.ErrorNoId(errStr)
		resp.Error = xyerror.Resp_BadInputData
		resp.Error.Desc = proto.String(errStr)
		err = xyerror.ErrBadInputData
		goto ErrHandle
	}

	failReason, err = api.SendToTransaction(uid, xybusiness.BusinessCode_QueryWallet, req, resp)
	if failReason != battery.ErrorCode_NoError {
		resp.Error = xyerror.ConstructError(failReason)
	}

ErrHandle:
	return
}

//根据好友成就信息来获取好友数据
// dbAccount *battery.DBAccount 好友账户信息
//return:
// friendData *battery.FriendData 好友数据
func (api *XYAPI) getFriendDataFromAccomplishment(userAccomplishment *battery.DBUserAccomplishment) (friendData *battery.FriendData) {
	friendData = &battery.FriendData{
		Uid:          proto.String(userAccomplishment.GetUid()),
		Selectroleid: proto.Uint64(api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).GetSelectRoleID(userAccomplishment.GetUid())),
	}

	//好友玩家数据只要获取以下指标
	for i, accomplishment := range userAccomplishment.GetAccomplishment() {
		friendData.Accomplishment = append(friendData.Accomplishment, new(battery.QuotaList))
		if i == int(battery.AccomplishmentType_AccomplishmentType_CheckPoint_Total) { //只取checkPointTotalAccomplishment的信息
			for _, quotaList := range accomplishment.GetList() {
				for _, quota := range quotaList.GetItems() {
					id := quota.GetId()
					if battery.QuotaEnum_Quota_AllCheckPointScore == id || //所有记忆点的总分 1006
						battery.QuotaEnum_Quota_AllCheckPointCharge == id || //所有记忆点的charge总数 1005
						battery.QuotaEnum_Quota_AllCheckPointStar == id || //所有记忆点的星星总数 4009
						battery.QuotaEnum_Quota_FarthestCheckPoint == id { //到达的最远记忆点 5004
						friendData.Accomplishment[i].Items = append(friendData.Accomplishment[i].Items, quota)
					}
				}
			}
		}
	}

	return
}

//判断uid信息是否可用:1.非空，2.uid存在
// uid string 玩家id
// isUidValid bool 是否可用，true 可用，false 不可用
func (api *XYAPI) isUidValid(uid string) (isUidValid bool) {

	isUidValid = true

	if uid == "" || !api.IsUserExist(uid) {
		isUidValid = false
		xylog.Error(uid, "Invalid user")
	}

	return
}

//判断uid信息是否可用:1.非空，2.uid存在
// uid string 玩家id
// isUidValid bool 是否可用，true 可用，false 不可用
func (api *XYAPI) isIdentityValid(identity uint64) (isValid bool) {

	isValid = true

	if !api.IsIdentityExist(identity) {
		isValid = false
	}

	return
}

////查询玩家账户信息
//// uids string 玩家id
//// accounts *[]*battery.DBAccount 保存玩家信息的列表
//// consistency mgo.Mode 一致性模型
//func (api *XYAPI) GetDBAccountDirect(uid string, selector interface{}, account *battery.DBAccount, consistency mgo.Mode) (err error) {
//	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).GetAccountDirect(uid, selector, account, consistency)
//	return
//}

////查询多个玩家的账户信息
//// uids []string 玩家uid列表
//// accounts *[]*battery.DBAccount 保存玩家信息的列表
//// consistency mgo.Mode 一致性模型
//func (api *XYAPI) GetDBAccountsDirect(uids []string, accounts *[]*battery.DBAccount, consistency mgo.Mode) (err error) {
//	selector := bson.M{}
//	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).GetAccountsDirect(uids, selector, accounts, consistency)
//	return
//}

////根据玩家查询玩家账户信息
//// identity uint64 玩家identity
//// accounts *[]*battery.DBAccount 保存玩家信息的列表
//// consistency mgo.Mode 一致性模型
//func (api *XYAPI) GetDBAccountDirectByIdentity(identity uint64, selector interface{}, account *battery.DBAccount, consistency mgo.Mode) (err error) {
//	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).GetAccountDirectByIdentity(identity, selector, account, consistency)
//	return
//}

//判断玩家名称是否需要更新，如果需要更新，就替换到结构体中的名称字段
// account *battery.Account 玩家信息结构体
// tpid battery.TPID 玩家的tpid信息
func (api *XYAPI) ShouldUpdateAccountName(account *battery.Account, tpid battery.TPID) (err error) {
	name := account.GetName()
	tpname := tpid.GetName()
	if tpname != "" && name != tpname {
		account.Name = proto.String(tpname)
	}

	return
}

//判断玩家是否存在
// uid string 玩家id
//return:
// true 玩家存在，false 玩家不存在
func (api *XYAPI) IsUserExist(uid string) bool {
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).IsUserExist(uid)
}

//判断玩家是否存在
// identity uint64 玩家identity
//return:
// true 玩家存在，false 玩家不存在
func (api *XYAPI) IsIdentityExist(identity uint64) bool {
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).IsIdentityExist(identity)
}

//根据identity查询玩家的uid
// identity uint64 玩家的uid
//return:
// uid string 玩家的uid标识
// err error 操作错误信息
func (api *XYAPI) GetUid(identity string) (uid string, err error) {
	//查询玩家uid
	account := &battery.DBAccount{}
	selector := bson.M{"uid": 1}
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).GetAccountDirectByIdentity(identity, selector, account, mgo.Strong)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("GetAccountDirectByIdentity(%d) failed : %v ", identity, err)
		return
	}

	uid = account.GetUid()

	return
}
