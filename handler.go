package sckio

import (
	"sync"

	goservice "github.com/baozhenglab/go-sdk"
	"github.com/baozhenglab/go-sdk/logger"
	socketio "github.com/googollee/go-socket.io"
)

// SocketIO Event Keys
const (
	SckDisconnected       = "disconnection"
	SckAuthenticate       = "SckAuthenticate"
	SckAuthenticated      = "SckAuthenticated"
	SckAuthenticateFailed = "SckAuthenticateFailed"
	SckChatMessage        = "SckChatMessage"
	SckChatRoomCreated    = "SckChatRoomCreated"
)

type socketHandler struct {
	// map user id to socket sessions
	uidSckMap map[uint32][]AppSocket
	locker    *sync.RWMutex
	server    *socketio.Server
}

func (hdl *socketHandler) BroadcastTo(room, event string, args ...interface{}) {
	hdl.server.BroadcastTo(room, event, args...)
}

func (hdl *socketHandler) UserDisconnect(uid uint32) {
	hdl.locker.Lock()
	delete(hdl.uidSckMap, uid)
	hdl.locker.Unlock()
}

func (hdl *socketHandler) UserConnect(uid uint32, as AppSocket) {
	hdl.locker.Lock()
	defer hdl.locker.Unlock()

	if s, ok := hdl.uidSckMap[uid]; ok {
		hdl.uidSckMap[uid] = append(s, as)
		return
	}
	hdl.uidSckMap[uid] = []AppSocket{as}
}

func (hdl *socketHandler) Sockets(uids []uint32) []AppSocket {
	hdl.locker.RLock()
	defer hdl.locker.RUnlock()

	sockets := make([]AppSocket, 0)

	for _, id := range uids {
		if s, ok := hdl.uidSckMap[id]; ok {
			sockets = append(sockets, s...)
		}
	}

	return sockets
}

func NewSckHdl() *socketHandler {
	return &socketHandler{
		uidSckMap: make(map[uint32][]AppSocket),
		locker:    new(sync.RWMutex),
	}
}

func (hdl *socketHandler) AddObservers(server *socketio.Server, sc goservice.ServiceContext,
	l logger.Logger, al HandlerUserJoin) func(socketio.Socket) {
	hdl.server = server

	return func(so socketio.Socket) {
		defer func() {
			if err := recover(); err != nil {
				l.Errorln(err)
			}
		}()

		as := newAppSocket(server, sc, l, so)
		as.locker = new(sync.RWMutex)

		l.Debugf("socket connected: %s", as)

		// Observers
		_ = as.On(SckDisconnected, func() {
			hdl.onDisconnected(as)
			al.OnUserDisConnect(sc, as)
		})
		_ = as.On(SckAuthenticate, func(data interface{}) {
			user, err := al.OnAuthentication(sc, data)
			if err != nil {
				al.OnAuthFail(sc, data, as)
			} else {
				as.SetCurrentUser(user)
				hdl.UserConnect(user.UserID(), as)
				al.OnAuthSuccessfully(sc, data, as)
			}
		})
	}
}
