package repository

import (
	"errors"
	"github.com/941112341/avalon/pkg/mygorm"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/jinzhu/gorm"
)

type UploadGroupKey struct {
	PSM     string
	Version string
}

func (k *UploadGroupKey) valid() error {
	if k == nil {
		return errors.New("key is nil")
	}
	if k.PSM == "" {
		return errors.New("psm is nil")
	}
	if k.Version == "" {
		return errors.New("version is nil")
	}
	return nil
}

type UploadUnionKey struct {
	UploadGroupKey
	Base string
}

func (k *UploadUnionKey) valid() error {
	if k == nil {
		return errors.New("key is nil")
	}
	if err := k.UploadGroupKey.valid(); err != nil {
		return err
	}
	if k.Base == "" {
		return errors.New("base is nil")
	}
	return nil
}

type UploadVo struct {
	mygorm.Model

	UploadUnionKey
	Content string
}

func (vo *UploadVo) valid() error {
	if vo == nil {
		return errors.New("vo is nil")
	}
	_, err := generic.NewThriftGroup(map[string]string{vo.Base: vo.Content})
	if err != nil {
		return err
	}
	return nil
}

func (vo *UploadVo) BeforeCreate(scope *gorm.Scope) error {
	if vo.Version == "" {
		vo.Version = inline.RandString(32)
	}
	return vo.Model.BeforeCreate(scope)
}

func (vo *UploadVo) BeforeUpdate(scope *gorm.Scope) error {
	if vo.Version == "" {
		vo.Version = inline.RandString(32)
	}
	return vo.Model.BeforeUpdate(scope)
}

func (*UploadVo) TableName() string {
	return "upload"
}
