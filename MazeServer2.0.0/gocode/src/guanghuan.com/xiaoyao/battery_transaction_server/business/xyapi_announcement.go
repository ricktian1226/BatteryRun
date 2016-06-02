// xyapi_announcement
package batteryapi

import (
	"guanghuan.com/xiaoyao/common/idgenerate"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

//公告id生成器定义
var DefAnnouncementIdGenerater *xyidgenerate.IdGenerater

func (api *XYAPI) OperationAnnouncement(req *battery.AnnouncementRequest, resp *battery.AnnouncementResponse) (err error) {

	uid := req.GetUid()

	//初始化返回值
	resp.Uid, resp.Cmd = req.Uid, req.Cmd
	resp.Error = xyerror.DefaultError()

	//获取操作码
	cmd := req.GetCmd()
	switch cmd {
	case battery.OP_CMD_QueryByTime:
		fallthrough
	case battery.OP_CMD_Query:
		err = api.QueryAnnouncement(req, resp)
	default:
		xylog.Error(uid, "Unkown announcement CMD %d", cmd)
		err = xyerror.DBErrFailedDueToClientError
	}

	return
}

//查询公告信息
func (api *XYAPI) QueryAnnouncement(req *battery.AnnouncementRequest, resp *battery.AnnouncementResponse) (err error) {

	uid := req.GetUid()

	//客户端上报的已有的公告列表
	m := make(map[uint64]bool, 0)
	for _, item := range req.Items {
		m[*item.Id] = false
	}

	xylog.Debug(uid, "Announcements on client : %v .", m)

	//从数据库中查询符合当前时间戳的公告并且有效的公告列表
	announcements := xybusinesscache.DefAnnouncementsCacheManager.AnnouncementsSlice()
	xylog.Debug(uid, "Announcements today : %v.", announcements)

	length := len(announcements)
	for i := 0; i < length; i++ {
		id := announcements[i].GetId()
		if _, ok := m[id]; ok {
			m[id] = true
			xylog.Debug(uid, "Announcement %d already exists in client,skip it.", id)
		} else {
			resp.Items = append(resp.Items, &announcements[i])
			xylog.Debug(uid, "Announcement %d does not exist in client,add it.", id)
		}
	}

	state := battery.ANNOUNCEMENT_STATE_ANNOUNCEMENT_STATE_INVALID
	for k, v := range m {
		//客户端上有的，但是没找到或者状态为invalid的，需要告诉客户端删除
		id := k
		if !v {
			t := &battery.Announcement{
				Id:    &id,
				State: &state,
			}
			resp.Items = append(resp.Items, t)
			xylog.Debug(uid, "Announcement %d is invalid.", k)
		}
	}

	return
}
