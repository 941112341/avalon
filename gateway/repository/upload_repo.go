package repository

import (
	"errors"
	"fmt"
	"github.com/941112341/avalon/gateway/database"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/sdk/inline"
	"strings"
)

func init() {
	_ = registry.Registry("UploadRepository", &uploadRepository{})
}

type UploadRepository interface {
	Insert(vo *UploadVo) error
	BatchInsert(vos []*UploadVo) error
	FindByKey(id *UploadUnionKey) (*UploadVo, error)
	FindGroup(id *UploadGroupKey) ([]*UploadVo, error)
	DeleteGroup(id *UploadGroupKey) error
}

type uploadRepository struct {
}

func (r uploadRepository) DeleteGroup(id *UploadGroupKey) error {
	if err := id.valid(); err != nil {
		return inline.PrependErrorFmt(err, "valid before deleted")
	}
	return database.DB.Model(&UploadVo{}).Where(&UploadVo{UploadUnionKey: UploadUnionKey{UploadGroupKey: *id}}).Update("Deleted", true).Error
}

func (r uploadRepository) BatchInsert(vos []*UploadVo) error {
	if len(vos) == 0 {
		return nil
	}
	if len(vos) > 100 {
		return errors.New("beyond 100")
	}
	unions := make([]string, 0)
	union := ` select ?, ?, ?, ?, ?, ?, ? `
	ifaces := make([]interface{}, 0)
	for i, vo := range vos {
		if err := vo.valid(); err != nil {
			return inline.PrependErrorFmt(err, "valid fail")
		}
		if i == 0 {
			union = ` select ? as id, ? as psm, ? as content, ? as base, ? as created, ? as updated, ? as version `
		}
		if err := vo.BeforeCreate(nil); err != nil {
			return inline.PrependErrorFmt(err, "vo %+v", vo)
		}
		ifaces = append(ifaces, vo.ID, vo.PSM, vo.Content, vo.Base, vo.Created, vo.Updated, vo.Version)
		unions = append(unions, union)
	}
	unionCompleted := strings.Join(unions, " union all ")

	raw := `
	insert into upload (id, psm, content, base, created, updated, version)
select b.* from upload  right join (
   	%s
    ) as b
on false = true
where not exists(
        select 1 from upload where deleted = 0 and version=? and psm = ?
    )

`
	rawCompleted := fmt.Sprintf(raw, unionCompleted)
	ifaces = append(ifaces, vos[0].Version, vos[0].PSM)
	inline.WithFields("raw", rawCompleted).Infoln("raw sql")
	return database.DB.Exec(rawCompleted, ifaces...).Error
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
