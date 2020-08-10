package common

import (
	"github.com/941112341/avalon/common/gen/base"
	"github.com/941112341/avalon/sdk/inline"
	"reflect"
)

func ConvertResponse(request, response interface{}, err error) (interface{}, error) {
	if reflect.ValueOf(response).IsNil() {
		elem := reflect.TypeOf(response).Elem()
		response = reflect.New(elem).Interface()
	}

	baseRespField := reflect.ValueOf(response).Elem().FieldByName("BaseResp")
	if baseRespField.IsNil() {
		var baseResp = &base.BaseResp{}
		if err != nil {
			aErr, ok := err.(inline.AvalonError)
			if !ok {
				baseResp = &base.BaseResp{Message: err.Error(), Code: int32(inline.Unknown)}

			} else {
				baseResp = &base.BaseResp{Message: aErr.Error(), Code: int32(aErr.Code())}
			}
		}
		baseRespField.Set(reflect.ValueOf(baseResp))
	}
	message := "success"
	if err != nil {
		message = "fail"
	}
	inline.WithFields("request", request, "response", response, "err", err).Infoln(message)

	return response, nil
}
