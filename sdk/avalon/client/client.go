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

	event := &Event{Ctx: ctx, Method: method, Args: args, Result: result}
	return c.LoadBalancer.Consume(event)
}

func NewClient(loadBalancer LoadBalancer) *Client {
	return &Client{LoadBalancer: loadBalancer}
}

func NewClientOptions(psm string) (*Client, error) {
	balancer, err := NewBalancerBuilder().PSM(psm).Build()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "build fail")
	}
	loadBalancer := NewLoadBalancerBuilder().Balancer(balancer).Build()
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

func (t TimeoutOption) LoadBalanceBuilder(builder *loadBalancerBuilder) {
	// todo

}

func WithTimeoutOption(timeout time.Duration) Option {
	return &TimeoutOption{timeout: timeout}
}
