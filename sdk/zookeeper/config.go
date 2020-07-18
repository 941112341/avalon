package zookeeper

import "time"

type ZkConfig struct {
	HostPorts      []string      `yaml:"hostPorts"`
	SessionTimeout time.Duration `yaml:"sessionTimeout"` // second
	Path           string        `yaml:"path"`
}
