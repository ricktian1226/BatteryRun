package xyhttpservice

import (
	"fmt"
	martini "github.com/codegangsta/martini"
	xylog "guanghuan.com/xiaoyao/common/log"
	"net"
	"net/http"
	"sync"
	"time"
)

type MartiniService struct {
	name          string
	isRunning     bool
	port          int
	host          string
	encrypt       int
	encoding      int
	msvr          *martini.Martini
	router        martini.Router
	http_listener net.Listener
	done          chan bool
	start         chan bool

	cur_request_count int // 当前的请求数量
	max_request       int // 允许同时处理的最大请求数
	request_mutex     sync.Mutex

	cur_timeout_request int // 当前连续超时的请求数
	max_timeout_request int // 允许连续超时的最大请求数
	max_request_time    int // 单个请求处理时间的警戒线 (ms)
	request_time_mutex  sync.Mutex
}

func NewMartiniService(name string, host string, port int) (svc *MartiniService) {
	svc = &MartiniService{
		name:   name,
		host:   host,
		port:   port,
		msvr:   martini.New(),
		router: martini.NewRouter(),
		done:   make(chan bool, 1),
	}

	return
}

func DefaultMartiniService(name string, host string, port int) (svc *MartiniService) {
	svc = NewMartiniService(name, host, port)
	svc.AddMartiniHandler(martini.Recovery())
	svc.AddMartiniHandler(svc.Logger())
	svc.AddRouter(HttpGet, "/ping", Ping())
	return
}

func (svc *MartiniService) Logger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		xylog.DebugNoId("[%s] Started %s %s", svc.Name(), req.Method, req.URL.Path)

		start := time.Now()
		rw := res.(martini.ResponseWriter)
		c.Next()

		xylog.DebugNoId("[%s] %s Completed %v %s in %v\n", svc.Name(), req.URL.Path, rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}

func Ping() martini.Handler {
	return func() (int, string) {
		return http.StatusOK, ""
	}
}

// 开启流量控制
// max_flow: 允许同时处理的最多请求数
// max_request_time: 单个请求处理时间的警戒线 (ms)
// max_timeout_request: 允许连续超时的请求数量
func (svc *MartiniService) EnableFlowControl(max_request int, max_request_time int, max_timeout_request int) {
	if max_request > 0 {
		svc.max_request = max_request
		xylog.InfoNoId("[%s] enable flow control: max reqeusts = %d", svc.Name(), svc.max_request)
		svc.AddMartiniHandler(svc.FlowControlHandler())
		if max_request_time > 0 {
			svc.max_request_time = max_request_time
			if max_timeout_request > max_request {
				xylog.WarningNoId("max timeout request(%d) > max requests(%d)", max_timeout_request, max_request)
				svc.max_timeout_request = max_request
			} else {
				svc.max_timeout_request = max_timeout_request
			}

			xylog.InfoNoId("[%s] enable flow control: timeout = %d ms, max timeout requests = %d",
				svc.Name(), svc.max_request_time, svc.max_timeout_request)
			svc.AddMartiniHandler(svc.RequestTimeHandler())
		}
	}
}

// 请求数控制
func (svc *MartiniService) FlowControlHandler() martini.Handler {
	return func(w http.ResponseWriter, c martini.Context) {
		//var (
		//	reject bool
		//)
		//svc.request_mutex.Lock()
		//svc.cur_request_count++
		//xylog.DebugNoId("Request pool: %d", svc.cur_request_count)
		//if svc.cur_request_count > svc.max_request {
		//	xylog.WarningNoId("[%s] too many requests: %d ! (max: %d), rejecting requests",
		//		svc.Name(), svc.cur_request_count, svc.max_request)
		//	reject = true
		//}
		//svc.request_mutex.Unlock()
		//if reject {
		//	svc.RejectRequest(w, "Too many requests, try again later")
		//} else {
		c.Next()
		//}
		//svc.request_mutex.Lock()
		//svc.cur_request_count--
		//xylog.DebugNoId("Request pool: %d", svc.cur_request_count)
		//svc.request_mutex.Unlock()
		return
	}
}

// 请求超时控制
func (svc *MartiniService) RequestTimeHandler() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, c martini.Context) {
		//var (
		//	start   time.Time
		//	dur     time.Duration
		//	reject  bool
		//	timeout bool
		//	uri     = r.URL.Path
		//)
		//svc.request_mutex.Lock()
		//xylog.DebugNoId("cur timeout request: %d", svc.cur_timeout_request)
		//if svc.cur_request_count > svc.max_timeout_request && // 当前在处理的请求 > 允许超时的最大请求数 (请求池)
		//	svc.cur_timeout_request > svc.max_timeout_request { // 当前超时的请求数 > 允许超时的请求数
		//	xylog.WarningNoId("[%s] slow processor, start rejecting request", svc.Name())
		//	reject = true
		//}
		//svc.request_mutex.Unlock()
		//if reject {
		//	svc.RejectRequest(w, "Request timeout, try again later")
		//} else {
		//start = time.Now()
		c.Next()
		//dur = time.Since(start)
		//if dur > time.Duration(svc.max_request_time)*time.Millisecond {
		//	xylog.WarningNoId("[%s] request <%s> takes %.03f ms", svc.Name(), uri, float64(dur)/float64(time.Millisecond))
		//	timeout = true
		//}
		//svc.request_time_mutex.Lock()
		//if timeout {
		//	svc.cur_timeout_request++
		//} else {
		//	if svc.cur_timeout_request > 0 {
		//		svc.cur_timeout_request--
		//	}
		//}
		//xylog.DebugNoId("cur timeout request: %d", svc.cur_timeout_request)
		//svc.request_time_mutex.Unlock()
		//}
		return
	}
}

// 拒绝请求 (客户端排队)
func (svc *MartiniService) RejectRequest(w http.ResponseWriter, msg string) {
	xylog.DebugNoId("Reject request")

	http.Error(w, msg, http.StatusRequestTimeout)
}

/////////////////////////////////////////////////////////////
//func (svc *MartiniService) SetEncryptMethod() {

//}

//func (svc *MartiniService) SetEncodingMethod() {

//}

// 添加 martini 处理函数
func (svc *MartiniService) AddMartiniHandler(h martini.Handler) {
	svc.msvr.Use(h)
}

// 添加 url的处理函数
func (svc *MartiniService) AddRouter(op HttpOp, url_pattern string, h ...martini.Handler) {
	if svc.router == nil {
		xylog.WarningNoId("Martini router is not initialized")
		return
	}
	switch op {
	case HttpAny:
		svc.router.Any(url_pattern, h...)
	case HttpGet:
		svc.router.Get(url_pattern, h...)
	case HttpPost:
		svc.router.Post(url_pattern, h...)
	case HttpDelete:
		svc.router.Delete(url_pattern, h...)
	case HttpPatch:
		svc.router.Patch(url_pattern, h...)
	case HttpOptions:
		svc.router.Options(url_pattern, h...)
	case HttpPut:
		svc.router.Put(url_pattern, h...)
	case HttpHead:
		svc.router.Head(url_pattern, h...)
	}
}

// 添加静态网页路径
func (svc *MartiniService) SetStaticFilePath(path string) {
	svc.AddMartiniHandler(martini.Static(path))
}

// 添加 404 错误的处理函数
func (svc *MartiniService) AddNotFoundHandler(h ...martini.Handler) {
	if svc.router == nil {
		xylog.WarningNoId("Martini router is not initialized")
		return
	}
	svc.router.NotFound(h...)
}

/////////////////////////////////////////////
func (svc *MartiniService) Name() string {
	return svc.name
}
func (svc *MartiniService) IsRunning() bool {
	return svc.isRunning
}
func (svc *MartiniService) Init() (err error) {
	return
}
func (svc *MartiniService) Start() (err error) {
	if svc.IsRunning() {
		return
	}
	svc.msvr.Action(svc.router.Handle)

	hp := fmt.Sprintf("%s:%d", svc.host, svc.port)

	l, err := net.Listen("tcp", hp)
	if err != nil {
		xylog.ErrorNoId("Can't listen to (%s): %s", hp, err.Error())
		return
	}
	svc.http_listener = l

	http_server := &http.Server{
		Addr:    hp,
		Handler: svc.msvr,
	}

	go func() {
		xylog.InfoNoId("MartiniServer [%s] started (%s)", svc.Name(), hp)
		http_server.Serve(svc.http_listener)
		xylog.InfoNoId("MartiniServer [%s] stopped", svc.Name())
		svc.done <- true
	}()

	svc.isRunning = true
	xylog.InfoNoId("MartiniService [%s] started", svc.Name())
	return
}
func (svc *MartiniService) Stop() (err error) {
	if svc.IsRunning() {
		if svc.http_listener != nil {
			svc.http_listener.Close()
			svc.http_listener = nil
			xylog.DebugNoId("waiting for martini server to stop")
			<-svc.done
		}
		svc.isRunning = false
	}
	xylog.InfoNoId("MartiniService [%s] stopped", svc.Name())
	return
}
