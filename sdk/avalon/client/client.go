package client

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

type DirectThriftClient struct {
	Hostport string
	Timeout  time.Duration
}

func (d *DirectThriftClient) Call(ctx context.Context, invoke *avalon.Invoke) error {
	tSocket, err := thrift.NewTSocketTimeout(d.Hostport, d.Timeout)
	if err != nil {
		return inline.PrependErrorFmt(err, "open socket %+v", *d)
	}
	defer tSocket.Close()
	transport, err := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()).GetTransport(tSocket)
	if err != nil {
		return inline.PrependErrorFmt(err, "get transport")
	}
	if err = transport.Open(); err != nil {
		return inline.PrependErrorFmt(err, "transport open")
	}
	input := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(transport)
	output := input
	client := thrift.NewTStandardClient(input, output)

	return client.Call(ctx, invoke.MethodName, invoke.Request.(thrift.TStruct), invoke.Response.(thrift.TStruct))
}

// may be pool
type ThriftClientFactory struct {
	Timeout string `default:"2s"`
	timeout time.Duration
}

func (t *ThriftClientFactory) Initial() error {
	t.timeout = inline.Parse(t.Timeout)
	return nil
}

func (t *ThriftClientFactory) Destroy() error {
	return nil
}

func (t *ThriftClientFactory) NewClient(hostport string) (interface{}, error) {
	return &DirectThriftClient{
		Hostport: hostport,
		Timeout:  t.timeout,
	}, nil
}
