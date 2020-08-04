package invoker

import (
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
)

type ArgsParser interface {
	Parse() (Args, error)
}

type ArgsMetaParser interface {
	FieldsParsers() []Args
	ArgsMeta() ArgsMeta
	ElemParsers() ArgsParser
	KVParsers() (ArgsParser, ArgsParser)
}

type FileMetaParser struct {
	ctx   generic.ThriftContext
	model generic.ThriftFieldModel
}

func (f *FileMetaParser) FieldsParsers() []Args {
	ctx, model := f.ctx, f.model
	structModel, err := ctx.Ptr(model.Base, model.StructTypeName)
	if err != nil {
		inline.WithFields("base", model.Base, "structName", model.StructTypeName).Errorln("ptr err %s", err)
		return nil
	}
	args := make([]Args, 0)
	for _, fieldModel := range structModel.FieldMap {
		parser := NewParser(fieldModel.TType, &FileMetaParser{
			ctx:   ctx,
			model: *fieldModel,
		})
		subArgs, err := parser.Parse()
		if err != nil {
			inline.WithFields("parser", parser).Errorln("parse fail err %s", err)
			continue
		}
		args = append(args, subArgs)
	}
	return args
}

func (f *FileMetaParser) ArgsMeta() ArgsMeta {
	return &f.model
}

func (f *FileMetaParser) ElemParsers() ArgsParser {
	elem := f.model.Elem()
	return NewParser(elem.TType, &FileMetaParser{
		ctx:   f.ctx,
		model: *elem,
	})
}

func (f *FileMetaParser) KVParsers() (ArgsParser, ArgsParser) {
	kelem, velem := f.model.KVElem()
	return NewParser(kelem.TType, &FileMetaParser{
			ctx:   f.ctx,
			model: *kelem,
		}), NewParser(velem.TType, &FileMetaParser{
			ctx:   f.ctx,
			model: *velem,
		})
}

func NewParser(ttypes thrift.TType, parser ArgsMetaParser) ArgsParser {
	switch ttypes {
	case thrift.LIST:
		return NewListParser(parser)
	case thrift.MAP:
		return NewMapParser(parser)
	case thrift.STRUCT:
		return NewStructParser(parser)
	default:
		return NewBaseParser(parser)
	}
}

type ReflectArgsMetaParser struct {
	Field reflect.StructField

	model *CommonModel
}

func (r *ReflectArgsMetaParser) FieldsParsers() []Args {
	parentType := inline.Redirect(r.Field.Type)
	args := make([]Args, 0)
	for i := 0; i < parentType.NumField(); i++ {
		field := parentType.Field(i)
		parser := NewParser(generic.Type2thrift(field.Type), &ReflectArgsMetaParser{
			Field: field,
		})
		arg, err := parser.Parse()
		if err != nil {
			inline.WithFields("parser", parser).Errorln("parse reflect fail err %s", err)
			continue
		}
		args = append(args, arg)
	}
	return args
}

func (r *ReflectArgsMetaParser) ArgsMeta() ArgsMeta {
	return r.getModel()
}

func (r *ReflectArgsMetaParser) getModel() *CommonModel {
	if r.model == nil {
		model, err := NewModel(r.Field)
		if err != nil {
			inline.WithFields("field", inline.ToJsonString(r.Field)).Errorln("new model err %s", err)
			return nil
		}
		r.model = model
	}
	return r.model
}

func (r *ReflectArgsMetaParser) ElemParsers() ArgsParser {
	model := r.getModel()
	subModel := model.Elem()
	return NewParser(subModel.ttype, &ReflectArgsMetaParser{model: subModel})
}

func (r *ReflectArgsMetaParser) KVParsers() (ArgsParser, ArgsParser) {
	model := r.getModel()
	k, v := model.KVElem()
	return NewParser(k.ttype, &ReflectArgsMetaParser{model: k}), NewParser(v.ttype, &ReflectArgsMetaParser{model: v})
}
