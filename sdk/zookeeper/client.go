package zookeeper

import (
	"github.com/941112341/avalon/sdk/collect"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type EventStatus int

const (
	Success EventStatus = iota
	GetFail
	ChildrenFail
)

type zkClient struct {
	Conn      *zk.Conn
	EventChan <-chan zk.Event
}

func NewClient(hostPorts []string, sessionTimeout time.Duration) (*ZkClient, error) {
	conn, eventChan, err := zk.Connect(hostPorts, sessionTimeout)
	if err != nil {
		return nil, errors.WithMessage(err, inline.JsonString(hostPorts))
	}
	cli := zkClient{Conn: conn, EventChan: eventChan}
	return &ZkClient{zkClient: cli}, nil
}

func (zkClient *zkClient) Close() {
	zkClient.Conn.Close()
}

type ZkClient struct {
	zkClient

	isClose bool
}

type Listener func(event Event)

type Event struct {
	Path    string
	OldData string
	Data    string
	Type    zk.EventType
	Stat    *zk.Stat
	Err     error
	Status  EventStatus
}

func (c *ZkClient) Close() {
	c.zkClient.Close()
	c.isClose = true
}

func (c *ZkClient) WatchTree(path string, listener Listener) {
	c.watchTree(path, listener, collect.NewSyncMap())
}

// 监听节点 todo 重构
func (c *ZkClient) watchTree(path string, listener Listener, dataMap *collect.SyncMap) {
	data, stat, ch, err := c.Conn.GetW(path)
	if err != nil {
		listener(Event{
			Path:   path,
			Err:    err,
			Status: GetFail,
		})
		return
	}
	dataStr := string(data)
	listener(Event{
		Path: path,
		Data: dataStr,
		Type: zk.EventNodeCreated,
		Stat: stat,
	})
	dataMap.Put(path, dataStr)

	child, _, ch, err := c.Conn.ChildrenW(path)
	if err != nil {
		listener(Event{
			Path:   path,
			Status: ChildrenFail,
			Err:    err,
		})
		return
	}
	for _, key := range child {
		subPath := path + "/" + key
		c.watchTree(subPath, listener, dataMap)
	}

	go func() {
		for !c.isClose {
			select {
			case event := <-ch:
				path := event.Path
				log.New().Debugf("server receive a event %s \n", inline.JsonString(event))
				switch event.Type {
				//case zk.EventNodeCreated: 不会触发，只用由exist触发
				//	log.New().Warnf("create event %s\n", inline.JsonString(event))
				case zk.EventNodeDataChanged:
					data, stat, ch, err = c.Conn.GetW(path)
					dataStr := string(data)
					if dataStr != "" {
						listener(Event{
							Path:    path,
							Data:    dataStr,
							OldData: dataMap.GetString(path),
							Type:    event.Type,
							Stat:    stat,
							Err:     err,
						})
					}
					dataMap.Put(path, dataStr)
				case zk.EventNodeDeleted:
					data, stat, ch, err = c.Conn.GetW(path)
					listener(Event{
						Path:    path,
						OldData: dataMap.GetString(path),
						Type:    event.Type,
						Stat:    stat,
						Err:     err,
					})
					dataMap.Delete(path)
				case zk.EventNodeChildrenChanged:
					child, stat, ch, err = c.Conn.ChildrenW(path)
					for _, key := range child {
						subPath := path + "/" + key
						if dataMap.Contains(subPath) {
							continue
						}
						c.watchTree(subPath, listener, dataMap)
					}
				}
			}
		}
	}()
}

func (c *ZkClient) ListenerTree(path string, maps *collect.SyncMap) {
	if maps == nil {
		return
	}

	c.WatchTree(path, func(event Event) {
		if event.Status > Success {
			log.New().WithField("event", event).Errorln("receive zk err event")
		}
		switch event.Type {
		case zk.EventNodeCreated, zk.EventNodeDataChanged:
			maps.Put(event.Path, event.Data)
		case zk.EventNodeDeleted:
			maps.Delete(event.Path)
		}
	})
}
