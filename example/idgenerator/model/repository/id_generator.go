package repository

import (
	"github.com/941112341/avalon/example/idgenerator/database"
	"github.com/941112341/avalon/example/idgenerator/registry"
)

const BizID = "idgenerator"

type idGeneratorRepository struct {
}

func (g idGeneratorRepository) FindByMaxIDBetween(left, right int64) (*IdGenerator, error) {
	var gen IdGenerator
	err := database.DBRead.Model(&gen).Where("max_id between ? and ? and biz_id <> ?", left, right, BizID).Order("max_id desc").First(&gen).Error
	if err != nil {
		return nil, err
	}
	return &gen, nil
}

func (g idGeneratorRepository) UpdateVersion(generator IdGenerator) (int64, error) {
	db := database.DB.Model(&generator).Where("version = ? and biz_id = ?", generator.Version, generator.BizID).Update(map[string]interface{}{
		"max_id": generator.MaxID, "length": generator.Length, "version": generator.Version + 1,
	})
	return db.RowsAffected, db.Error
}

func (g idGeneratorRepository) Save(generator IdGenerator) error {
	return database.DB.Save(&generator).Error
}

func (g idGeneratorRepository) Get() (*IdGenerator, error) {
	var gen IdGenerator
	if err := database.DB.Where(IdGenerator{BizID: BizID}).First(&gen).Error; err != nil {
		return nil, err
	}
	return &gen, nil
}

func init() {
	_ = registry.Registry("IdGeneratorRepository", &idGeneratorRepository{})
}
