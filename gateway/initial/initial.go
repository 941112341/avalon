package initial

import (
	"github.com/941112341/avalon/gateway/conf"
	"github.com/941112341/avalon/gateway/service"
)

/*
this package should only use for
*/

func InitAll(args ...interface{}) error {
	var err error
	err = conf.InitViper()
	if err != nil {
		return err
	}

	return service.Initial()
}

func InitAllForTest(args ...interface{}) error {

	return InitAll(args...)
}
