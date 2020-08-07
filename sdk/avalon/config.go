package avalon

import (
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"sync"
	"time"
)

type ClientConfig struct {
	Retry    int           `yaml:"retry",json:"retry,omitempty"`
	Wait     time.Duration `yaml:"wait",json:"wait,omitempty"` // 100ms
	HostPort string        `yaml:"hostPort",json:"hostPort,omitempty"`
	Timeout  time.Duration `yaml:"timeout",json:"timeout,omitempty"` // 1s
}

type ServerConfig struct {
	Port    int           `yaml:"port",json:"port,omitempty"`
	Timeout time.Duration `yaml:"timeout",json:"timeout,omitempty"` // 1s
}

type Config struct {
	PSM string `yaml:"psm",json:"psm,omitempty"`

	Client ClientConfig `yaml:"client",json:"client,omitempty"`

	Server ServerConfig `yaml:"server",json:"server,omitempty"`

	ZkConfig zookeeper.ZkConfig `yaml:"zkConfig",json:"zkConfig,omitempty"`
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
			inline.WithFields("err", err).Errorln("read fail")
		}
	})
	return cfg.newConfigWithPSM(d.psm)
}

func NewConfigBuilder(psm string, others ...ConfigSetter) *DefaultConfigBuilder {
	cfg := Config{
		PSM: "",
		Client: ClientConfig{
			Retry:   0,
			Wait:    100,
			Timeout: 1,
		},
		Server: ServerConfig{
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
		cfg = other(cfg)
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

type ConfigSetter func(cfg Config) Config

func TimeoutSetter(timeout time.Duration) ConfigSetter {
	return func(cfg Config) Config {
		cfg.Client.Timeout = timeout
		return cfg
	}
}
