package avalon

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
)

type iServer struct {
	opts      []Option
	config    *Config
	processor thrift.TProcessor

	tServer thrift.TServer
}

func NewServer(builder thrift.TProcessor, opts ...Option) *iServer {
	options := append([]Option{}, opts...)
	options = append(options, defaultOptions...)
	iServer := &iServer{
		opts:      options,
		config:    &Config{},
		processor: builder,
	}
	for _, option := range options {
		option(iServer.config)
	}

	return iServer
}

func (server *iServer) Start() error {

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
	return server.tServer.Stop()
}
