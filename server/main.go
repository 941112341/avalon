package main

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/idl/message"
	"github.com/941112341/avalon/sdk/inline"
	"os"
)

type Handler struct {
}

func (Handler Handler) MessageDispatcher(ctx context.Context, r *message.MessageRequest) (*message.MessageResponse, error) {
	return &message.MessageResponse{
		Body: []byte("hello world"),
	}, nil
}

func main() {
	_ = os.Setenv("base", "../base.yaml")

	err := message.Run(&Handler{})
	fmt.Println(inline.VString(err))
}
