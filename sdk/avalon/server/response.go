package server

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/inline"
	"reflect"
)

type ErrorWrapper struct {
	avalon.TodoBean
}

//type BaseResp struct {
//	Code int32 `thrift:"code,1" db:"code" json:"code"`
//	Message string `thrift:"message,2" db:"message" json:"message"`
//	Extra map[string]string `thrift:"extra,3" db:"extra" json:"extra"`
//}

func (e ErrorWrapper) Middleware(call avalon.Call) avalon.Call {
	return func(ctx context.Context, invoke *avalon.Invoke) error {
		field := reflect.ValueOf(invoke.Response).Elem().FieldByName("BaseResp")
		if field.IsNil() {
			fieldVal := reflect.New(field.Type().Elem())
			field.Set(fieldVal)
		}
		var code int32
		var message string

		err := call(ctx, invoke)
		if err != nil {
			aErr, ok := err.(inline.AvalonError)
			if ok {
				code = int32(aErr.Code())
			} else {
				code = -1 // unknown err
			}
			message = err.Error()
		}

		field.Elem().FieldByName("Code").Set(reflect.ValueOf(code))
		field.Elem().FieldByName("Message").Set(reflect.ValueOf(message))
		return nil
	}
}
