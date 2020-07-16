package avalon

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func MetricsMiddleware(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		t := time.Now()
		err := call(ctx, method, args, result)
		log.New().
			WithField("duration", time.Since(t).String()).
			WithField("err", err).
			WithField("psm", cfg.Psm).
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
		ctx = WithSession(ctx, &Session{HostPort: cfg.Client.HostPort, Attachments: map[string]interface{}{}})

		return call(ctx, method, args, result)
	}
}
