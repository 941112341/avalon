package zookeeper

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

type Watcher interface {
	Watch(event zk.Event)
}

// unThreadSave
type ZkNode struct {
	data     string
	children map[string]*ZkNode
	path     string
	parent   *ZkNode
}

func (n *ZkNode) Copy() *ZkNode {
	return &ZkNode{
		data:     n.data,
		children: n.children,
		path:     n.path,
		parent:   n.parent,
	}
}

func (n *ZkNode) Save(client *ZkClient, flag int32) error {
	ok, stat, err := client.Conn.Exists(n.path)
	if !ok && err != nil {
		return errors.WithMessage(err, n.path)
	}
	if !ok {
		_, err = client.Conn.Create(n.path, []byte(n.data), flag, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.New().WithFields(logrus.Fields{
				"err":  err,
				"cfg":  client.cfg,
				"path": n.path,
			}).Warnln("create fail")
		}
	} else {
		stat, err = client.Conn.Set(n.path, []byte(n.data), stat.Version)
		if err != nil {
			log.New().WithFields(logrus.Fields{
				"err":  err,
				"cfg":  client.cfg,
				"path": n.path,
				"data": n.data,
			}).Warnln("set fail")
		}
	}

	for dir, node := range n.children {
		if err := node.Save(client, flag); err != nil {
			return errors.WithMessage(err, dir)
		}
		node.parent = n
	}
	return nil
}

func (n *ZkNode) Delete(client *ZkClient) error {
	ok, stat, err := client.Conn.Exists(n.path)
	if !ok {
		return nil
	}
	if !ok && err != nil {
		return errors.WithMessage(err, n.path)
	}
	err = n.List(client, false)
	if err != nil {
		return errors.WithMessage(err, n.path)
	}
	for dir, node := range n.children {
		if err := node.Delete(client); err != nil {
			return errors.WithMessage(err, dir)
		}
	}

	return client.Conn.Delete(n.path, stat.Version)
}

func (n *ZkNode) Get(client *ZkClient) error {
	data, _, err := client.Conn.Get(n.path)
	if err != nil {
		return errors.WithMessage(err, n.path)
	}
	n.data = string(data)
	n.children = map[string]*ZkNode{}
	return nil
}

func (n *ZkNode) getW(client *ZkClient) (<-chan zk.Event, error) {
	log.New().WithField("path", n.path).Infof("watch data node")
	data, _, ch, err := client.Conn.GetW(n.path)
	if err != nil {
		return nil, errors.WithMessage(err, n.path)
	}
	n.data = string(data)
	n.children = map[string]*ZkNode{}
	return ch, nil
}

func (n *ZkNode) GetWL(client *ZkClient, watchers ...Watcher) error {
	return watchLoop(func() (<-chan zk.Event, error) {
		return n.getW(client)
	}, watchers...)
}

func watchLoop(loopFunc func() (<-chan zk.Event, error), watchers ...Watcher) (err error) {
	ch, err := loopFunc()
	if err != nil {
		return
	}
	go func() {
		defer func() { err = inline.RecoverErr() }()
		for err == nil {
			select {
			case event := <-ch:
				for _, watcher := range watchers {
					watcher.Watch(event)
				}
				ch, err = loopFunc()
			}
		}
		log.New().WithField("err", err).Errorln("watch err")
	}()
	return
}

// forAll = false only once
func (n *ZkNode) List(client *ZkClient, forAll bool) error {
	dirs, _, err := client.Conn.Children(n.path)
	if err != nil {
		return errors.WithMessage(err, n.path)
	}

	newChildren := make(map[string]*ZkNode)
	for _, dir := range dirs {
		if n.children[dir] != nil {
			newChildren[dir] = n.children[dir]
			continue
		}
		childNode := &ZkNode{
			path:   inline.JoinPath(n.path, dir),
			parent: n,
		}
		if err := childNode.Get(client); err != nil {
			return errors.WithMessage(err, childNode.path)
		}
		if forAll {
			if err := childNode.List(client, forAll); err != nil {
				return errors.WithMessage(err, childNode.path)
			}
		}
		newChildren[dir] = childNode
	}
	n.children = newChildren
	return nil
}

func (n *ZkNode) ListWL(client *ZkClient, forAll bool, watches ...Watcher) error {
	return watchLoop(func() (<-chan zk.Event, error) {
		log.New().WithField("path", n.path).Infof("watch children")
		dirs, _, ch, err := client.Conn.ChildrenW(n.path)
		if err != nil {
			return nil, err
		}
		newChildren := make(map[string]*ZkNode)
		for _, dir := range dirs {
			if n.children[dir] != nil {
				continue
			}
			childNode := &ZkNode{
				path:   inline.JoinPath(n.path, dir),
				parent: n,
			}
			if err := childNode.GetWL(client); err != nil {
				return nil, errors.WithMessage(err, childNode.path)
			}
			if forAll {
				if err := childNode.ListWL(client, forAll, watches...); err != nil {
					return nil, errors.WithMessage(err, childNode.path)
				}
			}
			newChildren[dir] = childNode
		}
		n.children = newChildren
		return ch, nil
	}, watches...)
}

type ZkNodeBuilder struct {
	node *ZkNode
}

func (b *ZkNodeBuilder) Data(data string) *ZkNodeBuilder {
	b.node.data = data
	return b
}

func (b *ZkNodeBuilder) Children(dir string, node *ZkNode) *ZkNodeBuilder {
	if b.node.children == nil {
		b.node.children = map[string]*ZkNode{}
	}
	b.node.children[dir] = node
	return b
}

func (b *ZkNodeBuilder) Build() *ZkNode {
	return b.node
}

func NewZkNodeBuilder(path string) *ZkNodeBuilder {
	return &ZkNodeBuilder{node: &ZkNode{path: path}}
}

func NewZkNodeBuilderWithNode(node *ZkNode) *ZkNodeBuilder {
	return &ZkNodeBuilder{node: node.Copy()}
}
