package model

import (
	"github.com/941112341/avalon/gateway/database"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/sdk/inline"
	"time"
)

func init() {
	_ = registry.Registry("UploadRepository", &uploadRepository{})
}

type UploadRepository interface {
	Insert(idlFile *IDLFile) error
	Update(idlFile *IDLFile) error
	FindByKey(id IDLFileID) (*IDLFile, error)
}

type uploadRepository struct {
}

func (uploadRepository) Insert(idlFile *IDLFile) error {
	return database.DB.Save(idlFile).Error
}

func (uploadRepository) Update(idlFile *IDLFile) error {
	db := idlFile
	version := db.Version
	dbResult := database.DB.Model(&IDLFile{}).Where(&IDLFile{Version: version, IDLFileID: idlFile.IDLFileID}).Update(map[string]interface{}{
		"content": idlFile.Content, "updated": time.Now(), "version": db.Version + 1,
	})
	if err := dbResult.Error; err != nil {
		return inline.PrependErrorFmt(err, "model:%+v", *db)
	}
	cnt := dbResult.RowsAffected
	if cnt == 0 {
		return inline.NewError(inline.ErrDBCas, "model update err:%+v", *db)
	}
	db.Version++

	return nil
}

func (uploadRepository) FindByKey(id IDLFileID) (*IDLFile, error) {
	upload := IDLFile{}
	if err := database.DBRead.Where(&IDLFile{IDLFileID: id}).Find(&upload).Error; err != nil {
		return nil, inline.PrependErrorFmt(err, "query = %s", inline.ToJsonString(id))
	}
	return &upload, nil
}
