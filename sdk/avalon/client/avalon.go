package client

import (
	"context"
	"errors"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
)

type IPDiscover interface {
	avalon.Bean
	GetHostports() []string
}

type Factory interface {
	avalon.Bean
	NewClient(hostport string) (interface{}, error)
}

type Caller interface {
	Call(ctx context.Context, invoke *avalon.Invoke) error
}

// implements
type AvalonClient struct {
	discover     IPDiscover
	loadBalancer LoadBalancer
	wrappers     avalon.WrapperComposite
	factory      Factory

	Retry int
}

func (f *AvalonClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	return inline.Retry(func() error {
		hostports := f.discover.GetHostports()
		target := f.loadBalancer.GetIP(hostports)
		if target == "" {
			return errors.New("no instance found")
		}
		o, err := f.factory.NewClient(target)
		if err != nil {
			return err
		}
		caller := o.(Caller)

		var call avalon.Call = caller.Call
		for _, wrapper := range f.wrappers {
			call = wrapper.Middleware(call)
		}

		return call(ctx, &avalon.Invoke{
			MethodName: method,
			Request:    args,
			Response:   result,
		})
	}, f.Retry, 0)
}

func (f *AvalonClient) Initial() error {
	if err := avalon.NewBean(f.discover).Initial(); err != nil {
		return err
	}
	if err := avalon.NewBean(f.loadBalancer).Initial(); err != nil {
		return err
	}
	if err := f.wrappers.Initial(); err != nil {
		return err
	}
	if err := avalon.NewBean(f.factory).Initial(); err != nil {
		return err
	}
	return nil
}

func (f *AvalonClient) Destroy() error {
	if err := f.discover.Destroy(); err != nil {
		return err
	}
	if err := f.loadBalancer.Destroy(); err != nil {
		return err
	}
	if err := f.wrappers.Destroy(); err != nil {
		return err
	}
	if err := f.factory.Destroy(); err != nil {
		return err
	}
	return nil
}

func (f *AvalonClient) SetDiscover(discover IPDiscover) *AvalonClient {
	f.discover = discover
	return f
}

func (f *AvalonClient) SetLoadBalancer(balancer LoadBalancer) *AvalonClient {
	f.loadBalancer = balancer
	return f
}

func (f *AvalonClient) AddWrapper(wrapper avalon.Wrapper) *AvalonClient {
	f.wrappers = append(f.wrappers, wrapper)
	return f
}

func (f *AvalonClient) SetCoreClient(factory Factory) *AvalonClient {
	f.factory = factory
	return f
}

func DefaultClient(psm string) *AvalonClient {
	return (&AvalonClient{}).SetDiscover(&ZkIPDiscover{PSM: psm}).
		SetLoadBalancer(&RandomBalancer{}).
		SetCoreClient(&ThriftClientFactory{})
}

func DefaultClientTimeout(psm string, timeout string) *AvalonClient {
	return (&AvalonClient{}).SetDiscover(&ZkIPDiscover{PSM: psm}).
		SetLoadBalancer(&RandomBalancer{}).
		SetCoreClient(&ThriftClientFactory{Timeout: timeout})
}
