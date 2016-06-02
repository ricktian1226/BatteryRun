package xyperf

import (
	xylog "guanghuan.com/xiaoyao/common/log"
	"sync"
	"time"
)

var perf *Perf

func InitPerf() {
	perf = NewPerf()
}

func NewPerf() (p *Perf) {
	logconfig := *(xylog.DefConfig)
	logconfig.Stdout = false

	p = &Perf{
		Statistic: Statistic{
			logger: xylog.NewLogger(&logconfig, 1000),
		},
		totolRequest:     0,
		totolTime:        0,
		requestIncrement: 0,
		timeIncrement:    0,
	}

	p.AppLogConfig(*p, &logconfig)

	return p
}

func AddTotalRequest(delta int64) {
	perf.AddTotalRequest(delta)
}

func AddTotalTime(delta int64) {
	perf.AddTotalTime(delta)
}

func AddRequestIncrement(delta int64) {
	perf.AddRequestIncrement(delta)
}

func AddTimeIncrement(delta int64) {
	perf.AddTimeIncrement(delta)
}

func PerfLog() {
	perf.Log()
}

//var timeRange *TimeRange

func InitTimeRange() {
	DefTimeRangeManager = NewTimeRangeManager()
}

type MAPLogId2TimeRange map[int]*TimeRange

type TimeRangeManager struct {
	timeranges MAPLogId2TimeRange
	mutex      sync.Mutex
}

var DefTimeRangeManager *TimeRangeManager

func NewTimeRangeManager() (p *TimeRangeManager) {
	p = &TimeRangeManager{
		timeranges: make(MAPLogId2TimeRange, 0),
	}

	//默认监控请求消息
	p.timeranges[DefLogId] = NewTimeRange(DefLogId)

	return
}

func NewTimeRange(logId int) (p *TimeRange) {
	logconfig := *(xylog.DefConfig)
	logconfig.Stdout = false
	logconfig.LogId = logId
	p = &TimeRange{
		Statistic: Statistic{
			logger: xylog.NewLogger(&logconfig, 1000),
		},
		Lt10:       0,
		Gt10Lt50:   0,
		Gt50Lt100:  0,
		Gt100Lt200: 0,
		Gt200Lt300: 0,
		Gt300:      0,
	}

	p.AppLogConfig(*p, &logconfig)

	return p
}

//func AddTimeRange(delta int64) {
//	timeRange.Add(delta)
//}

func TimeRangeLog() {
	DefTimeRangeManager.mutex.Lock()
	defer DefTimeRangeManager.mutex.Unlock()
	for _, timeRange := range DefTimeRangeManager.timeranges {
		//xylog.Debug("timeRange[%d] %v", id, timeRange)
		timeRange.Log()
	}

}

const DefLogId = 0

func Trace(logId int, begin *time.Time) {
	var timeRange *TimeRange
	DefTimeRangeManager.mutex.Lock()
	if _, ok := DefTimeRangeManager.timeranges[logId]; !ok {
		DefTimeRangeManager.timeranges[logId] = NewTimeRange(logId)
	}
	DefTimeRangeManager.mutex.Unlock()

	timeRange = DefTimeRangeManager.timeranges[logId]

	cost := time.Since(*begin)
	timeRange.Add(int64(cost / time.Millisecond))
	if logId == DefLogId {
		AddTotalRequest(1)
		AddTotalTime(int64(cost / time.Millisecond))
		AddRequestIncrement(1)
		AddTimeIncrement(int64(cost / time.Millisecond))
	}

	//AddTimeRange()
}
