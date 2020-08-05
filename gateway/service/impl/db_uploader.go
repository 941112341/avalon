package impl

import (
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/gateway/repository"
	"github.com/941112341/avalon/sdk/inline"
)

func init() {
	_ = registry.Registry("UploaderService", &DBUploader{})
}

type DBUploader struct {
	Repo repository.UploadRepository `inject:"UploadRepository"`
}

func (D *DBUploader) GroupContent(key repository.UploadGroupKey) (map[string]string, error) {
	groups, err := D.Repo.FindGroup(&key)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "find group key=%=v", key)
	}

	result := make(map[string]string)
	for _, group := range groups {
		result[group.Base] = group.Content
	}
	return result, nil
}
