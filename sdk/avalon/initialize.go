package avalon

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func RegisterService(cfg Config) error {
	zkCli, err := zookeeper.GetZkClientInstance(cfg.ZkConfig)
	if err != nil {
		return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
	}
	ip, err := inline.InetAddress()
	if err != nil {
		return err
	}

	hostPort := fmt.Sprintf("%s:%d", ip, cfg.Server.Port)
	node := zookeeper.NewZkNodeBuilder(inline.JoinPath(cfg.ZkConfig.Path, cfg.PSM, hostPort)).Build()
	err = node.Save(zkCli, zk.FlagEphemeral)
	if err != nil {
		return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
	}
	err = node.GetWL(zkCli, func(event zk.Event) {
		if event.Type != zk.EventNodeDeleted {
			return
		}
		node := zookeeper.NewZkNodeBuilder(event.Path).Build()
		_ = inline.RetryFun(func() error {
			return node.Save(zkCli, zk.FlagEphemeral)
		}, -1, time.Second, func(wt time.Duration, retry int) time.Duration {
			return time.Minute
		})
	})
	if err != nil {
		return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
	}
	inline.Infoln("register success", inline.NewPair("serviceName", cfg.PSM))
	return nil
}
