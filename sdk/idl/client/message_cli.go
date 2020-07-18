package client

import (
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/idl/message/base"
	"os"
)

var (
	DefaultMessageClient base.MessageService
)

func init() {
	os.Setenv("base", "../base.yaml")
	client, err := avalon.NewClient()
	if err != nil {
		//panic(err)
	}
	DefaultMessageClient = base.NewMessageServiceClient(client)
}
