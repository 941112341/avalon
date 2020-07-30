package model

import (
	"github.com/941112341/avalon/common/client"
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

	Repo UploadRepository `gorm:"-"`
}

func (IDLFile) TableName() string {
	return "upload"
}

func (i *IDLFile) Upload() error {
	old, err := i.Repo.FindByKey(i.IDLFileID)
	if err != nil {
		if !inline.IsErr(err, gorm.ErrRecordNotFound) {
			return err
		}
		i.ID = client.GenID()
		return i.Repo.Insert(i)
	} else {
		i.Version = old.Version
		return i.Repo.Update(i)
	}
}

func (i *IDLFile) Get() (*IDLFile, error) {
	result, err := i.Repo.FindByKey(i.IDLFileID)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "upload %s", inline.ToJsonString(i))
	}
	return result, nil
}
