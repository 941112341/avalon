package avalon

import (
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"sync"
	"time"
)

type Config struct {
	Psm string

	Client struct {
		Retry    int
		Wait     time.Duration // 100ms
		HostPort string
		Timeout  time.Duration // 1s
	}

	Server struct {
		Port    int
		Timeout time.Duration // 1s
	}

	ZkConfig zookeeper.ZkConfig
}

// default config
var _cfg = Config{
	Psm: "",
	Client: struct {
		Retry    int
		Wait     time.Duration
		HostPort string
		Timeout  time.Duration
	}{
		Retry:   0,
		Wait:    100,
		Timeout: 1,
	},
	Server: struct {
		Port    int
		Timeout time.Duration
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
