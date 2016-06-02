package xyutil

import (
	"fmt"
	//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//	xydb "guanghuan.com/xiaoyao/superbman_server/db/v1"
	//"log"
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"

	"guanghuan.com/xiaoyao/common/log"
)

const (
	SecondsPerDay = 60 * 60 * 24
)

// 当前时间 (秒)
func CurTimeSec() int64 {
	return time.Now().Unix()
}

// 当天的时间戳范围 (秒)
func CurTimeRangeSec() (int64, int64) {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix(), time.Date(year, month, day, 23, 59, 59, 0, time.Local).Unix()
}

// 当前时间 (纳秒)
func CurTimeNs() int64 {
	return time.Now().UnixNano()
}

// 当前时间 (微秒)
func CurTimeUs() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

// 当前时间 (毫秒)
func CurTimeMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// 当前时间的字符串格式：yyyymmddhhmmss
func CurTimeStr() (str string) {
	now := time.Now()
	str = fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	return
}

// 把数值时间转成文本时间
func ToStrTime(time_sec int64) (str string) {
	t := time.Unix(time_sec, 0)
	//	RFC3339 = "2006-01-02T15:04:05Z07:00"
	//	ChinaTime = "2006-01-02 14:04:05(08:00)""
	//	local, _ := time.LoadLocation("Local")
	str = t.Format(time.RFC3339)
	//	str = time.Parse(time.RFC3339, t)
	//	t, _ = time.ParseInLocation(time.RFC3339, str, local)
	//	str
	return
}

// 产生一个字符串id
func NewId() string {

	// 第一部分， 时间
	t := CurTimeNs()

	// 第二部分， 随机数
	rv := rand.Intn(1000)

	newid := fmt.Sprintf("%d%03d", t, rv)

	return newid
}

//|----------------------|-------|--------|
//         42 bit             10bits   12 bit
//
//高42位 时间戳，粒度为毫秒级别。137年后时间戳溢出
//中10位 服务节点id，其中的高4位为数据中心标识（gateway id）；低6位为节点标识（nodeid）。最多可以有16个gateway，每个gateway下可以挂64个nodeid
//低12位为计数器 最大支持每毫秒 4096个id分配
func NewUint64Id(timestamp, dcid, nodeid, counter int64) (newid uint64) {
	return uint64(timestamp<<22 + dcid<<18 + nodeid<<12 + counter)
}

var (
	ServerStartDate int64 // server start up date
)

func BytesToInt32(b []byte) (v int32, err error) {
	if len(b) > 4 {
		err = errors.New("[BytesToInt32] Invalid byte buffer")
		return
	}
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &v)
	return
}

//计算起止时间戳的间隔天数
// begin 开始时间戳
// end   结束时间戳
//return
// day 天数
func DayDiff(begin, end int64) (day int64) {
	timeBeginTmp := time.Unix(begin, 0)
	timeEndTmp := time.Unix(end, 0)
	timeBegin := time.Date(timeBeginTmp.Year(), timeBeginTmp.Month(), timeBeginTmp.Day(), 0, 0, 0, 0, time.Local).Unix()
	timeEnd := time.Date(timeEndTmp.Year(), timeEndTmp.Month(), timeEndTmp.Day(), 0, 0, 0, 0, time.Local).Unix()
	day = int64((timeEnd - timeBegin) / (24 * 60 * 60))
	xylog.DebugNoId("begin : %v, end : %v, day %d", timeBegin, timeEnd, day)
	return
}

//当天的起始时间戳
// now int64 当前时间戳
//return
// begin int64 当天0点的时间戳
// end   int64 当天24点的时间戳
func TodayTimeRange(now int64) (begin, end int64) {
	year, month, day := time.Unix(now, 0).Date()
	begin = time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
	end = begin + int64(SecondsPerDay)
	return
}
