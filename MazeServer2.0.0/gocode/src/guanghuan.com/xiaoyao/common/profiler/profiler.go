package xyprofiler

import (
	"fmt"
	xylog "guanghuan.com/xiaoyao/common/log"
	"sync"
	"time"
)

type Profiler struct {
	Name         string
	JobCount     int
	SuccessCount int
	ErrorCount   int
	TotalTime    time.Duration
	MaxTime      time.Duration
	TotalDataIn  int64
	MaxDataIn    int64
	TotalDataOut int64
	MaxDataOut   int64
	ProfMutex    sync.Mutex
}

type ProfilerMap struct {
	IsStart   bool
	StartTime time.Time
	StopTime  time.Time
	M         map[string]*Profiler
	Mutex     sync.Mutex
	P         *Profiler
}

func NewProfilerMap(init_cap int) (pm *ProfilerMap) {
	if init_cap <= 0 {
		init_cap = 10
	}
	pm = &ProfilerMap{
		M: make(map[string]*Profiler, init_cap),
		P: &Profiler{Name: "All Jobs"},
	}
	return
}

func (pm *ProfilerMap) IsStarted() bool {
	return pm.IsStart
}

func (pm *ProfilerMap) Stop() {
	pm.Mutex.Lock()
	defer pm.Mutex.Unlock()

	if pm.IsStart {
		pm.IsStart = false
		pm.StopTime = time.Now()
		xylog.InfoNoId("Profiler Stopped : %s", pm.StopTime.String())
	}
}

func (pm *ProfilerMap) Start() {
	pm.Mutex.Lock()
	defer pm.Mutex.Unlock()
	if !pm.IsStart {
		pm.IsStart = true
		pm.StartTime = time.Now()
		xylog.InfoNoId("Profiler Started : %s", pm.StartTime.String())
	}
}

func (pm *ProfilerMap) Reset() {
	pm.Mutex.Lock()
	for n, p := range pm.M {
		if p != nil {
			p.Reset()
		} else {
			delete(pm.M, n)
		}
	}
	pm.P.Reset()
	pm.StartTime = time.Now()
	xylog.Info("Profiler Reset : %s", pm.StartTime.String())
	pm.Mutex.Unlock()
}

func (pm *ProfilerMap) Profiler(name string) (p *Profiler) {
	pm.Mutex.Lock()
	if pm.M[name] == nil {
		p = &Profiler{Name: name}
		pm.M[name] = p
	} else {
		p = pm.M[name]
	}
	pm.Mutex.Unlock()
	return
}

func (pm *ProfilerMap) String() (str string) {
	pm.Mutex.Lock()
	m := pm.M

	//jobs := 0
	//jobs_ok := 0
	//jobs_fail := 0

	stop := pm.StopTime
	if pm.StopTime.Before(pm.StartTime) {
		stop = time.Now()
	}
	for n, p := range m {
		if p != nil {
			str = str + "\n" + n + ":" + p.String()
			//pm.P.AddJob(true, p.TotalTime, p.TotalDataOut, p.TotalDataIn)
			//jobs += p.JobCount
			//jobs_ok += p.SuccessCount
			//jobs_fail += p.ErrorCount
		} else {
			str = str + "\n" + n + ":" + "nil"
		}
	}
	//pm.P.JobCount = jobs
	//pm.P.SuccessCount = jobs_ok
	//pm.P.ErrorCount = jobs_fail

	str0 := fmt.Sprintf(` -------- profiler result --------
	Start : %s
	Stop  : %s
	Past  : %s
	[%d] profilers :%s`,
		pm.StartTime.String(),
		stop.String(),
		stop.Sub(pm.StartTime).String(),
		len(m),
		pm.P.String())
	pm.Mutex.Unlock()

	str = str0 + str
	return
}

func (pm *ProfilerMap) AddJobResult(job_name string, success bool, dur time.Duration, data_out int64, data_in int64) {
	if pm.IsStarted() {
		p := pm.Profiler(job_name)
		p.AddJob(success, dur, data_out, data_in)
		pm.P.AddJob(success, dur, data_out, data_in)
	}
}

func (p *Profiler) AddJob(success bool, dur time.Duration, data_out int64, data_in int64) {
	p.ProfMutex.Lock()
	p.JobCount++
	if success {
		p.SuccessCount++
	} else {
		p.ErrorCount++
	}
	p.TotalTime += dur
	if dur > p.MaxTime {
		p.MaxTime = dur
	}
	p.TotalDataIn += data_in
	if data_in > p.MaxDataIn {
		p.MaxDataIn = data_in
	}
	p.TotalDataOut += data_out
	if data_out > p.MaxDataOut {
		p.MaxDataOut = data_out
	}

	p.ProfMutex.Unlock()
}

func (p *Profiler) Reset() {
	p.ProfMutex.Lock()
	p.JobCount = 0
	p.ErrorCount = 0
	p.SuccessCount = 0
	p.TotalTime = 0
	p.MaxTime = 0
	p.TotalDataIn = 0
	p.MaxDataIn = 0
	p.TotalDataOut = 0
	p.MaxDataOut = 0
	p.ProfMutex.Unlock()
}
func (p *Profiler) String() (str string) {
	p.ProfMutex.Lock()

	if p.JobCount <= 0 {
		str = fmt.Sprintf(`
	Job Name      : %s
	Job Count     : %d (%d : %d)	
	`, p.Name, p.JobCount, 0, 0)
	} else {
		str = fmt.Sprintf(`
	Job Name      : %s
	Job Count     : %d (%d : %d)	

	Total time    : %0.3f s
	Max   time    : %0.3f ms
	Avg   time    : %0.3f ms

	Total Data    : %0.3f (kb)
		
	Total Data In : %0.3f (kb)
	Max   Data In : %d (byte)
	Avg   Data In : %d (byte)
	
	Total Data Out: %0.3f (kb)
	Max   Data Out: %d (byte)
	Ava   Data Out: %d (byte)
	-------------------
	`, p.Name,
			p.JobCount,
			p.SuccessCount,
			p.ErrorCount,

			p.TotalTime.Seconds(),
			p.MaxTime.Seconds()*1000,
			p.TotalTime.Seconds()*1000/float64(p.JobCount),

			float64(p.TotalDataIn+p.TotalDataOut)/1000,

			float64(p.TotalDataIn)/1000,
			p.MaxDataIn,
			p.TotalDataIn/int64(p.JobCount),

			float64(p.TotalDataOut)/1000,
			p.MaxDataOut,
			p.TotalDataOut/int64(p.JobCount))
	}
	p.ProfMutex.Unlock()

	return
}
