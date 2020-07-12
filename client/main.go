package main

import (
	"context"
	"github.com/941112341/avalon/sdk/idl/client"
	"github.com/941112341/avalon/sdk/idl/message"
	"github.com/941112341/avalon/sdk/log"
	"time"
)

func main() {
	for i := 0; i < 20; i++ {
		go c()
	}

	time.Sleep(20 * time.Second)
}

func c() {

	resp, err := client.DefaultMessageClient.MessageDispatcher(context.Background(), &message.MessageRequest{
		Header:     nil,
		Body:       nil,
		MethodName: "",
		URL:        "xx",
	})
	if err != nil {
		log.New().Errorln(err)
	} else {
		log.New().Infoln(string(resp.Body))
	}
}
