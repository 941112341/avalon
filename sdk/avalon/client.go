package avalon

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type Endpoint func(ctx context.Context, method string, args, result interface{}) error

type Middleware func(config Config, call Endpoint) Endpoint

type ConfigBuilder interface {
	Config() Config
}

type IClient struct {
	Middleware []Middleware
	builder    ConfigBuilder
}

func (c *IClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	var call Endpoint
	for i := len(c.Middleware) - 1; i >= 0; i-- {
		call = c.Middleware[i](c.builder.Config(), call)
	}

	return call(ctx, method, args, result)
}

func NewClientWithConfig(builder ConfigBuilder, middleware ...Middleware) *IClient {
	meddlers := []Middleware{
		scopeMiddleware,
		RetryMiddleware,
		DiscoverMiddleware,
		FixAddressMiddleware,
		MetricsMiddleware,
		metaMiddlewareClient,
		ThriftMiddleware,
	}
	return &IClient{
		builder:    builder,
		Middleware: append(meddlers, middleware...),
	}
}

func NewClient(psm string, others ...Config) *IClient {

	return NewClientWithConfig(NewConfigBuilder(psm, others...))
}
