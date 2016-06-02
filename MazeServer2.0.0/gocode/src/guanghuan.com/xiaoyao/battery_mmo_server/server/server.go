//server
package server

import (
	"fmt"
	xylogs "guanghuan.com/xiaoyao/common/log"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	mu       sync.RWMutex //用读写锁，应该性能会比互斥锁高些
	clients  map[ClientID]*Client
	gcid     ClientID //client id
	running  uint32
	listener net.Listener
	done     chan bool
	stats
}

func New(opts *Options) *Server {

	s := &Server{
		done:    make(chan bool, 1),
		clients: make(map[ClientID]*Client),
		running: 1,
	}

	s.handleSignals()

	return s
}

func PrintAndDie(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

func (s *Server) isRunning() bool {
	return 0 != atomic.LoadUint32(&s.running)
}

func (s *Server) handleSignals() {
	if GOpts.NoSigs {
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			xylogs.Error("Trapped Signal; %v", sig)
			xylogs.Error("Server Exiting..")
			os.Exit(0)
		}
	}()
}

func (s *Server) logPid() {
	pidStr := strconv.Itoa(os.Getpid())
	err := ioutil.WriteFile(GOpts.PidFile, []byte(pidStr), 0660)
	if err != nil {
		fmt.Printf("Could not write pidfile: %v\n", err)
	}
}

func (s *Server) Start() {

	if GOpts.PidFile != "" {
		s.logPid()
	}

	s.AcceptLoop()

}

func (s *Server) Shutdown() {
	s.mu.Lock()

	if 0 == atomic.LoadUint32(&s.running) {
		s.mu.Unlock()
		return
	}

	atomic.StoreUint32(&s.running, 0)

	// Copy off the clients
	clients := make(map[ClientID]*Client)
	for i, c := range s.clients {
		clients[i] = c
	}

	// Kick client AcceptLoop()
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}

	s.mu.Unlock()

	// Close client connections
	for _, c := range clients {
		c.Close()
	}

	// Block until the accept loops exit
	<-s.done
}

func (s *Server) AcceptLoop() {
	hp := fmt.Sprintf(":%d", GOpts.Port)
	l, e := net.Listen("tcp", hp)
	if e != nil {
		//todo : log critical
		xylogs.Error("net.Listen error,program will exit.")
		return
	}

	xylogs.Info("Battery_mmo_server ready.")

	s.listener = l

	tmpDelay := ACCEPT_MIN_SLEEP

	for s.isRunning() {
		conn, err := l.Accept()
		xylogs.Debug("Accept from " + conn.RemoteAddr().String())
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				//todo log Temporary client accept error
				xylogs.Error("Temporary client accept error")
				time.Sleep(tmpDelay)
				tmpDelay *= 2
				if tmpDelay > ACCEPT_MAX_SLEEP {
					tmpDelay = ACCEPT_MAX_SLEEP
				}
			} else if s.isRunning() {
				//todo log accept error
				xylogs.Error("Server accept error")
			}
			continue
		}
		tmpDelay = ACCEPT_MIN_SLEEP
		s.createClient(conn)
	}

	xylogs.Debug("battery_mmo_server exiting.")
	s.done <- true
}

//创建会话信息
func (s *Server) createClient(conn net.Conn) *Client {
	c := &Client{
		srv:      s,
		nc:       conn,
		inqueue:  make(chan MsgEntry, 100),
		outqueue: make(chan MsgEntry, 100),
		room:     nil,
	}

	// Initialize
	c.initClient()

	// Register with the server.
	s.mu.Lock()
	s.clients[c.cid] = c
	s.mu.Unlock()

	return c
}

//删除会话信息
func (s *Server) removeClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//删除客户端列表
	delete(s.clients, c.cid)
}
