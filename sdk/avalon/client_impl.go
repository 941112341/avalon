package avalon

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
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
			err = tSocket.SetTimeout(config.Timeout * time.Second)
			if err != nil {
				log.New().WithField("err", err.Error()).WithField("warn", config.Timeout).Warningln()
			}
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, err := transportFactory.GetTransport(tSocket)
		if err != nil {
			return errors.WithMessage(err, inline.ToJsonString(config))
		}
		err = transport.Open()
		if err != nil {
			return errors.WithMessage(err, inline.ToJsonString(config))
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
		log.New().
			WithField("duration", time.Since(t).String()).
			WithField("err", err).
			WithField("hostPort", config.HostPort).
			Info("call")
		return err
	}
}

func RetryMiddleware(config *ClientConfig, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		return inline.Retry(func() error {
			return call(ctx, method, args, result)
		}, config.Retry, config.Wait*time.Millisecond)
	}
}

func DiscoverMiddleware(config *ClientConfig, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		hostPortList := make([]string, 0)
		prefix := config.ZkConfig.Path + "/" + config.ServiceName
		DiscoverMap.Range(func(key, value interface{}) {
			ks, ok := key.(string)
			if ok {
				if strings.HasPrefix(ks, prefix) && ks != prefix {
					hostPort := strings.Replace(ks, prefix, "", 1)
					if strings.HasPrefix(hostPort, "/") {
						hostPort = hostPort[1:]
					}
					hostPortList = append(hostPortList, hostPort)
				}
			}
		})
		if len(hostPortList) == 0 {
			return errors.New("serviceName " + config.ServiceName + " service not found")
		}
		// 随机负载
		idx := rand.Intn(len(hostPortList))
		config.HostPort = hostPortList[idx]
		return call(ctx, method, args, result)
	}
}

func DebugMiddleware(config *ClientConfig, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		if config.ClientIp == "" {
			ip, err := inline.InetAddress()
			if err != nil {
				return errors.WithMessage(err, "get ip err")
			}
			config.ClientIp = ip
		}

		// 如果ip相同 使用本地环回地址 避免防火墙
		if strings.HasPrefix(config.HostPort, config.ClientIp) {
			config.HostPort = strings.Replace(config.HostPort, config.ClientIp, "localhost", 1)
		}
		return call(ctx, method, args, result)
	}
}
