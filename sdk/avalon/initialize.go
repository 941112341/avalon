package avalon

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
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
	node := zookeeper.NewZkNodeBuilder(inline.JoinPath(cfg.ZkConfig.Path, cfg.Psm, hostPort)).Build()
	err = node.Save(zkCli, zk.FlagEphemeral)
	if err != nil {
		return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
	}

	inline.Infoln("register success", inline.NewPair("serviceName", cfg.Psm))
	return nil
}
