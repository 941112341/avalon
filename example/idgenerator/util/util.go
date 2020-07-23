package util

import (
	"github.com/941112341/avalon/example/idgenerator/conf"
	"github.com/941112341/avalon/sdk/config"
)

func InitConfig(args ...interface{}) error {
	return config.Read(&conf.Config, "conf/config.yaml")
}
