package avalon

import (
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"sync"
	"time"
)

type Config struct {
	PSM string `yaml:"psm"`

	Client struct {
		Retry    int           `yaml:"retry"`
		Wait     time.Duration `yaml:"wait"` // 100ms
		HostPort string        `yaml:"hostPort"`
		Timeout  time.Duration `yaml:"timeout"` // 1s
	}

	Server struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"` // 1s
	} `yaml:"server"`

	ZkConfig zookeeper.ZkConfig `yaml:"zkConfig"`
}

// default config
var _cfg = Config{
	PSM: "",
	Client: struct {
		Retry    int           `yaml:"retry"`
		Wait     time.Duration `yaml:"wait"` // 100ms
		HostPort string        `yaml:"hostPort"`
		Timeout  time.Duration `yaml:"timeout"` // 1s
	}{
		Retry:   0,
		Wait:    100,
		Timeout: 1,
	},
	Server: struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"` // 1s
	}{
		Port:    8888,
		Timeout: 1,
	},
	ZkConfig: zookeeper.ZkConfig{
		SessionTimeout: 30,
		HostPorts:      []string{"localhost:2181"},
		Path:           "/host",
	},
}

var once sync.Once

func GetConfig() (Config, error) {
	var err error
	once.Do(func() {
		err = config.Read(&_cfg, inline.GetEnv("base", "base.yaml"))
		if err != nil {
			return
		}
	})
	return _cfg, err
}
