package xyidgenerate

import (
	"fmt"
	xyutil "guanghuan.com/xiaoyao/common/util"
	"sync"
	"sync/atomic"
	"time"
)

//id生成器的时间戳起点是2014-08-14 14:37:18
//用这个时间是因为代码是这个时间点写的~ o(∩_∩)o
//单位是毫秒
var DefIdGenerateBeginTimeStamp = IdGenerateBeginTimeStamp()

func IdGenerateBeginTimeStamp() int64 {
	t, err := time.Parse("2006-01-02 15:04:05", "2014-08-14 14:37:18")
	if err != nil {
		return 0
	}

	return t.UnixNano() / int64(time.Millisecond)
}

type IdGenerater struct {
	m         sync.Mutex
	name      string //业务名称
	from      int64  //起始时间
	dcid      int64  //数据中心标识
	nodeid    int64  //节点标识
	counter   int64  //计数器
	timestamp int64  //计数器有效时间戳
	//-------------------------------
	idsum          int64 //分配的id总数
	begintimestamp int64 //启动分配的时间戳
}

func NewIdGenerater(from, dcid, nodeid int64, name string) *IdGenerater {
	now := xyutil.CurTimeMs()
	return &IdGenerater{
		name:           name,
		from:           from,
		dcid:           dcid,
		nodeid:         nodeid,
		timestamp:      now,
		begintimestamp: now,
		counter:        0,
	}
}

//生成新的整型(uint64) id
func (ig *IdGenerater) NewID() (id uint64) {

	//如果时间戳变化，重置计数和时间戳
	now := xyutil.CurTimeMs()
	ig.m.Lock()
	if ig.timestamp < now {
		ig.counter = 0
		ig.timestamp = now
	}

	id = xyutil.NewUint64Id(ig.timestamp-ig.from, ig.dcid, ig.nodeid, ig.counter)

	ig.counter++
	ig.m.Unlock()

	//分配的id总数+1
	atomic.StoreInt64(&(ig.idsum), atomic.AddInt64(&(ig.idsum), 1))

	return
}

//线程安全的，不用加锁
func (ig *IdGenerater) String() string {
	idsum := atomic.LoadInt64(&(ig.idsum))
	timeDiff := atomic.LoadInt64(&(ig.timestamp)) - atomic.LoadInt64(&(ig.begintimestamp))
	var avg int64
	if timeDiff > 0 {
		avg = idsum / timeDiff
	}

	return fmt.Sprintf(`
	    %s (%d:%d)IdGenerater :
	        idsum : %d,
			time  : %d ms,
			avg   : %d /ms
	`, ig.name, ig.dcid, ig.nodeid, idsum, timeDiff, avg)
}
