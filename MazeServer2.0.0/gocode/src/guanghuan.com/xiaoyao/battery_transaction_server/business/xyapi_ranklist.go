package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xycache "guanghuan.com/xiaoyao/superbman_server/cache"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// 玩家全局排行榜查询
func (api *XYAPI) OperationGetGlobalRankList(req *battery.QueryGlobalRankRequest, resp *battery.QueryGlobalRankResponse) (err error) {
	var (
		uid = req.GetUid()
		//   sortType            = req.GetSortType() // 排序方式，目前没需求，全部按总分排序，时间紧有时间在加其他排序
		tpids               = make([]*battery.IDMap, 0)
		userAccomplishments = make([]*battery.DBUserAccomplishment, 0)
		users               = make([]*battery.FriendData, 0)
		mapUid2Tpid         = make(map[string]*battery.IDMap, 0)
		errStruct           = new(battery.Error)
		//  uids =make([]string,0)
	)
	platform := req.GetPlatform()
	api.SetDB(platform)

	*errStruct = *(xyerror.Resp_NoError)
	// 从缓存加载top排行榜
	gids := xycache.DefGlobalRankListManager.GlobalRank()
	xylog.DebugNoId("len :%v,gids :%v", len(gids), gids)
	if !api.isUserinGlobalRank(uid, gids) {
		// 排行榜数量不够，直接增加玩家自身，否则替换最后一名
		if len(gids) < DefConfigCache.Configs().GlobalRankSize {
			gids = append(gids, uid)
		} else {
			gids[len(gids)-1] = uid
		}

	}
	// 查询玩家tpid信息
	{
		selector := bson.M{"gid": 1, "sid": 1, "note": 1, "iconurl": 1}
		err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TPIDMAP).XYBusinessDB.QueryTpidsByUids(selector, gids, &tpids)
		if err != xyerror.ErrOK {
			errStruct.Code = battery.ErrorCode_QueryUserFriendsUidError.Enum()
			err = xyerror.ErrQueryUserFriendsUidError
			goto ErrHandle
		}
	}
	xylog.DebugNoId("tpids :%v", tpids)
	for _, tpid := range tpids {
		gid := tpid.GetGid()
		//sid := tpid.GetSid()
		if len(gid) > 0 {
			mapUid2Tpid[gid] = tpid
		}
	}
	err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_USERACCOMPLISHMENT).GetUsersAccoplishment(gids, &userAccomplishments, mgo.Monotonic)
	if err != xyerror.ErrOK {
		errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
		err = xyerror.ErrGetAccountByUidError
		goto ErrHandle
	}
	xylog.DebugNoId("len gids:%v,tpids:%v,accomplish:%v", len(gids), len(tpids), len(userAccomplishments))
	for _, userAccomplishment := range userAccomplishments {
		var (
			gid = userAccomplishment.GetUid()
			//sid string
			tpid *battery.IDMap
			ok   bool
		)

		//user := api.getFriendDataFromAccount(account)
		user := api.getFriendDataFromAccomplishment(userAccomplishment)
		//将好友id转换为第三方账户id
		if tpid, ok = mapUid2Tpid[gid]; ok {
			user.Uid = proto.String(tpid.GetSid())
			user.IconUrl = proto.String(tpid.GetIconUrl())
			user.Name = proto.String(tpid.GetNote())
		}
		users = append(users, user)

	}
	resp.Data = users
	resp.Source = req.Source
	resp.Uid = req.Uid
ErrHandle:
	resp.Error = errStruct
	xylog.Debug(uid, "errStruct : %v, resp.Error : %v", errStruct, resp.Error)
	return
}

// 判断玩家是否在全局排行榜内
func (api *XYAPI) isUserinGlobalRank(uid string, uids []string) bool {
	for _, tuid := range uids {
		if uid == tuid {
			return true
		}
	}
	return false
}

func (api *XYAPI) getFriendDataFromAccomplishment(userAccomplishment *battery.DBUserAccomplishment) (friendData *battery.FriendData) {
	friendData = &battery.FriendData{
		Uid:          proto.String(userAccomplishment.GetUid()),
		Selectroleid: proto.Uint64(api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_ROLEINFO).GetSelectRoleID(userAccomplishment.GetUid())),
	}
	xylog.DebugNoId("frienddata :%v", friendData)
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

// 更新玩家排行榜数据
func (api *XYAPI) updateUserRankList(uid string, checkpointsum *battery.DBUserCheckPoint, platform battery.PLATFORM_TYPE) {
	ranklist := &battery.RankLisk{
		Uid:          proto.String(uid),
		Score:        checkpointsum.Score,
		Charge:       checkpointsum.Charge,
		Star:         proto.Uint64(uint64(checkpointsum.GetCollectionsCount())),
		PlatformType: platform.Enum(),
	}

	err := api.GetCommonDB(xybusiness.Business_COMMON_COLLECTION_INDEX_RANKLIST).UpsertRankList(uid, ranklist)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("update user ranklist fail :%v", err)
		return
	}
	return
}
