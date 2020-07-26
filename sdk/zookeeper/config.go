package zookeeper

import "time"

type ZkConfig struct {
	HostPorts      []string      `yaml:"hostPorts",json:"hostPorts,omitempty"`
	SessionTimeout time.Duration `yaml:"sessionTimeout",json:"sessionTimeout,omitempty"` // second
	Path           string        `yaml:"path",json:"path,omitempty"`
}
