package avalon

import (
	"github.com/941112341/avalon/sdk/log"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Server interface {
	Start() error
	Stop() error
}

type thriftServer struct {
	processor thrift.TProcessor

	tServer thrift.TServer
	config  *ServerConfig
}

func NewThriftServer(builder thrift.TProcessor, config *ServerConfig) *thriftServer {
	return &thriftServer{
		config:    config,
		processor: builder,
	}
}

func (server *thriftServer) Start() error {

	serverTransport, err := thrift.NewTServerSocketTimeout(server.config.HostPort, server.config.Timeout*time.Second)
	if err != nil {
		return errors.WithMessage(err, "new socket")
	}
	log.New().WithFields(logrus.Fields{
		"hostPort": server.config.HostPort,
	}).Infof("server start")
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

func (server *thriftServer) Stop() error {
	return server.tServer.Stop()
}

type ServerWrapper func(cfg *ServerConfig, server Server) Server

type Bootstrap struct {
	server   Server
	cfg      *ServerConfig
	wrappers []ServerWrapper
	preHooks []func(cfg *ServerConfig) error
}

func (b *Bootstrap) Start() error {
	for _, hook := range b.preHooks {
		err := hook(b.cfg)
		if err != nil {
			return errors.Cause(err)
		}
	}
	s := b.server
	for _, wrapper := range b.wrappers {
		s = wrapper(b.cfg, s)
	}
	return s.Start()
}

func (b *Bootstrap) Stop() error {
	s := b.server
	for _, wrapper := range b.wrappers {
		s = wrapper(b.cfg, s)
	}
	return s.Stop()
}

func Wrap(cfg *ServerConfig, server Server) *Bootstrap {
	return &Bootstrap{
		server: server,
		cfg:    cfg,
		preHooks: []func(cfg *ServerConfig) error{
			func(cfg *ServerConfig) error {
				return startDiscover(cfg.ZkConfig)
			},
		},
		wrappers: []ServerWrapper{
			ServiceRegisterWrapper,
		},
	}
}
