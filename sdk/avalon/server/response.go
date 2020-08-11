package server

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"reflect"
)

type BaseResp struct {
	Code    int32  `thrift:"code,1" db:"code" json:"code"`
	Message string `thrift:"message,2" db:"message" json:"message"`
}

func ConvertResponseAdvance(call Call) Call {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := call(ctx, request)
		return ConvertErr2Resp(resp, err)
	}
}

func ConvertErr2Resp(resp interface{}, err error) (interface{}, error) {
	respType := reflect.ValueOf(resp)
	if respType.IsNil() {
		respType = reflect.New(respType.Type().Elem())
		resp = respType.Interface()
	}
	baseRespField := respType.Elem().FieldByName("BaseResp")
	if baseRespField.IsNil() {
		newBaseRespValue := reflect.New(baseRespField.Type().Elem()) // ptr
		baseRespField.Set(newBaseRespValue)
	}
	if err != nil {
		var baseResp BaseResp
		aErr, ok := err.(inline.AvalonError)
		if !ok {
			baseResp.Code = int32(inline.Unknown)
			baseResp.Message = err.Error()
		} else {
			baseResp.Code = int32(aErr.Code())
			baseResp.Message = aErr.Error()
		}
		baseRespField.Elem().FieldByName("Code").Set(reflect.ValueOf(baseResp.Code))
		baseRespField.Elem().FieldByName("Message").Set(reflect.ValueOf(baseResp.Message))
	}

	return resp, nil
}
