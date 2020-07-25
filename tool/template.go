package tool

type FileTemplate struct {
	Package     string
	Version     string
	IDLName     string
	ServiceName string
}

const generateTemplate = `
package {{.Package}}

import (
	"context"
	"github.com/941112341/avalon/common/gen/base"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/inline"
)

/*
	{{.Version}}
*/
type Handler struct {
	handler    {{.ServiceName}}
	cfg        avalon.Config
	middleware []avalon.Middleware
}


func Run(psm string, service {{.ServiceName}}, middleware ...avalon.Middleware) error {
	server := avalon.NewServer(psm, middleware...)
	handler := &Handler{
		handler:    service,
		cfg:        server.Config(),
		middleware: append(server.Middleware, middleware...),
	}
	err := server.Register(New{{.ServiceName}}Processor(handler))
	if err != nil {
		return err
	}
	return server.Start()
}

`

type MethodTemplate struct {
	MethodName string
	Request    string
	Response   string
}

const methodTemplate = `
	
func (h *Handler) {{.MethodName}}(ctx context.Context, request *{{.Request}}) (r *{{.Response}}, err error) {
	var call avalon.Endpoint = func(ctx context.Context, method string, _, _ interface{}) (err error) {
		defer func() {
			if iErr, ok := recover().(error); ok {
				inline.WithFields("requestID", avalon.RequestID(ctx), "err", iErr).Errorln("panic !!")
				err = iErr
			}
		}()
		r, err = h.handler.{{.MethodName}}(ctx, request)
		return err
	}

	for _, middleware := range h.middleware {
		call = middleware(h.cfg, call)
	}
	err = call(ctx, "{{.MethodName}}", request, r)
	if err != nil {
		aErr, ok := err.(*avalon.AErr)
		if ok {
			r = &{{.Response}}{BaseResp: &base.BaseResp{
				Code:    aErr.Code,
				Message: aErr.Error(),
			}}
		} else {
			r = &{{.Response}}{BaseResp: &base.BaseResp{
				Code:    avalon.UnknownErr,
				Message: err.Error(),
			}}
		}
	}
	if r == nil {
		r = &{{.Response}}{BaseResp: &base.BaseResp{
			Code: avalon.UnknownErr,
		}}
	}
	if r.BaseResp == nil {
		r.BaseResp = &base.BaseResp{}
	}
	return r, nil
}
`
