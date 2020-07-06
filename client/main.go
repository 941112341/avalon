package main

import (
	"context"
	"github.com/941112341/avalon/sdk/idl/client"
	"github.com/941112341/avalon/sdk/idl/message"
	"github.com/941112341/avalon/sdk/log"
)

func main() {
	resp, err := client.DefaultMessageClient.MessageDispatcher(context.Background(), &message.MessageRequest{
		Header:     nil,
		Body:       nil,
		MethodName: "",
		URL:        "xx",
	})
	if err != nil {
		panic(err)
	}
	log.New().Info(string(resp.Body))
}
