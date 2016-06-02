// Copyright 2012-2014 Apcera Inc. All rights reserved.

// +build !windows

package logger

import (
	"fmt"
	"log"
	"log/syslog"
	"net/url"
)

type SysLogger struct {
	writer *syslog.Writer
	debug  bool
	trace  bool
}

func NewSysLogger(debug, trace bool) *SysLogger {
	w, err := syslog.New(syslog.LOG_DAEMON|syslog.LOG_NOTICE, "gnatsd")
	if err != nil {
		log.Fatalf("error connecting to syslog: %q", err.Error())
	}

	return &SysLogger{
		writer: w,
		debug:  debug,
		trace:  trace,
	}
}

func NewRemoteSysLogger(fqn string, debug, trace bool) *SysLogger {
	network, addr := getNetworkAndAddr(fqn)
	w, err := syslog.Dial(network, addr, syslog.LOG_DEBUG, "gnatsd")
	if err != nil {
		log.Fatalf("error connecting to syslog: %q", err.Error())
	}

	return &SysLogger{
		writer: w,
		debug:  debug,
		trace:  trace,
	}
}

func getNetworkAndAddr(fqn string) (network, addr string) {
	u, err := url.Parse(fqn)
	if err != nil {
		log.Fatal(err)
	}

	network = u.Scheme
	if network == "udp" || network == "tcp" {
		addr = u.Host
	} else if network == "unix" {
		addr = u.Path
	} else {
		log.Fatalf("error invalid network type: %q", u.Scheme)
	}

	return
}

func (l *SysLogger) Noticef(format string, v ...interface{}) {
	l.writer.Notice(fmt.Sprintf(format, v...))
}

func (l *SysLogger) Fatalf(format string, v ...interface{}) {
	l.writer.Crit(fmt.Sprintf(format, v...))
}

func (l *SysLogger) Errorf(format string, v ...interface{}) {
	l.writer.Err(fmt.Sprintf(format, v...))
}

func (l *SysLogger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.writer.Debug(fmt.Sprintf(format, v...))
	}
}

func (l *SysLogger) Tracef(format string, v ...interface{}) {
	if l.trace {
		l.writer.Notice(fmt.Sprintf(format, v...))
	}
}
