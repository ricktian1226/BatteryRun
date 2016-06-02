// misson
package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"guanghuan.com/xiaoyao/common/idgenerate"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/performance"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xycache "guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	//"math/rand"
	"sort"
	"time"
)

//任务号生成器
var DefMissionIdGenerater *xyidgenerate.IdGenerater

//任务指标不需要记忆点完成，对于游戏外的任务指标，直接指定为false
const MissionQuotasNoNeedFinish = false

//查询任务信息消息接口
func (api *XYAPI) OperationQueryUserMission(req *battery.QueryUserMissionRequest, resp *battery.QueryUserMissionResponse) (err error) {
	var (
		uid               = req.GetUid()
		failReason        = battery.ErrorCode_NoError
		missionTypes      = req.GetTypes()
		missionCount      = int(req.GetMissionCount())
		dailyMissionCount = int(req.GetDailyMissionCount())
		now               = time.Now().Unix()
	)

	//获取请求的终端平台类型
	api.SetDB(req.GetPlatformType())

	//初始化resp
	resp.Uid = req.Uid
	resp.Error = xyerror.DefaultError()

	xylog.Debug(uid, "[%s] QueryUserMission missionTypes : %v, missionCount %d, dailyMissionCount %d", uid, missionTypes, missionCount, dailyMissionCount)

	//如果请求的条数大于上限，设置请求条数为上限
	if missionCount > DefConfigCache.Configs().MissionCountLimit && DefConfigCache.Configs().MissionCountLimit > 0 {
		missionCount = DefConfigCache.Configs().MissionCountLimit
	}
	if dailyMissionCount > DefConfigCache.Configs().DailyMissionCountLimit && DefConfigCache.Configs().DailyMissionCountLimit > 0 {
		dailyMissionCount = DefConfigCache.Configs().DailyMissionCountLimit
	}

	failReason, err = api.QueryUserMission(uid, missionTypes, missionCount, dailyMissionCount, &(resp.Entrys), now)
	if err != xyerror.ErrOK {
		//resp.Error = xyerror.Resp_QueryUserMissionError
		goto ErrHandle
	}

	xylog.Debug(uid, "QueryUserMission result : %v", resp.Entrys)

ErrHandle:
	resp.Error.Code = failReason.Enum()
	if failReason != battery.ErrorCode_NoError {
		resp.Entrys = nil
	}

	return
}

//查询玩家任务
// uid string 玩家id
// missionTypes []battery.MissionType 任务类型
// missionCount int 教学任务+主线任务数目
// dailyMissionCount int 日常任务数目
// entrys *[]*battery.DBUserMissionEntry 任务集合
// now int64 当前时间戳
func (api *XYAPI) QueryUserMission(uid string, missionTypes []battery.MissionType, missionCount int, dailyMissionCount int, entrys *[]*battery.UserMissionEntry, now int64) (failReason battery.ErrorCode, err error) {

	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_QUERYUSERMISSION, &begin)

	//从数据库获取任务信息
	dbEntrys := make([]*battery.DBUserMissionEntry, 0)
	failReason, err = api.QueryDBUserMission(uid, missionTypes, missionCount, dailyMissionCount, &dbEntrys, now)
	if failReason != xyerror.Resp_NoError.GetCode() || err != xyerror.ErrOK {
		return
	}

	//将数据库任务转换为消息任务
	for _, dbEntry := range dbEntrys {
		entry := &battery.UserMissionEntry{
			Type: dbEntry.GetType().Enum(),
		}

		for _, dbUserMission := range dbEntry.GetActive() {
			userMission := api.getUserMissionFromDBUserMission(uid, dbUserMission)
			if nil != userMission {
				entry.Active = append(entry.Active, userMission)
			}

		}

		for _, dbUserMission := range dbEntry.GetDoneNotCollect() {
			userMission := api.getUserMissionFromDBUserMission(uid, dbUserMission)
			if nil != userMission {
				entry.DoneNotCollect = append(entry.DoneNotCollect, userMission)
			}
		}

		for _, dbUserMission := range dbEntry.GetDoneCollect() {
			userMission := api.getUserMissionFromDBUserMission(uid, dbUserMission)
			if nil != userMission {
				entry.DoneCollect = append(entry.DoneCollect, userMission)
			}
		}

		*entrys = append(*entrys, entry)
	}

	return
}

//查询玩家任务
// uid string 玩家id
// missionTypes []battery.MissionType 查询任务类型
// missionCount int 教学任务+主线任务数目
// dailyMissionCount int 日常任务数目
// entrys *[]*battery.DBUserMissionEntry 任务集合
// now int64 当前时间戳
func (api *XYAPI) QueryDBUserMission(uid string, missionTypes []battery.MissionType, missionCount int, dailyMissionCount int, entrys *[]*battery.DBUserMissionEntry, now int64) (failReason battery.ErrorCode, err error) {

	//缓存下请求的任务类型，并做下校验
	setMissionTypesSet := make(Set, 0)
	for _, missionType := range missionTypes {
		switch missionType {
		case battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_Study, battery.MissionType_MissionType_MainLine:
			setMissionTypesSet[missionType] = empty{}
		default:
			xylog.Error(uid, "unkown MissionType %v", missionType)
		}
	}

	xylog.Debug(uid, "QueryUserMissions for setMissionTypesSet : %v", setMissionTypesSet)
	//没有一个可用的任务类型，玩蛋啊~~~
	if len(setMissionTypesSet) <= 0 {
		xylog.Error(uid, "[%s] invalid missionTypes %v", uid, missionTypes)
		failReason = xyerror.Resp_BadInputData.GetCode()
		err = xyerror.ErrBadInputData
		return
	}

	//日常任务
	if _, ok := setMissionTypesSet[battery.MissionType_MissionType_Daily]; ok {
		entry := new(battery.DBUserMissionEntry)
		failReason, err, entry = api.queryDailyUserMission(uid, now)
		if failReason != xyerror.Resp_NoError.GetCode() || err != xyerror.ErrOK {
			xylog.Warning(uid, "%v queryDBUserMissionsByType failed : %v", battery.MissionType_MissionType_Daily, err)
		} else {
			*entrys = append(*entrys, entry)
		}
	}

	//教学任务和主线任务，因为教学任务和主线任务是一起查询的，只要判断下是否查询教学任务就可以
	if _, ok := setMissionTypesSet[battery.MissionType_MissionType_Study]; ok {
		var studyEntry *battery.DBUserMissionEntry
		var mainLineEntry *battery.DBUserMissionEntry
		failReason, err, studyEntry, mainLineEntry = api.queryStudyAndMainLineUserMission(uid, missionCount, now)
		if failReason != xyerror.Resp_NoError.GetCode() || err != xyerror.ErrOK {
			xylog.Warning(uid, "[%s] %v queryStudyAndMainLineUserMission failed : %v", uid, battery.MissionType_MissionType_Daily, err)
		} else {
			*entrys = append(*entrys, studyEntry)
			*entrys = append(*entrys, mainLineEntry)
		}
	}

	//xylog.Debug("[%s] All missions queryDBUserMissionsByType entrys : %v", uid, entrys)
	return
}

//查询玩家日常任务
// uid string 玩家id
// now int64  当前时间戳
//return:
// entry *battery.DBUserMissionEntry 日常任务集合
func (api *XYAPI) queryDailyUserMission(uid string, now int64) (failReason battery.ErrorCode, err error, entry *battery.DBUserMissionEntry) {
	//查找当天激活的日常任务
	var (
		missionTypes      = []battery.MissionType{battery.MissionType_MissionType_Daily}
		missionStates     = make([]battery.UserMissionState, 0)
		dailyUserMissions = make([]*battery.DBUserMission, 0)
	)

	//查询玩家当天所有的日常任务
	err = api.queryDBUserMissionByTypesAndStatesAndTimestamp(uid, missionTypes, missionStates, now, &dailyUserMissions)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		xylog.Error(uid, "queryDBUserMissionByTypesAndStatesAndTimestamp failed : %v", err)
		failReason = xyerror.Resp_QueryUserMissionError.GetCode()
		return
	}

	xylog.Debug(uid, "queryDailyUserMission current : %v", dailyUserMissions)

	setActiveDailyUserMission := make(Set, 0)
	setDoneDailyUserMission := make(Set, 0)

	if len(dailyUserMissions) > 0 {
		//按照任务状态进行分类
		doneNotCollectDailyUserMissions := make([]*battery.DBUserMission, 0) //完成未收集任务
		doneCollectDailyUserMissions := make([]*battery.DBUserMission, 0)    //完成已收集任务
		activeDailyUserMissions := make([]*battery.DBUserMission, 0)         //活跃任务
		for _, dailyUserMission := range dailyUserMissions {
			switch dailyUserMission.GetState() {
			case battery.UserMissionState_UserMissionState_Active:
				activeDailyUserMissions = append(activeDailyUserMissions, dailyUserMission)
				setActiveDailyUserMission[dailyUserMission.GetOriginalmid()] = empty{}
			case battery.UserMissionState_UserMissionState_DoneNotCollect:
				doneNotCollectDailyUserMissions = append(doneNotCollectDailyUserMissions, dailyUserMission)
				setDoneDailyUserMission[dailyUserMission.GetOriginalmid()] = empty{}
			case battery.UserMissionState_UserMissionState_DoneCollected:
				doneCollectDailyUserMissions = append(doneCollectDailyUserMissions, dailyUserMission)
				setDoneDailyUserMission[dailyUserMission.GetOriginalmid()] = empty{}
			default: //do nothing
			}
		}

		dailyUserMissions = []*battery.DBUserMission{} //清空

		if len(activeDailyUserMissions) > 0 {
			//如果有激活未完成的任务，则直接返回该任务（只返回一个，玩家理论上不会有多于一个活跃的日常任务）
			dailyUserMissions = append(dailyUserMissions, activeDailyUserMissions[0])
		} else if len(doneNotCollectDailyUserMissions) > 0 {
			//如果有完成未领取的任务，则直接返回该任务（只返回一个，玩家理论上不会有多于一个完成未领取的日常任务）
			dailyUserMissions = append(dailyUserMissions, doneNotCollectDailyUserMissions[0])
		} else if len(doneCollectDailyUserMissions) >= DefConfigCache.Configs().DailyMissionCountLimit {
			//如果完成领取的任务数已经达到当天可完成的日常任务上限，则返回最后一个完成领取的任务
			finalMission := doneCollectDailyUserMissions[0]
			length := len(doneCollectDailyUserMissions)
			for i := 1; i < length; i++ {
				if finalMission.GetTimestamps()[0].GetTimestamp() < doneCollectDailyUserMissions[i].GetTimestamps()[0].GetTimestamp() {
					finalMission = doneCollectDailyUserMissions[i]
				}
			}

			dailyUserMissions = append(dailyUserMissions, finalMission)
		}
	}

	//没有可返回的日常任务，则激活一个
	if len(dailyUserMissions) == 0 {
		missionCountRest := 1
		failReason, err = api.newUserMission(uid, now, missionTypes, &setActiveDailyUserMission, setDoneDailyUserMission, &dailyUserMissions, &missionCountRest)
		if err != xyerror.ErrOK || failReason != battery.ErrorCode_NoError {
			xylog.Error(uid, "newUserMission for dailyUserMission failed : %v", err)
		}
	}

	entry = new(battery.DBUserMissionEntry)
	entry.Type = battery.MissionType_MissionType_Daily.Enum()

	for _, m := range dailyUserMissions {
		switch m.GetState() {
		case battery.UserMissionState_UserMissionState_Active:
			entry.Active = append(entry.Active, m)
		case battery.UserMissionState_UserMissionState_DoneNotCollect:
			entry.DoneNotCollect = append(entry.DoneNotCollect, m)
		case battery.UserMissionState_UserMissionState_DoneCollected:
			entry.DoneCollect = append(entry.DoneCollect, m)
		default:
		}
	}

	return
}

//查询教学任务和主线任务
// uid string 玩家id
// missionCount int 任务总数
// now int64 当前时间戳
//return:
// studyEntry *battery.DBUserMissionEntry 教学任务集合
// mainLineEntry *battery.DBUserMissionEntry 主线任务集合
func (api *XYAPI) queryStudyAndMainLineUserMission(uid string, missionCount int, now int64) (failReason battery.ErrorCode, err error, studyEntry *battery.DBUserMissionEntry, mainLineEntry *battery.DBUserMissionEntry) {
	var (
		userDoneMissions = make([]*battery.DBUserMission, 0)
		missionTypes     = []battery.MissionType{battery.MissionType_MissionType_Study, battery.MissionType_MissionType_MainLine}
		missionStates    = []battery.UserMissionState{battery.UserMissionState_UserMissionState_DoneNotCollect, battery.UserMissionState_UserMissionState_DoneCollected}
	)

	studyEntry = &battery.DBUserMissionEntry{
		Type: battery.MissionType_MissionType_Study.Enum(),
	}

	mainLineEntry = &battery.DBUserMissionEntry{
		Type: battery.MissionType_MissionType_MainLine.Enum(),
	}

	//已完成未领取教学/主线任务
	//查询用户已完成的任务（包含已完成未领取和已完成已领取任务）
	err = api.queryDBUserMissionByTypesAndStatesAndTimestamp(uid, missionTypes, missionStates, 0, &userDoneMissions)
	if err != xyerror.ErrOK {
		return
	}

	setDoneMissionId := make(Set, 0)
	for _, v := range userDoneMissions {
		state := v.GetState()
		id := v.GetOriginalmid()
		setDoneMissionId[id] = empty{}
		switch state {
		case battery.UserMissionState_UserMissionState_DoneNotCollect:
			switch v.GetType() {
			case battery.MissionType_MissionType_Study:
				studyEntry.DoneNotCollect = append(studyEntry.DoneNotCollect, v)
			case battery.MissionType_MissionType_MainLine:
				mainLineEntry.DoneNotCollect = append(mainLineEntry.DoneNotCollect, v)
			}
		default:
			//do nothing
		}
	}

	//完成未领取的任务数
	var missionCountRest = missionCount // - len(studyEntry.DoneNotCollect) + len(mainLineEntry.DoneNotCollect)
	if missionCountRest < 0 {
		xylog.Warning(uid, "[%s] missionDoneNotCollectCount(%d) < missionCount(%d), are you kidding?", uid, len(studyEntry.DoneNotCollect)+len(mainLineEntry.DoneNotCollect), missionCount)
		return
	}

	setActiveMissionId := make(Set, 0)
	if missionCountRest > 0 {
		//---------- 先统一处理一下，查找活跃任务，把过期的任务设置一下状态 ----------
		missionStates = []battery.UserMissionState{battery.UserMissionState_UserMissionState_Active}
		var (
			userActiveMissions []*battery.DBUserMission
			activeMissions     []*battery.DBUserMission
		)
		err = api.queryDBUserMissionByTypesAndStatesAndTimestamp(uid, missionTypes, missionStates, 0, &userActiveMissions)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "queryDBUserMissionByTypesAndStatesAndTimestamp failed : %v", err)
			failReason = xyerror.Resp_QueryUserMissionError.GetCode()
			return
		}

		xylog.Debug(uid, "%v userActiveMissions before %d : %v", missionTypes, len(userActiveMissions), userActiveMissions)

		if len(userActiveMissions) > 0 {
			expiredMissions := make([]uint64, 0)
			api.checkUserMissions(uid, userActiveMissions, &activeMissions, &expiredMissions, now)
			xylog.Debug(uid, "%v ActiveMissions after check %d : %v\nExpiredMissions after check %d: %v", missionTypes, len(activeMissions), activeMissions, len(expiredMissions), expiredMissions)
			//刷新过期任务的状态
			api.setUserMissionState(uid, expiredMissions, battery.UserMissionState_UserMissionState_Expired, now)

			//保存一下活跃任务
			for _, m := range activeMissions {
				id := m.GetOriginalmid()
				setActiveMissionId[id] = empty{}

				switch m.GetType() {
				case battery.MissionType_MissionType_Study:
					studyEntry.Active = append(studyEntry.Active, m)
				case battery.MissionType_MissionType_MainLine:
					mainLineEntry.Active = append(mainLineEntry.Active, m)
				default:
					//do nothing
				}
			}

			missionCountRest -= len(studyEntry.Active) + len(mainLineEntry.Active)
			if missionCountRest < 0 {
				xylog.Warning(uid, "[%s] missionActiveCount(%d) > missionCountRest(%d), are you kidding?", uid, len(studyEntry.Active)+len(mainLineEntry.Active), len(studyEntry.Active)+len(mainLineEntry.Active)+missionCountRest)
				return
			}
		}

		//---------- 激活教学任务 ---------- todelete 没有教学任务概念
		//if missionCountRest > 0 {
		//	//查找任务信息
		//	missionTypes = []battery.MissionType{battery.MissionType_MissionType_Study}
		//	failReason, err = api.newUserMission(uid, now, missionTypes, &setActiveMissionId, setDoneMissionId, &(studyEntry.Active), &missionCountRest)
		//	if failReason != xyerror.Resp_NoError.GetCode() || err != xyerror.ErrOK {
		//		return
		//	}
		//}

		//---------- 激活主线任务 ----------
		if missionCountRest > 0 {
			//查找玩家已完成的主线任务列表
			ids := make([]uint64, 0)
			ids, err = api.queryDBUserDoneCollectionMission(uid, battery.MissionType_MissionType_MainLine)
			for _, id := range ids {
				setDoneMissionId[id] = empty{}
			}

			xylog.Debug(uid, "setDoneMissionId : %v", setDoneMissionId)

			//查找任务信息
			missionTypes = []battery.MissionType{battery.MissionType_MissionType_MainLine}
			failReason, err = api.newUserMission(uid, now, missionTypes, &setActiveMissionId, setDoneMissionId, &(mainLineEntry.Active), &missionCountRest)
			if failReason != xyerror.Resp_NoError.GetCode() || err != xyerror.ErrOK {
				return
			}
		}
	}

	return
}

//生成玩家任务
// uid string 玩家id
// now int64 当前时间戳
// missionTypes []battery.MissionType 任务类型列表
// setActiveMissionId *Set 活跃任务列表
// setDoneMissionId Set 完成任务列表
// activeMissions *[]*battery.DBUserMission 活跃任务列表
// missionCountRest *int 剩余任务数
func (api *XYAPI) newUserMission(uid string, now int64, missionTypes []battery.MissionType, setActiveMissionId *Set, setDoneMissionId Set, activeMissions *[]*battery.DBUserMission, missionCountRest *int) (failReason battery.ErrorCode, err error) {
	mapMissionItems := make(map[uint32][]*battery.MissionItem, 0)
	var prioritySlice sort.IntSlice
	xylog.Debug(uid, "queryMission for %v", missionTypes)

	err = api.queryMission(missionTypes, &mapMissionItems, &prioritySlice, now)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "queryMission failed : %v", err)
		failReason = xyerror.Resp_QueryMissionError.GetCode()
		return
	} else {
		//xylog.Debug("[%s] queryMission result : %v", uid, mapMissionItems)
	}

LOOP1:
	//按照优先级查找可激活任务
	for _, priority := range prioritySlice {
		if missionItems, ok := mapMissionItems[uint32(priority)]; ok {
			for _, m := range missionItems {

				//任务数目够了，就跳出
				if *missionCountRest == 0 {
					break LOOP1
				}

				id := m.GetId()
				//已经存在激活关联任务，跳过
				if _, ok := (*setActiveMissionId)[id]; ok {
					xylog.Debug(uid, "mission %d already exist in active user mission. skip it.", id)
					continue
				}

				//如果任务已经完成，跳过
				if _, ok := setDoneMissionId[id]; ok {
					xylog.Debug(uid, "mission %d already exist in done user mission. skip it.", id)
					continue
				}

				//判断任务的前置条件是否完成
				if !(api.isMissionReady(uid, setDoneMissionId, m)) {
					continue
				}

				userMission := api.getDBUserMissionFromMissionItem(uid, m, now)
				*activeMissions = append(*activeMissions, userMission)
				xylog.Debug(uid, "addUserMission : %v", userMission)
				err = api.addUserMission(userMission)
				if err != nil {
					xylog.Error(uid, "addUserMission failed : %v", err)
					continue
				}

				//设置任务已激活
				(*setActiveMissionId)[id] = empty{}

				//剩余数减1
				(*missionCountRest)--
			}
		}
	}

	return
}

//调整日常任务的指标时间
// mission *battery.DBUserMission  任务指针
// now int64 当前时间戳
func (api *XYAPI) setDailyExpireTimeStamp(mission *battery.DBUserMission, now int64) {
	var (
		nowTmp           = time.Unix(now, 0)
		year, month, day = nowTmp.Date()
		hour             = nowTmp.Hour()
		timestamp        int64
		quotaNum         int
	)

	hour = DefConfigCache.Configs().DailyMissionRefreshHour
	if hour >= DefConfigCache.Configs().DailyMissionRefreshHour { //8点后（可配置，参照配置项"dailymissionrefreshhour"）领到的任务，过期时间为第二天8点
		now += OneDaySeconds
		nowTmp = time.Unix(now, 0)
		year, month, day = nowTmp.Date()
	}

	timestamp = time.Date(year, month, day, hour, 0, 0, 0, time.Local).Unix()

	quotaNum = len(mission.GetQuotaStats())
	for j := 0; j < quotaNum; j++ {
		mission.QuotaStats[j].CycleType = battery.QuotaCycleType_QuotaCycleType_DAY.Enum()
		mission.QuotaStats[j].CycleValue = proto.Int64(timestamp)
	}
}

//查询当前玩家激活的任务列表
// uid string 玩家id
// missionTypes []battery.MissionType 任务类型列表
//return:
// err error 错误
// missions []*battery.DBUserMission 查到的任务列表
func (api *XYAPI) queryDBActiveUserMissions(uid string, missionTypes []battery.MissionType, now int64) (err error, dbUserMissions []*battery.DBUserMission) {
	states := []battery.UserMissionState{battery.UserMissionState_UserMissionState_Active}
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERMISSION).QueryUserMissionByTypesAndStatesAndTimestamp(uid, missionTypes, states, now, &dbUserMissions)
	return
}

//通过数据库查询的玩家任务获取消息的玩家任务结构体
func (api *XYAPI) getUserMissionFromDBUserMission(uid string, dbUserMission *battery.DBUserMission) (userMission *battery.UserMission) {
	originalMid := dbUserMission.GetOriginalmid()
	config := xycache.DefMissionCacheManager.Mission(originalMid)
	if nil == config {
		xylog.Error(uid, "get missionconfig of %d failed", originalMid)
		return nil
	}

	userMission = &battery.UserMission{
		Mid:         proto.Uint64(dbUserMission.GetMid()),
		Originalmid: proto.Uint64(originalMid),
		Rewards:     config.GetRewards(),
		Tipid:       config.Tipid,
		TipDesc:     config.TipDesc,
	}

	for _, q := range dbUserMission.GetQuotaStats() {
		quota := &battery.MissionQuotaStat{
			Id:        q.GetId().Enum(),
			GoalValue: proto.Uint64(q.GetGoalValue()),
			DoneValue: proto.Uint64(q.GetDoneValue()),
		}
		userMission.QuotaStats = append(userMission.QuotaStats, quota)
	}
	return
}

//根据任务配置信息生成玩家任务信息
func (api *XYAPI) getDBUserMissionFromMissionItem(uid string, missionItem *battery.MissionItem, now int64) (dbUserMission *battery.DBUserMission) {
	dbUserMission = &battery.DBUserMission{
		Uid:         proto.String(uid),
		Type:        missionItem.GetType().Enum(),
		Mid:         proto.Uint64(DefMissionIdGenerater.NewID()),
		Originalmid: proto.Uint64(missionItem.GetId()),
		Timestamp:   proto.Int64(0),
	}

	dbUserMission.QuotaStats = make([]*battery.DBMissionQuotaStat, 0)
	quotas := missionItem.GetQuotas()
	for _, quota := range quotas {
		quotaStat := &battery.DBMissionQuotaStat{
			Id:              quota.GetId().Enum(),
			CycleType:       quota.GetCycleType().Enum(),
			GoalValue:       proto.Uint64(quota.GetGoalValue()),
			DoneValue:       proto.Uint64(0),
			ConstraintEqual: proto.Bool(quota.GetConstraintEqual()),
		}

		//根据周期类型调整周期值
		quotaStat.CycleValueLimit = proto.Int64(quota.GetCycleValue()) //局数目标值设置
		switch quota.GetCycleType() {
		case battery.QuotaCycleType_QuotaCycleType_ROUND: //如果是局数周期值，直接保存局数
			quotaStat.CycleValue = proto.Int64(0) //局数当前值初始化为0
		case battery.QuotaCycleType_QuotaCycleType_SECOND:
			quotaStat.CycleValue = proto.Int64(now + quota.GetCycleValue())
		default: // battery.QuotaCycleType_QuotaCycleType_DAY ，默认为天周期，转换为对应的自然日0点时间
			quotaStat.CycleValue = proto.Int64(api.getQuotaStatTimestamp(now, quota.GetCycleValue()))
		}

		dbUserMission.QuotaStats = append(dbUserMission.QuotaStats, quotaStat)
	}

	if missionItem.GetType() == battery.MissionType_MissionType_Daily {
		//日常任务的过期时间在第二天的8:00，任务指标过期时间重设一下
		api.setDailyExpireTimeStamp(dbUserMission, now)
	}

	timestamp := &battery.UserMissionTimestamp{
		State:     battery.UserMissionState_UserMissionState_Active.Enum(),
		Timestamp: proto.Int64(now),
	}

	dbUserMission.Timestamps = make([]*battery.UserMissionTimestamp, 0)
	dbUserMission.Timestamps = append(dbUserMission.Timestamps, timestamp)
	dbUserMission.State = battery.UserMissionState_UserMissionState_Active.Enum()
	return
}

//生成默认的玩家任务信息
// uid string 玩家id
// id uint64 任务id
// missionType battery.MissionType 任务类型列表
// autoCollect bool 是否自动领取
// now int64 当前时间戳
func (api *XYAPI) defaultDBUserMission(uid string, id uint64, missionType battery.MissionType, autoCollect bool, now int64) (userMission *battery.DBUserMission) {
	userMission = &battery.DBUserMission{
		Uid:         proto.String(uid),
		Type:        missionType.Enum(),
		Mid:         proto.Uint64(DefMissionIdGenerater.NewID()),
		Originalmid: proto.Uint64(id),
		State:       battery.UserMissionState_UserMissionState_Active.Enum(),
		Timestamp:   proto.Int64(0),
	}

	timestamp := &battery.UserMissionTimestamp{
		State:     battery.UserMissionState_UserMissionState_Active.Enum(),
		Timestamp: proto.Int64(now),
	}

	userMission.Timestamps = make([]*battery.UserMissionTimestamp, 0)
	userMission.Timestamps = append(userMission.Timestamps, timestamp)

	return
}

//刷新玩家的任务指标状态
// uid string 玩家id
// missionTypes []battery.MissionType 任务类型
// quotas []*battery.Quota 指标列表
// now int64 当前时间戳
// isFinish bool 记忆点是否完成
func (api *XYAPI) updateUserMissionsQuotas(uid string, missionTypes []battery.MissionType, quotas []*battery.Quota, now int64, isFinish bool) (err error) {

	xylog.Debug(uid, "quotas : %v", quotas)

	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_UPDATEUSERMISSIONSTAT, &begin)

	//查询用户当前激活的任务列表
	activeUserMissions := []*battery.DBUserMission{}
	//在根据指标更新任务状态的时候，时间戳设置为0，查到所有当前活跃的任务
	err, activeUserMissions = api.queryDBActiveUserMissions(uid, missionTypes, 0)
	if err != xyerror.ErrOK {
		xylog.Debug(uid, "queryDBActiveUserMissions of %v failed : %v", missionTypes, err)
		return
	}

	xylog.Debug(uid, "queryDBActiveUserMissions result : %v", activeUserMissions)

	//针对每个任务更新指标值和任务状态
	doneCollectedIds := make(map[battery.MissionType][]uint64, 0)
	for _, mission := range activeUserMissions {
		api.updateUserMissionQuotas(uid, mission, quotas, now, &doneCollectedIds, isFinish)
	}

	//xylog.Debug("[%s] doneCollectedIds : %v", uid, doneCollectedIds)

	err = api.appendDoneCollected(uid, doneCollectedIds)

	return
}

//刷新玩家任务状态
// uid string 玩家id
// mission *battery.DBUserMission 玩家任务指针
// quotas []*battery.Quota 指标集合
// now int64 当前时间戳
func (api *XYAPI) updateUserMissionQuotas(uid string, mission *battery.DBUserMission, quotas []*battery.Quota, now int64, doneCollectedIds *map[battery.MissionType][]uint64, isFinish bool) (err error) {

	xylog.Debug(uid, "updateUserMissionQuotas[%v].", quotas)

	//修改一下指标值
	for _, quota := range quotas {

		id, value := quota.GetId(), quota.GetValue()

		//如果指标值小于等于0，直接跳过
		if value <= 0 {
			xylog.Debug(uid, "quota[%v].Value[%d] <= 0, skip.", id, value)
			continue
		}

		//如果指标不需要记录，直接跳过
		if !api.isNeedToBeCarved(id, isFinish) {
			xylog.Debug(uid, "quota[%d] no need to be carved.", id)
			continue
		}

		length := len(mission.QuotaStats)
		for i := 0; i < length; i++ {
			if id == mission.QuotaStats[i].GetId() {
				//指标值累加
				*(mission.QuotaStats[i].DoneValue) += value

				//如果指标周期是局，把指标值加一
				if mission.QuotaStats[i].GetCycleType() == battery.QuotaCycleType_QuotaCycleType_ROUND {
					*(mission.QuotaStats[i].CycleValue) += 1
				}
				break
			}
		}
	}

	//判断是否完成，如果已经完成设置一下任务状态，并记录下时间戳
	if api.isUserMissionComplete(mission) {
		//先加上doneNotCollect状态信息
		userMissionTimestamp := &battery.UserMissionTimestamp{
			State:     battery.UserMissionState_UserMissionState_DoneNotCollect.Enum(),
			Timestamp: proto.Int64(now),
		}
		mission.Timestamps = append(mission.Timestamps, userMissionTimestamp)

		//根据任务的AutoCollect标识，进行任务完成的奖励发放操作
		originalMid := mission.GetOriginalmid()
		config := xycache.DefMissionCacheManager.Mission(originalMid)
		if nil == config {
			xylog.Error(uid, "get missionconfig of %d failed", originalMid)
		} else {
			//xylog.Debug("[%s] missionconfig : %v", uid, config)
			if config.GetAutoCollect() { //如果任务是自动领取奖励，则开始发放奖励
				//设置doneCollect状态
				userMissionTimestamp := &battery.UserMissionTimestamp{
					State:     battery.UserMissionState_UserMissionState_DoneCollected.Enum(),
					Timestamp: proto.Int64(now),
				}
				mission.Timestamps = append(mission.Timestamps, userMissionTimestamp)
				mission.State = battery.UserMissionState_UserMissionState_DoneCollected.Enum()
				mission.Timestamp = proto.Int64(now)

				xylog.Debug(uid, "before updateUserMission : %v", mission)

				//发放奖励
				for _, propItem := range config.GetRewards() {
					api.GainProp(uid, nil, propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
				}

				// 把任务添加到donecollected列表中
				missionType := mission.GetType()
				if _, ok := (*doneCollectedIds)[missionType]; !ok {
					(*doneCollectedIds)[missionType] = make([]uint64, 0)
				}
				(*doneCollectedIds)[missionType] = append((*doneCollectedIds)[missionType], originalMid)

			} else { //如果任务是不是自动领取奖励，则只设置任务状态
				mission.State = battery.UserMissionState_UserMissionState_DoneNotCollect.Enum()
			}
		}
	}

	//校验一下任务状态，对于过期需要再重启的，重启一下
	failReason, _, isUpdate := api.checkUserMission(uid, mission, now)
	if failReason != battery.ErrorCode_NoError {
		return
	}

	//没刷新过任务状态，就刷新一下
	if !isUpdate {
		err = api.updateUserMission(mission)
	}

	return
}

//获取指标的过期时间戳
// now 当前时间戳
// dayCount 指标完成天数限制
func (api *XYAPI) getQuotaStatTimestamp(now, dayCount int64) int64 {
	secs := now + dayCount*OneDaySeconds
	timeTmp := time.Unix(secs, 0)
	year, month, day := timeTmp.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
}

//判断任务是否就绪
// uid string 玩家id
// setUserCompleteMissionId Set 已经完成任务id列表
// mission *battery.MissionItem 任务指针
func (api *XYAPI) isMissionReady(uid string, setUserCompleteMissionId Set, mission *battery.MissionItem) bool {

	//前置任务是否完成，如果有一个未完成，该任务就无法激活
	relatedMissions := mission.GetRelatedMissions()
	for _, v := range relatedMissions {
		if _, ok := setUserCompleteMissionId[v]; !ok {
			xylog.Debug(uid, "mission %d no ready because of mission %d no done", mission.GetId(), v)
			return false
		}
	}

	//前置道具是否拥有，如果有一个未拥有，该任务就无法激活
	for _, v := range mission.GetRelatedProps() {
		if !api.OwnProp(uid, v) {
			xylog.Debug(uid, "mission %d no ready because of prop %d no exist", mission.GetId(), v)
			return false
		}
	}

	return true
}

//玩家任务是否完成
// mission *battery.DBUserMission 任务指针
func (api *XYAPI) isUserMissionComplete(mission *battery.DBUserMission) bool {

	for _, missionQuota := range mission.GetQuotaStats() {
		done, goal := missionQuota.GetDoneValue(), missionQuota.GetGoalValue()
		if (missionQuota.GetConstraintEqual() && done != goal) || //不符合严格一致要求
			done < goal { //非严格一致要求，但是达到目标值
			return false
		}
	}

	return true
}

//查询玩家doneCollectedMission信息
// uid string 玩家id
// missionType battery.MissionType 任务类型
func (api *XYAPI) queryDBUserDoneCollectionMission(uid string, missionType battery.MissionType) (ids []uint64, err error) {
	var userDoneCollectedMission *battery.DBUserDoneCollectedMission
	userDoneCollectedMission, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERDONECOLLECTIONMISSION).QueryUserDoneCollectedMission(uid, missionType)

	if err == xyerror.ErrOK {
		ids = userDoneCollectedMission.GetIds()
		xylog.Debug(uid, "userDoneCollectedMission : %v", userDoneCollectedMission)
	}

	return
}

//根据类型+状态+时间戳查找玩家任务
// uid string 玩家id
// missionTypes []battery.MissionType 任务类型
// missionStates []battery.UserMissionState  任务状态
// now int64 任务时间戳，例如：输入为2014-09-23的时间戳，则查询的任务为2014-09-23当天激活的任务（主要用于日常任务的查询条件，教学任务和主线任务填0即可）
// userMissions *[]*battery.DBUserMission 输出参数，查找到的玩家任务
//return:
// err error
func (api *XYAPI) queryDBUserMissionByTypesAndStatesAndTimestamp(uid string, missionTypes []battery.MissionType, missionStates []battery.UserMissionState, now int64, userMissions *[]*battery.DBUserMission) (err error) {
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERMISSION).QueryUserMissionByTypesAndStatesAndTimestamp(uid, missionTypes, missionStates, now, userMissions)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "QueryUserMissionByTypesAndStates : %v", err)
	}
	return
}

//查询任务的详细信息
// uid string 玩家id
// missionId uint64 任务id
// userMission *battery.DBUserMission 玩家任务保存指针
func (api *XYAPI) queryUserMissionDetail(uid string, missionId uint64, userMission *battery.DBUserMission) (err error) {
	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_QUERYUSERMISSIONDETAIL, &begin)

	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERMISSION).QueryUserMissionDetail(uid, missionId, userMission)
	return
}

func (api *XYAPI) addUserMission(userMission *battery.DBUserMission) (err error) {
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERMISSION).AddUserMission(userMission)
}

func (api *XYAPI) updateUserMission(userMission *battery.DBUserMission) (err error) {
	begin := time.Now()
	defer xyperf.Trace(LOGTRACE_UPDATEUSERMISSION, &begin)
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERMISSION).UpdateUserMission(userMission)
}

//校验单个任务
// userMission *battery.DBUserMission 任务指针
// now int64 当前的时间戳
func (api *XYAPI) checkUserMission(uid string, mission *battery.DBUserMission, now int64) (failReason battery.ErrorCode, isExpired bool, isUpdate bool) {
	failReason, isExpired, isUpdate = battery.ErrorCode_NoError, false, false

	originalMid := mission.GetOriginalmid()
	config := xycache.DefMissionCacheManager.Mission(originalMid)
	if nil == config {
		xylog.Error(uid, "get missionconfig of %d failed", originalMid)
		failReason = battery.ErrorCode_QueryMissionError
		return
	}

	quotaStats := mission.GetQuotaStats()
	if nil != quotaStats {
		for _, stat := range quotaStats {
			cycleType := stat.GetCycleType()
			cycleValue := stat.GetCycleValue()
			done, goal := stat.GetDoneValue(), stat.GetGoalValue()
			switch cycleType {
			case battery.QuotaCycleType_QuotaCycleType_DAY, battery.QuotaCycleType_QuotaCycleType_SECOND:
				if now > cycleValue { //时间过期
					isExpired = true
				}
			case battery.QuotaCycleType_QuotaCycleType_ROUND:
				if (cycleValue >= stat.GetCycleValueLimit() && done < goal) || //局数到达指标限定值且指标未完成
					(cycleValue >= stat.GetCycleValueLimit() && stat.GetConstraintEqual() && done != goal) { //局数到达指标限定值且匹配
					isExpired = true
				}
			default:
				//QuotaCycleType_QuotaCycleType_UNKOWN , do nothing如果任务的周期类型是未知，则什么都不做，任务一直存在直到完成
			}
		}
	}

	if isExpired {
		if config.GetExpiredRestart() { //任务过期重启
			xylog.Debug(uid, "ExpiredRestart mission : %v", mission)
			api.resetUserMission(mission, now) //重置任务
			//更新到数据库
			err := api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERMISSION).UpdateUserMission(mission)
			if err != xyerror.ErrOK { //更新失败，处理下一条
				xylog.Error(uid, "UpdateUserMission %v failed : %v", *mission, err)
				failReason = battery.ErrorCode_UpdateUserMissionError
				return
			}

			isUpdate = true
		}
	}

	return
}

//校验任务，活跃的任务保存在activeMission中，过期的任务保存在expiredMissions中
// userMissions []*battery.DBUserMission 所有的玩家任务
// activeMission *[]*battery.DBUserMission 活跃的玩家任务
// expiredMissions *[]uint64 过期的任务id列表
// now int64 当前的时间戳
func (api *XYAPI) checkUserMissions(uid string, userMissions []*battery.DBUserMission, activeMission *[]*battery.DBUserMission, expiredMissions *[]uint64, now int64) {

	for _, mission := range userMissions {

		//校验下任务的状态，如果需要重启的，就重启下
		failReason, isExpired, _ := api.checkUserMission(uid, mission, now)
		if failReason != battery.ErrorCode_NoError {
			continue
		}

		if !isExpired { //未过期的仍旧保留
			*activeMission = append(*activeMission, mission)
		} else {
			//过期的保存下来用于刷新数据库中的任务状态
			*expiredMissions = append(*expiredMissions, mission.GetMid())
		}
	}

}

//重置玩家任务
func (api *XYAPI) resetUserMission(mission *battery.DBUserMission, now int64) {
	quotaStatsLen := len(mission.GetQuotaStats())
	for i := 0; i < quotaStatsLen; i++ {
		switch mission.QuotaStats[i].GetCycleType() {
		case battery.QuotaCycleType_QuotaCycleType_ROUND: //如果是局为统计周期单位，则将已玩局数设置为0
			mission.QuotaStats[i].CycleValue = proto.Int64(0)
		case battery.QuotaCycleType_QuotaCycleType_SECOND: //如果是秒为统计周期单位，则当前时间加上对应秒数
			mission.QuotaStats[i].CycleValue = proto.Int64(now + mission.QuotaStats[i].GetCycleValueLimit())
		default: // battery.QuotaCycleType_QuotaCycleType_DAY ，默认为天周期，转换为对应的自然日0点时间
			mission.QuotaStats[i].CycleValue = proto.Int64(api.getQuotaStatTimestamp(now, mission.QuotaStats[i].GetCycleValueLimit()))
		}
		mission.QuotaStats[i].DoneValue = proto.Uint64(0)
	}
	timestamp := &battery.UserMissionTimestamp{
		State:     battery.UserMissionState_UserMissionState_Reset.Enum(),
		Timestamp: proto.Int64(now),
	}
	mission.Timestamps = append(mission.Timestamps, timestamp)
	mission.State = battery.UserMissionState_UserMissionState_Active.Enum()
}

//设置玩家任务状态
// uid string 玩家id
// missions []uint64 任务列表
// state battery.UserMissionState 任务状态
// now int64 当前时间戳
func (api *XYAPI) setUserMissionState(uid string, missions []uint64, state battery.UserMissionState, now int64) (err error) {

	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERMISSION).UpdateUserMissionState(uid, missions, state, now)
	if err != nil {
		xylog.Error(uid, "UpdateUserMissionState %v to %v failed : %v", missions, state, err)
		return
	}
	return
}

//查询任务信息，按照优先级做好排序的map
// missionTypes []battery.MissionType 任务类型列表
// now int64 当前时间戳
//return:
// mapMissionItems *map[uint32][]*battery.MissionItem 保存任务信息的map
// prioritySlice *sort.IntSlice 优先级列表
func (api *XYAPI) queryMission(missionTypes []battery.MissionType, mapMissionItems *map[uint32][]*battery.MissionItem, prioritySlice *sort.IntSlice, now int64) (err error) {

	missionItems := xycache.DefMissionCacheManager.TypeMissionItems(missionTypes, now)

	for _, missionItem := range missionItems {
		priority := missionItem.GetPriority()
		if (*mapMissionItems)[priority] == nil {
			(*mapMissionItems)[priority] = make([]*battery.MissionItem, 0)
			*prioritySlice = append(*prioritySlice, int(priority))
		}

		(*mapMissionItems)[priority] = append((*mapMissionItems)[priority], missionItem)
	}

	sort.Sort(*prioritySlice)

	return
}

//根据任务类型和当前时间戳查询当前可激活的任务信息
// missionTypes []battery.MissionType 任务类型列表
// now int64 当前时间戳
//return:
// missionItems *[]*battery.MissionItem
func (api *XYAPI) queryMissionSlice(missionTypes []battery.MissionType, missionItems *[]*battery.MissionItem, now int64) {
	*missionItems = xycache.DefMissionCacheManager.TypeMissionItems(missionTypes, now)
}

//确认任务消息
func (api *XYAPI) OperationConfirmUserMission(req *battery.ConfirmUserMissionRequest, resp *battery.ConfirmUserMissionResponse) (err error) {
	var (
		uid                  = req.GetUid()
		missionType          = req.GetType()
		mid                  = req.GetMid()
		now                  = time.Now().Unix()
		doneCollectedIds     = make(map[battery.MissionType][]uint64, 0)
		originalMid          uint64
		config               *battery.MissionItem
		userMissionTimestamp *battery.UserMissionTimestamp
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//初始化resp
	resp.Uid = req.Uid
	resp.Type = req.Type
	resp.Mid = req.Mid
	resp.Error = xyerror.DefaultError()

	//查询任务信息是否存在
	userMission := new(battery.DBUserMission)
	err = api.queryUserMissionDetail(uid, mid, userMission)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "queryUserMission req (%v) failed : %v", req, err)
		resp.Error.Code = battery.ErrorCode_QueryUserMissionError.Enum()
		goto ErrHandle
	}

	//校验一下任务
	if !api.isUserMissionValid(uid, mid, missionType, battery.UserMissionState_UserMissionState_DoneNotCollect, userMission) {
		xylog.Error(uid, "isUserMissionValid %v : false", userMission)
		resp.Error.Code = battery.ErrorCode_QueryUserMissionError.Enum()
		goto ErrHandle
	}

	//发放奖品
	originalMid = userMission.GetOriginalmid()
	config = xycache.DefMissionCacheManager.Mission(originalMid)
	if nil == config {
		xylog.Error(uid, "get missionconfig of %d failed", originalMid)
		resp.Error.Code = battery.ErrorCode_QueryMissionError.Enum()
		goto ErrHandle
	}

	for _, propItem := range config.GetRewards() {
		api.GainProp(uid, nil, propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
	}

	//修改任务状态
	//err = api.setUserMissionState(uid, missions, battery.UserMissionState_UserMissionState_DoneCollected, now)
	//if err != xyerror.ErrOK {
	//	xylog.Error("[%s] setUserMissionState(%v) failed", uid, missions)
	//	goto ErrHandle
	//}

	//将任务id增加到donecollected列表
	doneCollectedIds[missionType] = []uint64{originalMid}
	err = api.appendDoneCollected(uid, doneCollectedIds)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "appendDoneCollected(%v) failed", doneCollectedIds)
		goto ErrHandle
	}

	//设置doneCollect状态
	userMissionTimestamp = &battery.UserMissionTimestamp{
		State:     battery.UserMissionState_UserMissionState_DoneCollected.Enum(),
		Timestamp: proto.Int64(now),
	}
	userMission.Timestamps = append(userMission.Timestamps, userMissionTimestamp)
	userMission.State = battery.UserMissionState_UserMissionState_DoneCollected.Enum()
	userMission.Timestamp = proto.Int64(now)

	err = api.updateUserMission(userMission)

ErrHandle:

	return
}

//任务增加到doneCollected列表
// uid string 玩家id
// missionType battery.MissionType 任务类型
// doneCollectedIds []uint64 donecollected任务id列表
func (api *XYAPI) appendDoneCollected(uid string, doneCollectedIds map[battery.MissionType][]uint64) (err error) {

	//xylog.Debug("[%s] doneCollectedIds : %v", uid, doneCollectedIds)

	for missionType, ids := range doneCollectedIds {

		if missionType == battery.MissionType_MissionType_Daily { //日常任务不需要保存
			continue
		}

		var userDoneCollectedMission *battery.DBUserDoneCollectedMission
		userDoneCollectedMission, err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERDONECOLLECTIONMISSION).QueryUserDoneCollectedMission(uid, missionType)
		if err != xyerror.ErrOK {
			if err == xyerror.ErrNotFound { //如果没有找到，就创建一条
				userDoneCollectedMission = api.defaultDBUserDoneCollectedMission(uid, missionType)
				//userDoneCollectedMission.Ids = ids
			} else { //数据库错误
				xylog.Error(uid, "QueryUserDoneCollectedMission failed : %v", err)
				return
			}
		}
		userDoneCollectedMission.Ids = append(userDoneCollectedMission.Ids, ids...)

		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERDONECOLLECTIONMISSION).UpsertUserDoneCollectedMission(userDoneCollectedMission)
		if err != xyerror.ErrOK {
			xylog.Error(uid, "UpsertUserDoneCollectedMission failed : %v", err)
			return
		}

		xylog.Debug(uid, "userDoneCollectedMission : %v", userDoneCollectedMission)
	}
	return
}

//生成默认的玩家DoneCollectedMission信息列表
// uid string 玩家id
// missionType battery.MissionType 任务类型
func (api *XYAPI) defaultDBUserDoneCollectedMission(uid string, missionType battery.MissionType) *battery.DBUserDoneCollectedMission {
	return &battery.DBUserDoneCollectedMission{
		Uid:  proto.String(uid),
		Type: missionType.Enum(),
		Ids:  []uint64{},
	}
}

//校验玩家任务是否是合法的
// uid string 玩家id
// mid uint64 任务id
// missionType battery.MissionType 任务类型
// state battery.UserMissionState 任务状态
// userMission *battery.DBUserMission 玩家任务信息指针
//return:
// true 合法，false 非法
func (api *XYAPI) isUserMissionValid(uid string, mid uint64, missionType battery.MissionType, state battery.UserMissionState, userMission *battery.DBUserMission) bool {
	return userMission.GetUid() == uid && userMission.GetMid() == mid && userMission.GetType() == missionType && userMission.GetState() == state
}

//增加指标到指标集
// id battery.QuotaEnum指标id
// value uint64 指标值
// quotas *[]*battery.Quota 指标集合
func (api *XYAPI) addQuota(id battery.QuotaEnum, value uint64, quotas *[]*battery.Quota) {
	*quotas = append(*quotas, &battery.Quota{
		Id:    id.Enum(),
		Value: proto.Uint64(value),
	})
}

//判断任务指标是否需要跳过
// id battery.QuotaEnum 任务指标id
// isFinish bool 记忆点是否完成
//return:
// true  指标需要记录
// false 指标不需要记录
func (api *XYAPI) isNeedToBeCarved(id battery.QuotaEnum, isFinish bool) bool {
	if !isFinish { //记忆点未完成，需要以记忆点完成为前提的指标跳过
		for _, quotaId := range api.Config.Configs().QuotasNeedFinish {
			if quotaId == int32(id) {
				return false
			}
		}
	}
	return true
}
