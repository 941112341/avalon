package model

import "net/http"

type MapperRules interface {
	GetApplication(request *http.Request) (Application, error)
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
}

type MapperRule interface {
	GetType() MapperRuleType
	Order() int
	Match(request *http.Request) (ApplicationKey, bool)
}

type MapperRuleType int32
