package zookeeper

import "time"

type ZkConfig struct {
	HostPorts      []string      `yaml:"HostPorts"`
	SessionTimeout time.Duration `yaml:"SessionTimeout"` // 单位second
	Path           string        `yaml:"Path"`
}
