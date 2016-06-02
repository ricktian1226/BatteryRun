// useridentity
// 玩家identity管理器
package xybusinesscache

import (
	"fmt"
	"sync"

	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

//前缀列表
//var Prefix [26]byte = {'A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z'}
var Prefix = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

// UserIdentityManager 玩家identity管理器定义
type UserIdentityManager struct {
	lock    sync.Mutex
	counter []int32 //玩家数目计数器
	prefix  string  //identity前缀
}

var DefUserIdentityManager = NewUserIdentityManager()

func NewUserIdentityManager() *UserIdentityManager {
	return &UserIdentityManager{}
}

func (m *UserIdentityManager) InitWhileStart(dcId, nodeId int) bool {
	if dcId > len(Prefix) || nodeId > len(Prefix) {
		xylog.ErrorNoId("UserIdentityManager.Init failed, error dcId(%d) nodeId(%d)", dcId, nodeId)
		return false
	}

	m.prefix = fmt.Sprintf("%c%c", Prefix[dcId], Prefix[nodeId])

	if err := m.load(); err != xyerror.ErrOK {
		xylog.ErrorNoId("UserIdentityManager.load failed : %v", err)
		return false
	}

	xylog.DebugNoId("UserIdentityManager : %v", m)

	return true
}

// load 加载当且节点的useridentitycounter信息
func (m *UserIdentityManager) load() (err error) {
	var userIdentityCounter = &battery.DBUserIdentityCounter{}

	m.lock.Lock()
	defer m.lock.Unlock()

	err = DefCacheDB.LoadUserIdentityCounter(m.prefix, userIdentityCounter)
	if err != xyerror.ErrOK {
		return
	}

	m.counter = userIdentityCounter.GetCounter()

	return
}

const (
	BASE_IDENTITY_SEGMENT = 10000
)

// 生成玩家identity
func (m *UserIdentityManager) Spawn(platform battery.PLATFORM_TYPE) (identity string, err error) {
	var counter int32
	counter, err = m.add(platform)
	if err != xyerror.ErrOK {
		return
	}

	// identity加了一个BASE_IDENTITY_SEGMENT是为了保留一些靓号段 :=)
	identity = fmt.Sprintf("%s%d", m.prefix, BASE_IDENTITY_SEGMENT+counter)

	return
}

// 增加计数
func (m *UserIdentityManager) add(platform battery.PLATFORM_TYPE) (counter int32, err error) {
	var (
		index = int(platform)
	)

	m.lock.Lock()
	defer m.lock.Unlock()

	// tothink 假如觉得每次增加计数都写数据库，写操作过于频繁，可以优化为增加N个玩家或多长时间未刷新后，写一次库，以此来减少写库操作。
	// 但是需要考虑服务异常crash后，内存的计数来不及持久化到数据库的情况，对内存中的计数做修正。
	err = DefCacheDB.IncreaseUserIdentityCounter(m.prefix, index)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("DefCacheDB.IncreaseUserIdentityCounter for (%s, %d) failed : %v", m.prefix, index, err)
		return
	}
	(m.counter[index])++
	counter = m.counter[index]
	return

}
