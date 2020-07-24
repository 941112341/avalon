package base

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
)

/*
	1.0.0
*/
type Handler struct {
	handler    MessageService
	cfg        avalon.Config
	middleware []avalon.Middleware
}

func Run(service MessageService, middleware ...avalon.Middleware) error {
	server, err := avalon.NewServer(middleware...)
	if err != nil {
		return err
	}
	handler := &Handler{
		handler:    service,
		cfg:        server.builder,
		middleware: append(server.Middleware, middleware...),
	}
	err = server.Register(NewMessageServiceProcessor(handler))
	if err != nil {
		return err
	}
	return server.Start()
}

func (h *Handler) MessageDispatcher(ctx context.Context, request *MessageRequest) (r *MessageResponse, err error) {
	defer func() {
		if iErr, ok := recover().(error); ok {
			err = iErr
		}
	}()
	var call avalon.Endpoint = func(ctx context.Context, method string, _, _ interface{}) error {
		r, err = h.handler.MessageDispatcher(ctx, request)
		return err
	}

	for _, middleware := range h.middleware {
		call = middleware(h.cfg, call)
	}
	err = call(ctx, "MessageDispatcher", request, r)
	if err != nil {
		aErr, ok := err.(*avalon.AErr)
		if ok {
			r = &MessageResponse{BaseResp: &BaseResp{
				Code:    aErr.Code,
				Message: aErr.Error(),
			}}
		} else {
			r = &MessageResponse{BaseResp: &BaseResp{
				Code:    avalon.UnknownErr,
				Message: err.Error(),
			}}
		}
	}
	if r == nil {
		r = &MessageResponse{BaseResp: &BaseResp{
			Code: avalon.UnknownErr,
		}}
	}
	if r.BaseResp == nil {
		r.BaseResp = &BaseResp{}
	}
	return r, nil
}
