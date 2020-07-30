package util

import (
	"github.com/941112341/avalon/gateway/conf"
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/inline"
)

func InitConfig(args ...interface{}) error {
	return config.Read(&conf.Config, inline.GetEnv("conf", "./conf/config.yaml"))
}
