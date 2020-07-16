package main

import (
	"context"
	"github.com/941112341/avalon/sdk/idl/client"
	"github.com/941112341/avalon/sdk/idl/message"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"time"
)

func main() {
	for i := 0; i < 150; i++ {
		go c()
	}

	time.Sleep(50 * time.Second)
}

func c() {

	resp, err := client.DefaultMessageClient.MessageDispatcher(context.Background(), &message.MessageRequest{
		Header:     nil,
		Body:       nil,
		MethodName: "",
		URL:        "xx",
	})
	if err != nil {
		log.New().Errorln(inline.VString(err))
	} else {
		log.New().Infoln(string(resp.Body))
	}
}
