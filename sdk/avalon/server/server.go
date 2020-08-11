package server

import (
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"time"
)

type Server struct {
	Timeout  time.Duration
	PSM      string
	Hostport both.Hostport
}

func (s *Server) Run(processor thrift.TProcessor) error {
	hostport := s.Hostport.Port()
	serverTransport, err := thrift.NewTServerSocketTimeout(hostport, s.Timeout)
	if err != nil {
		return errors.WithMessage(err, "new socket")
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	tServer := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	return tServer.Serve()
}

type serverBuilder struct {
	s *Server
}

func Builder() *serverBuilder {
	return &serverBuilder{s: &Server{Timeout: time.Second}}
}

func (b *serverBuilder) Timeout(t time.Duration) *serverBuilder {
	b.s.Timeout = t
	return b
}

func (b *serverBuilder) PSM(psm string) *serverBuilder {
	b.s.PSM = psm
	return b
}

func (b *serverBuilder) Hostport(hostport both.Hostport) *serverBuilder {
	b.s.Hostport = hostport
	return b
}

func (b *serverBuilder) Build() *Server {
	selfBase = &both.Base{
		Psm:   b.s.PSM,
		IP:    inline.GetIP(),
		Time:  0,
		Extra: map[string]string{},
		Base:  nil,
	}
	return b.s
}
