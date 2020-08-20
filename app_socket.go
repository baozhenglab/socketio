package sckio

import (
	"fmt"
	"sync"

	goservice "github.com/baozhenglab/go-sdk"
	"github.com/baozhenglab/go-sdk/logger"
	"github.com/baozhenglab/sdkcm"
	socketio "github.com/googollee/go-socket.io"
)

type appSocket struct {
	locker *sync.RWMutex
	once   *sync.Once
	sio    *socketio.Server
	sc     goservice.ServiceContext
	cu     sdkcm.Requester
	logger logger.Logger
	Socket
	// Map room object id to speed up check socket has in a specific room
	roomMap map[string]bool
}

func (as appSocket) String() string {
	if as.cu != nil {
		return fmt.Sprintf("sck_id: %s - user_id: %d", as.Id(), as.cu.UserID())
	}

	return fmt.Sprintf("sck_id: %s - user_id: not logged in", as.Id())
}

func newAppSocket(sio *socketio.Server, sc goservice.ServiceContext, l logger.Logger, socket Socket) *appSocket {
	return &appSocket{
		locker:  new(sync.RWMutex),
		once:    new(sync.Once),
		sio:     sio,
		sc:      sc,
		logger:  l,
		Socket:  socket,
		roomMap: make(map[string]bool),
	}
}

func (as *appSocket) BroadcastToRoom(room, event string, args ...interface{}) {
	as.sio.BroadcastTo(room, event, args...)
}
func (as *appSocket) ServiceContext() goservice.ServiceContext { return as.sc }
func (as *appSocket) Logger() logger.Logger                    { return as.logger }
func (as *appSocket) CurrentUser() sdkcm.Requester             { return as.cu }
func (as *appSocket) SetCurrentUser(r sdkcm.Requester)         { as.cu = r }

func (as *appSocket) Join(room string) error {
	if err := as.Socket.Join(room); err != nil {
		return err
	}

	as.locker.Lock()
	as.roomMap[room] = true
	as.locker.Unlock()
	return nil
}

// Custom function for check socket has in room
func (as *appSocket) HasInRoom(room string) bool {
	as.locker.RLock()
	result, ok := as.roomMap[room]
	if !ok {
		return ok
	}
	as.locker.RUnlock()
	return result
}
