package avalon

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func ThriftMiddleware(cfg Config, _ Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		tArgs, ok := args.(thrift.TStruct)
		if !ok {
			return errors.New(fmt.Sprintf("args %+v is not TStruct", args))
		}
		tResult, ok := result.(thrift.TStruct)
		if !ok {
			return errors.New(fmt.Sprintf("result %+v is not TStruct", args))
		}

		hostPort := GetHostPort(ctx)
		if hostPort == "" {
			return errors.New("host port is empty")
		}

		tSocket, err := thrift.NewTSocket(hostPort)
		if err != nil {
			return errors.Wrap(err, hostPort)
		}
		if cfg.Client.Timeout != time.Duration(0) {
			err = tSocket.SetTimeout(cfg.Client.Timeout * time.Second)
			if err != nil {
				log.New().WithField("err", err.Error()).WithField("warn", cfg.Client.Timeout).Warningln()
			}
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, err := transportFactory.GetTransport(tSocket)
		if err != nil {
			return errors.WithMessage(err, inline.ToJsonString(cfg))
		}
		err = transport.Open()
		if err != nil {
			return errors.WithMessage(err, inline.ToJsonString(cfg))
		}
		defer transport.Close()
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := thrift.NewTStandardClient(protocolFactory.GetProtocol(transport), protocolFactory.GetProtocol(transport))
		return client.Call(ctx, method, tArgs, tResult)
	}
}

func MetricsMiddleware(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		t := time.Now()
		err := call(ctx, method, args, result)
		log.New().
			WithField("duration", time.Since(t).String()).
			WithField("err", err).
			WithField("hostPort", cfg.Client.HostPort).
			Info("call")
		return err
	}
}

func RetryMiddleware(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		return inline.Retry(func() error {
			return call(ctx, method, args, result)
		}, cfg.Client.Retry, cfg.Client.Wait*time.Millisecond)
	}
}

func FixAddressMiddleware(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		ip, err := inline.InetAddress()
		if err != nil {
			return errors.WithMessage(err, "get ip err")
		}
		hostPort := GetHostPort(ctx)
		if strings.HasPrefix(hostPort, ip) {
			SetHostPort(ctx, strings.Replace(hostPort, ip, "localhost", 1))
		}
		return call(ctx, method, args, result)
	}
}

func CreateSessionMiddleware(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		ctx = WithSession(ctx, &Session{HostPort: cfg.Client.HostPort})

		return call(ctx, method, args, result)
	}
}
