package repository

import (
	"github.com/941112341/avalon/gateway/database"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/t-tiger/gorm-bulk-insert/v2"
)

type MapperRepository interface {
	AllMapper() (MapperList, error)
	AddMapper(mapperList MapperList) error
	DelMapper(mapperList MapperList) error
}

func init() {
	_ = registry.Registry("MapperRepository", &mapperRepository{})
}

type mapperRepository struct {
}

func (mapperRepository) AllMapper() (MapperList, error) {
	mappers := make(MapperList, 0)

	if err := database.DBRead.Model(&MapperVo{}).Where("deleted = 0").Order("type asc").Order("updated desc").Find(&mappers).Error; err != nil {
		return nil, err
	}
	return mappers, nil
}

func (mapperRepository) AddMapper(mapperList MapperList) error {
	if mapperList.IsEmpty() {
		return nil
	}
	ifaces := make([]interface{}, 0)
	for _, mapper := range mapperList {
		if err := mapper.BeforeCreate(nil); err != nil {
			return inline.PrependErrorFmt(err, "init mapper err %+v", mapper)
		}

		ifaces = append(ifaces, mapper)
	}

	return gormbulk.BulkInsert(database.DB, ifaces, 1000)
}

func (mapperRepository) DelMapper(mapperList MapperList) error {
	// by ids
	if mapperList.IsEmpty() {
		return nil
	}
	ids := make([]int64, 0)
	for _, mapper := range mapperList {
		ids = append(ids, mapper.ID)
	}

	return database.DB.Model(&MapperVo{}).Where("id in (?)", ids).Update("Deleted", true).Error
}
