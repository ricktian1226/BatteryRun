package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"errors"
	"fmt"
	nats "github.com/nats-io/nats"
	batterydb "guanghuan.com/xiaoyao/battery_app_server/db"
	xyconf "guanghuan.com/xiaoyao/common/conf"
	"guanghuan.com/xiaoyao/common/db"
	//xyidgenerate "guanghuan.com/xiaoyao/common/idgenerate"
	xylog "guanghuan.com/xiaoyao/common/log"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//"guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"math/rand"
	"strconv"
	"time"
)

type XYAPI struct {
	Config   *ConfigCache
	platform battery.PLATFORM_TYPE //请求的平台类型
}

//var DefReceiptIdGenerater *xyidgenerate.IdGenerater

var apiConfigUtil xyconf.ApiConfigUtil

func NewXYAPI() *XYAPI {
	return &XYAPI{
		Config:   DefConfigCache,
		platform: battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN,
	}
}

//根据终端平台类型设置业务数据库指针，在业务请求时调用
// platform battery.PLATFORM_TYPE 终端平台类型
func (api *XYAPI) SetDB(platform battery.PLATFORM_TYPE) {
	api.platform = platform
}

//获取业务数据库会话
// index int 业务索引
func (api *XYAPI) GetDB(index int) *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, api.platform)
	return dbInterface.(*batterydb.BatteryDB)
}
func (api *XYAPI) GetXYDB(index int) *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, api.platform)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}

func (api *XYAPI) GetCommonDB(index int) *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return dbInterface.(*batterydb.BatteryDB)
}

func (api *XYAPI) GetCommonXYDB(index int) *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}

func (api *XYAPI) GetCommonXYBusinessDB(index int) *xybusiness.XYBusinessDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYBusinessDB)
}

func (api *XYAPI) GetLogDB() *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return dbInterface.(*batterydb.BatteryDB)
}

func (api *XYAPI) GetLogXYDB() *xydb.XYDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOG, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	return &(dbInterface.(*batterydb.BatteryDB).XYDB)
}

func (api *XYAPI) FormMsg(uid string, uidSuffix uint32, businessCode int, reqData *[]byte) {

	msgHeader := make([]byte, 0)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uidSuffix)
	msgHeader = append(msgHeader, buf...)
	binary.LittleEndian.PutUint32(buf, uint32(businessCode))
	msgHeader = append(msgHeader, buf...)
	xylog.Debug(uid, "FormMsg uidSuffix %d businessCode %d msgHeader %v", uidSuffix, businessCode, msgHeader)
	*reqData = append(msgHeader, *reqData...) //拼上消息头
}

func (api *XYAPI) SendToTransaction(uid string, businessCode int, reqData proto.Message, respData proto.Message) (failReason battery.ErrorCode, err error) {
	var (
		subject   string
		uidSuffix uint32
		req_data  []byte
		resp_data []byte
		reply     *nats.Msg
	)

	failReason = battery.ErrorCode_NoError

	if uid != "" {
		if DefBannedUserManager.IsUidBanned(uid) { //校验玩家是否被封号，如果被封号则直接返回错误码
			xylog.Debug(uid, "Uid  is banned")
			failReason = battery.ErrorCode_UserBannedError
			return
		}
	}

	//根据uid获取下游的transaction服务节点subject
	failReason, err = api.getTransactionNode(uid, &subject, &uidSuffix)
	if failReason != battery.ErrorCode_NoError || err != xyerror.ErrOK {
		return
	}

	//xylog.Debug(uid, "[%s] Get uidSuffix : %d, businessCode : %v, transactionNode : %s", uid, uidSuffix, businessCode, subject)

	//下发请求，等待返回
	req_data, err = proto.Marshal(reqData)
	if err != nil {
		failReason = battery.ErrorCode_ServerError
		xylog.Error(uid, "[app iapvalidate] proto.Marshal failed : %v", err)
		return
	}
	//xylog.Debug("[%s] req_data before formmsg: %v", uid, req_data)

	//拼上uidSuffix和businessCode
	api.FormMsg(uid, uidSuffix, businessCode, &req_data)

	//xylog.Debug("[%s] req_data after formmsg: %v", uid, req_data)

	reply, err = xynatsservice.Nats_service.Request(subject, req_data, time.Duration(api.Config.Configs().NatsTimeOut)*time.Second)

	if err != nil {
		xylog.Error(uid, "<%s> Error: %s", subject, err.Error())
		failReason = battery.ErrorCode_SendToTransactionError
		return
	} else {
		if reply != nil && len(reply.Data) > 0 {
			//xylog.Debug("[%s] reply: %v", uid, *reply)
			resp_data = reply.Data
			err = proto.Unmarshal(resp_data, respData)
			if err != xyerror.ErrOK {
				failReason = battery.ErrorCode_ServerError
				xylog.Error(uid, "[app iapvalidate] proto.Unmarshal failed : %v", err)
				return
			}
		} else {
			xylog.Error(uid, "[app iapvalidate] no reply data")
			err = errors.New("no reply data")
		}
	}
	return
}

//广播消息接口
// subject string 消息subject
// msg []byte 消息体
func (api *XYAPI) NatsPublish(subject string, msg []byte) (failReason battery.ErrorCode) {
	err := xynatsservice.Nats_service.Publish(subject, msg)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("NatsPublish to %s failed : %v", subject, err)
		failReason = battery.ErrorCode_ServerError
		return
	}
	return
}

//根据uid获取transaction node路由
// uid string 玩家id
// transactionNode *string
// uidSuffix *uint32 包含uid后9位十进制整型的路由信息，需要传给transaction，作为channel路由的依据
func (api *XYAPI) getTransactionNode(uid string, transactionNode *string, uidSuffix *uint32) (failReason battery.ErrorCode, err error) {
	var (
		i int
		n = len(uid)
	)

	//配置错误，直接返回错误
	if api.Config.Configs().TransactionNodeCount <= 0 {
		failReason = battery.ErrorCode_ServerError
		s := fmt.Sprintf("error TransactionNodeCount %d <= 0", api.Config.Configs().TransactionNodeCount)
		xylog.ErrorNoId(s)
		err = errors.New(s)
		return
	}

	//如果uid不足9位，直接随机一个node
	if n < 9 {
		i = rand.Intn(api.Config.Configs().TransactionNodeCount)
		*uidSuffix = uint32(rand.Intn(int(1e9))) //取的是uid的后9位数字，所以是1e9以内
	} else {
		sub := uid[len(uid)-9:] //取uid的后9位数字转换为。之所以为9，因为9位数字表示的整型范围为[0,999999999],数字够大，又足够安全
		i, err = strconv.Atoi(string(sub))
		if err != nil {
			failReason = battery.ErrorCode_ServerError
			s := fmt.Sprintf("can't get transactionNode string cause' sub is %s", string(sub))
			xylog.ErrorNoId(s)
			err = errors.New(s)
			return
		}
		*uidSuffix = uint32(i)
	}

	//取uid后9位整型数的高5位作为transaction node路由依据
	node := (int(*uidSuffix) / xybusiness.BASE_UID_BAND) % api.Config.Configs().TransactionNodeCount
	xylog.DebugNoId("(uidSuffix(%d) / xybusiness.BASE_UID_BAND(%d)) api.Config.TransactionNodeCount(%d) : %d", *uidSuffix, xybusiness.BASE_UID_BAND, api.Config.Configs().TransactionNodeCount, node)
	*transactionNode = fmt.Sprintf("transaction%03d", node)

	return
}
