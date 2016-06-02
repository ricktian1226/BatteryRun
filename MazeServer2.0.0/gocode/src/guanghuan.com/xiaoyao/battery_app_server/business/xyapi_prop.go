// xyapi_prop
package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	//"guanghuan.com/xiaoyao/common/cache"
	"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
)

func (api *XYAPI) OperationQueryPropRes(req *battery.QueryPropResRequest, resp *battery.QueryPropResResponse) (err error) {

	var (
		uid    = req.GetUid()
		errStr string
	)

	//初始化resp
	resp.Uid = req.Uid

	if !api.isUidValid(uid) {
		errStr = fmt.Sprintf("[%s] GetAccountByUidError", uid)
		xylog.ErrorNoId(errStr)
		err = xyerror.ErrGetAccountByUidError
		resp.Error = xyerror.Resp_GetAccountByUidError
		goto ErrHandle
	}

	err = api.queryPropRes(uid, &(resp.Props))
	if err != xyerror.ErrOK {
		errStr = fmt.Sprintf("[%s] QueryPropsFromCacheError : %v", uid, err)
		xylog.ErrorNoId(errStr)
		resp.Error = xyerror.Resp_QueryPropsFromCacheError
		goto ErrHandle
	}

ErrHandle:
	//resp.Error.Desc = proto.String(errStr)
	return
}

//查询道具缓存信息
// uid string 玩家id
// props *[]*battery.Prop 返回的道具信息列表
func (api *XYAPI) queryPropRes(uid string, props *[]*battery.Prop) (err error) {

	mapProps := xybusinesscache.DefPropCacheManager.Props()
	xylog.Debug(uid, "mapProps : %v", *mapProps)
	mapProps.Print()
	if nil != mapProps {
		for k, v := range *mapProps {
			prop := &battery.Prop{
				Id:    proto.Uint64(k),
				Type:  v.Type.Enum(),
				Items: v.Items,
			}
			*props = append(*props, prop)
		}
	} else {
		err = xyerror.ErrQueryPropsFromCacheError
	}

	return
}
