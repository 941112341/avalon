package client

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"time"
)

type ZkIPDiscover struct {
	namespace string

	Session      []string `default:"localhost:2181"`
	Timeout      string   `default:"2s"`
	RegisterTime string   `default:"10s"`

	Path string `viper:"server.discover.path"`

	timeout time.Duration
	client  *zookeeper.ZkClient

	hostports []string
	PSM       string
}

func (z *ZkIPDiscover) Key() string {
	return "zookeeper"
}

func (z *ZkIPDiscover) Initial() error {
	z.timeout = inline.Parse(z.Timeout)

	cli, err := zookeeper.NewClient(zookeeper.ZkConfig{
		HostPorts:      z.Session,
		SessionTimeout: z.timeout,
	})
	if err != nil {
		return err
	}
	z.client = cli

	node := zookeeper.NewZkNodeBuilder(path.Join(z.Path, z.PSM)).Build()
	err = node.ListWL(cli, false, func(event zk.Event) {
		switch event.Type {
		case zk.EventNodeCreated, zk.EventNodeDeleted:
			z.doRefresh(event, node, cli)
		}
	})
	z.hostports = node.GetChildrenKey()
	return err
}

func (z *ZkIPDiscover) doRefresh(event zk.Event, node *zookeeper.ZkNode, cli *zookeeper.ZkClient) {
	inline.WithFields("path", event.Path).Infoln("receive a register event")
	if err := node.List(cli, false); err != nil {
		inline.WithFields("path", event.Path).Errorln("refresh host list fail")
	} else {
		inline.WithFields("hostports", node.GetChildrenKey()).Infoln("refresh host list success")
		z.hostports = node.GetChildrenKey()
	}
}

func (z *ZkIPDiscover) Destroy() error {
	z.client.Conn.Close()
	return nil
}

func (z *ZkIPDiscover) GetHostports() []string {
	return z.hostports
}
