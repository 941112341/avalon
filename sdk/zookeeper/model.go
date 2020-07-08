package zookeeper

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
)

type Watcher interface {
	Watch(event zk.Event, oldNode, newNode *ZkNode)
}

type ZkNode struct {
	data     string
	children map[string]*ZkNode
	path     string
	parent   *ZkNode

	watches      []Watcher
	isPersistent bool
	varLock      sync.RWMutex

	ch <-chan zk.Event
}

// it's default Synchronous, should use 'go' if you want Asynchronous
func watchLoop(cfg *ZkConfig, ch <-chan zk.Event, loopFunc func(event zk.Event, cfg *ZkConfig) (<-chan zk.Event, error)) error {
	var err error
	select {
	case event := <-ch:
		ch, err = loopFunc(event, cfg)
		if err != nil {
			return errors.WithMessage(err, inline.JsonString(event))
		}
	}
	return watchLoop(cfg, ch, loopFunc)
}

func nodeLoop(cfg *ZkConfig, zkCli *ZkClient, builder func(cfg *ZkConfig, zkCli *ZkClient) (*ZkNode, error)) (*ZkNode, error) {
	node, err := builder(cfg, zkCli)
	if err != nil {
		return nil, err
	}

	go func() {
		watchLoop(cfg, node.ch, func(event zk.Event, cfg *ZkConfig) (<-chan zk.Event, error) {
			newNode, err := builder(cfg, zkCli)
			if err != nil {
				return nil, err
			}
			for _, watch := range newNode.watches {
				watch.Watch(event, node, newNode)
			}
			node = newNode
		})
	}()
	return node, nil
}

func NewZkNode(cfg *ZkConfig, watchers ...Watcher) (*ZkNode, error) {
	zkCli, err := GetZkClientInstance(cfg)
	if err != nil {
		return nil, err
	}
	data, _, ch, err := zkCli.Conn.GetW(cfg.Path)
	node := &ZkNode{
		path:         cfg.Path,
		isPersistent: true,
		data:         string(data),
	}
	go func() {
		err := watchLoop(cfg, ch, func(event zk.Event, cfg *ZkConfig) (<-chan zk.Event, error) {
			newData, _, ch, err := zkCli.Conn.GetW(cfg.Path)
			if err != nil {
				// maybe net error，todo
				return nil, err
			}
			newNode := &ZkNode{
				data:         string(newData),
				path:         cfg.Path,
				watches:      watchers,
				isPersistent: true,
			}
			for _, watcher := range watchers {
				watcher.Watch(event, node, newNode)
			}
			return ch, nil
		})
		log.New().WithField("err", err).Errorln("stop loop with error")
	}()

	return node, nil
}
