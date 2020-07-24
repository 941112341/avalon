package initial

import (
	"github.com/941112341/avalon/example/idgenerator/database"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/example/idgenerator/util"
	"os"
	"path"
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

func InitAllForTest(args ...interface{}) error {
	GOPATH := os.Getenv("GOPATH")
	dir := "src/github.com/941112341/avalon/example/idgenerator/conf/config.yaml"
	base := "src/github.com/941112341/avalon/example/idgenerator/base.yaml"
	conf := path.Join(GOPATH, dir)
	base = path.Join(GOPATH, base)
	_ = os.Setenv("conf", conf)
	_ = os.Setenv("base", base)

	return InitAll(args...)
}
