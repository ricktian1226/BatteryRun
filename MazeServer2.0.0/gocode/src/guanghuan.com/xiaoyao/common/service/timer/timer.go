// xytimer 定时器接口定义
package xytimer

import (
	"errors"
	//"strconv"
	//"strings"
	"time"

	xylog "guanghuan.com/xiaoyao/common/log"
)

//定时任务类型
const (
	TIMER_TYPE_FIXED    = iota //定点任务
	TIMER_TYPE_INTERVAL        //固定时间间隔任务
)

type TimerHandler func()

// TimerMoment 定义了定点时间结构
type TimerMoment struct {
	Hour, Minute, Second int
}

// TimerSingleOption 定义了单个定时任务时间相关的参数
type TimerSingleOption struct {
	Type        int   //定时任务类型
	Interval    int64 //间隔定义
	TimerMoment       //时刻定义
}

// TimerOption 定义了逻辑定时任务的相关参数
type TimerOption struct {
	Type     int           //定时任务类型
	Moments  []TimerMoment //时刻列表，可以定义在多个时刻触发
	Interval int64         //间隔时间，单位:秒
}

var (
	ErrMomentStringError = errors.New("moments string error")
	ErrTimeTypeError     = errors.New("time type is not TIME_TYPE_FIXED/TIME_TYPE_INTERVAL")
)

// InitTimer 初始化定时器
//
func InitTimer(timerOption TimerOption, handler TimerHandler) (err error) {
	var (
		moments       []int64
		singleOptions []TimerSingleOption
	)

	moments, singleOptions, err = GetInitialAwakeTime(timerOption)
	if err != nil {
		return
	}

	xylog.DebugNoId("GetInitialAwakeTime : %v", moments)

	for i := 0; i < len(moments); i++ {
		singleOption := singleOptions[i]
		timer := time.NewTimer(time.Second * time.Duration(moments[i]))

		xylog.DebugNoId("NewTimer : awake in %d seconds", moments[i])
		//每个定时器启动一个goroutine来侦听
		go func() {
			for range timer.C {
				FuncTimer(singleOption, handler, timer)
			}
		}()
	}
	return
}

// FuncTimer 定时器函数
// singleOption TimerSingleOption 定时器配置
// handler TimerHandler 定时处理函数
func FuncTimer(singleOption TimerSingleOption, handler TimerHandler, timer *time.Timer) {
	//
	handler()

	//计算下次唤醒时间
	nextTime, err := GetNextAwakeTime(singleOption)
	if err != nil {
		xylog.ErrorNoId("GetNextAwakeTime failed : %v", err)
		return
	}

	xylog.DebugNoId("nextTime %d seconds later", nextTime)

	//timer.AfterFunc(time.Second*time.Duration(nextTime), func() {
	//	FuncTimer(singleOption, handler)
	//})

	//重新设置定时器时间
	timer.Reset(time.Second * time.Duration(nextTime))
}

//GetInitialAwakeTime 获取初始唤醒时间列表
func GetInitialAwakeTime(opt TimerOption) ([]int64, []TimerSingleOption, error) {
	switch opt.Type {
	case TIMER_TYPE_FIXED:
		return getInitialAwakeTimeFixed(opt)
	case TIMER_TYPE_INTERVAL:
		return getInitialAwakeTimeInterval(opt)
	default:
		return nil, nil, ErrTimeTypeError
	}
}

// GetNextAwakeTime 获取下次唤醒时间
// timeSingleOption TimerSingleOption
func GetNextAwakeTime(timeSingleOption TimerSingleOption) (int64, error) {
	nowTime := time.Now().Unix()
	switch timeSingleOption.Type {
	case TIMER_TYPE_FIXED:
		return getNextDayAwakeTime(nowTime, timeSingleOption.Hour, timeSingleOption.Minute, timeSingleOption.Second), nil
	case TIMER_TYPE_INTERVAL:
		return getNextAwakeTimeInterval(nowTime, timeSingleOption.Interval), nil
	default:
		return 0, ErrTimeTypeError
	}
}

// getInitialAwakeTimeFixed 获取定点任务的初始唤醒时间
// opt TimerOption 时间参数
//returns:
// moments []int64 唤醒的时间戳列表
func getInitialAwakeTimeFixed(opt TimerOption) (moments []int64, singleOptions []TimerSingleOption, err error) {
	moments, singleOptions = make([]int64, 0), make([]TimerSingleOption, 0)

	for _, m := range opt.Moments {
		//比较一下当前时间和当天的唤醒时间，根据过没过当天唤醒决定定时器的唤醒时间
		now := time.Now()
		year, month, day := now.Date()
		todayAwakeTime := time.Date(year, month, day, m.Hour, m.Minute, m.Second, 0, time.Local).Unix()
		nowTime := now.Unix()

		if nowTime >= todayAwakeTime { //已经过了当天的唤醒时间，设置为第二天唤醒时间
			moments = append(moments, getNextDayAwakeTime(nowTime, m.Hour, m.Minute, m.Second))
		} else { //未过当天的唤醒时间，设置为当天唤醒时间
			moments = append(moments, todayAwakeTime-nowTime)
		}

		singleOption := TimerSingleOption{
			Type: opt.Type,
		}
		singleOption.Hour = m.Hour
		singleOption.Minute = m.Minute
		singleOption.Second = m.Second
		singleOptions = append(singleOptions, singleOption)
	}
	return
}

//getNextDayAwakeTime 获取第二天的唤醒时间
// now time.Time
// hour int 唤醒的小时
// minute int 唤醒分钟
// second int 唤醒秒
//returns:
// nextDayAwakeTime int64 第二天唤醒时间
func getNextDayAwakeTime(nowTime int64, hour, minute, second int) (nextDayAwakeTime int64) {
	nextDayTime := nowTime + 3600*24
	year, month, day := time.Unix(nextDayTime, 0).Date()
	nextDayAwakeTime = time.Date(year, month, day, hour, minute, second, 0, time.Local).Unix() - nowTime
	return
}

// getInitialAwakeTimeInterval 获取间隔任务的初始唤醒时间
// opt TimerOption 时间参数
//returns:
// moments []int64 唤醒的时间戳列表
func getInitialAwakeTimeInterval(opt TimerOption) (moments []int64, singleOptions []TimerSingleOption, _ error) {
	moments, singleOptions = make([]int64, 0), make([]TimerSingleOption, 0)
	moments = append(moments, opt.Interval)
	singleOption := TimerSingleOption{
		Type:     opt.Type,
		Interval: opt.Interval,
	}
	singleOptions = append(singleOptions, singleOption)

	return moments, singleOptions, nil
}

// getNextAwakeTimeInterval 获取间隔定时器的下次唤醒时间
// nowTime int64 当前时间戳
// interval int64 时间间隔
func getNextAwakeTimeInterval(nowTime int64, interval int64) int64 {
	return interval
}
