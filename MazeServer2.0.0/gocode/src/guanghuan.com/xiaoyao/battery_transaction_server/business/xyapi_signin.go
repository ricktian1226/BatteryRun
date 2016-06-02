// xyapi_signin
package batteryapi

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/cache"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
	"time"
)

//处理旧版本签到查询请求
func (api *XYAPI) OperationQuerySignIn(req *battery.QuerySignInRequest, resp *battery.QuerySignInResponse) (err error) {
	var (
		uid                = req.GetUid()
		now                int64
		activitys          []*battery.DBUserSignInActivity
		activeActivitys    []*battery.DBUserSignInActivity
		setExceptActivitys = make(Set, 0)
		signInActivitys    = make([]*battery.SignInActivity, 0)
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//初始化返回
	resp.Error = xyerror.DefaultError()

	//xylog.Debug(uid, "[%s] QuerySignIn request begin", uid)
	//defer xylog.Debug(uid, "[%s] QuerySignIn request end", uid)

	//查询用户签到活动信息
	activitys, err = api.queryUserSignInActivitys(uid)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			xylog.Debug(uid, "UserSignInActivitys nofound.")
		} else {
			xylog.Error(uid, "UserSignInActivitys failed : %v", err)
			resp.Error.Code = battery.ErrorCode_QueryUserSignInRecordError.Enum()
			goto ErrHandle
		}
	}

	//查找需要添加的签到活动
	for _, activity := range activitys {
		setExceptActivitys[activity.GetId()] = empty{}
		if activity.GetState() == battery.UserMissionState_UserMissionState_Active { //正在进行的活动，加入返回列表中
			activeActivitys = append(activeActivitys, activity)
		}
	}

	now = time.Now().Unix()
	signInActivitys, err = api.querySignInActivity(uid, now, setExceptActivitys)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "querySignInActivity from cache failed : %v", err)
		resp.Error.Code = battery.ErrorCode_QuerySignInActivitysFromCacheError.Enum()
		goto ErrHandle
	}

	for _, signInActivity := range signInActivitys {
		newUserSignInActivity := api.defaultUserSignInActivity(uid, now, signInActivity)
		err = api.addUserSignInActivity(newUserSignInActivity)
		if err == xyerror.ErrOK { //插入正常，返回消息中也加上新的签到活动
			activeActivitys = append(activeActivitys, newUserSignInActivity)
		}
	}

	//在返回消息中带上所有的签到活动
	for _, activeActivitys := range activeActivitys {
		//根据签到活动类型重算下当前的完成值
		value := api.calcActivityDoneValue(uid, now, activeActivitys.GetType(), activeActivitys.GetDoneValue(), activeActivitys.GetLatestTimestamp())
		if value == 0 { //同一天，跳过该活动
			continue
		}

		//生成新签到活动信息
		userSignInActivity := &battery.UserSignInActivity{
			Id:        activeActivitys.Id,
			Type:      activeActivitys.Type,
			DoneValue: proto.Uint32(value),
			GoalValue: activeActivitys.GoalValue,
			Items:     activeActivitys.Items,
		}
		//添加到返回信息中
		resp.Items = append(resp.Items, userSignInActivity)
	}

ErrHandle:

	resp.Uid = proto.String(uid)

	return
}

//处理新版本签到查询请求
func (api *XYAPI) OperationQuerySignIn2(req *battery.NewQuerySignInRequest, resp *battery.NewQuerySignInResponse) (err error) {
	var (
		uid                = req.GetUid()
		now                int64
		activitys          []*battery.DBUserSignInActivity
		activeActivitys    []*battery.DBUserSignInActivity
		setExceptActivitys = make(Set, 0)
		signInActivitys    = make([]*battery.SignInActivity, 0)
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//初始化返回
	resp.Error = xyerror.DefaultError()

	//xylog.Debug(uid, "[%s] QuerySignIn request begin", uid)
	//defer xylog.Debug(uid, "[%s] QuerySignIn request end", uid)

	//查询用户签到活动信息
	activitys, err = api.queryUserSignInActivitys(uid)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			xylog.Debug(uid, "UserSignInActivitys nofound.")
		} else {
			xylog.Error(uid, "UserSignInActivitys failed : %v", err)
			resp.Error.Code = battery.ErrorCode_QueryUserSignInRecordError.Enum()
			goto ErrHandle
		}
	}

	//查找需要添加的签到活动
	for _, activity := range activitys {
		setExceptActivitys[activity.GetId()] = empty{}
		if activity.GetState() == battery.UserMissionState_UserMissionState_Active { //正在进行的活动，加入返回列表中
			activeActivitys = append(activeActivitys, activity)
		}
	}

	now = time.Now().Unix()
	signInActivitys, err = api.querySignInActivity(uid, now, setExceptActivitys)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "querySignInActivity from cache failed : %v", err)
		resp.Error.Code = battery.ErrorCode_QuerySignInActivitysFromCacheError.Enum()
		goto ErrHandle
	}

	for _, signInActivity := range signInActivitys {
		newUserSignInActivity := api.defaultUserSignInActivity(uid, now, signInActivity)
		err = api.addUserSignInActivity(newUserSignInActivity)
		if err == xyerror.ErrOK { //插入正常，返回消息中也加上新的签到活动
			activeActivitys = append(activeActivitys, newUserSignInActivity)
		}
	}

	//在返回消息中带上所有的签到活动
	for _, activeActivitys := range activeActivitys {
		//根据签到活动类型重算下当前的完成值
		value, diff := api.calcActivityDoneValue2(uid, now, activeActivitys.GetType(), activeActivitys.GetDoneValue(), activeActivitys.GetLatestTimestamp())

		//生成新签到活动信息
		userSignInActivity := &battery.NewUserSignInActivity{
			Id:        activeActivitys.Id,
			Type:      activeActivitys.Type,
			DoneValue: proto.Uint32(value),
			GoalValue: activeActivitys.GoalValue,
			Items:     activeActivitys.Items,
		}

		if diff < 1 { // 设置当天已签到
			userSignInActivity.IsSignin = proto.Bool(true)
		}

		//添加到返回信息中
		resp.Items = append(resp.Items, userSignInActivity)
	}

ErrHandle:

	resp.Uid = proto.String(uid)

	return
}

//计算签到完成值
//args:
//now int64 当前时间戳
//activityType battery.SignInActivityType 签到活动类型
//doneValue uint32 玩家当前完成值
//latestTimestamp int64 玩家最近签到时间戳
//return:
//value uint32 玩家新的完成值（用于前端显示是第几天）
func (api *XYAPI) calcActivityDoneValue2(uid string, now int64, activityType battery.SignInActivityType, doneValue uint32, latestTimestamp int64) (value uint32, diff int64) {
	diff = xyutil.DayDiff(latestTimestamp, now)

	switch activityType {
	case battery.SignInActivityType_SignInActivityType_UnContinuous: //累计型
		if diff >= 1 { //跨天+1
			value = doneValue + 1
		} else if diff < 1 { //同一天

			value = doneValue
		}

	default: // battery.SignInActivityType_SignInActivityType_Continuous//连续型
		if diff > 1 { //间隔多天重置为第一天登录
			value = 1
		} else if diff == 1 { //连续
			value = doneValue + 1
		} else if diff < 1 { //同一天

			value = doneValue
		}
	}

	xylog.Debug(uid, "diff(%d) return value(%d)", diff, value)

	return
}

//计算旧版本签到完成值
func (api *XYAPI) calcActivityDoneValue(uid string, now int64, activityType battery.SignInActivityType, doneValue uint32, latestTimestamp int64) (value uint32) {
	diff := xyutil.DayDiff(latestTimestamp, now)

	switch activityType {
	case battery.SignInActivityType_SignInActivityType_UnContinuous: //累计型
		if diff >= 1 { //跨天+1
			value = doneValue + 1
		} else if diff < 1 { //同一天
			//do nothing
		}

	default: // battery.SignInActivityType_SignInActivityType_Continuous//连续型
		if diff > 1 { //间隔多天重置为第一天登录
			value = 1
		} else if diff == 1 { //连续
			value = doneValue + 1
		} else if diff < 1 { //同一天
			//do nothing
		}
	}

	xylog.Debug(uid, "diff(%d) return value(%d)", diff, value)

	return
}

//查询玩家当前的签到活动信息
//uid string 玩家id
//return:
//activitys []*battery.UserSignInActivity  玩家当前签到活动列表
func (api *XYAPI) queryUserSignInActivitys(uid string) (activitys []*battery.DBUserSignInActivity, err error) {
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSIGNINACTIVITY).QueryUserSignInActivitys(uid, &activitys)
	if err != nil {
		xylog.Error(uid, "QueryUserSignInActivitys failed : %v", err)
		return
	}
	return
}

func (api *XYAPI) queryUserSignInActivityDetail(uid string, id uint64) (activity *battery.DBUserSignInActivity, err error) {
	activity = new(battery.DBUserSignInActivity)
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSIGNINACTIVITY).QueryUserSignInActivityDetail(uid, id, activity)
	if err != nil {
		xylog.Error(uid, "QueryUserSignInActivitys failed : %v", err)
		return
	}
	return
}

//增加玩家签到活动信息
//activity *battery.UserSignInActivity 新的签到活动指针
func (api *XYAPI) addUserSignInActivity(activity *battery.DBUserSignInActivity) (err error) {
	return api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSIGNINACTIVITY).AddUserSignInActivity(activity)
}

//生成玩家签到活动信息
//uid string 玩家id
//now int    当前时间戳
//activity *battery.SignInActivity 对应的签到活动信息指针
//*battery.UserSignInActivity 生成的玩家签到活动信息指针
func (api *XYAPI) defaultUserSignInActivity(uid string, now int64, activity *battery.SignInActivity) *battery.DBUserSignInActivity {

	Timestamp := &battery.UserMissionTimestamp{
		State:     battery.UserMissionState_UserMissionState_Active.Enum(),
		Timestamp: &now,
	}

	return &battery.DBUserSignInActivity{
		Uid:             proto.String(uid),
		Id:              activity.Id,
		Type:            activity.Type,
		DoneValue:       proto.Uint32(0), //默认未签到
		GoalValue:       activity.GoalValue,
		LatestTimestamp: proto.Int64(0), //默认未签到过
		State:           battery.UserMissionState_UserMissionState_Active.Enum(),
		Timestamps:      []*battery.UserMissionTimestamp{Timestamp},
		Items:           activity.Items,
	}
}

//查询签到配置信息
// now int64 当前时间戳
// exceptId Set 需要过滤的签到活动id集合
func (api *XYAPI) querySignInActivity(uid string, now int64, exceptId Set) (signInActivitys []*battery.SignInActivity, err error) {

	signInActivitys = make([]*battery.SignInActivity, 0)

	activitys := xybusinesscache.DefSignInCacheManager.Activitys()

	if nil == activitys {
		err = xyerror.ErrQuerySignInActivitysFromCacheError
		return
	}

	for _, activity := range *activitys {
		//玩家已经存在对应活动，跳过
		if _, ok := exceptId[activity.GetId()]; ok {
			xylog.Debug(uid, "activity %v in already in usersigninactivity.skip it", activity)
			continue
		}

		//不在签到活动期间内，跳过
		if activity.GetBeginTime() > now || activity.GetEndTime() < now {
			xylog.Debug(uid, "activity %v no in timerange (%d,%d)", activity, activity.GetBeginTime(), activity.GetEndTime())
			continue
		}

		//玩家未存在的活动，添加
		signInActivitys = append(signInActivitys, activity)
		xylog.Debug(uid, "add signin activity %v", activity)
	}

	return
}

//处理玩家签到请求
// req *battery.SignInRequest  签到请求指针
// resp *battery.SignInResponse 签到回应指针
func (api *XYAPI) OperationSignIn(req *battery.SignInRequest, resp *battery.SignInResponse) (err error) {
	var (
		uid      = req.GetUid()
		id       = req.GetId()
		activity *battery.DBUserSignInActivity
		value    uint32
		now      int64
		errStr   string
	)

	//获取请求的终端平台类型
	platform := req.GetPlatformType()
	api.SetDB(platform)

	//初始化resp
	resp.Uid = req.Uid
	resp.Id = req.Id
	resp.Error = xyerror.DefaultError()

	xylog.Debug(uid, "SignIn request begin", uid)
	defer xylog.Debug(uid, "SignIn request end", uid)

	//查找玩家当前的活动详情
	activity, err = api.queryUserSignInActivityDetail(uid, id)
	if err != xyerror.ErrOK {
		if err == xyerror.ErrNotFound {
			errStr = fmt.Sprintf("[%s] UserSignInActivitys nofound.", uid)
		} else {
			errStr = fmt.Sprintf("[%s] UserSignInActivitys failed : %v", uid, err)
		}
		xylog.Error(uid, errStr)
		resp.Error.Code = battery.ErrorCode_QueryUserSignInRecordError.Enum()
		goto ErrHandle
	}

	//校验活动状态，如果不是Active状态，则报错返回
	if activity.GetState() != battery.UserMissionState_UserMissionState_Active {
		xylog.Error(uid, "SignInActivitys state %v, can't signin.", activity.GetState())
		resp.Error.Code = battery.ErrorCode_SignInActivitysStateError.Enum()
		goto ErrHandle
	}

	//计算玩家签到活动完成值
	now = time.Now().Unix()
	value = api.calcActivityDoneValue(uid, now, activity.GetType(), activity.GetDoneValue(), activity.GetLatestTimestamp())
	if value == 0 {
		//同一天签到，告警
		xylog.Warning(uid, "[%s] signin at the same day, something wrong.", uid)
		resp.Error.Code = battery.ErrorCode_DuplicateSignInError.Enum()
		goto ErrHandle
	}
	activity.DoneValue = proto.Uint32(value)
	activity.LatestTimestamp = proto.Int64(now)

	//先分发当天奖励
	err = api.gainSignInAward(uid, int(value), activity.GetItems())
	if err != xyerror.ErrOK {
		xylog.Error(uid, "gainSignInAward failed : %v", err)
		resp.Error.Code = battery.ErrorCode_GainSignInAwardError.Enum()
		goto ErrHandle
	}

	//判断玩家签到活动是否完成
	if activity.GetDoneValue() >= activity.GetGoalValue() {
		//签到活动完成，刷新活动状态
		api.setUserSignInState(now, battery.UserMissionState_UserMissionState_DoneCollected, activity)
		//发放奖励
		err = api.gainSignInAward(uid, int(0), activity.GetItems())
		if err != xyerror.ErrOK {
			xylog.Error(uid, "gainSignInAward failed : %v", err)
			resp.Error.Code = battery.ErrorCode_GainSignInAwardError.Enum()
			goto ErrHandle
		}
	}

	//更新玩家活动信息
	err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_USERSIGNINACTIVITY).UpsertUserSignInActivity(uid, activity)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "UpsertUserSignInActivity failed : %v", err)
		resp.Error.Code = battery.ErrorCode_UpsertUserSignInRecordError.Enum()
		goto ErrHandle
	}

ErrHandle:

	return
}

func (api *XYAPI) gainSignInAward(uid string, index int, items []*battery.SignInItem) (err error) {
	if index >= len(items) { //如果签到天数超过奖励设置天数，则报错
		xylog.Error(uid, "index(%d) >= len(activity.GetItems())(%d) while get signInItems of items %v", index, len(items), items)
		err = xyerror.ErrGetSignInItemsError
		return
	} else {
		//分发奖品
		//再分发目标达成奖励
		awards := (items)[index].GetAward()
		for _, propItem := range awards {
			err = api.GainProp(uid, nil, propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
			if err != xyerror.ErrOK {
				return
			}
		}
	}

	return
}

//设置签到状态
func (api *XYAPI) setUserSignInState(now int64, state battery.UserMissionState, activity *battery.DBUserSignInActivity) {
	activity.State = state.Enum()
	timestamp := &battery.UserMissionTimestamp{
		State:     state.Enum(),
		Timestamp: proto.Int64(now),
	}

	activity.Timestamps = append(activity.Timestamps, timestamp)
}
