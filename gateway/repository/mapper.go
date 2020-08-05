package repository

import (
	"github.com/941112341/avalon/pkg/mygorm"
	"github.com/jinzhu/gorm"
)

type MapperList []MapperVo

func (l MapperList) IsEmpty() bool {
	return len(l) == 0
}

type MapperVo struct {
	mygorm.Model
	URL     string
	Type    int16
	Domain  string
	PSM     string
	Base    string
	Method  string
	Version string
}

func (vo *MapperVo) BeforeSave(scope *gorm.Scope) error {
	return vo.Model.BeforeCreate(scope)
}

func (vo *MapperVo) BeforeUpdate(scope *gorm.Scope) error {
	return vo.Model.BeforeUpdate(scope)
}

func (*MapperVo) TableName() string {
	return "mapper"
}
