// basic server for all kinds of services
package xyserver

import (
	xylog "guanghuan.com/xiaoyao/common/log"
	xyservice "guanghuan.com/xiaoyao/common/service"
	xydbservice "guanghuan.com/xiaoyao/common/service/db"
	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	xyversion "guanghuan.com/xiaoyao/common/version"
	"os"
	"os/signal"
	//"runtime"
	"time"
)

type Options struct {
	NoSigs    bool
	StopOnErr bool
}

type ServiceMap map[string]xyservice.Service
type ServiceList []xyservice.Service

type Server struct {
	start     time.Time
	Name      string
	isRunning bool
	done      chan bool
	exit      chan bool
	opts      *Options
	services  ServiceList
	version   *xyversion.Version
}

const (
	DefaultDBServiceName      string = "DB"
	DefaultNatsServiceName    string = "NATS"
	DefaultMartiniServiceName string = "Martini"
)

func New(name string, opts *Options) *Server {
	s := &Server{
		Name:      name,
		start:     time.Now(),
		done:      make(chan bool, 1),
		opts:      opts,
		isRunning: false,
		services:  ServiceList(make([]xyservice.Service, 0, 10)),
		version:   xyversion.New(1, 0, 0),
	}
	return s
}

func (s *Server) RegisterService(svr_name string, svr xyservice.Service) {
	if svr != nil {
		s.services = append(s.services, svr)
		xylog.InfoNoId("Register Service: %s", svr_name)
		if s.isRunning {
			svr.Start()
		}
	}
}

func (s *Server) QuickRegService(svr xyservice.Service) {
	s.RegisterService(svr.Name(), svr)
}

func (s *Server) Init() (err error) {
	xylog.InfoNoId("========================")
	xylog.InfoNoId("Server [%s] init start ...", s.Name)
	for k, svr := range s.services {
		xylog.InfoNoId("------------------------")
		xylog.InfoNoId("#%02d.[%s] init ...", k+1, svr.Name())
		err = svr.Init()
		if err != nil {
			xylog.ErrorNoId("#%02d.[%s] init failed, ERROR: %s",
				k+1, svr.Name(), err.Error())
			if s.opts.StopOnErr {
				xylog.InfoNoId("Stop on error")
				break
			} else {
				xylog.InfoNoId("Ignore error")
				err = nil
			}
		} else {
			xylog.InfoNoId("#%02d.[%s] init done", k+1, svr.Name())
		}
	}
	xylog.InfoNoId("Server [%s] init done", s.Name)
	return
}

func (s *Server) Start() (err error) {
	err = s.Init()
	if err != nil {
		if s.opts.StopOnErr {
			xylog.InfoNoId("Stop on error")
			return
		} else {
			xylog.InfoNoId("Ignore error")
		}
	}
	xylog.InfoNoId("========================")
	xylog.InfoNoId("Server [%s] is starting ...", s.Name)
	s.handleSignals()
	for k, svr := range s.services {
		xylog.InfoNoId("------------------------")
		xylog.InfoNoId("#%02d.[%s] starting ...", k+1, svr.Name())
		err = svr.Start()
		if err != nil {
			xylog.ErrorNoId("#%02d.[%s] failed to start, ERROR: %s",
				k+1, svr.Name(), err.Error())
			if s.opts.StopOnErr {
				xylog.InfoNoId("Stop on error")
				break
			} else {
				xylog.InfoNoId("Ignore error")
				err = nil
			}
		} else {
			xylog.InfoNoId("#%02d.[%s] started", k+1, svr.Name())
		}
	}

	if err == nil {
		s.isRunning = true
	}
	//	if !s.opts.NoSigs {
	//go s.Run()
	//	}

	//runtime.Goexit()
	return
}

func (s *Server) Shutdown() {
	xylog.InfoNoId("========================")
	xylog.InfoNoId("Server [%s] stopping ...", s.Name)
	for i := len(s.services) - 1; i >= 0; i-- {
		if s.services != nil {
			svr := s.services[i]
			xylog.InfoNoId("------------------------")
			xylog.InfoNoId("#%02d.[%s] stopping ...", i+1, svr.Name())
			(svr).Stop()
			xylog.InfoNoId("#%02d.[%s] stopped", i+1, svr.Name())
		}
	}
	xylog.InfoNoId("Server [%s] stopped", s.Name)
	xylog.Close()
	os.Exit(0)
}

func (s *Server) Run() {
	xylog.InfoNoId("Server [%s] is running now", s.Name)
	xylog.InfoNoId("========================")
	<-s.done
	xylog.InfoNoId("Server [%s] is about to stop", s.Name)
	s.Shutdown()
}

// Signal Handling
func (s *Server) handleSignals() {
	if s.opts.NoSigs {
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		//for sig := range c {
		//	//			xylog.Debug("Trapped Signal; %-v", sig)

		//	xylog.Info("Server Exiting..")
		//	s.SendShutdownSignal()
		//}
		<-c
		xylog.InfoNoId("Server Exiting..")
		s.SendShutdownSignal()
	}()
}

func (s *Server) SendShutdownSignal() {
	s.done <- true
}

// Add Default DB Service
func (s *Server) EnableDBService(db_url string, db_name string) (svc *xydbservice.DBService) {
	svc = xydbservice.NewXYDBService(DefaultDBServiceName, db_url, db_name)
	s.RegisterService(svc.Name(), svc)
	return svc
}

// Add Default Nats Service
func (s *Server) EnableNatsService(nsurl string) (svc *xynatsservice.NatsService) {
	svc = xynatsservice.NewNatsService(DefaultNatsServiceName, nsurl)
	s.RegisterService(svc.Name(), svc)
	return svc
}

// Add Default Martini Http service
func (s *Server) EnableMartiniService(host string, port int) (svc *xyhttpservice.MartiniService) {
	svc = xyhttpservice.DefaultMartiniService(DefaultMartiniServiceName, host, port)
	s.RegisterService(svc.Name(), svc)
	return svc
}

func (s *Server) Version() string {
	return s.version.String()
}
func (s *Server) SetVersion(x, y, z int32) string {
	s.version.Set(x, y, z)
	return s.version.String()
}
