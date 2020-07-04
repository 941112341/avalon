package main

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
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
	processor := message.NewMessageServiceProcessor(&Handler{})
	server := avalon.NewServer(processor)
	err := server.Start()
	if err != nil {
		panic(err)
	}
	defer server.Stop()
}
