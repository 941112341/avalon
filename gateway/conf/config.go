package conf

import "context"
import conf "github.com/941112341/avalon/sdk/config"

type config struct {
	Https struct {
		Port int `yaml:"port"`
	} `yaml:"https"`
	Http struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	Database struct {
		DB     string `yaml:"DB"`
		DBRead string `yaml:"DBRead"`
	} `yaml:"Database"`
}

var Config config

func InitConfig(ctx context.Context) error {
	return conf.Read(&Config, "conf/config.yaml")
}
