package avalon

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"time"
)

const (
	MetaKey = "__META_KEY__"
)

type Base struct {
	IP        string `json:"ip"`
	PSM       string `json:"psm"`
	Time      int64
	Extra     map[string]string
	RequestID string

	Base *Base
}

func GetBase(ctx context.Context) *Base {
	base, _ := ctx.Value(MetaKey).(*Base)
	return base
}

func SetBase(ctx context.Context, base *Base) context.Context {
	return context.WithValue(ctx, MetaKey, base)
}

func RequestID(ctx context.Context) string {
	base := GetBase(ctx)
	if base == nil {
		return "nil"
	}
	return base.RequestID
}

func metaMiddlewareClient(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		scope := GetScope(ctx).find(FromCrossRPC)
		preMeta := GetBase(ctx)

		extra := make(map[string]string)
		scope.cache.Range(func(key, value interface{}) bool {
			k, v := key.(string), value.(string)
			if k == HostPortKey {
				return true
			}
			extra[k] = v
			return true
		})
		meta := &Base{
			IP:    inline.GetIP(),
			PSM:   cfg.PSM,
			Time:  time.Now().Unix(),
			Extra: extra,
			Base:  preMeta,
		}
		if preMeta == nil || preMeta.RequestID == "" {
			meta.RequestID = inline.RandString(32)
		} else {
			meta.RequestID = preMeta.RequestID
		}
		r, err := inline.GetField(args, "Request")
		if err != nil {
			return errors.Wrap(err, "Get requset fail")
		}
		err = inline.SetFieldJSON(r, "Base", meta)
		if err != nil {
			return errors.Wrap(err, "set context error")
		}

		return call(ctx, method, args, result)
	}
}

func metaMiddlewareServer(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		i, err := inline.GetField(args, "Base")
		if err != nil {
			return err
		}

		var base Base
		err = inline.Copy(i, &base)
		if err != nil {
			return errors.Wrap(err, "copy")
		}
		if base.RequestID == "" {
			base.RequestID = inline.RandString(32)
		}
		ctx = SetBase(ctx, &base)
		return call(ctx, method, args, result)
	}
}
