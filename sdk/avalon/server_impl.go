package avalon

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type serverImpl struct {
	start func(cfg *ServerConfig) error
	stop  func(cfg *ServerConfig) error
}

func (s *serverImpl) Start(cfg *ServerConfig) error {
	return s.start(cfg)
}

func (s *serverImpl) Stop(cfg *ServerConfig) error {
	return s.stop(cfg)
}

type thriftServer struct {
	processor thrift.TProcessor
	tServer   thrift.TServer
}

func NewThriftServer(builder thrift.TProcessor) *thriftServer {
	return &thriftServer{
		processor: builder,
	}
}

func (server *thriftServer) Start(cfg *ServerConfig) error {
	serverTransport, err := thrift.NewTServerSocketTimeout(cfg.HostPort, cfg.Timeout*time.Second)
	if err != nil {
		return errors.WithMessage(err, "new socket")
	}
	log.New().WithFields(logrus.Fields{
		"hostPort": cfg.HostPort,
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

func (server *thriftServer) Stop(config *ServerConfig) error {
	return server.tServer.Stop()
}

func ServiceRegisterWrapper(coreServer Server) Server {
	return &serverImpl{
		start: func(cfg *ServerConfig) error {
			zkCli, err := zookeeper.GetZkClientInstance(&cfg.ZkConfig)
			if err != nil {
				return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
			}
			ip, err := inline.InetAddress()
			if err != nil {
				return err
			}
			idx := strings.LastIndex(cfg.HostPort, ":")
			port := cfg.HostPort[idx:]
			hostPort := ip + port
			node := zookeeper.NewZkNodeBuilder(inline.JoinPath(cfg.Path, cfg.ServiceName, hostPort)).Build()
			err = node.Save(zkCli, zk.FlagEphemeral)
			if err != nil {
				return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
			}

			log.New().WithField("serviceName", cfg.ServiceName).Infoln("register success")

			return coreServer.Start(cfg)
		},
		stop: coreServer.Stop,
	}
}
