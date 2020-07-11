package zookeeper

import "time"

type ZkConfig struct {
	HostPorts      []string      `yaml:"HostPorts",default:"localhost:127.0.0.1:2181"`
	SessionTimeout time.Duration `yaml:"SessionTimeout",default:"30"` // second
	Path           string        `yaml:"Path",default:"/host"`
}
