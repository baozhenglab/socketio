package sckio

import (
	"flag"
	"fmt"

	goservice "github.com/baozhenglab/go-sdk"
	"github.com/baozhenglab/go-sdk/logger"
	socketio "github.com/googollee/go-socket.io"
)

type Config struct {
	Name           string
	MaxConnection  int
	TransportNames string
}

type sckServer struct {
	Config
	io     *socketio.Server
	logger logger.Logger
}

func New(name string) goservice.PrefixConfigure {
	return &sckServer{
		Config: Config{Name: name},
	}
}

func (s *sckServer) GetPrefix() string {
	return s.Config.Name
}

func (s *sckServer) Get() interface{} {
	return s
}

func (s *sckServer) Name() string {
	return s.Config.Name
}

func (s *sckServer) InitFlags() {
	pre := s.GetPrefix()
	flag.IntVar(&s.MaxConnection, fmt.Sprintf("%s-max-connection", pre), 2000, "socket max connection")
	flag.StringVar(&s.TransportNames, fmt.Sprintf("%s-tranports-name", pre), "websocket", "List tranport name, example: websocket,polling")

	s.logger = logger.GetCurrent().GetLogger("io.socket")
}
