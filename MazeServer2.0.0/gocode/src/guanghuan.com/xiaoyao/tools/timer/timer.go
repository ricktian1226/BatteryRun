// timer
package main

import (
	"time"
)

func main() {
	Init()
	go func() {
		UpsertSysLottoInfo()
		timer := time.NewTicker(time.Minute * time.Duration(DefConfig.InfoInterval))
		for {
			select {
			case <-timer.C:
				UpsertSysLottoInfo()
			}
		}
	}()

	go func() {
		UpsertSysFreeRest()
		timer := time.NewTicker(time.Minute * time.Duration(DefConfig.FreeInterval))
		for {
			select {
			case <-timer.C:
				UpsertSysFreeRest()
			}
		}
	}()

	ch := make(chan int, 1)
	<-ch

}
