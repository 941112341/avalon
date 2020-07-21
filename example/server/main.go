package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/941112341/avalon/sdk/idl/message/base"
	"github.com/941112341/avalon/sdk/inline"
	"os"
)

type Handler struct {
}

func (Handler Handler) MessageDispatcher(ctx context.Context, r *base.MessageRequest) (*base.MessageResponse, error) {
	inline.Infoln("message dispatcher", inline.NewPairs("req", inline.ToJsonString(r))...)

	return nil, errors.New("test err")
}

func main() {
	_ = os.Setenv("base", "../base.yaml")

	err := base.Run(&Handler{})
	fmt.Println(inline.VString(err))
}
