package avalon

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

const (
	MetaKey = "__META_KEY__"
)

type Base struct {
	IP        string
	PSM       string
	Time      int64
	Extra     map[string]string
	RequestID string

	Base *Base
}

func metaMiddlewareClient(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		fld := reflect.ValueOf(args).Elem().FieldByName("Base")
		if !fld.IsValid() {
			return call(ctx, method, args, result)
		}

		scope := GetScope(ctx).find(FromCrossRPC)
		preMeta, _ := ctx.Value(MetaKey).(*Base)
		psm, _ := scope.Get(PSMKey, true)

		extra := make(map[string]string)
		scope.cache.Range(func(key, value interface{}) bool {
			k, v := key.(string), value.(string)
			if k == HostPortKey || k == PSMKey {
				return true
			}
			extra[k] = v
			return true
		})
		meta := &Base{
			IP:    inline.GetIP(),
			PSM:   psm,
			Time:  time.Now().Unix(),
			Extra: extra,
			Base:  preMeta,
		}
		if preMeta == nil || preMeta.RequestID == "" {
			meta.RequestID = inline.RandString(32)
		} else {
			meta.RequestID = preMeta.RequestID
		}
		err := inline.SetField(args, "Base", meta)
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
		base, ok := i.(Base)
		if !ok {
			return errors.WithMessage(err, inline.ToJsonString(i))
		}
		if base.RequestID == "" {
			base.RequestID = inline.RandString(32)
		}
		ctx = context.WithValue(ctx, MetaKey, &base)
		return call(ctx, method, args, result)
	}
}
