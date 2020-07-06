package zookeeper

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

func TestListener(t *testing.T) {
	cli, err := NewClient([]string{"localhost:2181"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for {
		data, stat, eventChan, err := cli.Conn.ChildrenW("/host")
		fmt.Println(data, stat, err)
		select {
		case event := <-eventChan:
			fmt.Printf("event %s\n", inline.JsonString(event))
		}
	}
	/*for true {
		_, _, ch, _ := cli.Conn.GetW("/host/test")
		select {
		case ev := <-ch:
			fmt.Printf("event %s \n", inline.JsonString(ev))
		}
	}*/

}

func TestWrite(t *testing.T) {
	cli, err := NewClient([]string{"localhost:2181"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	data, err := cli.Conn.Create("/host/job", []byte("localhost:8888"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	fmt.Println(data, err)
	time.Sleep(10 * time.Second)
}

func TestClient(t *testing.T) {
	cli, err := NewClient([]string{"localhost:2181"}, time.Minute)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	cli.WatchTree("/host", func(event Event) {
		log.New().Infof("receive a event: %s", inline.JsonString(event))
	})

	time.Sleep(10 * time.Minute)
}
