package client

import (
	"context"
	"errors"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

type Client struct {
	LoadBalancer LoadBalancer
}

func (c *Client) Call(ctx context.Context, method string, args, result thrift.TStruct) error {

	event := Event{Ctx: ctx, Method: method, Args: args, Result: result}
	return c.LoadBalancer.Consume(event)
}

type ClientBuilder struct {
	psm           string
	clientTimeout time.Duration
	config        zookeeper.ZkConfig
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		psm:           "",
		clientTimeout: time.Second,
		config: zookeeper.ZkConfig{
			HostPorts:      []string{"localhost:2181"},
			SessionTimeout: 10 * time.Second,
		},
	}
}

func (b *ClientBuilder) ZkConfig(config zookeeper.ZkConfig) *ClientBuilder {
	b.config = config
	return b
}

func (b *ClientBuilder) PSM(psm string) *ClientBuilder {
	b.psm = psm
	return b
}

func (b *ClientBuilder) Timeout(timeout time.Duration) *ClientBuilder {
	b.clientTimeout = timeout
	return b
}

func (b *ClientBuilder) valid() error {
	if b.psm == "" {
		return errors.New("psm is nil")
	}
	return nil
}

func (b *ClientBuilder) Build() (*Client, error) {
	config := b.config
	builder := NewBalancerBuilder()
	builder.Timeout(config.SessionTimeout).Session(config.HostPorts)
	balancer, err := builder.PSM(b.psm).Build()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "loadbalancer get fail")
	}
	loadBalancer := NewLoadBalancer(balancer, b.clientTimeout)
	return &Client{LoadBalancer: loadBalancer}, nil
}
