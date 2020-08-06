package service

import "github.com/941112341/avalon/gateway/repository"

type MapperService interface {
	FetchMapperList() (repository.MapperList, error)
	AddMapperRule(mapper *MapperData) error
}

type MapperData struct {
	URL     string
	Type    int16
	Domain  string
	PSM     string
	Base    string
	Method  string
	Version string
}
