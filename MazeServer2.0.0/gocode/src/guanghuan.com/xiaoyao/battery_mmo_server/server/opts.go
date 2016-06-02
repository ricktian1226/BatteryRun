// opts
package server

import (
	"fmt"
)

type Options struct {
	Host         string `json:"addr"`
	Port         int    `json:"port"`
	ConfigFile   string `json:"cfgfile"`
	MaxConn      int    `json:"maxconnections"`
	PingInterval int    `json:"pinginterval"`
	PingTimeOut  int    `json:"pingtimeout"`
	PidFile      string `json:"pidfile"`
	NoSigs       bool   `json:"nosigs"`
	RoundTime    int    `json:"roundtime"`
	RoundDelay   int    `json:"rounddelay"`
}

func (opts *Options) String() string {
	return fmt.Sprintf(`Options : host : %s 
	            port : %d
				maxconnection : %d
				pinginterval : %d
				pingtimeout : %d
				pidfile : %s
				nosigs : %v
				roundtime : %d
				rounddelay : %d`,
		opts.Host,
		opts.Port,
		opts.PingInterval,
		opts.PingTimeOut,
		opts.PidFile,
		opts.NoSigs,
		opts.RoundTime,
		opts.RoundDelay)
}

var GOpts = &Options{}
