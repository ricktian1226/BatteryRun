package batteryapi

import (
	//"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

func (api *XYAPI) OperationMemCache(req *battery.MemCacheRequest, resp *battery.MemCacheResponse) (err error) {

	uid, optype, key, value := req.GetUid(), req.GetOpType(), req.GetOpKey(), req.GetOpValue()
	errCode := battery.ErrorCode_NoError

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	resp.OpKey, resp.PlatformType = req.OpKey, req.PlatformType

	switch {
	case optype == battery.MemCacheOperationType_MemCacheOperationType_GET &&
		len(req.Units) == 0:
		errCode, value = api.GetMemCacheValue(uid, key)
		if errCode == battery.ErrorCode_NoError {
			resp.OpValue = &value
		}

	case optype == battery.MemCacheOperationType_MemCacheOperationType_SET &&
		len(req.Units) == 0:
		errCode = api.SetMemCacheValue(uid, value, key)
		resp.OpValue = req.OpValue

	case optype == battery.MemCacheOperationType_MemCacheOperationType_GET &&
		len(req.Units) > 0:
		resp.Units = req.Units
		errCode = api.GetMemCacheValues(uid, &(resp.Units))

	case optype == battery.MemCacheOperationType_MemCacheOperationType_SET &&
		len(req.Units) > 0:
		resp.Units = req.Units
		errCode = api.SetMemCacheValues(uid, &(resp.Units))
	}

	resp.Error = &battery.Error{
		Code: errCode.Enum(),
	}

	return
}

func (api *XYAPI) GetMemCacheValue(uid string, key battery.MemCacheEnum) (errCode battery.ErrorCode, value string) {
	////todelete
	//db := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MEMCACHE)
	//if db == nil {
	//	errCode = battery.ErrorCode_SetMemCacheError
	//	return
	//}

	var err error
	//先查下brcommondb.memcache下有没有，有的话，直接返回
	value, err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MEMCACHE).GetMemCacheValue(uid, key, api.platform, mgo.Strong)
	if err != xyerror.ErrOK {
		errCode = battery.ErrorCode_GetMemCacheError
		return
	}

	xylog.Debug(uid, "GetMemCacheValue(%v) platform(%v) : %s", key, api.platform, value)

	return
}

func (api *XYAPI) GetMemCacheValues(uid string, units *[]*battery.MemCacheUnit) (errCode battery.ErrorCode) {

	mapKey2Values, keys := make(map[battery.MemCacheEnum]*battery.DBMemCache, 0), make([]battery.MemCacheEnum, 0)
	for _, u := range *units {
		keys = append(keys, u.GetKey())
	}

	values, err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MEMCACHE).GetMemCacheValues(uid, keys, api.platform, mgo.Strong)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		errCode = battery.ErrorCode_GetMemCacheError
		return
	}

	xylog.Debug(uid, "GetMemCacheValues : %v", values)

	for _, v := range values {
		mapKey2Values[v.GetKey()] = v
	}

	//返回查询结果
	for i, _ := range *units {
		if v, ok := mapKey2Values[(*units)[i].GetKey()]; ok {
			(*units)[i].Value = v.Value
			(*units)[i].Result = battery.ErrorCode_NoError.Enum()
		} else {
			(*units)[i].Result = battery.ErrorCode_RecordNotExisting.Enum()
		}
	}

	xylog.Debug(uid, "GetMemCacheValues result : %v", units)

	return
}

func (api *XYAPI) SetMemCacheValue(uid, value string, key battery.MemCacheEnum) (errCode battery.ErrorCode) {
	////todelete
	//db := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MEMCACHE)
	//if db == nil {
	//	errCode = battery.ErrorCode_SetMemCacheError
	//	return
	//}

	err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MEMCACHE).SetMemCacheValue(uid, value, key, api.platform, mgo.Strong)
	if err != xyerror.ErrOK {
		errCode = battery.ErrorCode_SetMemCacheError
		return
	}
	xylog.Debug(uid, "SetMemCacheValue(%v) : %s succeed", key, value)
	return
}

func (api *XYAPI) SetMemCacheValues(uid string, units *[]*battery.MemCacheUnit) (errCode battery.ErrorCode) {

	for _, unit := range *units {
		key, value := unit.GetKey(), unit.GetValue()
		err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MEMCACHE).SetMemCacheValue(uid, value, key, api.platform, mgo.Strong)
		if err != xyerror.ErrOK {
			unit.Result = battery.ErrorCode_SetMemCacheError.Enum()
			xylog.Debug(uid, "SetMemCacheValue(%v) : %s failed", key, value)
		} else {
			unit.Result = battery.ErrorCode_NoError.Enum()
			xylog.Debug(uid, "SetMemCacheValue(%v) : %s succeed", key, value)
		}
	}

	return
}
