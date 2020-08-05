package repository

import (
	"github.com/941112341/avalon/gateway/database"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/sdk/inline"
)

func init() {
	_ = registry.Registry("UploadRepository", &uploadRepository{})
}

type UploadRepository interface {
	Insert(vo *UploadVo) error
	FindByKey(id *UploadUnionKey) (*UploadVo, error)
	FindGroup(id *UploadGroupKey) ([]*UploadVo, error)
}

type uploadRepository struct {
}

func (r uploadRepository) FindGroup(id *UploadGroupKey) ([]*UploadVo, error) {
	uploads := make([]*UploadVo, 0)
	if err := database.DBRead.Where(&UploadVo{UploadUnionKey: UploadUnionKey{UploadGroupKey: *id}}).Find(&uploads).Error; err != nil {
		return nil, inline.PrependErrorFmt(err, "query = %s", inline.ToJsonString(id))
	}
	return uploads, nil
}

func (uploadRepository) Insert(UploadVo *UploadVo) error {
	return database.DB.Create(UploadVo).Error
}

func (uploadRepository) FindByKey(id *UploadUnionKey) (*UploadVo, error) {
	upload := UploadVo{}
	if err := database.DBRead.Where(&UploadVo{UploadUnionKey: *id}).First(&upload).Error; err != nil {
		return nil, inline.PrependErrorFmt(err, "query = %s", inline.ToJsonString(id))
	}
	return &upload, nil
}
