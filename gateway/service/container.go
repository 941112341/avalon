package service

import "github.com/941112341/avalon/gateway/registry"

func init() {
	_ = registry.Registry("", ServiceContainer)
}

var ServiceContainer = &Container{}

type Container struct {
	UploadService UploadService `inject:"UploadService"`
	MapperService MapperService `inject:"MapperService"`
}
