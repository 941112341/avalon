package client

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/inline"
	"reflect"
)

func SetBaseArgs(ctx context.Context, args interface{}) {
	req := reflect.ValueOf(args).Elem().FieldByName("Request").Interface()
	SetBaseConsistent(ctx, req)
}

func SetBaseConsistent(ctx context.Context, req interface{}) {
	base := both.GetBase(ctx)
	reqValue := reflect.ValueOf(req)
	baseField := reqValue.Elem().FieldByName("Base")
	if baseField.IsNil() {
		newBaseField := reflect.New(baseField.Type().Elem())
		baseField.Set(newBaseField)
	}

	inline.MustCopy(base, baseField.Interface())
}
