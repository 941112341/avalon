package client

import (
	"errors"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"math/rand"
	"path"
	"time"
)

type Balancer interface {
	Choice() (string, error)
	Choices() []string
}

// 完全随机器
type DiscoverBalancer struct {
	PSM     string
	Session []string
	Timeout time.Duration

	Cache []string
}

func (d *DiscoverBalancer) Choices() []string {
	return d.Cache
}

func (d *DiscoverBalancer) Choice() (string, error) {
	if len(d.Cache) == 0 {
		return "", errors.New("discover error")
	}
	idx := rand.Intn(len(d.Cache))
	return d.Cache[idx], nil
}

func (d *DiscoverBalancer) Refresh() error {
	cli, err := zookeeper.NewClient(zookeeper.ZkConfig{
		HostPorts:      d.Session,
		SessionTimeout: d.Timeout,
	})

	if err != nil {
		return inline.PrependErrorFmt(err, "connect zk fail %+v", *d)
	}

	node := zookeeper.NewZkNodeBuilder(path.Join("/host", d.PSM)).Build()
	err = node.ListWL(cli, false, func(event zk.Event) {
		switch event.Type {
		case zk.EventNodeCreated, zk.EventNodeDeleted:
			d.doRefresh(event, node, cli)
		}
	})
	if err != nil {
		return inline.PrependErrorFmt(err, "get host fail")
	}
	d.Cache = node.GetChildrenKey()
	return nil
}

func (d *DiscoverBalancer) doRefresh(event zk.Event, node *zookeeper.ZkNode, cli *zookeeper.ZkClient) {
	inline.WithFields("path", event.Path).Infoln("receive a register event")
	if err := node.List(cli, false); err != nil {
		inline.WithFields("path", event.Path).Errorln("refresh host list fail")
	} else {
		inline.WithFields("hostports", node.GetChildrenKey()).Infoln("refresh host list success")
		d.Cache = node.GetChildrenKey()
	}
}

type BalancerBuilder struct {
	balancer *DiscoverBalancer
}

func NewBalancerBuilder() *BalancerBuilder {
	return &BalancerBuilder{balancer: &DiscoverBalancer{Session: []string{"localhost:2181"}, Timeout: time.Second * 10}}
}

func (b *BalancerBuilder) PSM(psm string) *BalancerBuilder {
	b.balancer.PSM = psm
	return b
}

func (b *BalancerBuilder) Timeout(timeout time.Duration) *BalancerBuilder {
	b.balancer.Timeout = timeout
	return b
}

func (b *BalancerBuilder) Session(session []string) *BalancerBuilder {
	b.balancer.Session = session
	return b
}

func (b *BalancerBuilder) Build() (*DiscoverBalancer, error) {
	if err := b.balancer.Refresh(); err != nil {
		return nil, err
	}
	return b.balancer, nil
}
