package impl

import (
	"errors"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/json-iterator/go"
)

const MapperKey = "mapper"

var (
	ErrArgsKeyExists = errors.New("mapper key exists")
)

// default json
type BaseConverter struct {
	mapperArgs map[string]interface{}
}

func (b *BaseConverter) ConvertRequest(request *model.HttpRequest) (interface{}, error) {
	raw := make(map[string]interface{})
	err := jsoniter.UnmarshalFromString(request.Body, &raw)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "unmarshal body %s", request.Body)
	}
	if _, ok := raw[MapperKey]; ok {
		return nil, inline.PrependErrorFmt(ErrArgsKeyExists, "body %s", request.Body)
	} else {
		raw[MapperKey] = b.mapperArgs
	}
	return map[string]interface{}{
		"request": raw,
	}, nil
}

func (b *BaseConverter) ConvertResponse(data interface{}) (*model.HttpResponse, error) {
	switch d := data.(type) {
	case string:
		return DefaultHttpResponse(d), nil
	case error:
		return &model.HttpResponse{
			HTTPCode: 200, Headers: DefaultHttpResponseHeader(), Body: d.Error(),
		}, nil
	case inline.AvalonError:
		return &model.HttpResponse{
			HTTPCode: int(d.Code()),
			Headers:  DefaultHttpResponseHeader(),
			Body:     "error happend",
		}, nil
	default:
		return DefaultHttpResponse(inline.ToJsonString(data)), nil
	}
}

func DefaultHttpResponse(body string) *model.HttpResponse {
	return &model.HttpResponse{
		HTTPCode: 200,
		Headers:  DefaultHttpResponseHeader(),
		Body:     body,
	}

}

func DefaultHttpResponseHeader() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}
