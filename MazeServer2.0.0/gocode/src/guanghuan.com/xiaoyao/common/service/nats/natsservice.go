package xynatsservice

import (
	proto "code.google.com/p/goprotobuf/proto"
	nats "github.com/nats-io/nats"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyservice "guanghuan.com/xiaoyao/common/service"
	"reflect"
	"strings"
	"time"
)

type ApiHandler func(req proto.Message, resp proto.Message) (err error)

type MsgHandler struct {
	Subject  string
	ReqType  reflect.Type
	RespType reflect.Type
	Handler  ApiHandler
	Nats     *NatsService
}

type MsgHandlerMap map[string]*MsgHandler

func (mm MsgHandlerMap) AddHandler(subject string, req interface{}, resp interface{}, h ApiHandler) {
	m := map[string]*MsgHandler(mm)
	m[subject] = &MsgHandler{
		Subject:  subject,
		ReqType:  reflect.TypeOf(req),
		RespType: reflect.TypeOf(resp),
		Handler:  h,
	}
	xylog.InfoNoId("Adding handler: subj=[%s], req=%s, resp=%s", subject, reflect.TypeOf(req).Name(), reflect.TypeOf(resp).Name())
}

func (mm MsgHandlerMap) GetHandler(subject string) (h *MsgHandler) {
	m := map[string]*MsgHandler(mm)
	return m[subject]
}
func (mm MsgHandlerMap) Size() int {
	m := map[string]*MsgHandler(mm)
	return len(m)
}

type MsgCodeHandler struct {
	Code     uint32
	ReqType  reflect.Type
	RespType reflect.Type
	Handler  ApiHandler
	Nats     *NatsService
}

type MsgCodeHandlerMap map[uint32]*MsgCodeHandler

func (mm MsgCodeHandlerMap) AddHandler(code uint32, req interface{}, resp interface{}, h ApiHandler) {
	mm[code] = &MsgCodeHandler{
		Code:     code,
		ReqType:  reflect.TypeOf(req),
		RespType: reflect.TypeOf(resp),
		Handler:  h,
	}
	//xylog.Info("Adding handler: code=[%d], req=%s, resp=%s", code, reflect.TypeOf(req).Name(), reflect.TypeOf(resp).Name())
}

func (mm MsgCodeHandlerMap) GetHandler(code uint32) (h *MsgCodeHandler) {
	return mm[code]
}
func (mm MsgCodeHandlerMap) Size() int {
	return len(mm)
}

type NatsSubHandler struct {
	subject string
	handler nats.MsgHandler
	//subcription *nats.Subscription
	subcription nats.Subscription
}

type NatsQueueHandler struct {
	queue string
	group string
	//subcription *nats.Subscription
	subcription nats.Subscription
	handler     nats.MsgHandler
}

// a wraper class of NATS operations
type NatsService struct {
	xyservice.DefaultService
	server_url string
	nc         *nats.Conn
	//subscribers     []*NatsSubHandler
	subscribers []NatsSubHandler
	//queue_receivers []*NatsQueueHandler
	queue_receivers []NatsQueueHandler
}

//global var nats_service
var Nats_service *NatsService

func NewNatsService(name string, server_url string) (svc *NatsService) {
	svc = &NatsService{
		DefaultService: *xyservice.NewDefaultService(name),
		server_url:     server_url,
		//subscribers:     make([]*NatsSubHandler, 0, 10),
		//queue_receivers: make([]*NatsQueueHandler, 0, 10),
		subscribers:     make([]NatsSubHandler, 0, 10),
		queue_receivers: make([]NatsQueueHandler, 0, 10),
	}
	return
}

// 单向发送消息
func (svc *NatsService) Publish(subject string, msg []byte) (err error) {
	if svc.nc != nil && !svc.nc.IsClosed() {
		//		svc.nc.Subscribe(subject, handler)
		err = svc.nc.Publish(subject, msg)
		//		xylog.Debug("publish message: %s", subject)
	} else {
		xylog.WarningNoId("nats connection is not ready!")
	}
	return
}

// 发送消息，并指定返回的消息地址
func (svc *NatsService) PublishRequest(subject string, reply string, msg []byte) (err error) {
	if svc.nc != nil && !svc.nc.IsClosed() {
		err = svc.nc.PublishRequest(subject, reply, msg)
	}
	return
}

// 发送消息，并返回消息
func (svc *NatsService) Request(subject string, msg []byte, timeout time.Duration) (reply_msg *nats.Msg, err error) {
	if svc.nc != nil && !svc.nc.IsClosed() {
		reply_msg, err = svc.nc.Request(subject, msg, timeout)
		//		if err == nil {
		//			err = reply_handler(reply_msg)
		//		}
	}
	return
}

// 订阅指定消息
func (svc *NatsService) AddSubscriber(subject string, handler nats.MsgHandler) (err error) {
	if subject != "" && handler != nil {
		//h := &NatsSubHandler{
		//	subject: subject,
		//	handler: handler,
		//}
		h := NatsSubHandler{
			subject: subject,
			handler: handler,
		}

		svc.subscribers = append(svc.subscribers, h)
		if svc.nc != nil && !svc.nc.IsClosed() {
			var tmp *nats.Subscription
			//h.subcription, err = svc.nc.Subscribe(subject, handler)
			tmp, err = svc.nc.Subscribe(subject, handler)
			h.subcription = *tmp
			xylog.DebugNoId("add subscriber: %s", subject)
		}
	}
	return
}

// 订阅指定消息队列
func (svc *NatsService) AddQueueSubscriberGroup(queue string, group string, handler nats.MsgHandler) (err error) {
	if queue != "" && group != "" && handler != nil {
		//h := &NatsQueueHandler{
		//	queue:   queue,
		//	group:   group,
		//	handler: handler,
		//}
		h := NatsQueueHandler{
			queue:   queue,
			group:   group,
			handler: handler,
		}

		svc.queue_receivers = append(svc.queue_receivers, h)
		if svc.nc != nil && !svc.nc.IsClosed() {
			var tmp *nats.Subscription
			//h.subcription, err = svc.nc.QueueSubscribe(queue, group, handler)
			tmp, err = svc.nc.QueueSubscribe(queue, group, handler)
			h.subcription = *tmp
			xylog.DebugNoId("add queue group: %s (%s)", queue, group)
		}
	}
	return
}
func (svc *NatsService) AddQueueSubscriber(queue string, handler nats.MsgHandler) {
	svc.AddQueueSubscriberGroup(queue, queue, handler)
}

// 订阅指定的请求
func (svc *NatsService) AddRequestHandler(request string, handler nats.MsgHandler) {
	svc.AddQueueSubscriber(request, handler)
}

//func (svc *NatsService) Name() string {
//	return svc.defsvc.Name()
//}

func (svc *NatsService) Init() (err error) {
	return svc.DefaultService.Init()
}

// 启动服务
func (svc *NatsService) Start() (err error) {
	if svc.IsRunning() {
		xylog.DebugNoId("Service[%s] already running", svc.Name())
		return
	}

	//修改下nats的默认配置，先简单地写死，后续改成配置项
	nats.DefaultOptions.MaxReconnect = 200
	nats.DefaultOptions.ReconnectWait = time.Duration(2 * time.Second)

	//xylog.Debug("MaxReconnect %d, ReconnectWait %d seconds", nats.DefaultOptions.MaxReconnect, nats.DefaultOptions.ReconnectWait)

	opts := nats.DefaultOptions
	opts.Name = svc.DefaultService.Name()
	opts.Servers = strings.Split(svc.server_url, ",")
	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}
	xylog.InfoNoId("Connecting to NATS server: %s, opts : %v", svc.server_url, opts)
	svc.nc, err = opts.Connect()
	if err != nil {
		xylog.ErrorNoId("Can't connect: %v\n", err)
		return
	}
	for _, sub := range svc.subscribers {
		//if sub != nil {
		var tmp *nats.Subscription
		//sub.subcription, err = svc.nc.Subscribe(sub.subject, sub.handler)
		tmp, err = svc.nc.Subscribe(sub.subject, sub.handler)
		sub.subcription = *tmp
		xylog.InfoNoId("add subscriber: %s", sub.subject)
		//}
	}
	for _, queue := range svc.queue_receivers {
		//if queue != nil {
		//queue.subcription, err = svc.nc.QueueSubscribe(queue.queue, queue.group, queue.handler)
		var tmp *nats.Subscription
		tmp, err = svc.nc.QueueSubscribe(queue.queue, queue.group, queue.handler)
		queue.subcription = *tmp
		xylog.InfoNoId("add queue group: %s (%s)", queue.queue, queue.group)
		//}
	}
	err = svc.DefaultService.Start()
	return
}

// 停止服务
func (svc *NatsService) Stop() (err error) {
	if svc.IsRunning() {
		if svc.nc != nil && !svc.nc.IsClosed() {
			for _, sub := range svc.subscribers {
				//if sub != nil {
				sub.subcription.Unsubscribe()
				xylog.InfoNoId("remove subscriber: %s", sub.subject)
				//}
			}
			for _, queue := range svc.queue_receivers {
				//if queue != nil {
				queue.subcription.Unsubscribe()
				xylog.InfoNoId("remove queue group: %s (%s)", queue.queue, queue.group)
				//}
			}
			svc.nc.Close()
		}
		err = svc.DefaultService.Stop()
		//		svc.isRunning = false
	}
	return
}

// 服务状态
//func (svc *NatsService) IsRunning() bool {
//	return svc.defsvc.IsRunning()
//}
