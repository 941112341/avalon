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
		Wait     time.Duration `default:"100"` // 100ms
		HostPort string
		Timeout  time.Duration `default:"1"` // 1s
	}

	Server struct {
		Port int `default:"8888"`
	}

	ZkConfig zookeeper.ZkConfig
}

var _cfg Config
var once sync.Once

func (cfg *Config) Initial() {
	inline.SetDefaultValue(cfg)
}

func GetConfig() (Config, error) {
	var err error
	once.Do(func() {
		err = config.Read(&_cfg, inline.GetEnv("base", "base.yaml"))
		if err != nil {
			return
		}
		_cfg.Initial()
	})
	return _cfg, err
}
