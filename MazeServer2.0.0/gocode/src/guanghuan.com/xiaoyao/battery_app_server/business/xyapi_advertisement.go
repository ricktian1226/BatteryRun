// xyapi_announcement
package batteryapi

import (
	"fmt"

	"code.google.com/p/goprotobuf/proto"

	"guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xybusinesscache "guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
)

func (api *XYAPI) OperationAdvertisement(req *battery.AdvertisementRequest, resp *battery.AdvertisementResponse) (err error) {

	var (
		uid    = req.GetUid()
		errStr string
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	resp.Uid = req.Uid
	resp.PlatformType = req.PlatformType
	resp.Error = xyerror.DefaultError()

	// 查询用户是否可用
	if !api.isUidValid(uid) {
		errStr = fmt.Sprintf("[%s] invalid user.", uid)
		xylog.ErrorNoId(errStr)
		err = xyerror.ErrGetAccountByUidError
		resp.Error.Code = battery.ErrorCode_GetAccountByUidError.Enum()
		goto ErrHandle
	}

	//选取广告信息
	api.ElectAdvertisement(req.GetAdvertisementSpaceId(), &(resp.Items))

ErrHandle:
	return
}

// ElectAdvertisement 选取广告
// ids []uint32 广告位id列表
// items *[]*battery.AdvertisementItem 广告信息列表指针
func (api *XYAPI) ElectAdvertisement(ids []uint32, items *[]*battery.AdvertisementItem) {

	if len(ids) <= 0 {
		return
	}

	xylog.DebugNoId("ElectAdvertisement for %v", ids)

	*items = make([]*battery.AdvertisementItem, 0)
	for _, id := range ids {

		ret, enable, flags, advertisement := xybusinesscache.DefAdvertisementManager.Elect(id)
		if ret && advertisement != nil {
			advertisementItem := &battery.AdvertisementItem{
				AdvertisementSpaceId: proto.Uint32(id),
				Item:                 advertisement,
				Enable:               proto.Bool(enable),
				Flags:                flags,
			}
			*items = append(*items, advertisementItem)
		}
	}

	xylog.DebugNoId("Advertisement items : %v", *items)
}
