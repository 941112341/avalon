package avalon

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"time"
)

type initial func(cfg Config) error

type Bootstrap interface {
	Start() error
	Stop() error
	Register(handler interface{}) error
	Config() Config
}

type IServer struct {
	server     thrift.TServer
	builder    ConfigBuilder
	initials   []initial
	Middleware []Middleware
}

func (s *IServer) Config() Config {
	return s.builder.Config()
}

func (s *IServer) Start() error {
	cfg := s.builder.Config()
	inline.Infoln("start server", inline.NewPair("port", cfg.Server.Port))
	return s.server.Serve()
}

func (s *IServer) Stop() error {
	return s.server.Stop()
}

// now support only tProcessor
func (s *IServer) Register(handler interface{}) error {
	processor, ok := handler.(thrift.TProcessor)
	if !ok {
		return fmt.Errorf("handler is not tProcessor %b", handler)
	}
	cfg := s.builder.Config()
	for idx, initial := range s.initials {
		if err := initial(cfg); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("index[%d]", idx))
		}
	}

	hostPort := fmt.Sprintf(":%d", cfg.Server.Port)
	timeout := time.Second
	if cfg.Server.Timeout != 0 {
		timeout = cfg.Server.Timeout * time.Second
	}
	serverTransport, err := thrift.NewTServerSocketTimeout(hostPort, timeout)
	if err != nil {
		return errors.WithMessage(err, "new socket")
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	tServer := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	s.server = tServer
	return nil
}

func NewServerWithConfig(cfg ConfigBuilder, middleware ...Middleware) *IServer {
	return &IServer{
		builder: cfg,
		initials: []initial{
			RegisterService,
		},
		Middleware: append([]Middleware{

			metaMiddlewareServer,
		}, middleware...),
	}
}

func NewServer(psm string, middleware ...Middleware) *IServer {

	return NewServerWithConfig(NewConfigBuilder(psm), middleware...)
}
