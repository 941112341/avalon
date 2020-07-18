package avalon

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type Endpoint func(ctx context.Context, method string, args, result interface{}) error

type Middleware func(config Config, call Endpoint) Endpoint

type IClient struct {
	Middleware []Middleware
	Cfg        Config
}

func (c *IClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	var call Endpoint
	for i := len(c.Middleware) - 1; i >= 0; i-- {
		call = c.Middleware[i](c.Cfg, call)
	}

	return call(ctx, method, args, result)
}

func NewClientWithConfig(cfg Config, middleware ...Middleware) *IClient {
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
		Cfg:        cfg,
		Middleware: append(meddlers, middleware...),
	}
}

func NewClient(middleware ...Middleware) (*IClient, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return NewClientWithConfig(cfg, middleware...), nil
}
