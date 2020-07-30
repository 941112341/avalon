package database

import (
	"github.com/941112341/avalon/gateway/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DBRead *gorm.DB
	DB     *gorm.DB
)

func InitDatabase() (err error) {
	DBRead, err = gorm.Open("mysql", conf.Config.Database.DBRead)
	if err != nil {
		return
	}

	DB, err = gorm.Open("mysql", conf.Config.Database.DBRead)
	return
}
