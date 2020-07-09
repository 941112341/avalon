package zookeeper

import "time"

type ZkConfig struct {
	HostPorts      []string      `yaml:"HostPorts"`
	SessionTimeout time.Duration `yaml:"SessionTimeout"` // second
	Path           string        `yaml:"path"`
}
