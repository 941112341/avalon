package model

type Gateway interface {
	getMapperRules() (MapperRules, error)

	Registry()
	Transfer(request HttpRequest) (*HttpResponse, error)
}

type HttpRequest struct {
	Headers    map[string]string
	Body       string
	URL        string
	HTTPMethod string
}

type HttpResponse struct {
	HTTPCode int
	Headers  map[string]string
	Body     string
}
