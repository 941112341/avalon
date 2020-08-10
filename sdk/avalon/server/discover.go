package server

import (
	"errors"
	"fmt"
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"time"
)

type Discover interface {
	Register() error
	both.Hostport
}

type ZkDiscover struct {
	session  []string
	timeout  time.Duration
	psm      string
	hostport string
	port     string
}

func (z *ZkDiscover) Port() string {
	return z.port
}

func (z *ZkDiscover) Hostport() string {
	return z.hostport
}

func (z *ZkDiscover) valid() error {
	if z.psm == "" {
		return errors.New("psm is empty")
	}
	if z.hostport == "" {
		return errors.New("host port is nil")
	}
	return nil
}

func (z *ZkDiscover) nodePath() string {
	return path.Join("/host", z.psm, z.hostport)
}

type zkDiscoverBuilder struct {
	discover *ZkDiscover
}

func NewZkDiscoverBuilder() *zkDiscoverBuilder {
	return &zkDiscoverBuilder{discover: &ZkDiscover{timeout: 10 * time.Second, session: []string{"localhost:2181"}}}
}

func (b *zkDiscoverBuilder) Session(session []string) *zkDiscoverBuilder {
	b.discover.session = session
	return b
}

func (b *zkDiscoverBuilder) Timeout(t time.Duration) *zkDiscoverBuilder {
	b.discover.timeout = t
	return b
}

func (b *zkDiscoverBuilder) PSM(psm string) *zkDiscoverBuilder {
	b.discover.psm = psm
	return b
}

func (b *zkDiscoverBuilder) Hostport(hostport string) *zkDiscoverBuilder {
	b.discover.hostport = hostport
	return b
}

func (b *zkDiscoverBuilder) Port(port int) *zkDiscoverBuilder {
	ip := inline.GetIP()
	b.discover.hostport = fmt.Sprintf("%s:%d", ip, port)
	b.discover.port = fmt.Sprintf(":%d", port)
	return b
}

func (b *zkDiscoverBuilder) Build() *ZkDiscover {
	if err := b.discover.Register(); err != nil {
		panic(err)
	}
	return b.discover
}

func (z *ZkDiscover) Register() error {
	err := z.register()
	if err != nil {
		return err
	}
	inline.WithFields("psm", z.psm, "timeout", z.timeout).Infoln("start register")

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

func (z *ZkDiscover) register() error {
	if err := z.valid(); err != nil {
		return inline.PrependErrorFmt(err, "zk valid err")
	}

	cli, err := zookeeper.GetZkClientInstance(zookeeper.ZkConfig{
		HostPorts:      z.session,
		SessionTimeout: z.timeout,
	})
	if err != nil {
		return inline.PrependErrorFmt(err, "get zk client %+v", *z)
	}

	node := zookeeper.NewZkNodeBuilder(z.nodePath()).Build()
	if err := node.Save(cli, zk.FlagEphemeral); err != nil {
		return inline.PrependErrorFmt(err, "save node fail %s", z.nodePath())
	}

	return nil
}
