package invoker

import (
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"strings"
)


type FieldModelParser struct {
	ctx generic.ThriftContext

	fieldModel *generic.ThriftFieldModel
}

func (f *FieldModelParser) Parse() (Args, error) {
	fieldModel := f.fieldModel
	var arg Args
	switch f.fieldModel.Type {
	case thrift.STRUCT:
		/*arg = &StructArgs{
			TypeName:   inline.Ucfirst(fieldModel.FieldName),
			jsonPath:   fieldModel.FieldName,
			thriftName: fieldModel.FieldName,
			LazyFields: LazyCacheArgs{},
			optional:   fieldModel.Optional,
			any:        nil,
			ID:         fieldModel.Idx,
		}*/

		ss := strings.Split(f.fieldModel.StructTypeName, ".")
		var base, structName string
		switch len(ss) {
		case 1:
			base = fieldModel.Base
			structName = ss[1]
		case 2:
			base, structName = ss[0], ss[1]
		}

		reference, err := f.ctx.Ptr(base, structName)
		if err != nil {
			return nil, inline.PrependErrorFmt(err, "ptr base %s, structName %s", base, structName)
		}
		parser := NewStructModelParser(f.ctx, reference)
		return parser.Parse()
	//case thrift.LIST
	}

	return arg, nil
}



