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

func ThriftMiddleware(config *Config, _ Call) Call {
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
				log.NewLoggerWithRotate().WithField("err", err.Error()).WithField("warn", config.Timeout).Warningln()
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

func MetricsMiddleware(config *Config, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		t := time.Now()
		err := call(ctx, method, args, result)
		log.NewLoggerWithRotate().WithField("duration", time.Since(t).String()).WithField("err", err).Info("call")
		return err
	}
}

func RetryMiddleware(config *Config, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		return inline.Retry(func() error {
			return call(ctx, method, args, result)
		}, config.Retry, 100*time.Millisecond)
	}
}

func DownstreamMiddleware(config *Config, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		err := call(ctx, method, args, result)
		if err != nil && config.Downstream != nil {
			return config.Downstream(ctx, method, args, result, err)
		}
		return errors.Cause(err)
	}
}

func ConfigMiddleware(config *Config, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		*config = *config.Get(method)
		return call(ctx, method, args, result)
	}
}
