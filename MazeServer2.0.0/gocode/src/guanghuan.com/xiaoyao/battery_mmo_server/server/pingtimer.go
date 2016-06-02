package server

import (
	xylogs "guanghuan.com/xiaoyao/common/log"
	"sync"
	"time"
)

//连接过期处理器接口类
type DierInterface interface {
	Die()
}

type PingTimer struct {
	mu    sync.Mutex
	ptmr  *time.Timer   //ping定时器
	pTime time.Time     //记录最近ping时间戳
	dier  DierInterface //心跳死亡后的执行函数
}

//注册连接过期的处理器
func (t *PingTimer) SetDier(pi DierInterface) {
	t.dier = pi
	t.pTime = time.Now()
}

func (t *PingTimer) Process() {
	bDie := false

	t.mu.Lock()

	t.ptmr = nil

	//如果客户端的时间戳超过最大超时时长，则认为该客户端过期，销毁处理
	xylogs.Debug("time.Since(t.pTime) : %d , PingTimeOut : %d", time.Since(t.pTime), (time.Duration)(GOpts.PingTimeOut)*time.Second)

	if t.pTime.Nanosecond() != 0 && (time.Since(t.pTime) > (time.Duration)(GOpts.PingTimeOut)*time.Second) {
		bDie = true
	}

	t.mu.Unlock()

	if bDie {
		xylogs.Debug("PingTimer : Client should go to hell.")
		t.dier.Die()
		t.Clear()
	} else {
		//重置心跳检测定时器
		xylogs.Debug("PingTimer : Client should not go to hell.")
		t.Set()
	}
}

func (t *PingTimer) Set() {
	t.mu.Lock()
	defer t.mu.Unlock()
	//xylogs.Debug("Set PingTimer.")
	t.ptmr = time.AfterFunc(time.Duration(GOpts.PingInterval)*time.Second, func() { t.Process() })
}

func (t *PingTimer) RefleshTime() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.pTime = time.Now()
}

func (t *PingTimer) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.ptmr == nil {
		return
	}
	t.ptmr.Stop()
	t.ptmr = nil
}
