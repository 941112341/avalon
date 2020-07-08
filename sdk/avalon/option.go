package avalon

import (
	"github.com/941112341/avalon/sdk/zookeeper"
	"time"
)

type Config struct {
	Timeout     time.Duration `yaml:"Timeout"`
	HostPort    string        `yaml:"HostPort"`
	ServiceName string        `yaml:"ServiceName"`

	zookeeper.ZkConfig `yaml:"ZkConfig"`
}

type ClientConfig struct {
	Retry    int           `yaml:"Retry"`
	Wait     time.Duration `yaml:"Wait"`
	ClientIp string        `yaml:"ClientIp"`
	Config   `yaml:"Config"`
}

type ServerConfig struct {
	Config `yaml:"Config"`
}
