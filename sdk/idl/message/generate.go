package message

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
)

/*
	this file should use tools generate
*/
type Handler struct {
	handler    MessageService
	cfg        avalon.CallConfig
	middleware []avalon.Middleware
}

func (h *Handler) MessageDispatcher(ctx context.Context, request *MessageRequest) (r *MessageResponse, err error) {
	var call avalon.Endpoint = func(ctx context.Context, method string, _, _ interface{}) error {
		resp, err := h.handler.MessageDispatcher(ctx, request)
		if err != nil {
			return err
		}
		*r = *resp
		return nil
	}

	for _, middleware := range h.middleware {
		call = middleware(h.cfg, call)
	}
	err = call(ctx, "MessageDispatcher", request, r)
	return
}

func Run(service MessageService, middleware ...avalon.Middleware) error {
	server, err := avalon.NewServer(middleware...)
	if err != nil {
		return err
	}
	handler := &Handler{
		handler:    service,
		cfg:        avalon.NewCallConfig(server.Cfg),
		middleware: append(server.Middleware, middleware...),
	}
	err = server.Register(handler)
	if err != nil {
		return err
	}
	return server.Start()
}
