package initial

import (
	"github.com/941112341/avalon/example/idgenerator/database"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/example/idgenerator/util"
)

/*
this package should only use for

*/

func InitAll(args ...interface{}) error {
	var err error
	err = util.InitConfig()
	if err != nil {
		return err
	}
	err = database.InitDatabase()
	if err != nil {
		return err
	}
	err = registry.InitInject()
	if err != nil {
		return err
	}
	return nil
}
