// xyperf
package xyperf

import (
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	"os"
	"sync/atomic"
)

type StatisticInterface interface {
	LogFile(lc *xylog.LoggerConfig)
}

type Statistic struct {
	logger *xylog.XYLogger
}

func (s *Statistic) add(data *int64, delta int64) {
	if delta < 1 {
		delta = 1
	}
	atomic.AddInt64(data, delta)
}

func (s *Statistic) AppLogConfig(si StatisticInterface, lc *xylog.LoggerConfig) {
	if lc == nil {
		lc = xylog.DefConfig
	}
	logger := s.logger.Logger()
	logger.SetLevel((int)(xylog.TraceLevel)) //set log level to trace

	if lc.Verbose {
		logger.EnableFuncCallDepth(true)
		logger.SetLogFuncCallDepth(4)
	} else {
		logger.EnableFuncCallDepth(false)
		logger.SetLogFuncCallDepth(0)
	}

	var strconfig string
	if lc.Stdout {
		strconfig = fmt.Sprintf(`{"level":%v}`, lc.Level)
		logger.SetLogger("console", strconfig)
	} else {
		si.LogFile(lc)
		strconfig = fmt.Sprintf(`{"filename":"%s/%s","maxlines":%v,"maxsize":%v,"daily":%v,"maxdays":%v,"rotate":%v}`,
			lc.Path,
			lc.Filename,
			lc.Maxlines,
			lc.Maxsize,
			lc.Daily,
			lc.Maxdays,
			lc.Rotate) //,
		//logconfig.Level)
		logger.SetLogger("file", strconfig)
	}

	xylog.DebugNoId("%s\n", strconfig)
}

//the func from interface
func (s Statistic) LogFile(lc *xylog.LoggerConfig) {
	//if lc.Filename == "" {
	if lc.NodeId >= 0 {
		lc.Filename = fmt.Sprintf("perf_%s_%d.log", lc.AppName, lc.NodeId)
	} else {
		lc.Filename = fmt.Sprintf("perf_%s.%d.log", lc.AppName, os.Getpid())
	}
	//}
}

type Perf struct {
	Statistic
	totolRequest int64 //total requests from the moment server start
	totolTime    int64 //total time cost for process
	/*	request          int64 //total request during the period
		time             int64 //total time during the period*/
	requestIncrement int64 //request increment per (maybe... minute)
	timeIncrement    int64 //time cost increment per (maybe... minute)
}

func (p *Perf) AddTotalRequest(delta int64) {
	p.add(&p.totolRequest, delta)
}

func (p *Perf) AddTotalTime(delta int64) {
	p.add(&p.totolTime, delta)
}

/*func (p *Perf) AddRequest(delta int64) {
	p.add(&p.request, delta)
}

func (p *Perf) AddTime(delta int64) {
	p.add(&p.time, delta)
}*/

func (p *Perf) AddRequestIncrement(delta int64) {
	p.add(&p.requestIncrement, delta)
}

func (p *Perf) AddTimeIncrement(delta int64) {
	p.add(&p.timeIncrement, delta)
}

func (p *Perf) Truncate() {
	atomic.SwapInt64(&p.requestIncrement, 0)
	atomic.SwapInt64(&p.timeIncrement, 0)
}

func (p *Perf) Log() {
	requestIncrement, timeIncrement := atomic.LoadInt64(&(p.requestIncrement)), atomic.LoadInt64(&(p.timeIncrement))
	var rt, tr int64
	if requestIncrement <= 0 {
		tr = 0
	} else {
		tr = timeIncrement / requestIncrement
	}

	if timeIncrement <= 0 {
		rt = 0
	} else {
		rt = requestIncrement / 60
	}

	p.logger.Log(xylog.TraceLevel, `TR(total request):%v	TT(total time):%vms	RI(request increment):%v	TI(time increment):%vms	R/S(request per second):%v	T/R(time per request):%vms`,
		atomic.LoadInt64(&(p.totolRequest)),
		atomic.LoadInt64(&(p.totolTime)),
		atomic.LoadInt64(&(p.requestIncrement)),
		atomic.LoadInt64(&(p.timeIncrement)),
		rt,
		tr)
	p.Truncate()
}

type TimeRange struct {
	Statistic
	Lt1        int64 //(0,1ms]
	Lt10       int64 //(1,10ms]
	Gt10Lt50   int64 //(10ms,50ms]
	Gt50Lt100  int64 //(50ms,100ms]
	Gt100Lt200 int64 //(100ms,200ms]
	Gt200Lt300 int64 //(200ms,300ms]
	Gt300      int64 //(300ms,)
}

func (p *TimeRange) Truncate() {
	atomic.SwapInt64(&p.Lt1, 0)
	atomic.SwapInt64(&p.Lt10, 0)
	atomic.SwapInt64(&p.Gt10Lt50, 0)
	atomic.SwapInt64(&p.Gt50Lt100, 0)
	atomic.SwapInt64(&p.Gt100Lt200, 0)
	atomic.SwapInt64(&p.Gt200Lt300, 0)
	atomic.SwapInt64(&p.Gt300, 0)
}

func (p *TimeRange) Add(delta int64) {
	switch {
	case delta >= 0 && delta <= 1:
		p.add(&p.Lt1, 1)
	case delta > 1 && delta <= 10:
		p.add(&p.Lt10, 1)
	case delta > 10 && delta <= 50:
		p.add(&p.Gt10Lt50, 1)
	case delta > 50 && delta <= 100:
		p.add(&p.Gt50Lt100, 1)
	case delta > 100 && delta <= 200:
		p.add(&p.Gt100Lt200, 1)
	case delta > 200 && delta <= 300:
		p.add(&p.Gt200Lt300, 1)
	case delta > 300:
		p.add(&p.Gt300, 1)
	}
}

//the func from interface
func (p TimeRange) LogFile(lc *xylog.LoggerConfig) {
	//if lc.Filename == "" {
	if lc.NodeId >= 0 {
		lc.Filename = fmt.Sprintf("perf_timerange_%s_%d", lc.AppName, lc.NodeId)
	} else {
		lc.Filename = fmt.Sprintf("perf_timerange_%s.%d", lc.AppName, os.Getpid())
	}

	if lc.LogId >= 0 {
		lc.Filename += fmt.Sprintf("_%d", lc.LogId)
	}

	lc.Filename += ".log"

	xylog.DebugNoId("logfile : %s", lc.Filename)
	//}
}

func (p *TimeRange) Log() {
	p.logger.Log(xylog.TraceLevel, `(0,1ms]:%v	(1,10ms]:%v	(10ms,50ms]:%v	(50ms,100ms]:%v	(100ms,200ms]%v	(200ms,300ms]:%v	(300ms,...]:%v`,
		atomic.LoadInt64(&(p.Lt1)),
		atomic.LoadInt64(&(p.Lt10)),
		atomic.LoadInt64(&(p.Gt10Lt50)),
		atomic.LoadInt64(&(p.Gt50Lt100)),
		atomic.LoadInt64(&(p.Gt100Lt200)),
		atomic.LoadInt64(&(p.Gt200Lt300)),
		atomic.LoadInt64(&(p.Gt300)))
	p.Truncate()
}
