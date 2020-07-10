package zookeeper

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

var ZkClientInstance *ZkClient
var once sync.Once

func GetZkClientInstance(cfg *ZkConfig) (*ZkClient, error) {
	var err error
	once.Do(func() {
		ZkClientInstance, err = NewClient(cfg)
	})
	return ZkClientInstance, err
}

func NewClient(cfg *ZkConfig) (*ZkClient, error) {
	conn, eventChan, err := zk.Connect(cfg.HostPorts, cfg.SessionTimeout*time.Second)
	if err != nil {
		return nil, errors.WithMessage(err, inline.ToJsonString(cfg.HostPorts))
	}
	return &ZkClient{Conn: conn, EventChan: eventChan, cfg: cfg}, nil
}

type ZkClient struct {
	Conn      *zk.Conn
	EventChan <-chan zk.Event
	cfg       *ZkConfig
}
