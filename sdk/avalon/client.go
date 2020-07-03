package avalon

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type Call func(ctx context.Context, method string, args, result interface{}) error

type Middleware func(config *Config, call Call) Call

type iClient struct {
	opts   []Option
	config *Config

	middleware []Middleware
	call       Call
}

func NewClient(opts ...Option) thrift.TClient {
	client := &iClient{
		opts:   opts,
		config: &Config{MethodConfig: map[string]*Config{}},
		middleware: []Middleware{
			ConfigMiddleware, DownstreamMiddleware, RetryMiddleware, MetricsMiddleware, ThriftMiddleware,
		},
	}
	for _, opt := range opts {
		opt(client.config)
	}
	for _, opt := range defaultOptions {
		opt(client.config)
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
