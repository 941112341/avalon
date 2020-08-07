package zookeeper

import (
	"fmt"
	"testing"
	"time"
)

func TestModel(t *testing.T) {
	cli, err := NewClient(ZkConfig{
		HostPorts:      []string{"localhost:2181"},
		SessionTimeout: 10 * time.Second,
		Path:           "/host",
	})

	if err != nil {
		panic(err)
	}
	//err = cli.Conn.Delete("/host/example.jiangshihao.test/10.86.124.245:8889", 0)
	d, stat, err := cli.Conn.Get("/host/example.jiangshihao.test/10.86.124.245:8888")
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
	fmt.Println(stat)
}
