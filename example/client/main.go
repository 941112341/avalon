package main

import (
	"context"
	"github.com/941112341/avalon/sdk/idl/client"
	"github.com/941112341/avalon/sdk/idl/message/base"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"time"
)

func main() {
	for i := 0; i < 3; i++ {
		go c()
	}

	time.Sleep(10 * time.Second)
}

func c() {
	resp, err := client.DefaultMessageClient.MessageDispatcher(context.Background(), &base.MessageRequest{
		Header:     nil,
		Body:       nil,
		MethodName: "",
		URL:        "xx",
	})
	if err != nil {
		log.New().Errorln(inline.VString(err))
	} else {
		log.New().Infoln(inline.ToJsonString(resp))
	}
}
