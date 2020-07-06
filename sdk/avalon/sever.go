package avalon

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"sort"
)

type Lifecycle interface {
	PreStart(config *ServerConfig) error
	PreStop(config *ServerConfig) error
	Order() int
}

type Lifecycles []Lifecycle

func (ls Lifecycles) Len() int {
	return len(ls)
}

func (ls Lifecycles) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls Lifecycles) Less(i, j int) bool {
	return ls[i].Order() < ls[j].Order()
}

type iServer struct {
	config    *ServerConfig
	processor thrift.TProcessor

	tServer thrift.TServer
	hooks   []Lifecycle
}

func NewServer(builder thrift.TProcessor) *iServer {
	return NewServerWithConfig(builder, defaultServerConfig)
}

func NewServerWithConfig(builder thrift.TProcessor, config *ServerConfig, hooks ...Lifecycle) *iServer {
	sort.Sort(Lifecycles(hooks))
	return &iServer{
		config:    config,
		processor: builder,
		hooks:     hooks,
	}
}

func (server *iServer) Start() error {
	for _, hook := range server.hooks {
		err := hook.PreStart(server.config)
		if err != nil {
			return errors.WithMessage(err, inline.VString(hook))
		}
	}

	serverTransport, err := thrift.NewTServerSocketTimeout(server.config.HostPort, server.config.Timeout)
	if err != nil {
		return errors.WithMessage(err, "new socket")
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	tServer := thrift.NewTSimpleServer4(server.processor, serverTransport, transportFactory, protocolFactory)
	err = tServer.Serve()
	if err != nil {
		return errors.WithMessage(err, "server")
	}
	server.tServer = tServer
	return nil
}

func (server *iServer) Stop() error {
	err := server.tServer.Stop()

	for _, hook := range server.hooks {
		err := hook.PreStop(server.config)
		if err != nil {
			return errors.WithMessage(err, inline.VString(hook))
		}
	}
	return err
}
