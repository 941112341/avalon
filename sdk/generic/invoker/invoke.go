package invoker

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
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
	thrift.TStruct
	BindValue(any interface{}) error
	Data() interface{}
	IsSkip() bool
	Meta() ArgsMeta
}

type ArgsMeta interface {
	TypeName() string
	ID() int16
	JsonPath() string
	ThriftName() string
	Type() thrift.TType
	Optional() bool
}

type argsMeta struct {
	typeName   string
	id         int16
	optional   bool
	jsonPath   string
	thriftName string
	ttype      thrift.TType
}

func (a *argsMeta) Optional() bool {
	return a.optional
}

func (a *argsMeta) TypeName() string {
	return a.typeName
}

func (a *argsMeta) ID() int16 {
	return a.id
}

func (a *argsMeta) JsonPath() string {
	return a.jsonPath
}

func (a *argsMeta) ThriftName() string {
	return a.thriftName
}

func (a *argsMeta) Type() thrift.TType {
	return a.ttype
}

func FileInvoker(base []string, method string) (Invoker, error) {
	grp, err := generic.NewThriftGroupBase(base)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "base %+v", base)
	}
	return CreateInvoker(grp, "", "", method)
}

func CreateInvoker(ctx generic.ThriftContext, base, service, method string) (Invoker, error) {
	model, ok := ctx.GetMethod(base, service, method)
	if !ok {
		return nil, inline.Error("base %s, service %s, method %s not found", base, service, method)
	}
	request := StructArgs{
		ArgsMeta: &argsMeta{
			typeName: fmt.Sprintf("%s_arg", method),
			id:       0,
		},
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				parser := NewParser(thrift.STRUCT, &FileMetaParser{
					ctx: ctx,
					model: generic.ThriftFieldModel{
						Base:           base,
						FieldName:      "request",
						Idx:            1,
						TType:          thrift.STRUCT,
						OptionalVar:    false,
						StructTypeName: model.Request,
					},
				})
				arg, err := parser.Parse()
				if err != nil {
					inline.Errorln("parse err " + err.Error())
					return nil
				}
				return []Args{arg}
			},
		},
	}

	response := StructArgs{
		ArgsMeta: &argsMeta{
			id:       0,
			typeName: fmt.Sprintf("%s_result", method),
		},
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				parser := NewParser(thrift.STRUCT, &FileMetaParser{
					ctx: ctx,
					model: generic.ThriftFieldModel{
						Base:           base,
						FieldName:      "success",
						Idx:            0,
						TType:          thrift.STRUCT,
						OptionalVar:    false,
						StructTypeName: model.Response,
					},
				})
				arg, err := parser.Parse()
				if err != nil {
					inline.Errorln("parse err " + err.Error())
					return nil
				}
				return []Args{arg}
			},
		},
	}

	return NewInvoker(method, &request, &response), nil
}

func ReflectInvoker(method string, args, result thrift.TStruct) (Invoker, error) {
	argType := reflect.TypeOf(args).Elem()
	request := StructArgs{
		ArgsMeta: &argsMeta{
			typeName: fmt.Sprintf("%s_arg", method),
			id:       0,
		},
		skip: false,
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				parser := NewParser(thrift.STRUCT, &ReflectArgsMetaParser{
					Field: argType.Field(0),
				})
				arg, err := parser.Parse()
				if err != nil {
					inline.Errorln("parse err " + err.Error())
					return nil
				}
				return []Args{arg}
			},
		},
	}

	resultType := reflect.TypeOf(result).Elem()
	response := StructArgs{
		ArgsMeta: &argsMeta{
			id:       0,
			typeName: fmt.Sprintf("%s_result", method),
		},
		skip: false,
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				parser := NewParser(thrift.STRUCT, &ReflectArgsMetaParser{
					Field: resultType.Field(0),
				})
				arg, err := parser.Parse()
				if err != nil {
					inline.Errorln("parse err " + err.Error())
					return nil
				}
				return []Args{arg}
			},
		},
	}

	return NewInvoker(method, &request, &response), nil
}
