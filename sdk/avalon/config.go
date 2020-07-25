package avalon

import (
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	jsoniter "github.com/json-iterator/go"
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

type DefaultConfigBuilder struct {
	once sync.Once
	cfg  *Config
	psm  string
}

func (d *DefaultConfigBuilder) Config() Config {
	cfg := d.cfg
	var err error
	d.once.Do(func() {
		err = config.Read(cfg, inline.GetEnv("base", "base.yaml"))
		if err != nil {
			inline.WithFields("err", err.Error()).Errorln("read fail")
		}
	})
	return cfg.newConfigWithPSM(d.psm)
}

func NewConfigBuilder(psm string, others ...Config) *DefaultConfigBuilder {
	cfg := Config{
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
	for _, other := range others {
		str := inline.ToJsonString(other)
		_ = jsoniter.UnmarshalFromString(str, &cfg)
	}
	return &DefaultConfigBuilder{
		once: sync.Once{},
		cfg:  &cfg,
		psm:  psm,
	}
}

func (cfg *Config) newConfigWithPSM(psm string) Config {
	return Config{
		PSM:      psm,
		Client:   cfg.Client,
		ZkConfig: cfg.ZkConfig,
		Server:   cfg.Server,
	}
}
