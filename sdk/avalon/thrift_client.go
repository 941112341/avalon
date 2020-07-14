package avalon

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/collect"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var poolMap sync.Map

type consumer struct {
	socket          *thrift.TSocket
	transport       thrift.TTransport
	protocalFactory thrift.TProtocolFactory
}

func (c *consumer) Close() error {
	return c.socket.Close()
}

func (c *consumer) Do(ctx context.Context, args ...interface{}) error {
	client := thrift.NewTStandardClient(c.protocalFactory.GetProtocol(c.transport), c.protocalFactory.GetProtocol(c.transport))
	method := args[1].(string)
	tArgs, tResult := args[2].(thrift.TStruct), args[3].(thrift.TStruct)

	return client.Call(ctx, method, tArgs, tResult)
}

type factory struct {
	hostPort         string
	timeout          time.Duration
	transportFactory thrift.TTransportFactory
	protocalFactory  thrift.TProtocolFactory
}

func (c *factory) Create() (collect.Consumer, error) {
	tSocket, err := thrift.NewTSocket(c.hostPort)
	if err != nil {
		return nil, errors.Wrap(err, c.hostPort)
	}
	if c.timeout != time.Duration(0) {
		err = tSocket.SetTimeout(c.timeout)
		if err != nil {
			inline.Warningln("set timeout err", inline.NewPairs("err", err, "timeout", c.timeout)...)
		}
	}
	transport, err := c.transportFactory.GetTransport(tSocket)
	if err != nil {
		return nil, errors.Wrap(err, c.hostPort)
	}
	if err = transport.Open(); err != nil {
		return nil, errors.Wrap(err, c.hostPort)
	}

	return &consumer{
		transport:       transport,
		socket:          tSocket,
		protocalFactory: c.protocalFactory,
	}, nil
}

func newFactory(hostPort string, timeout time.Duration, transportFactory thrift.TTransportFactory,
	protocalFactory thrift.TProtocolFactory) collect.ConsumerFactory {
	return &factory{
		hostPort:         hostPort,
		timeout:          timeout,
		transportFactory: transportFactory,
		protocalFactory:  protocalFactory,
	}
}

func NewFactory(hostPort string, timeout time.Duration) collect.ConsumerFactory {
	return newFactory(hostPort, timeout, thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()),
		thrift.NewTBinaryProtocolFactoryDefault())
}

func ThriftMiddleware(cfg Config, _ Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		hostPort := GetHostPort(ctx)
		if hostPort == "" {
			return errors.New("host port is empty")
		}

		any, ok := poolMap.Load(hostPort)
		if !ok {
			any = collect.NewPool(time.Hour, 10, 20, NewFactory(hostPort, cfg.Client.Timeout))
			poolMap.Store(hostPort, any)
		}
		pool := any.(collect.Pool)
		t := time.Now()
		consumer, err := pool.GetConsumerBlock(cfg.Client.Timeout)
		inline.Infoln("get consumer spend", inline.NewPairs("time", time.Since(t).String(), "hostPort", hostPort)...)
		if err != nil {
			return errors.Wrap(err, "get consumer")
		}
		return consumer.Do(ctx, method, args, result)
	}
}
