package server

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/inline"
	"time"
)

type LogWrapper struct {
	avalon.TodoBean
}

func (l *LogWrapper) Middleware(call avalon.Call) avalon.Call {
	return func(ctx context.Context, invoke *avalon.Invoke) error {
		start := time.Now()
		err := call(ctx, invoke)

		inline.WithFields("request", invoke.Request, "response", invoke.Response, "err", err,
			"spend", time.Since(start).String()).Infoln("finish %s", invoke.MethodName)
		return err
	}
}
