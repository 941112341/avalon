package impl

import (
	"context"
	"errors"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

const MapperKey = "mapper"

var (
	ErrArgsKeyExists = errors.New("mapper key exists")
)

// default json
type BaseConverter struct {
	mapperArgs map[string]interface{}
}

func (b *BaseConverter) ConvertRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	raw := make(map[string]interface{})
	err = jsoniter.Unmarshal(body, &raw)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "unmarshal body %s", request.Body)
	}
	for k, v := range b.mapperArgs {
		old, ok := raw[k]
		if ok {
			inline.WithFields("key", k, "old", old).Warnln("parameter overwrite")
		}
		raw[k] = v
	}

	base := both.GetBase(ctx)
	if base == nil {
		return nil, errors.New("base not found")
	}
	raw["base"] = base
	return map[string]interface{}{
		"request": raw,
	}, nil
}

func (b *BaseConverter) ConvertResponse(ctx context.Context, data interface{}) (*model.HttpResponse, error) {
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
		return DefaultHttpResponse(data), nil
	}
}

func DefaultHttpResponse(body interface{}) *model.HttpResponse {
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
