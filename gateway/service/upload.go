package service

import (
	"github.com/941112341/avalon/gateway/repository"
)

type UploadService interface {
	GroupContent(key repository.UploadGroupKey) (map[string]string, error)
	SaveGroupContent(request *SaveGroupContentRequest) error
}

type SaveGroupContentRequest struct {
	PSM     string
	Version string

	Data map[string]string
}
