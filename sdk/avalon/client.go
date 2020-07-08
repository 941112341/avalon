package avalon

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
)

type Call func(ctx context.Context, method string, args, result interface{}) error

type Middleware func(config *ClientConfig, call Call) Call

type iClient struct {
	config *ClientConfig

	middleware []Middleware
}

func NewClientWithConfig(config *ClientConfig) thrift.TClient {
	return &iClient{
		config: config,
		middleware: []Middleware{
			RetryMiddleware, DiscoverMiddleware, DebugMiddleware, MetricsMiddleware, ThriftMiddleware,
		},
	}
}

func (client *iClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	cfg := &ClientConfig{}
	if err := inline.Copy(client.config, cfg); err != nil {
		return errors.WithMessage(err, "copy fail")
	}

	var call Call
	for i := len(client.middleware) - 1; i >= 0; i-- {
		call = client.middleware[i](cfg, call)
	}

	return call(ctx, method, args, result)
}
