package xycache

import (
	"sync/atomic"
)

//缓存下标管理器
type CacheBase struct {
	index int32
}

//获取主缓存下标
func (c *CacheBase) Major() int32 {
	return atomic.LoadInt32(&c.index)
}

//获取备缓存下标
func (c *CacheBase) Secondary() int32 {
	return (atomic.LoadInt32(&c.index) + 1) % 2
}

//切换缓存
func (c *CacheBase) Switch() {
	atomic.StoreInt32(&c.index, (atomic.LoadInt32(&c.index)+1)%2)
}
