package model

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/jinzhu/gorm"
	"time"
)

type IDLFileID struct {
	PSM  string
	Base string
}

type IDLFile struct {
	ID int64
	IDLFileID

	Content string
	Version int
	Deleted *bool
	Created time.Time
	Updated time.Time
}

func (IDLFile) TableName() string {
	return "upload"
}

func (i *IDLFile) Upload(repo UploadRepository) error {
	old, err := repo.FindByKey(i.IDLFileID)
	if err != nil {
		if !inline.IsErr(err, gorm.ErrRecordNotFound) {
			return err
		}
		return repo.Insert(i)
	} else {
		i.Version = old.Version
		return repo.Update(i)
	}
}

func (i *IDLFile) Get(repo UploadRepository) (*IDLFile, error) {
	result, err := repo.FindByKey(i.IDLFileID)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "upload %s", inline.ToJsonString(i))
	}
	return result, nil
}
