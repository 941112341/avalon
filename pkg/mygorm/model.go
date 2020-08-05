package mygorm

import (
	"github.com/941112341/avalon/common/client"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/jinzhu/gorm"
	"time"
)

type Model struct {
	ID      int64
	Deleted *bool
	Created time.Time
	Updated time.Time
}

func (m *Model) BeforeCreate(*gorm.Scope) error {
	if m.ID == 0 {
		m.ID = client.GenID()
	}
	if m.Deleted == nil {
		m.Deleted = inline.BoolPtr(false)
	}
	if m.Created.IsZero() {
		m.Created = time.Now()
	}
	if m.Updated.IsZero() {
		m.Updated = time.Now()
	}
	return nil
}

func (m *Model) BeforeUpdate(*gorm.Scope) error {
	if m.Updated.IsZero() {
		m.Updated = time.Now()
	}
	return nil
}
