package repository

import (
	"github.com/941112341/avalon/example/idgenerator/database"
	"github.com/941112341/avalon/example/idgenerator/registry"
)

const bizID = "idgenerator"

type idGeneratorRepository struct {
}

func (g idGeneratorRepository) UpdateVersion(generator IdGenerator) (int64, error) {
	db := database.DB.Model(&generator).Where("version = ? and biz_id = ?", generator.Version, generator.BizID).Update(map[string]interface{}{
		"max_id": generator.MaxID, "length": generator.Length, "version": generator.Version + 1,
	})
	return db.RowsAffected, db.Error
}

func (g idGeneratorRepository) Save(generator IdGenerator) error {
	return database.DB.Model(&generator).Save(generator).Error
}

func (g idGeneratorRepository) Get() (*IdGenerator, error) {
	var gen IdGenerator
	if err := database.DB.Where(IdGenerator{BizID: bizID}).First(&gen).Error; err != nil {
		return nil, err
	}
	return &gen, nil
}

func init() {
	_ = registry.Registry("IdGeneratorRepository", &idGeneratorRepository{})
}
