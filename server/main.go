package main

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/config"
	"github.com/941112341/avalon/sdk/idl/message"
)

type Handler struct {
}

func (Handler Handler) MessageDispatcher(ctx context.Context, r *message.MessageRequest) (*message.MessageResponse, error) {
	return &message.MessageResponse{
		Body: []byte("hello world"),
	}, nil
}

func main() {
	cfg := &avalon.ServerConfig{}
	config.Read(cfg, "../sdk/config/config.server.yaml")
	processor := message.NewMessageServiceProcessor(&Handler{})
	server := avalon.NewThriftServer(processor)

	bootstrap := avalon.NewBootstrap(server)
	err := bootstrap.Start(cfg)
	if err != nil {
		panic(err)
	}
	defer bootstrap.Stop(cfg)
}
