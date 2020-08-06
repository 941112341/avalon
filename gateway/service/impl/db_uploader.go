package impl

import (
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/gateway/repository"
	"github.com/941112341/avalon/gateway/service"
	"github.com/941112341/avalon/pkg/mygorm"
	"github.com/941112341/avalon/sdk/inline"
)

func init() {
	_ = registry.Registry("UploaderService", &DBUploaderService{})
}

type DBUploaderService struct {
	Repo repository.UploadRepository `inject:"UploadRepository"`
}

func (D *DBUploaderService) SaveGroupContent(request *service.SaveGroupContentRequest) error {

	vos := make([]*repository.UploadVo, 0)
	for base, content := range request.Data {
		vos = append(vos, &repository.UploadVo{
			Model: mygorm.Model{},
			UploadUnionKey: repository.UploadUnionKey{
				UploadGroupKey: repository.UploadGroupKey{
					PSM:     request.PSM,
					Version: request.Version,
				},
				Base: base,
			},
			Content: content,
		})
	}
	return D.Repo.BatchInsert(vos)
}

func (D *DBUploaderService) GroupContent(key repository.UploadGroupKey) (map[string]string, error) {
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
