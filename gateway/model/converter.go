package model

import (
	"context"
	"net/http"
)

type Converter interface {
	ConvertRequest(ctx context.Context, request *http.Request) (interface{}, error)
	ConvertResponse(ctx context.Context, data interface{}) (*HttpResponse, error)
}
