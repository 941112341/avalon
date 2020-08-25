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
	Timeout      string   `default:"10s"`
	RegisterTime string   `default:"10s"`

	Path string `viper:"server.discover.path" default:"/host"`

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
	err = node.ListWL(cli, false, &watcher{
		z: z,
	})
	z.hostports = node.GetChildrenKey()
	return err
}

func (z *ZkIPDiscover) doRefresh() {
	node := zookeeper.NewZkNodeBuilder(z.path()).Build()
	inline.WithFields("path", z.path()).Infoln("receive a register event")
	if err := node.List(z.client, false); err != nil {
		inline.WithFields("path", z.path()).Errorln("refresh host list fail")
	} else {
		inline.WithFields("hostports", node.GetChildrenKey()).Infoln("refresh host list success")
		z.hostports = node.GetChildrenKey()
	}
}

func (z *ZkIPDiscover) path() string {
	return path.Join(z.Path, z.PSM)
}

func (z *ZkIPDiscover) Destroy() error {
	z.client.Conn.Close()
	return nil
}

func (z *ZkIPDiscover) GetHostports() []string {
	if len(z.hostports) == 0 {
		z.doRefresh()
	}
	return z.hostports
}

type watcher struct {
	z *ZkIPDiscover
}

func (w *watcher) WatchEvent(event zk.Event) {
	switch event.Type {
	case zk.EventNodeCreated, zk.EventNodeDeleted:
		w.z.doRefresh()
	}
}

func (w *watcher) WatchError(err error) {
	w.z.hostports = nil
	inline.WithFields("err", err).Errorln("watch %s instance fail", w.z.PSM)
}
