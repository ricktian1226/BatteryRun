package xyservice

import (
	xylog "guanghuan.com/xiaoyao/common/log"
)

type Service interface {
	Name() string
	IsRunning() bool
	Init() (err error)
	Start() (err error)
	Stop() (err error)
}

type DefaultService struct {
	name       string
	is_running bool
}

func NewDefaultService(name string) (svc *DefaultService) {
	svc = &DefaultService{
		name: name,
	}
	return svc
}
func (svc *DefaultService) Name() string {
	return svc.name
}

func (svc *DefaultService) Init() (err error) {
	xylog.InfoNoId("Service [%s] init", svc.Name())
	return
}

func (svc *DefaultService) Start() (err error) {
	xylog.InfoNoId("Service [%s] start", svc.Name())
	svc.is_running = true
	return
}

func (svc *DefaultService) Stop() (err error) {
	if svc.is_running {
		xylog.InfoNoId("Service [%s] stop", svc.Name())
		svc.is_running = false
	}
	return
}

func (svc *DefaultService) IsRunning() bool {
	return svc.is_running
}
