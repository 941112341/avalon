package avalon

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
)

type PreHook func(server *iServer) error

type PostHook func(server *iServer) error

type iServer struct {
	config    *ServerConfig
	processor thrift.TProcessor

	preHooks  []PreHook
	postHooks []PostHook
	tServer   thrift.TServer
}

func NewServer(builder thrift.TProcessor) *iServer {
	return NewServerWithConfig(builder, defaultServerConfig)
}

func NewServerWithConfig(builder thrift.TProcessor, config *ServerConfig) *iServer {
	return &iServer{
		config:    config,
		processor: builder,
	}
}

func (server *iServer) Start() error {
	for _, hook := range server.preHooks {
		err := hook(server)
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

	for _, hook := range server.postHooks {
		err := hook(server)
		if err != nil {
			return errors.WithMessage(err, inline.VString(hook))
		}
	}
	return err
}
