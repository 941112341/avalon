package repository

import (
	"github.com/941112341/avalon/pkg/mygorm"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/jinzhu/gorm"
)

type UploadGroupKey struct {
	PSM     string
	Version string
}

type UploadUnionKey struct {
	UploadGroupKey
	Base string
}

type UploadVo struct {
	mygorm.Model

	UploadUnionKey
	Content string
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
