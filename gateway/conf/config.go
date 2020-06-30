package conf

import "context"
import conf "github.com/941112341/avalon/sdk/config"

type config struct {
	Https struct {
		Port int
	}
	Http struct {
		Port int
	}
}

var Config config

func InitConfig(ctx context.Context) error {
	return conf.Read(&Config, "conf/config.yaml")
}
