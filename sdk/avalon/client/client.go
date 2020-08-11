package client

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

type Client struct {
	LoadBalancer LoadBalancer
}

func (c *Client) Call(ctx context.Context, method string, args, result thrift.TStruct) error {

	SetBaseArgs(ctx, args)
	event := &Event{Ctx: ctx, Method: method, Args: args, Result: result}
	return c.LoadBalancer.Consume(event)
}

func NewClient(loadBalancer LoadBalancer) *Client {
	return &Client{LoadBalancer: loadBalancer}
}

func NewClientOptions(psm string, options ...Option) (*Client, error) {
	builder := NewBalancerBuilder()
	for _, option := range options {
		option.BalanceBuilder(builder)
	}
	balancer, err := builder.PSM(psm).Build()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "build fail")
	}

	balancerBuilder := NewLoadBalancerBuilder()
	for _, option := range options {
		option.LoadBalanceBuilder(balancerBuilder)
	}
	loadBalancer := balancerBuilder.Balancer(balancer).Build()
	return &Client{LoadBalancer: loadBalancer}, nil
}

type Option interface {
	BalanceBuilder(builder *BalancerBuilder)
	LoadBalanceBuilder(builder *loadBalancerBuilder)
}

type TimeoutOption struct {
	timeout time.Duration
}

func (t TimeoutOption) BalanceBuilder(builder *BalancerBuilder) {
	builder.Timeout(t.timeout)
}

func (t TimeoutOption) LoadBalanceBuilder(builder *loadBalancerBuilder) {}

func WithTimeoutOption(timeout time.Duration) Option {
	return &TimeoutOption{timeout: timeout}
}
