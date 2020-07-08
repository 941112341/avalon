package avalon

import (
	"github.com/941112341/avalon/sdk/collect"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/pkg/errors"
	"time"
)

var DiscoverMap *collect.SyncMap
var ZkClient *zookeeper.ZkClient

// must after initialClient
func InitialDiscover(config zookeeper.ZkConfig) error {
	if DiscoverMap == nil {
		DiscoverMap = collect.NewSyncMap()
		if err := initialZkClient(config); err != nil {
			return errors.Cause(err)
		}
		ZkClient.ListenerTree(config.Path, DiscoverMap)
	}
	return nil
}

func initialZkClient(config zookeeper.ZkConfig) error {
	if ZkClient == nil {
		cli, err := zookeeper.NewClient(config.HostPorts, config.SessionTimeout*time.Second)
		if err != nil {
			return errors.WithMessage(err, "zk connection fail")
		}
		ZkClient = cli
	}
	return nil
}
