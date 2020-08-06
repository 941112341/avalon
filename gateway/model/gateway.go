package model

import (
	"context"
	"github.com/941112341/avalon/gateway/service"
	"net/http"
)

type Gateway interface {
	GetMapperRules() (MapperRules, error)
	ClearRules()
	Transfer(ctx context.Context, request *http.Request) (*HttpResponse, error)

	AddMapper(ctx context.Context, request *service.MapperData) error
	AddUploader(ctx context.Context, request *service.SaveGroupContentRequest) error
}

type HttpResponse struct {
	HTTPCode int
	Headers  map[string]string
	Body     string
}
