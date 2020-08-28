package main

import (
	"fmt"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/avalon/client"
	"testing"
)

func TestIit(t *testing.T) {
	bean := &client.ZkIPDiscover{}
	avalon.NewBean(bean).Initial()
	fmt.Println(bean)
}
