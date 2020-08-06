package repository

import (
	"github.com/941112341/avalon/gateway/database"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/sdk/inline"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

func init() {
	_ = registry.Registry("UploadRepository", &uploadRepository{})
}

type UploadRepository interface {
	Insert(vo *UploadVo) error
	BatchInsert(vos []*UploadVo) error
	FindByKey(id *UploadUnionKey) (*UploadVo, error)
	FindGroup(id *UploadGroupKey) ([]*UploadVo, error)
}

type uploadRepository struct {
}

func (r uploadRepository) BatchInsert(vos []*UploadVo) error {
	if len(vos) == 0 {
		return nil
	}
	iface := make([]interface{}, 0)
	for _, vo := range vos {
		if err := vo.BeforeCreate(nil); err != nil {
			return inline.PrependErrorFmt(err, "vo %+v", vo)
		}
		iface = append(iface, *vo)
	}

	return gormbulk.BulkInsert(database.DB, iface, 1000)
}

func (r uploadRepository) FindGroup(id *UploadGroupKey) ([]*UploadVo, error) {
	if err := id.valid(); err != nil {
		return nil, err
	}
	uploads := make([]*UploadVo, 0)
	if err := database.DBRead.Where(&UploadVo{UploadUnionKey: UploadUnionKey{UploadGroupKey: *id}}).Find(&uploads).Error; err != nil {
		return nil, inline.PrependErrorFmt(err, "query = %s", inline.ToJsonString(id))
	}
	return uploads, nil
}

func (uploadRepository) Insert(UploadVo *UploadVo) error {
	if err := UploadVo.valid(); err != nil {
		return err
	}
	return database.DB.Create(UploadVo).Error
}

func (uploadRepository) FindByKey(id *UploadUnionKey) (*UploadVo, error) {
	if err := id.valid(); err != nil {
		return nil, err
	}
	upload := UploadVo{}
	if err := database.DBRead.Where(&UploadVo{UploadUnionKey: *id}).First(&upload).Error; err != nil {
		return nil, inline.PrependErrorFmt(err, "query = %s", inline.ToJsonString(id))
	}
	return &upload, nil
}
