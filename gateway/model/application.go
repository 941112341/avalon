package model

import "context"

type ApplicationKey interface {
	GetApplication() Application
}

type Application interface {
	Invoker(ctx context.Context, request *HttpRequest) (*HttpResponse, error)
}
