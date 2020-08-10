package client

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/collect"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"math/rand"
	"time"
)

type Event struct {
	Method       string
	Args, Result thrift.TStruct
	Ctx          context.Context
}

type TClient struct {
	socket        *thrift.TSocket
	input, output thrift.TProtocol
	UUID          string
}

func (t *TClient) Consume(e collect.Event) error {
	event, ok := e.(*Event)
	if !ok {
		return fmt.Errorf("event %+v", e)
	}

	time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	inline.WithFields("uuid", t.UUID, "event", e).Infoln("handler")
	client := thrift.NewTStandardClient(t.input, t.output)
	return client.Call(event.Ctx, event.Method, event.Args, event.Result)
}

func (t *TClient) Shutdown() error {
	return t.socket.Close()
}

type TClientFactory struct {
	transportFactory thrift.TTransportFactory
	protocalFactory  thrift.TProtocolFactory
	hostport         string
	timeout          time.Duration
}

func NewClientFactory(hostport string, timeout time.Duration, transportFactory thrift.TTransportFactory,
	protocalFactory thrift.TProtocolFactory) *TClientFactory {
	return &TClientFactory{
		transportFactory: transportFactory,
		protocalFactory:  protocalFactory,
		hostport:         hostport,
		timeout:          timeout,
	}
}

func NewDefaultFactory(hostPort string, timeout time.Duration) collect.ConsumerFactory {
	return NewClientFactory(hostPort, timeout, thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()),
		thrift.NewTBinaryProtocolFactoryDefault())
}

func (T *TClientFactory) CreateConsumer() (collect.Consumer, error) {
	tSocket, err := thrift.NewTSocketTimeout(T.hostport, T.timeout)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "open socket %+v", *T)
	}
	transport, err := T.transportFactory.GetTransport(tSocket)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "get transport")
	}
	if err = transport.Open(); err != nil {
		return nil, inline.PrependErrorFmt(err, "transport open")
	}
	return &TClient{
		socket: tSocket,
		input:  T.protocalFactory.GetProtocol(transport),
		output: T.protocalFactory.GetProtocol(transport),
		UUID:   inline.RandString(32),
	}, nil
}

func NewPool(timeout time.Duration, min, max int, factory collect.ConsumerFactory) (collect.Consumer, error) {
	return collect.NewManagerBuilder().
		Timeout(timeout).
		Max(int64(max)).
		Min(int64(min)).
		Factory(factory).
		Build()
}
