package client

import (
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/idl/message"
)

var (
	DefaultMessageClient message.MessageService
)

func init() {
	DefaultMessageClient = &message.MessageServiceClient{C: avalon.NewClient()}
}
