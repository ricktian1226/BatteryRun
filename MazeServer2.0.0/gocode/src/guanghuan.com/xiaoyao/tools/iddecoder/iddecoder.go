package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	TIMESTAMP_MASK uint64 = 0xffffffffffc00000 //42bit
	DcID_MASK      uint64 = 0x00000000003c0000 //4bit
	NodeID_MASK    uint64 = 0x000000000003f000 //6bit
	COUNTER_MASK   uint64 = 0x0000000000000fff //12bit
)

const (
	SHIFT_BITS_TIMESTAMP = 64 - 42
	SHIFT_BITS_DcID      = 64 - 42 - 4
	SHIFT_BITS_NodeID    = 64 - 42 - 4 - 6
	SHIFT_BITS_Counter   = 64 - 42 - 4 - 6 - 12
)

func main() {
	var d uint64
	flag.Uint64Var(&d, "d", 0, "number to decode")
	flag.Parse()

	if d > 0 {
		timestamp := (d & TIMESTAMP_MASK) >> SHIFT_BITS_TIMESTAMP
		dcid := (d & DcID_MASK) >> SHIFT_BITS_DcID
		nodeid := (d & NodeID_MASK) >> SHIFT_BITS_NodeID
		counter := (d & COUNTER_MASK) >> SHIFT_BITS_Counter
		timestamp *= uint64(time.Millisecond)
		begin, _ := time.Parse("2006-01-02 15:04:05", "2014-08-14 14:37:18")
		nsec := begin.UnixNano()
		timestamp += uint64(nsec)

		timestr := time.Unix(int64(timestamp)/int64(time.Second), int64(timestamp)%int64(time.Second)).String()

		fmt.Printf("id stands : time : %s , dcid : %d, nodeid : %d, counter : %d",
			timestr,
			dcid,
			nodeid,
			counter)
	}
}
