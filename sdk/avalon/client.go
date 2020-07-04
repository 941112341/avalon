package avalon

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type Call func(ctx context.Context, method string, args, result interface{}) error

type Middleware func(config *ClientConfig, call Call) Call

type iClient struct {
	config *ClientConfig

	middleware []Middleware
	call       Call
}

func NewClient() thrift.TClient {
	return NewClientWithConfig(defaultClientConfig)
}

func NewClientWithConfig(config *ClientConfig) thrift.TClient {
	client := &iClient{
		config: config,
		middleware: []Middleware{
			RetryMiddleware, DiscoverMiddleware, MetricsMiddleware, ThriftMiddleware,
		},
	}
	var call Call
	for i := len(client.middleware) - 1; i >= 0; i-- {
		call = client.middleware[i](client.config, call)
	}
	client.call = call

	return client
}

func (client *iClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	return client.call(ctx, method, args, result)
}
