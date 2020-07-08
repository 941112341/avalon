package client

import (
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/idl/message"
)

var (
	DefaultMessageClient message.MessageService
)

func init() {
	cfg := &avalon.ClientConfig{}
	err := config.Read(cfg, "../sdk/config/config.client.yaml")
	if err != nil {
		panic(err)
	}
	DefaultMessageClient = &message.MessageServiceClient{C: avalon.NewClientWithConfig(cfg)}
}
