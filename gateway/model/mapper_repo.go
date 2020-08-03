package model

import (
	"github.com/941112341/avalon/gateway/database"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/t-tiger/gorm-bulk-insert/v2"
)

type MapperRepository interface {
	AllMapper() (*MapperList, error)
	AddMapper(mapperList *MapperList) error
	DelMapper(mapperList *MapperList) error
}

func init() {
	_ = registry.Registry("MapperRepository", &mapperRepository{})
}

type mapperRepository struct {
}

func (mapperRepository) AllMapper() (*MapperList, error) {
	mappers := make([]*Mapper, 0)

	if err := database.DBRead.Model(&Mapper{}).Where("deleted = 0").Find(&mappers).Error; err != nil {
		return nil, err
	}
	return &MapperList{Mappers: mappers, Absolute: map[string]*Mapper{}}, nil
}

func (mapperRepository) AddMapper(mapperList *MapperList) error {
	if mapperList.IsEmpty() {
		return nil
	}
	ifaces := make([]interface{}, 0)
	for _, mapper := range mapperList.Mappers {
		ifaces = append(ifaces, *mapper)
	}

	return gormbulk.BulkInsert(database.DB, ifaces, 1000)
}

func (mapperRepository) DelMapper(mapperList *MapperList) error {
	// by ids
	if mapperList.IsEmpty() {
		return nil
	}
	ids := make([]int64, 0)
	for _, mapper := range mapperList.Mappers {
		ids = append(ids, mapper.ID)
	}

	return database.DB.Where("id in (?)", ids).Delete(&Mapper{}).Error
}
