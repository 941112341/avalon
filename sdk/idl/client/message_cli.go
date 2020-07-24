package client

import (
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/idl/message/base"
)

var (
	DefaultMessageClient base.MessageService
)

func init() {
	client, err := avalon.NewClient("avalon.test.idgenerator")
	if err != nil {
		//panic(err)
	}
	DefaultMessageClient = base.NewMessageServiceClient(client)
}
