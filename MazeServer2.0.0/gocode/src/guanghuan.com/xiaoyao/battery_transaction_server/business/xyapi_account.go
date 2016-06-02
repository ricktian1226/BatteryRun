package batteryapi

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"guanghuan.com/xiaoyao/common/idgenerate"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	xymoney "guanghuan.com/xiaoyao/superbman_server/money"
	xybusiness "guanghuan.com/xiaoyao/superbman_server/server"
)

const (
	UserVerifyFail = "0"  // 用户验证失败
	SignVerifyFail = "-1" // 签名验证失败
	URLOverTime    = "-2" // url失效
)

//玩家登录请求
//保存玩家账户信息和状态的结构体
type AccountWithFlag struct {
	account *battery.DBAccount
	bChange bool
}

func (accountWithFlag *AccountWithFlag) SetChange() {
	accountWithFlag.bChange = true
}

func (accountWithFlag *AccountWithFlag) Print() {
	xylog.Debug(accountWithFlag.account.GetUid(), "account : %v, change : %t", accountWithFlag.account, accountWithFlag.bChange)
}

//是否延迟刷新玩家账户数据
const ACCOUNT_UPDATE_DELAY = true
const ACCOUNT_UPDATE_NO_DELAY = false

var DefIdentityIdGenerater *xyidgenerate.IdGenerater

const (
	SDKLoginKey_Token = "token" // sdk验证字符串
	SDKLoginKey_Time  = "time"  // sdk验证时间
)

func (api *XYAPI) OperationLogin(req *battery.LoginRequest, resp *battery.LoginResponse, addr string) (err error) {

	var (
		devid, errStr   string
		loginId         = req.GetLoginId()
		account         *battery.DBAccount
		isNew           = false
		op              int32
		iconUrl         = req.GetIconUrl()
		now             = xyutil.CurTimeSec()
		idMap           *battery.IDMap
		accomplishments = make([]*battery.Accomplishment, 0)
		loginType       battery.LoginType
		errStruct       = xyerror.DefaultError()
		loginIp         = req.GetLoginIp()
		statisticsDevid = req.GetStatisticsDeviceId()

		sdkPlatform        = req.GetSdkPlatform() // sdk 接入第三方渠道
		token              = req.GetToken()
		name        string = loginId.GetName()
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	if req.DeviceId != nil {
		devid = req.GetDeviceId().GetId()
	}

	// 接入平台 旧版本接入平台为空，不做处理
	if sdkPlatform != "" {
		keys, values := make([]string, 0), make(map[string]string, 0)

		values[SDKLoginKey_Time] = xyutil.CurTimeStr()
		keys = append(keys, SDKLoginKey_Time)

		values[SDKLoginKey_Token] = token
		keys = append(keys, SDKLoginKey_Token)

		sort.Strings(keys)
		xylog.DebugNoId("sdklogin verifykey :%v", keys)

		var encodeStr string
		for i, k := range keys {
			if i == 0 {
				encodeStr += k + "=" + values[k]
			} else {
				encodeStr += "&" + k + "=" + values[k]
			}
		}
		urlEncode := url.QueryEscape(encodeStr)
		urlEncode += fmt.Sprintf("&%s", DefConfigCache.Configs().Appkey)
		xylog.DebugNoId("urlencode before md5 :%v", urlEncode)
		h := md5.New()
		io.WriteString(h, urlEncode)
		sign := fmt.Sprintf("%x", h.Sum(nil))
		sign = strings.ToLower(sign)
		encodeStr += "&sign=" + sign
		requrl := fmt.Sprintf("%s%s/user_check?%v", DefConfigCache.Configs().SDKLoginUrl, sdkPlatform, encodeStr)
		xylog.DebugNoId("requrl :%v", requrl)
		resp, err := http.Get(requrl)
		if err != nil {
			xylog.ErrorNoId("query sdk fail :%v", err)
			errStruct.Code = battery.ErrorCode_SDKQueryURLFail.Enum()
			goto ErrorHandle
		}
		xylog.DebugNoId("resp : %v", resp)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		uuid := string(body)
		xylog.DebugNoId("resp body :%v", string(body))

		// 非人性化接口，错误码和内容共用字段！！
		switch string(uuid) {
		case UserVerifyFail:
			fallthrough
		case SignVerifyFail:
			fallthrough
		case URLOverTime:
			errStruct.Code = battery.ErrorCode_SDKLoginVerifyFail.Enum()
			goto ErrorHandle
		default:
			if uuid != loginId.GetId() {
				errStruct.Code = battery.ErrorCode_RegistTpidError.Enum()
				goto ErrorHandle
			}
		}

	}

	//根据loginId查询玩家的tpid信息
	idMap, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetIdMapByTpid(loginId)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound { //数据库报不存在，说明是新玩家
			isNew = true
		} else { //数据库错误
			errStr = fmt.Sprintf("[%s] GetIdMapByTpid failed : %v", loginId.GetId(), err)
			xylog.ErrorNoId(errStr)
			errStruct.Code = battery.ErrorCode_QueryTpidError.Enum()
			goto ErrorHandle
		}
	}
	if idMap.Note != nil {
		name = idMap.GetNote()
	}

	xylog.DebugNoId("idmap:%v   ", idMap)
	if !isNew { //老玩家登录
		xylog.Debug(idMap.GetGid(), "%v : %s -> %s exists", idMap.GetSource(), idMap.GetSid(), idMap.GetGid())
		op = ACCOUNT_OP_LOGIN
		account, accomplishments = api.oldGuyLogin(idMap, loginId, loginIp, devid, statisticsDevid, iconUrl, now, errStruct, &loginType)
		xylog.DebugNoId("idmap", idMap.Note)

	} else { //新玩家登录
		xylog.DebugNoId("%v : %s no exists, will create", idMap.GetSource(), idMap.GetSid())
		op = ACCOUNT_OP_NEW
		account = api.newGuyLogin(loginId, loginIp, devid, statisticsDevid, iconUrl, now, errStruct)
		loginType = battery.LoginType_LoginType_New

	}

ErrorHandle:

	resp.Error = errStruct

	if errStruct.GetCode() == battery.ErrorCode_NoError {
		//resp.Uid = proto.String(account.GetUid())
		resp.Data = api.getAccountFromDBAccount(account)
		resp.Data.Accomplishment = accomplishments
		resp.Data.Name = proto.String(name)
		resp.ServerTime = proto.Int64(now)
		resp.LoginType = loginType.Enum()

		// ToDelete:
		////登陆成功
		////判断是否有当天的登陆礼包邮件
		////有弹出公告 并发送礼包邮件到系统邮箱
		//api.AddLoginGiftMail(account.GetUid())

		////通知apns，使能devicetoken
		//notification := &xyapn.APNNotification{
		//	DeviceToken: account.Deviceid,
		//	Cmd:         xyapn.APNNotificationCMD_EnableDeviceToken.Enum(),
		//	Uid:         resp.Uid,
		//}
		//go xyapn.Send(notification)

		err = xyerror.ErrOK

	} else {
		resp.Data = nil
		resp.Uid = proto.String("0")
		resp.ServerTime = proto.Int64(0)
	}

	go api.AddAccountLog(addr, account.GetUid(), loginId.GetId(), loginId.GetSource(), loginId.GetName(), op, req.GetVersion(), resp.Error.GetCode(), errStr)

	return
}

//老伙计登录
// loginId *battery.TPID 第三方账户信息
// devid string 设备信息
//return:
// errStruct *battery.Error 处理结果信息
// account *battery.DBAccount 玩家账户信息
//func (api *XYAPI) oldGuyLogin(loginId *battery.TPID, devid, iconUrl string) (errStruct *battery.Error, account *battery.DBAccount) {
func (api *XYAPI) oldGuyLogin(idMap *battery.IDMap, loginId *battery.TPID, loginIp, devid, statisticsDevid, iconUrl string, now int64, errStruct *battery.Error, loginType *battery.LoginType) (account *battery.DBAccount, accomplishments []*battery.Accomplishment) {
	var (
		errStr             string
		err                error
		accountWithFlag    AccountWithFlag
		userAccomplishment = &battery.DBUserAccomplishment{}
		//pgid               = idMap.GetPgid()
		uid = idMap.GetGid()
	)

	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_OLDLOGIN, &begin)

	errStruct.Code = battery.ErrorCode_NoError.Enum()

	//获取玩家账户信息详情
	account = new(battery.DBAccount)
	accountWithFlag.account = account
	accountWithFlag.bChange = false
	err = api.GetDBAccount(uid, &accountWithFlag)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound { //没找到，说明是已有账户，但是在该平台下没有账户,创建账户信息
			account = api.addNewGuyPlatformInfo(loginId, loginIp, devid, statisticsDevid, uid, now, errStruct)
			*loginType = battery.LoginType_LoginType_New_Platform
			goto ErrHandle
		}

		errStr = fmt.Sprintf("[%s] GetDBAccount failed: %s", uid, err.Error())
		xylog.Error(uid, errStr)
		errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
		goto ErrHandle
	}

	//获取玩家的成就信息
	api.getUserAccomplishment(uid, userAccomplishment, errStruct)
	if errStruct.GetCode() != battery.ErrorCode_NoError {
		goto ErrHandle
	} else {
		accomplishments = userAccomplishment.GetAccomplishment()
	}

	if devid != "" && account.GetDeviceid() != devid {
		account.Deviceid = proto.String(devid)
		accountWithFlag.SetChange()
	}
	if statisticsDevid != "" && account.GetStatisticsDeviceid() != statisticsDevid {
		account.StatisticsDeviceid = proto.String(statisticsDevid)
		accountWithFlag.SetChange()
	}
	if loginIp != "" && account.GetLastLoginIp() != loginIp {
		account.LastLoginIp = proto.String(loginIp)
		accountWithFlag.SetChange()
	}
	//判断是否需要刷新tpid信息，需要的话就刷新
	if api.ShouldUpdateTpid(loginId, iconUrl, idMap) {
		err = api.UpsertTPID(idMap)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "UpsertTPID failed : %v", err)
			errStruct.Code = battery.ErrorCode_UpsertTpidError.Enum()
			goto ErrHandle
		}
	}

	////向下兼容，如果没有identity的就补上一个
	//if 0 == account.GetIdentity() {
	//	account.Identity = proto.Uint64(DefIdentityIdGenerater.NewID())
	//	accountWithFlag.SetChange()
	//}

	//调用玩家账户信息刷新函数，需要刷新的话就刷新
	err = api.UpdateAccountWithFlag(&accountWithFlag)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_UpdateAccountError.Enum()
		goto ErrHandle
	}

	*loginType = battery.LoginType_LoginType_Old

ErrHandle:

	return
}

//更新玩家账户信息
func (api *XYAPI) UpdateAccountWithFlag(accountWithFlag *AccountWithFlag) (err error) {
	//begin := time.Now()
	//defer xyperf.Trace(LOGTRACE_UPDATEDBACCOUNT, &begin)

	if accountWithFlag.bChange {
		//更新玩家账户信息
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).UpdateAccount(accountWithFlag.account)
		if err != xyerror.ErrOK {
			xylog.Error(accountWithFlag.account.GetUid(), "UpdateAccount failed: %v", err)
			return
		}
	}
	return
}

//新玩家登录，针对特定平台
// 这种登录类型，只需增加玩家账户信息
// loginId *battery.TPID 第三方账户信息
// devid string 设备信息
//return:
// errStruct *battery.Error 处理结果信息
// account *battery.DBAccount 玩家账户信息
func (api *XYAPI) addNewGuyPlatformInfo(loginId *battery.TPID, loginIp, devid, statisticsDevid, uid string, now int64, errStruct *battery.Error) (account *battery.DBAccount) {
	//增加玩家的账户信息
	var (
		err  error
		info *battery.SysLottoInfo
	)
	account, err = api.AddDBAccount(devid, statisticsDevid, loginIp, uid, loginId)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("[%s] adding new account failed: %v", loginId.GetId(), err)
		errStruct = xyerror.Resp_UpdateAccountError
		goto ErrHandle
	}

	//增加玩家的系统抽奖信息
	err = api.addNewSysLottoInfo(uid, &info, now)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "addNewSysLottoInfo failed: %v", err)
		errStruct.Code = battery.ErrorCode_AddUserSysLottoInfoError.Enum()
		goto ErrHandle
	}

	//增加玩家的角色背包
	_, err = api.addDefaultUserRoleInfo(uid)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "addDefaultUserRoleInfo failed: %v", err)
		errStruct.Code = battery.ErrorCode_AddUserRoleError.Enum()
		goto ErrHandle
	}

	//增加拼图信息
	// todo: 增加一个可配新账户拥有道具配置。以后有新账户需要拥有的道具，可以在该配置信息中配置。
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_JIGSAW).AddJigsaw(uid, JIGSAW_FOR_NEWACCOUNT)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "AddJigsaw failed: %v", err)
		errStruct.Code = battery.ErrorCode_AddJigsawError.Enum()
		goto ErrHandle
	}

ErrHandle:
	return
}

//默认的源账户gid为空字符串
const DEFAULT_PGID = ""

//全新玩家登录，增加平台互通数据信息
// loginId *battery.TPID 玩家第三方账户信息
// now int64 当前时间戳
// errStruct *battery.Error 业务错误信息
//return:
// uid string 玩家账户信息
func (api *XYAPI) addNewGuyCommonInfo(loginId *battery.TPID, now int64, iconUrl string, errStruct *battery.Error) (uid string) {
	//增加玩家的第三方账户信息
	// 通过登录的注册方式，其源
	uid, err := api.RegisterTPID(loginId, DEFAULT_PGID, iconUrl)
	if err != xyerror.DBErrOK {
		xylog.ErrorNoId("[%s] RegisterTPID:%v", loginId.GetId(), err)
		errStruct.Code = battery.ErrorCode_RegistTpidError.Enum()
		return
	}

	api.addNewUserAccomplishment(uid, now, errStruct)
	if errStruct.GetCode() != battery.ErrorCode_NoError {
		xylog.ErrorNoId("[%s] addNewUserAccomplishment:%v", loginId.GetId(), err)
		return
	}

	return
}

//增加新玩家的成就信息
func (api *XYAPI) addNewUserAccomplishment(uid string, now int64, errStruct *battery.Error) {

	userAccomplishment := api.defaultDBUserAccomplishment(uid, now)

	err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERACCOMPLISHMENT).AddUserAccomplishment(userAccomplishment)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_AddNewUserAccomplishmentError.Enum()
		return
	}

	return
}

//获取玩家的成就信息
// uid string 玩家标识
// userAccomplishment *battery.DBUserAccomplishment 玩家成就信息指针
//return
//errStruct *battery.Error 错误信息
func (api *XYAPI) getUserAccomplishment(uid string, userAccomplishment *battery.DBUserAccomplishment, errStruct *battery.Error) {
	err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERACCOMPLISHMENT).GetUserAccomplishment(uid, userAccomplishment, mgo.Strong)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "GetUserAccomplishment failed : %v", err)
		errStruct.Code = battery.ErrorCode_GetUserAccomplishmentError.Enum()
		return
	}

	return
}

//新玩家登录，全新的玩家登录
// 这种登录类型，需要增加玩家第三方账户信息、玩家账户信息
// loginId *battery.TPID 第三方账户信息
// devid string 设备信息
//return:
// errStruct *battery.Error 处理结果信息
// account *battery.DBAccount 玩家账户信息
func (api *XYAPI) newGuyLogin(loginId *battery.TPID, loginIp, devid, statisticsDevid, iconUrl string, now int64, errStruct *battery.Error) (account *battery.DBAccount) {

	xylog.DebugNoId("[%s] a new account %v", loginId.GetId(), loginId)

	//增加玩家的平台互通信息
	uid := api.addNewGuyCommonInfo(loginId, now, iconUrl, errStruct)
	if errStruct.GetCode() != battery.ErrorCode_NoError {
		goto ErrHandle
	}

	account = api.addNewGuyPlatformInfo(loginId, loginIp, devid, statisticsDevid, uid, now, errStruct)

	//发放登录礼包
	api.newAccountGainProps(uid, loginId.GetSource(), errStruct)

ErrHandle:

	return
}

//查询玩家钱包信息
func (api *XYAPI) OperationQueryWallet(req *battery.QueryWalletRequest, resp *battery.QueryWalletResponse) (err error) {

	var (
		uid             = req.GetUid()
		timeleft        = int32(-1)
		errStruct       = xyerror.DefaultError()
		accountWithFlag *AccountWithFlag
	)

	//设置玩家终端的平台类型
	api.SetDB(req.GetPlatformType())

	//初始化resp
	resp.Uid = proto.String(uid)

	//查询刷新体力前的玩家账户信息
	account := &battery.DBAccount{}
	err = api.GetDBAccountDirect(uid, account, mgo.Strong)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
		goto ErrHandle
	}

	//刷新玩家体力
	accountWithFlag = &AccountWithFlag{
		account: account,
		bChange: false,
	}
	_, timeleft, err = api.getCurrentStamina(uid, accountWithFlag, 0)
	if err == xyerror.ErrOK {
		resp.Timeleft = proto.Int32(timeleft)
	} else {
		errStruct.Code = battery.ErrorCode_QueryStaminaError.Enum()
		goto ErrHandle
	}

	//调用玩家账户信息刷新函数，需要刷新的话就刷新
	err = api.UpdateAccountWithFlag(accountWithFlag)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_UpdateAccountError.Enum()
		goto ErrHandle
	}

	//返回最新的钱包信息
	resp.Wallet = xymoney.GetMoneyItemsFromMoneys(account.Wallet)

ErrHandle:
	resp.Error = errStruct

	return
}

// 查询玩家的数据 (附带刷新当前的体力值)
func (api *XYAPI) GetDBAccount(uid string, accountWithFlag *AccountWithFlag) (err error) {

	err = api.GetDBAccountDirect(uid, accountWithFlag.account, mgo.Strong)
	if err == xyerror.ErrOK {
		//刷新一下体力值
		_, _, err = api.getCurrentStamina(uid, accountWithFlag, 0)
	}

	xylog.Debug(uid, "[%s] GetDBAccount err : %v", uid, err)

	return
}

// 查询玩家的数据（accountWithFlag可能为空）
// uid string 玩家标识
// accountWithFlag **AccountWithFlag 玩家账户结构体指针的指针
func (api *XYAPI) fillAccountWithFlag(uid string, accountWithFlag **AccountWithFlag) (err error) {
	if nil == *accountWithFlag { //为空的话，就查一下
		account := new(battery.DBAccount)
		err = api.GetDBAccountDirect(uid, account, mgo.Strong)
		if err == xyerror.ErrOK {
			*accountWithFlag = &AccountWithFlag{
				account: account,
				bChange: false,
			}
		} else {
			xylog.Error(uid, "GetDBAccountDirect failed : %v", err)
		}
	}

	return
}

//根据tpid获取玩家id
// tpid battery.TPID 第三方平台玩家标识
//func (api *XYAPI) GetUid(tpid *battery.TPID, dbIndex int) (uid string, err error) {
//	uid, err = api.GetDB(dbIndex).GetGidByTpid(tpid)
//	return
//}

//根据第三方id获取uid
// sids []string 第三方id列表
//return:
// uids []string uid列表
// mapUid2Sid map[string]string uid->sid的映射列表
func (api *XYAPI) GetUidsFromSids(uid string, sids []string, source battery.ID_SOURCE) (uids []string, mapUid2Sid map[string]string, err error) {
	tpids := make([]*battery.IDMap, 0)
	selector := bson.M{"gid": 1, "sid": 1}
	err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).XYBusinessDB.QueryTpidsBySids(selector, sids, source, &tpids)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "QueryFriendsTpid for %v failed : %v", sids, err)
		return
	}

	uids = make([]string, 0)
	mapUid2Sid = make(map[string]string, 0)
	for _, tpid := range tpids {
		gid := tpid.GetGid()
		sid := tpid.GetSid()
		if len(gid) > 0 {
			mapUid2Sid[gid] = sid
			uids = append(uids, gid)
		}
	}
	return
}

//直接查询玩家的数据（不刷新玩家体力）
// uid string 玩家id
// account *battery.DBAccount 返回的玩家信息
func (api *XYAPI) GetDBAccountDirect(uid string, account *battery.DBAccount, consistency mgo.Mode) (err error) {
	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_GETDBACCOUNT, &begin)
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).GetAccountDirect(uid, selector, account, consistency)
}

//直接查询玩家的devicetoken
// uid string 玩家id
// account *battery.DBAccount 返回的玩家信息
func (api *XYAPI) GetDeviceToken(uid string, consistency mgo.Mode) (deviceToken, indentityString string, err error) {
	account := &battery.DBAccount{}
	selector := bson.M{"deviceid": 1, "identitystring": 1}
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).GetAccountDirect(uid, selector, account, consistency)
	if err == xyerror.ErrOK {
		deviceToken = account.GetDeviceid()
		indentityString = account.GetIdentityString()
	}

	return
}

func (api *XYAPI) GetStatisticsDevid(uid string, consistency mgo.Mode) (deviceToken, indentityString string, err error) {
	account := &battery.DBAccount{}
	selector := bson.M{"statisticsdeviceid": 1, "identitystring": 1}
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).GetAccountDirect(uid, selector, account, consistency)
	if err == xyerror.ErrOK {
		deviceToken = account.GetStatisticsDeviceid()
		indentityString = account.GetIdentityString()
	}
	return
}

//通过数据库账户信息获取消息用账户信息
// dbAccount *battery.DBAccount 数据库中的账户信息
//return:
// account *battery.Account 返回前台的账户信息
func (api *XYAPI) getAccountFromDBAccount(dbAccount *battery.DBAccount) (account *battery.Account) {
	account = &battery.Account{
		Uid:      proto.String(dbAccount.GetUid()),
		Deviceid: proto.String(dbAccount.GetDeviceid()),
		//Identity: proto.Uint64(dbAccount.GetIdentity()),
		IdentityString: proto.String(dbAccount.GetIdentityString()),
	}

	for _, money := range dbAccount.GetWallet() {
		moneyItem := &battery.MoneyItem{}
		moneyItem.Type = money.GetType().Enum()
		moneyItem.Amount = proto.Uint32(money.GetIapamount() + money.GetOapamount() + money.GetGainamount())
		account.Wallet = append(account.Wallet, moneyItem)
	}
	return
}

//增加玩家账户数据
// devid string 玩家设备token
// uid string 玩家标识
// loginId *battery.TPID 玩家第三方账户信息
//returns:
// account *battery.DBAccount 数据库中玩家账户信息
// err error 返回错误
func (api *XYAPI) AddDBAccount(devid, statisticsDevid, loginIp, uid string, loginId *battery.TPID) (account *battery.DBAccount, err error) {

	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_ADDACCOUNT, &begin)

	xylog.DebugNoId("[%s] New uid: [%s]->[%s], name:%s", loginId.GetId(), uid, loginId.GetId(), loginId.GetName())

	account, err = api.defaultDBAccount()
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("defaultDBAccount failed : %v", err)
		return
	}

	account.Uid = proto.String(uid)

	if loginIp != "" {
		account.CreateIp = proto.String(loginIp)
		account.LastLoginIp = proto.String(loginIp)
	}

	if devid != "" {
		account.Deviceid = proto.String(devid)
	}

	if statisticsDevid != "" {
		account.StatisticsDeviceid = proto.String(statisticsDevid)
	}

	if account.GetCreateDate() == 0 {
		account.CreateDate = proto.Int64(xyutil.CurTimeSec())
	}

	xylog.Debug(uid, "account : %s", account.String())

	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT).AddAccount(account)
	if err != xyerror.DBErrOK {
		xylog.Error(uid, "Adding account failed: %v", err)
	}

	return
}

//生成默认的玩家账户数据信息
func (api *XYAPI) defaultDBAccount() (account *battery.DBAccount, err error) {

	account = &battery.DBAccount{
		Uid: proto.String(""),
		//Identity:              proto.Uint64(DefIdentityIdGenerater.NewID()),
		Deviceid:              proto.String(""),
		StaminaLastUpdateTime: proto.Int64(DefaultStaminaLastUpdateTime),
		IdentityString:        proto.String(""),
		StatisticsDeviceid:    proto.String(""),
	}

	//生成玩家的identity
	*(account.IdentityString), err = xybusinesscache.DefUserIdentityManager.Spawn(api.platform)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("DefUserIdentityManager.Spawn for platform(%v) failed : %v", api.platform, err)
		return
	}

	//初始化钱包
	account.Wallet = make([]*battery.Money, battery.MoneyType_MoneyType_max)
	length := int(battery.MoneyType_MoneyType_max)
	for i := 0; i < length; i++ {
		account.Wallet[i] = api.defaultMoney(i)
	}

	return
}

//生成默认的玩家成就信息
func (api *XYAPI) defaultDBUserAccomplishment(uid string, now int64) (dbUserAccomplishment *battery.DBUserAccomplishment) {
	dbUserAccomplishment = &battery.DBUserAccomplishment{
		Uid:        proto.String(uid),
		CreateDate: proto.Int64(now),
	}

	//初始化成就
	length := int(battery.AccomplishmentType_AccomplishmentType_max)
	dbUserAccomplishment.Accomplishment = make([]*battery.Accomplishment, battery.AccomplishmentType_AccomplishmentType_max)
	for i := 0; i < length; i++ {
		dbUserAccomplishment.Accomplishment[i] = api.defaultAccomplishment()
	}

	return
}

//生成默认的代币数据
// i int 代币类型
//return:
// money *battery.Money 生成的代币指针
func (api *XYAPI) defaultMoney(i int) (money *battery.Money) {
	var gainamount uint32 = 0
	//体力和钻石填写默认
	switch i {
	case int(battery.MoneyType_stamina):
		gainamount = DefConfigCache.Configs().DefaultStamina
	case int(battery.MoneyType_diamond):
		gainamount = DefConfigCache.Configs().DefaultDiamond
	case int(battery.MoneyType_chip):
		gainamount = DefConfigCache.Configs().DefaultChip
	case int(battery.MoneyType_coin):
		gainamount = DefConfigCache.Configs().DefaultCoin
	default:
		//do nothing
	}

	monteyType := battery.MoneyType(i)

	return &battery.Money{
		Type:       &monteyType,
		Iapamount:  proto.Uint32(0),
		Oapamount:  proto.Uint32(0),
		Gainamount: proto.Uint32(gainamount),
	}
}

//默认的成绩
func (api *XYAPI) defaultAccomplishment() (accomplishment *battery.Accomplishment) {
	accomplishment = &battery.Accomplishment{}
	return
}

//增加玩家登录日志信息
// addr string 玩家手机ip
// uid string 玩家标识
// sid string 玩家第三方账户标识
// source battery.ID_SOURCE 第三方账户标识
// name string 玩家第三方账户名称
// op int32 登录类型
// ver string 客户端版本信息
// fail_reason battery.ErrorCode 操作错误码
// desc string 描述信息
func (api *XYAPI) AddAccountLog(addr, uid, sid string, source battery.ID_SOURCE, name string, op int32, ver int32, fail_reason battery.ErrorCode, desc string) (err error) {
	l := battery.AccountLog{
		Uid:           proto.String(uid),
		Address:       proto.String(addr),
		Sid:           proto.String(sid),
		Name:          proto.String(name),
		Source:        source.Enum(),
		OpType:        proto.Int32(op),
		Timestamp:     proto.Int64(xyutil.CurTimeSec()),
		ClientVersion: proto.Int32(ver),
		Error: &battery.Error{
			Code: fail_reason.Enum(),
			Desc: &desc,
		},
		OpDateStr: proto.String(xyutil.CurTimeStr()),
	}
	err = api.GetLogDB().AddAccountLog(l)
	return
}

//刷新玩家的碎片信息
// uid string 玩家标识
// accountWithFlag *AccountWithFlag 玩家账户信息
// amount uint32 数目
// delay bool 是否延迟刷新
func (api *XYAPI) updateUserDataWithChip(uid string, accountWithFlag *AccountWithFlag, amount uint32, delay bool, moneySubType battery.MoneySubType) (err error) {
	if delay {
		return api.updateUserDataWithMoney(uid, accountWithFlag, battery.MoneyType_chip, moneySubType, amount)
	} else {
		return api.updateUserDataWithMoneyDirect(uid, battery.MoneyType_chip, moneySubType, amount)
	}

}

//刷新玩家的金币信息
// uid string 玩家标识
// accountWithFlag *AccountWithFlag 玩家账户信息
// amount uint32 数目
// delay bool 是否延迟刷新
func (api *XYAPI) updateUserDataWithCoin(uid string, accountWithFlag *AccountWithFlag, amount uint32, delay bool, moneySubType battery.MoneySubType) (err error) {
	if delay {
		return api.updateUserDataWithMoney(uid, accountWithFlag, battery.MoneyType_coin, moneySubType, amount)
	} else {
		return api.updateUserDataWithMoneyDirect(uid, battery.MoneyType_coin, moneySubType, amount)
	}
}

//刷新玩家的钻石信息
// uid string 玩家标识
// accountWithFlag *AccountWithFlag 玩家账户信息
// amount uint32 数目
// delay bool 是否延迟刷新
func (api *XYAPI) updateUserDataWithDiamond(uid string, accountWithFlag *AccountWithFlag, amount uint32, delay bool, moneySubType battery.MoneySubType) (err error) {
	if delay {
		return api.updateUserDataWithMoney(uid, accountWithFlag, battery.MoneyType_diamond, moneySubType, amount)
	} else {
		return api.updateUserDataWithMoneyDirect(uid, battery.MoneyType_diamond, moneySubType, amount)
	}
}

//刷新玩家的徽章信息
// uid string 玩家标识
// accountWithFlag *AccountWithFlag 玩家账户信息
// amount uint32 数目
// delay bool 是否延迟刷新
func (api *XYAPI) updateUserDataWithBadge(uid string, accountWithFlag *AccountWithFlag, amount uint32, delay bool, moneySubType battery.MoneySubType) (err error) {
	if delay {
		return api.updateUserDataWithMoney(uid, accountWithFlag, battery.MoneyType_badge, moneySubType, amount)
	} else {
		return api.updateUserDataWithMoneyDirect(uid, battery.MoneyType_badge, moneySubType, amount)
	}
}

//刷新玩家的代币信息，延后刷新
// uid string 玩家标识
// accountWithFlag *AccountWithFlag 玩家账户信息
// amount uint32 数目
// mainType battery.MoneyType 代币主类型
// subType battery.MoneySubType 代币子类型
func (api *XYAPI) updateUserDataWithMoney(uid string, accountWithFlag *AccountWithFlag, mainType battery.MoneyType, subType battery.MoneySubType, amount uint32) (err error) {
	err = xymoney.Add(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT), amount, mainType, subType, accountWithFlag.account, ACCOUNT_UPDATE_DELAY)
	if err == xyerror.ErrOK {
		accountWithFlag.SetChange()
	}
	return
}

//刷新玩家的代币信息，直接刷新
// uid string 玩家标识
// mainType battery.MoneyType 代币主类型
// subType battery.MoneySubType 代币子类型
// amount uint32 数目
func (api *XYAPI) updateUserDataWithMoneyDirect(uid string, mainType battery.MoneyType, subType battery.MoneySubType, amount uint32) (err error) {
	return xymoney.Add(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT), amount, mainType, subType, nil, ACCOUNT_UPDATE_NO_DELAY)
}
