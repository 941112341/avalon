package service

import (
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/gateway/repository"
)

type UploadService interface {
	GroupContent(key repository.UploadGroupKey) (map[string]string, error)
}

var UploadBuilder = &UploadServiceBuilder{}

func init() {
	_ = registry.Registry("", UploadBuilder)
}

type UploadServiceBuilder struct {
	UploadService UploadService `inject:"UploadService"`
}
