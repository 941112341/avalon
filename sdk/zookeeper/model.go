package zookeeper

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

type Watcher func(event zk.Event)

// unThreadSave
type ZkNode struct {
	data     string
	children map[string]*ZkNode
	path     string
	parent   *ZkNode
}

func (n *ZkNode) GetNode(path string) (*ZkNode, error) {
	if !strings.HasPrefix(path, n.path) {
		return nil, errors.New(path + " data not found")
	}
	if path == n.path {
		return n, nil
	}

	arr := strings.SplitN(path, "/", 3)
	if len(arr) < 3 {
		return nil, errors.New(path + " not invalid")
	}
	name := arr[1]
	child, ok := n.children[name]
	if !ok {
		return nil, errors.New(path + " data not found")
	}
	return child.GetNode(arr[2])
}

// parent only has one child
func (n *ZkNode) GetParent() *ZkNode {
	if n.parent != nil {
		return n.parent
	}
	idx := strings.LastIndex(n.path, "/")
	if idx < 1 {
		return nil
	}
	key := n.path[idx+1:]
	return &ZkNode{path: n.path[:idx], children: map[string]*ZkNode{
		key: n,
	}}
}

func (n *ZkNode) GetData() string {
	return n.data
}

func (n *ZkNode) GetChildrenKey() []string {
	keys := make([]string, 0)
	for s, _ := range n.children {
		keys = append(keys, s)
	}
	return keys
}

func (n *ZkNode) GetChildrenMap(forAll bool) map[string]string {
	m := make(map[string]string)
	for _, node := range n.children {
		if forAll {
			subMap := node.GetChildrenMap(forAll)
			for subKey, subData := range subMap {
				m[subKey] = subData
			}
		}
		m[node.path] = node.data
	}
	return m
}

func (n *ZkNode) Copy() *ZkNode {
	return &ZkNode{
		data:     n.data,
		children: n.children,
		path:     n.path,
		parent:   n.parent,
	}
}

func (n *ZkNode) Exist(client *ZkClient) (bool, *zk.Stat, error) {
	return client.Conn.Exists(n.path)
}

func (n *ZkNode) Save(client *ZkClient, flag int32) error {
	ok, stat, err := n.Exist(client)
	if err != nil {
		return errors.Wrap(err, "exist")
	}
	if !ok {
		_, err = client.Conn.Create(n.path, []byte(n.data), flag, zk.WorldACL(zk.PermAll))
		if err != nil {
			if strings.Contains(err.Error(), "zk: node does not exist") {
				return n.GetParent().Save(client, 0)
			}
			return errors.Wrap(err, "create node fail")
		}
	} else if n.data != "" {
		stat, err = client.Conn.Set(n.path, []byte(n.data), stat.Version)
		return errors.Wrap(err, "set value fail")
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
					watcher(event)
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
			return nil, errors.Wrap(err, "ChildrenW")
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
				return nil, errors.Wrap(err, childNode.path)
			}
			if forAll {
				if err := childNode.ListWL(client, forAll, watches...); err != nil {
					return nil, errors.Wrap(err, childNode.path)
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
