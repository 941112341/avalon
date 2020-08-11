package server

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/inline"
	"reflect"
)

var selfBase *both.Base

func SetBase2Consistent(call Call) Call {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ctx = SetBaseConsistent(ctx, request)
		return call(ctx, request)
	}
}

func SetBaseConsistent(ctx context.Context, request interface{}) context.Context {
	requestValue := reflect.ValueOf(request)
	base := requestValue.Elem().FieldByName("Base").Interface()

	var baseEntity both.Base
	inline.MustCopy(base, &baseEntity)
	ctx = both.SetBase(ctx, selfBase)
	return both.SetBase(ctx, &baseEntity)
}
