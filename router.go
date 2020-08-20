package sckio

import (
	"log"
	"strings"

	goservice "github.com/baozhenglab/go-sdk"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func (s *sckServer) Router(sc goservice.ServiceContext, al HandlerUserJoin) func(*gin.Engine) {
	return func(engine *gin.Engine) {
		transportsName := strings.Split(strings.TrimSpace(s.TransportNames), ",")
		server, err := socketio.NewServer(transportsName)
		if err != nil {
			log.Fatal(err)
		}

		op := NewSckHdl()

		server.SetMaxConnection(s.MaxConnection)
		s.io = server

		_ = s.io.On("connection", op.AddObservers(server, sc, s.logger, al))

		engine.GET("/socket.io/", gin.WrapH(server))
		engine.POST("/socket.io/", gin.WrapH(server))
	}
}
