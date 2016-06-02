package batteryapi

import (
	"time"

	proto "code.google.com/p/goprotobuf/proto"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"

	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// 绑定账户消息处理
func (api *XYAPI) OperationBind(req *battery.BindRequest, resp *battery.BindResponse) (err error) {
	var (
		uid       = req.GetUid()
		dst       = req.GetTarget()
		platform  = req.GetPlatformType()
		errStruct = xyerror.DefaultError()
	)

	api.SetDB(platform)
	resp.Error = xyerror.DefaultError()
	resp.PlatformType = req.PlatformType
	resp.Uid = req.Uid

	//查询源账户信息
	srcTpid, err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetIdMapByGid(uid)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "GetIdMapByGid %v failed : %v", srcTpid, err)
		resp.Error.Code = battery.ErrorCode_DBError.Enum()
		return
	}

	//查询目标账户
	isExist, err := api.IsTPIDRegistered(dst)
	if err != xyerror.ErrOK { //数据库
		if err != xyerror.ErrNotFound {
			xylog.Error(uid, "IsTPIDRegistered failed : %v", err)
			resp.Error.Code = battery.ErrorCode_DBError.Enum()
			return
		}
	}

	//如果目标账户已经存在，则无法再进行重复绑定
	if isExist {
		xylog.Error(uid, "TPID %v already exist", dst)
		resp.Error.Code = battery.ErrorCode_TargetTpidAlreadyExistError.Enum()
		return
	}

	//创建新的目标账户tpid信息
	err = api.RegisterTPIDForBind(dst, uid)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "RegisterTPID failed : %v", err)
		resp.Error.Code = battery.ErrorCode_RegistTpidError.Enum()
		return
	}

	//删除源账户信息
	err = api.RemoveTPID(srcTpid.GetGid(), srcTpid.GetSource())
	if err != xyerror.ErrOK {
		xylog.Error(uid, "RegisterTPID failed : %v", err)
		resp.Error.Code = battery.ErrorCode_RemoveTpidError.Enum()
		return
	}

	//发放登录礼包
	api.newAccountGainProps(uid, dst.GetSource(), errStruct)

	resp.Error = errStruct

	return
}

//func (api *XYAPI) GidToGCid(gid string) (gcid string, err error) {
//	return api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetSidByGid(gid, battery.ID_SOURCE_SRC_GAMECENTER)
//}

//func (api *XYAPI) GidToSid(gid string, src battery.ID_SOURCE) (sid string, err error) {
//	//if api.IsTPIDSupported(src) {
//	sid, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetSidByGid(gid, src)
//	//} else {
//	//	err = xyerror.ErrNotSupport
//	//}
//	return
//}

// NoteByGid 根据玩家uid获取玩家昵称
// gid string 玩家标识
//returns:
// note string 玩家昵称
// err error 操作错误
func (api *XYAPI) NoteByGid(gid string) (note string, err error) {
	note, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetNoteByGid(gid)
	return
}

func (api *XYAPI) SidToGid(sid string, src battery.ID_SOURCE) (gid string, err error) {
	//if api.IsTPIDSupported(src) {
	gid, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetGidBySid(sid, src)
	//} else {
	//	err = xyerror.ErrNotSupport
	//}
	return
}

func (api *XYAPI) IsTPIDRegistered(tpid *battery.TPID) (isExisting bool, err error) {
	isExisting, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).IsTpidRegistered(tpid)
	return
}

//func (api *XYAPI) IsTPIDSupported(source battery.ID_SOURCE) (isSupported bool) {
//	switch source {
//	case battery.ID_SOURCE_SRC_GAMECENTER, battery.ID_SOURCE_SRC_SINA_WEIBO:
//		isSupported = true
//	}
//	return
//}

// 注册tpid信息
// tpid *battery.TPID 待注册的tpid信息
// sourceGid string 源账户gid
// iconUrl string 账户头像图片下载链接
func (api *XYAPI) RegisterTPID(tpid *battery.TPID, sourceGid, iconUrl string) (gid string, err error) {
	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_ADDTPID, &begin)

	idMap := &battery.IDMap{
		Gid:        proto.String(xyutil.NewId()),
		Source:     tpid.GetSource().Enum(),
		Sid:        proto.String(tpid.GetId()),
		Note:       proto.String(tpid.GetName()),
		CreateDate: proto.Int64(xyutil.CurTimeSec()),
		IconUrl:    proto.String(iconUrl),
		//State:      battery.USER_ACCOUNT_STATE_USER_ACCOUNT_STATE_UNBIND.Enum(), //默认账户都是未绑定的
		//Pgid:       proto.String(sourceGid),                                     //源账户gid
	}

	gid, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).AddTPID(idMap)

	return
}

// 注册tpid信息
// tpid *battery.TPID 待注册的tpid信息
// sourceGid string 源账户gid
// iconUrl string 账户头像图片下载链接
func (api *XYAPI) RegisterTPIDForBind(tpid *battery.TPID, sourceGid string) (err error) {
	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_ADDTPID, &begin)

	idMap := &battery.IDMap{
		Gid:        proto.String(sourceGid), //绑定时新生成的
		Source:     tpid.GetSource().Enum(),
		Sid:        proto.String(tpid.GetId()),
		Note:       proto.String(tpid.GetName()),
		CreateDate: proto.Int64(xyutil.CurTimeSec()),
		//IconUrl:    proto.String(iconUrl),
		//State: battery.USER_ACCOUNT_STATE_USER_ACCOUNT_STATE_UNBIND.Enum(), //默认账户都是未绑定的
		//Pgid:       proto.String(sourceGid),                                     //源账户gid
	}

	_, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).AddTPID(idMap)

	return
}

//是否需要刷新tpid信息
// tpid信息在老玩家登录时，可能会存在name或者iconurl的刷新，通过该函数进行判断
func (api *XYAPI) ShouldUpdateTpid(tpid *battery.TPID, iconUrl string, idMap *battery.IDMap) (change bool) {

	//玩家名称改变
	if tpid.GetName() != "" {
		if tpid.GetName() != idMap.GetNote() {
			idMap.Note = tpid.Name
			change = true
		}
	}

	//玩家头像url改变
	if iconUrl != idMap.GetIconUrl() {
		idMap.IconUrl = proto.String(iconUrl)
		change = true
	}
	return
}

//刷新玩家第三方信息
func (api *XYAPI) UpsertTPID(idMap *battery.IDMap) error {
	return api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).UpsertTPID(idMap)
}

// 删除玩家第三方账户信息
func (api *XYAPI) RemoveTPID(gid string, source battery.ID_SOURCE) error {
	return api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).RemoveTPID(gid, source)
}

// 玩家起名
func (api *XYAPI) OperationCreatName(req *battery.CreatNameRequest, resp *battery.CreatNameResponse) (err error) {
	var (
		uid   = req.GetUid()
		name  = req.GetName()
		idMap = battery.IDMap{}
	)
	resp.Name = req.Name
	resp.Uid = req.Uid
	resp.Error = xyerror.DefaultError()
	resp.PlatformType = req.PlatformType

	idMap, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).GetIdMapByGid(uid)
	if err != xyerror.ErrOK {
		resp.Error.Code = battery.ErrorCode_QueryTpidError.Enum()
		return
	}
	idMap.Note = proto.String(name)
	api.updateUserName(&idMap)
	return

}

func (api *XYAPI) updateUserName(idMap *battery.IDMap) (err error) {
	return api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).UpsertTPIDName(idMap)
}
