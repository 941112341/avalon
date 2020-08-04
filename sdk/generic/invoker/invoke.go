package invoker

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
)

type Invoker interface {
	Invoke(ctx context.Context, client thrift.TClient, data interface{}) (interface{}, error)
}

type IInvoker struct {
	Method   string
	Request  Args
	Response Args
}

func (i *IInvoker) Invoke(ctx context.Context, client thrift.TClient, data interface{}) (interface{}, error) {
	err := i.Request.BindValue(data)
	if err != nil {
		return nil, err
	}
	err = client.Call(ctx, i.Method, i.Request, i.Response)
	if err != nil {
		return nil, err
	}
	return i.Response.Data(), nil
}

func NewInvoker(method string, request, response Args) Invoker {
	return &IInvoker{
		Method:   method,
		Request:  request,
		Response: response,
	}
}

type Args interface {
	SubArgs
	BindValue(any interface{}) error
	Data() interface{}
	JSONPath() string
	Index() int16
	ThriftName() string
	IsSkip() bool
}

type SubArgs interface {
	thrift.TStruct
	GetType() thrift.TType
}

func CreateInvoker(ctx generic.ThriftContext, base, service, method string) (Invoker, error) {
	model, ok := ctx.GetMethod(base, service, method)
	if !ok {
		return nil, inline.Error("base %s, service %s, method %s not found", base, service, method)
	}
	request := StructArgs{
		TypeName: fmt.Sprintf("%s_arg", method),
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				parser := NewParser(ctx, generic.ThriftFieldModel{
					Base:           base,
					FieldName:      "request",
					Idx:            1,
					Type:           thrift.STRUCT,
					Optional:       false,
					StructTypeName: model.Request,
				})
				arg, err := parser.Parse()
				if err != nil {
					inline.Errorln("parse err " + err.Error())
					return nil
				}
				return []Args{arg}
			},
		},
		ID: 0,
	}

	response := StructArgs{
		TypeName: fmt.Sprintf("%s_result", method),
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				parser := NewParser(ctx, generic.ThriftFieldModel{
					Base:           base,
					FieldName:      "success",
					Idx:            0,
					Type:           thrift.STRUCT,
					Optional:       false,
					StructTypeName: model.Response,
				})
				arg, err := parser.Parse()
				if err != nil {
					inline.Errorln("parse err " + err.Error())
					return nil
				}
				return []Args{arg}
			},
		},
		ID: 0,
	}

	return NewInvoker(method, &request, &response), nil
}
