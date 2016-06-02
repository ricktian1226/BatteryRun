// xyperftimer
package xyperf

import (
	"runtime/debug"
	"time"
)

func StartPerfTimer() {
	go func() {
		timer := time.NewTicker(time.Minute)
		for {
			select {
			case <-timer.C:
				PerfLog()
			}
		}
	}()
}

func StartTimeRangeTimer() {
	go func() {
		timer := time.NewTicker(time.Minute)
		for {
			select {
			case <-timer.C:
				TimeRangeLog()
			}
		}
	}()
}

func StartFreeOSMemoryTimer() {
	go func() {
		timer := time.NewTicker(time.Minute * 10)
		for {
			select {
			case <-timer.C:
				debug.FreeOSMemory()
			}
		}
	}()
}
