// room
package server

import (
	proto "code.google.com/p/goprotobuf/proto"
	xymmopb "guanghuan.com/xiaoyao/battery_mmo_server/pb"
	xylogs "guanghuan.com/xiaoyao/common/log"
	"sync"
	"sync/atomic"
	"time"
)

type RoomID uint32

var gRoom = NewRoom()

type Room struct {
	mu         sync.RWMutex
	rid        RoomID               //room的标识
	clients    map[ClientID]*Client //room对应的client列表
	start      bool                 //是否已经开始
	startTime  int64                //游戏开始时间
	roundTimer *time.Timer          //回合计时器
}

func NewRoom() *Room {
	return &Room{
		rid:     1,
		clients: make(map[ClientID]*Client),
		start:   false,
	}
}

//增加
func (r *Room) AddClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	//校验一下是否重复插入
	if _, ok := r.clients[c.cid]; ok {
		xylogs.Debug("Duplicate Client#%d attend to insert into Room#%d", c.cid, r.rid)
	} else {
		xylogs.Debug("Client#%d insert into Room#%d", c.cid, r.rid)
		r.clients[c.cid] = c
		c.JoinRoom(r)
	}
}

func (r *Room) RemoveClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeClient(c)
}

//内部使用，非多线程安全的
func (r *Room) removeClient(c *Client) {
	delete(r.clients, c.cid)
	//如果房间中没有玩家了，就停止房间
	if 0 >= len(r.clients) {
		r.start = false
	}
}

//本轮结束，广播给各个玩家，销毁各个玩家的房间信息
func (r *Room) TimeOut() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.start = false
	r.roundTimer.Stop()
	r.roundTimer = nil

	var msg MsgEntry
	for _, c := range r.clients {
		msg.Cmd = (uint16)(xymmopb.CMD_CODE_CMD_S_GameOver)
		msg.MsgBuf = nil
		c.outqueue <- msg
		r.removeClient(c)
		c.LeaveRoom()
	}
}

func (r *Room) Broadcast(msg *MsgEntry) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	xylogs.Debug("Room#%d There are %d clients", r.rid, len(r.clients))
	for _, c := range r.clients {
		c.outqueue <- *msg
	}
}

func (r *Room) Start() (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	//房间游戏启动策略(暂时是：玩家超过两个的时候才开始)
	//后续策略：玩家超过两个，或者超过一定时间后自动开始
	xylogs.Debug("Room#%d start?(%v) ,len(r.clients) : %d ", r.rid, r.start, len(r.clients))
	if r.start == false && len(r.clients) >= 2 {
		//增加一个定时器，时间到了以后就通知各个客户端游戏结束
		r.roundTimer = time.AfterFunc(time.Duration(GOpts.RoundTime+GOpts.RoundDelay)*(time.Second), func() { r.TimeOut() })
		r.start = true
		r.startTime = (time.Now().UnixNano() + (int64)(GOpts.RoundDelay)*((int64)(time.Second))) / (int64)(time.Microsecond) //开始时间，单位微秒
	}

	//房间游戏才向玩家广播开始消息
	if r.start {
		req := &xymmopb.Msg_S_GameStart{
			Map:       &xymmopb.Map{Id: proto.Uint32(1)}, //map默认设为1
			Time:      proto.Uint64((uint64)(r.startTime)),
			Timedelay: proto.Uint64((uint64)(time.Duration(GOpts.RoundDelay) * (time.Second / time.Microsecond))),
		}

		//把目前房间的玩家列表向所有玩家广播
		for cid, _ := range r.clients {
			uid := proto.Uint32(uint32(cid))
			req.Players = append(req.Players, &xymmopb.Player{
				Uid: uid,
			})
		}
		xylogs.Debug("Recently players will send to client : %v.", req.Players)
		var msgOut MsgEntry
		var err error
		msgOut.MsgBuf, err = proto.Marshal(req)
		if err != nil {
			xylogs.Debug("proto.Marshal failed.")
			return false, err
		}

		msgOut.Cmd = (uint16)(xymmopb.CMD_CODE_CMD_S_GameStart)

		for _, c := range r.clients {
			c.outqueue <- msgOut
		}
	}

	return r.start, nil
}

func (r *Room) Empty() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.clients) <= 0
}

func (r *Room) Full() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.clients) > 4 //先默认定为4个玩家一个房间
}

type RoomMgr struct {
	mu    sync.Mutex
	rooms map[RoomID]*Room
	grid  uint32 //room id
}

func (rmgr *RoomMgr) CreateRoom() *Room {
	room := &Room{
		rid:     (RoomID)(atomic.AddUint32(&rmgr.grid, 1)),
		clients: make(map[ClientID]*Client),
	}

	return room
}

func (rmgr *RoomMgr) AddClient(c *Client) {
	addOk := false
	rmgr.mu.Lock()
	defer rmgr.mu.Unlock()
	for _, r := range rmgr.rooms {
		//给room加上锁
		r.mu.Lock()
		defer r.mu.Unlock()

		if !r.Full() {
			r.AddClient(c)
			addOk = true
			break
		}
	}

	if !addOk {
	}
}
