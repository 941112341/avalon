package avalon

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"time"
)

type initial func(cfg CallConfig) error

type Bootstrap interface {
	Start() error
	Stop() error
	Register(handler interface{}) error
}

type IServer struct {
	server     thrift.TServer
	Cfg        Config
	initials   []initial
	Middleware []Middleware
}

func (s *IServer) Start() error {
	inline.Infoln("start server", inline.NewPair("port", s.Cfg.Server.Port))
	return s.server.Serve()
}

func (s *IServer) Stop() error {
	return s.server.Stop()
}

// now support only tProcessor
func (s *IServer) Register(handler interface{}) error {
	cfg := NewCallConfig(s.Cfg)
	for idx, initial := range s.initials {
		if err := initial(cfg); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("index[%d]", idx))
		}
	}

	processor, ok := handler.(thrift.TProcessor)
	if !ok {
		return fmt.Errorf("handler is not tProcessor %b", handler)
	}
	serverTransport, err := thrift.NewTServerSocketTimeout(cfg.HostPort, cfg.Timeout*time.Second)
	if err != nil {
		return errors.WithMessage(err, "new socket")
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	tServer := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	s.server = tServer
	return nil
}

func NewServerWithConfig(cfg Config, middleware ...Middleware) *IServer {
	return &IServer{
		Cfg: cfg,
		initials: []initial{
			RegisterService,
		},
		Middleware: middleware,
	}
}

func NewServer(middleware ...Middleware) (*IServer, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return NewServerWithConfig(cfg, middleware...), nil
}
