package model

import (
	"context"
	"net/http"
)

type ApplicationKey interface {
	GetApplication() Application
}

type Application interface {
	Invoker(ctx context.Context, request *http.Request) (*HttpResponse, error)
}
