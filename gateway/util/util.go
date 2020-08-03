package util

import (
	"github.com/941112341/avalon/gateway/conf"
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/inline"
	"strings"
)

func InitConfig(args ...interface{}) error {
	return config.Read(&conf.Config, inline.GetEnv("conf", "./conf/config.yaml"))
}

func StandardURL(url string) string {
	idx := strings.Index(url, "?")
	if idx > -1 {
		return StandardURL(url[:idx])
	}
	return strings.Trim(url, "/")
}
