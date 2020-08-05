package model

type Converter interface {
	ConvertRequest(request *HttpRequest) (interface{}, error)
	ConvertResponse(data interface{}) (*HttpResponse, error)
}
