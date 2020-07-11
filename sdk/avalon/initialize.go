package avalon

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

func RegisterService(cfg CallConfig) error {
	zkCli, err := zookeeper.GetZkClientInstance(cfg.ZkConfig)
	if err != nil {
		return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
	}
	ip, err := inline.InetAddress()
	if err != nil {
		return err
	}
	idx := strings.LastIndex(cfg.HostPort, ":")
	port := cfg.HostPort[idx:]
	hostPort := ip + port
	node := zookeeper.NewZkNodeBuilder(inline.JoinPath(cfg.Path, cfg.Psm, hostPort)).Build()
	err = node.Save(zkCli, zk.FlagEphemeral)
	if err != nil {
		return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
	}

	inline.Infoln("register success", inline.NewPair("serviceName", cfg.Psm))
	return nil
}
