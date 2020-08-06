package model

import "net/http"

type Converter interface {
	ConvertRequest(request *http.Request) (interface{}, error)
	ConvertResponse(data interface{}) (*HttpResponse, error)
}
