package registry

import "github.com/facebookgo/inject"

var graph inject.Graph

func Registry(Name string, value interface{}) error {
	return graph.Provide(&inject.Object{Value: value, Name: Name})
}

func InitInject() error {
	return graph.Populate()
}
