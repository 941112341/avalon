package server

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"time"
)

type Zookeeper struct {
	Session      []string `default:"localhost:2181"`
	Timeout      string   `default:"2s"`
	RegisterTime string   `default:"10s"`
	PSM          string   `viper:"server.psm"`
	Port         int      `viper:"server.port"`
	Path         string   `viper:"server.discover.path" default:"/host"`

	timeout      time.Duration
	registerTime time.Duration
	client       *zookeeper.ZkClient
}

func (z *Zookeeper) Key() string {
	return "zookeeper"
}

func (z *Zookeeper) path() string {
	return path.Join(z.Path, z.PSM, z.hostport())
}

func (z *Zookeeper) Initial() (err error) {
	z.timeout = inline.Parse(z.Timeout)
	z.registerTime = inline.Parse(z.Timeout)

	z.client, err = zookeeper.NewClient(zookeeper.ZkConfig{
		HostPorts:      z.Session,
		SessionTimeout: z.timeout,
	})
	return
}

func (z *Zookeeper) Destroy() error {
	z.client.Conn.Close()
	return nil
}

func (z *Zookeeper) BeforeRun() error {
	return z.Register()
}

func (z *Zookeeper) Register() error {
	err := z.register()
	if err != nil {
		return err
	}
	inline.WithFields("path", z.path()).Infoln("start register")

	ticker := time.NewTicker(z.timeout)
	go func() {
		for range ticker.C {
			err := z.register()
			if err != nil {
				inline.WithFields("err", err).Errorln("register err")
			}
		}
	}()
	return nil
}

func (z *Zookeeper) register() error {

	cli := z.client

	node := zookeeper.NewZkNodeBuilder(z.path()).Build()
	if err := node.Save(cli, zk.FlagEphemeral); err != nil {
		return inline.PrependErrorFmt(err, "save node fail %s", z.path())
	}

	return nil
}

func (z *Zookeeper) hostport() string {
	return fmt.Sprintf("%s:%d", inline.GetIP(), z.Port)
}
