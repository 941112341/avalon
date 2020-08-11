package server

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
)

func Log(call Call) Call {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := call(ctx, request)
		inline.WithFields("request", request, "response", resp, "err", err).Infoln("finish")
		return resp, err
	}
}
