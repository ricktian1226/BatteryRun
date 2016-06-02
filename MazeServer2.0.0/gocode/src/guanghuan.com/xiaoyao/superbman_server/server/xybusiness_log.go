// xybusiness_log
package xybusiness

import (
	"code.google.com/p/goprotobuf/proto"
	"guanghuan.com/xiaoyao/common/db"
	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

var (
	DB_TABLE_OPERATION_LOG = "operationlog"
)

func defaultOperationLog(uid, req, resp string, businessCode int, optype int32, err *battery.Error) *battery.OperationLog {
	return &battery.OperationLog{
		Uid:          proto.String(uid),
		Req:          proto.String(req),
		Resp:         proto.String(resp),
		BusinessCode: proto.Int32(int32(businessCode)),
		Optype:       proto.Int32(optype),
		Opdate:       proto.String(xyutil.CurTimeStr()),
		Error:        err,
	}
}

func AddOperationLog(uid, req, resp string, db *xydb.XYDB, businessCode int, optype int32, err *battery.Error) error {
	return db.AddData(DB_TABLE_OPERATION_LOG, defaultOperationLog(uid, req, resp, businessCode, optype, err))
}

//初始化打开调试开关的玩家id信息
// db *XYBusinessDB 数据库指针
func LoadDebugUsers(db *XYBusinessDB) (err error) {
	var ids []interface{}
	ids, err = db.QueryDebugUsers()
	if err == nil {
		xylog.DefIdManager.Load(ids)
	} else {
		xylog.ErrorNoId("QueryDebugUsers failed : %v", err)
	}
	return
}
