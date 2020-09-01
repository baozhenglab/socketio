package sckio

import (
	"net/http"

	goservice "github.com/baozhenglab/go-sdk"
	"github.com/baozhenglab/go-sdk/logger"
	"github.com/baozhenglab/sdkcm"
	"github.com/gin-gonic/gin"
)

type Socket interface {
	Id() string
	Rooms() []string
	Request() *http.Request
	On(event string, f interface{}) error
	Emit(event string, args ...interface{}) error
	Join(room string) error
	Leave(room string) error
	Disconnect()
	BroadcastTo(room, event string, args ...interface{}) error
}

type SocketIOService interface {
	Router(goservice.ServiceContext, HandlerUserJoin) func(*gin.Engine)
}

type AppSocket interface {
	ServiceContext() goservice.ServiceContext
	Logger() logger.Logger
	CurrentUser() sdkcm.Requester
	SetCurrentUser(sdkcm.Requester)
	BroadcastToRoom(room, event string, args ...interface{})
	String() string
	Socket
}

type HandlerUserJoin interface {
	OnAuthentication(sc goservice.ServiceContext, as AppSocket, data interface{}) (sdkcm.Requester, error)
	OnAuthFail(sc goservice.ServiceContext, data interface{}, as AppSocket)
	OnAuthSuccessfully(sc goservice.ServiceContext, data interface{}, as AppSocket)
	OnUserDisConnect(sc goservice.ServiceContext, as AppSocket)
}
