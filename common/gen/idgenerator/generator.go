
package idgenerator

import (
	"context"
	"github.com/941112341/avalon/common/gen/base"
	"github.com/941112341/avalon/sdk/avalon"
)

/*
	1.0.0
*/
type Handler struct {
	handler    IDGenerator
	cfg        avalon.Config
	middleware []avalon.Middleware
}


func Run(service IDGenerator, middleware ...avalon.Middleware) error {
	server, err := avalon.NewServer(middleware...)
	if err != nil {
		return err
	}
	handler := &Handler{
		handler:    service,
		cfg:        server.Cfg,
		middleware: append(server.Middleware, middleware...),
	}
	err = server.Register(NewIDGeneratorProcessor(handler))
	if err != nil {
		return err
	}
	return server.Start()
}



	
func (h *Handler) GenIDs(ctx context.Context, request *IDRequest) (r *IDResponse, err error) {
	defer func() {
		if iErr, ok := recover().(error); ok {
			err = iErr
		}
	}()
	var call avalon.Endpoint = func(ctx context.Context, method string, _, _ interface{}) error {
		r, err = h.handler.GenIDs(ctx, request)
		return err
	}

	for _, middleware := range h.middleware {
		call = middleware(h.cfg, call)
	}
	err = call(ctx, "GenIDs", request, r)
	if err != nil {
		aErr, ok := err.(*avalon.AErr)
		if ok {
			r = &IDResponse{BaseResp: &base.BaseResp{
				Code:    aErr.Code,
				Message: aErr.Error(),
			}}
		} else {
			r = &IDResponse{BaseResp: &base.BaseResp{
				Code:    avalon.UnknownErr,
				Message: err.Error(),
			}}
		}
	}
	if r == nil {
		r = &IDResponse{BaseResp: &base.BaseResp{
			Code: avalon.UnknownErr,
		}}
	}
	if r.BaseResp == nil {
		r.BaseResp = &base.BaseResp{}
	}
	return r, nil
}
