package zookeeper

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

var cfg = &ZkConfig{
	HostPorts: []string{
		"192.168.0.109:2181",
	},
	SessionTimeout: 100,
	Path:           "/test",
}

func TestZkNode_Save(t *testing.T) {
	zkCli, err := GetZkClientInstance(cfg)
	if err != nil {
		panic(err)
	}

	node1 := NewZkNodeBuilder("/test/jiangshihao").Data("hello world2").Build()
	node2 := NewZkNodeBuilder("/test/jiangshihao2").Data("hello world_update").Build()
	node := NewZkNodeBuilder("/test").Children("jiangshihao", node1).Children("jiangshihao2", node2).Build()

	err = node.Save(zkCli, 0)
	if err != nil {
		panic(err)
	}
}

type LogiWatchers struct {
	Id string
}

func (l LogiWatchers) Watch(event zk.Event) {
	log.New().WithField("id", l.Id).Infof(inline.ToJsonString(event))
}

func TestZkNode_Get(t *testing.T) {
	zkCli, err := GetZkClientInstance(cfg)
	if err != nil {
		panic(err)
	}

	//go func() {
	//	for true {
	//
	//		select {
	//		case event:= <- zkCli.EventChan:
	//			log.New().WithField("event", event).Infof("cli connect receive")
	//		}
	//	}
	//}()

	node := NewZkNodeBuilder("/test").Build()
	err = node.GetWL(zkCli, LogiWatchers{})
	if err != nil {
		panic(err)
	}
	err = node.GetWL(zkCli, LogiWatchers{})
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Minute)
}

func TestZkNode_ListWL(t *testing.T) {
	zkCli, err := GetZkClientInstance(cfg)
	if err != nil {
		panic(err)
	}
	node := NewZkNodeBuilder("/test").Build()
	err = node.ListWL(zkCli, true, LogiWatchers{})
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Minute)
}

func TestChildListen(t *testing.T) {
	zkCli, err := GetZkClientInstance(cfg)
	if err != nil {
		panic(err)
	}
	_, _, _, _ = zkCli.Conn.GetW("/test")
	_, _, ch2, _ := zkCli.Conn.ChildrenW("/test")

	/*go func() {
		select {
		case event := <- ch:
			log.New().WithField("name", "getw").Infoln(inline.ToJsonString(event))
		}
	}()*/

	go func() {
		for {
			select {
			case event := <-ch2:
				log.New().WithField("name", "children").Infoln(inline.ToJsonString(event))
				_, _, ch2, _ = zkCli.Conn.ChildrenW("/test")
			}
		}
	}()

	time.Sleep(100 * time.Second)
}

func TestZkNode_Delete(t *testing.T) {

	zkCli, err := GetZkClientInstance(cfg)
	if err != nil {
		panic(err)
	}

	node := NewZkNodeBuilder("/test").Build()
	err = node.Delete(zkCli)
	if err != nil {
		panic(err)
	}
}
