package model

import "net/http"

type MapperRules interface {
	GetApplication(request *http.Request) (Application, error)
}

type MapperRule interface {
	GetType() MapperRuleType
	Match(request *http.Request) (ApplicationKey, bool)
}

type MapperRuleType int32
