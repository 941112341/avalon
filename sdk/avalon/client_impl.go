package avalon

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"time"
)

func ThriftMiddleware(config *ClientConfig, _ Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		tArgs, ok := args.(thrift.TStruct)
		if !ok {
			return errors.New(fmt.Sprintf("args %+v is not TStruct", args))
		}
		tResult, ok := result.(thrift.TStruct)
		if !ok {
			return errors.New(fmt.Sprintf("result %+v is not TStruct", args))
		}
		if config.HostPort == "" {
			return errors.New("host port is empty")
		}
		tSocket, err := thrift.NewTSocket(config.HostPort)
		if err != nil {
			return errors.Wrap(err, config.HostPort)
		}
		if config.Timeout != time.Duration(0) {
			err = tSocket.SetTimeout(config.Timeout)
			if err != nil {
				log.New().WithField("err", err.Error()).WithField("warn", config.Timeout).Warningln()
			}
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, err := transportFactory.GetTransport(tSocket)
		if err != nil {
			return errors.WithMessage(err, inline.JsonString(config))
		}
		err = transport.Open()
		if err != nil {
			return errors.WithMessage(err, inline.JsonString(config))
		}
		defer transport.Close()
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := thrift.NewTStandardClient(protocolFactory.GetProtocol(transport), protocolFactory.GetProtocol(transport))
		return client.Call(ctx, method, tArgs, tResult)
	}
}

func MetricsMiddleware(config *ClientConfig, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		t := time.Now()
		err := call(ctx, method, args, result)
		log.New().WithField("duration", time.Since(t).String()).WithField("err", err).Info("call")
		return err
	}
}

func RetryMiddleware(config *ClientConfig, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		return inline.Retry(func() error {
			return call(ctx, method, args, result)
		}, config.Retry, 100*time.Millisecond)
	}
}

func DiscoverMiddleware(config *ClientConfig, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		if config.HostPort == "" {
			config.HostPort = "localhost:8888" // todo zk
		}
		return call(ctx, method, args, result)
	}
}
