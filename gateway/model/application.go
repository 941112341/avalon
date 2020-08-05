package model

import "context"

type ApplicationKey interface {
	GetApplication() Application
}

type Application interface {
	Invoker(ctx context.Context, request *HttpRequest) (*HttpResponse, error)
}

type Converter interface {
	ConvertRequest(request *HttpRequest) (interface{}, error)
	ConvertResponse(data interface{}) (*HttpResponse, error)
}
