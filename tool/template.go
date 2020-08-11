package tool

type FileTemplate struct {
	Package     string
	Version     string
	IDLName     string
	ServiceName string
}

const generateTemplate = `
package {{.Package}}

/*
	{{.Version}}
*/
import (
    "context"
    "github.com/941112341/avalon/sdk/avalon/server"
)

type AvalonHandler struct {

    advances []server.Advance
    handler  {{.ServiceName}}
}

func NewAvalonHandler(handler {{.ServiceName}}) {{.ServiceName}} {
    return &AvalonHandler{
        advances: server.DefaultAdvance(),
        handler:  handler,
    }
}

`

type MethodTemplate struct {
	MethodName string
	Request    string
	Response   string
}

const methodTemplate = `
func (a *AvalonHandler) {{.MethodName}}(ctx context.Context, request *{{.Request}}) (r *{{.Response}}, err error) {
    var call server.Call = func(ctx context.Context, request interface{}) (interface{}, error) {
        return a.handler.{{.MethodName}}(ctx, request.(*{{.Request}}))
    }
    for _, advance := range a.advances {
        call = advance(call)
    }
    resp, err := call(ctx, request)
    if err != nil {
        return nil, err
    }
    return resp.(*{{.Response}}), nil
}
`
