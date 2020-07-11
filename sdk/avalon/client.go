package avalon

import (
	"context"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

type CallConfig struct {
	Retry int
	Wait  time.Duration // millSec
	Psm   string

	// common
	HostPort string
	zookeeper.ZkConfig
	Timeout time.Duration // sec
}

func NewCallConfig(c Config) CallConfig {
	return CallConfig{
		HostPort: c.Client.HostPort,
		Retry:    c.Client.Retry,
		Wait:     c.Client.Wait,
		Psm:      c.Psm,
		Timeout:  c.Client.Timeout,
		ZkConfig: c.ZkConfig,
	}
}

type Endpoint func(ctx context.Context, method string, args, result interface{}) error

type Middleware func(config CallConfig, call Endpoint) Endpoint

type IClient struct {
	Middleware []Middleware
	Cfg        Config
}

func (c *IClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	cfg := NewCallConfig(c.Cfg)
	var call Endpoint
	for i := len(c.Middleware) - 1; i >= 0; i-- {
		call = c.Middleware[i](cfg, call)
	}

	return call(ctx, method, args, result)
}

func NewClientWithConfig(cfg Config, middleware ...Middleware) *IClient {
	meddlers := []Middleware{
		RetryMiddleware, DiscoverMiddleware, FixAddressMiddleware, MetricsMiddleware, ThriftMiddleware,
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
