package invoker

import (
	"github.com/941112341/avalon/sdk/generic"
	"github.com/apache/thrift/lib/go/thrift"
)

type ArgsParser interface {
	Parse() (Args, error)
}

func NewParser(ctx generic.ThriftContext, model generic.ThriftFieldModel) ArgsParser {
	switch model.Type {
	case thrift.LIST:
		return NewListParser(ctx, model)
	case thrift.MAP:
		return NewMapParser(ctx, model)
	case thrift.STRUCT:
		return NewStructParser(ctx, model)
	default:
		return NewBaseParser(model)
	}
}