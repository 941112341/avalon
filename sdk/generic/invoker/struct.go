package invoker

import (
	"fmt"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/json-iterator/go"
)

type StructArgs struct {
	TypeName   string
	jsonPath   string
	thriftName string
	LazyFields LazyCacheArgs `json:",omitempty"` // struct
	optional   bool
	ID         int16
	skip       bool

}

func (s *StructArgs) Index() int16 {
	return s.ID
}

func (s *StructArgs) ThriftName() string {
	return s.thriftName
}

func (s *StructArgs) JSONPath() string {
	return s.jsonPath
}

func (s *StructArgs) Data() interface{} {
	m := make(map[string]interface{})
	for _, args := range s.LazyFields.fields() {
		m[args.JSONPath()] = args.Data()
	}
	return m
}

func (s *StructArgs) findSubArg(idx int16) (Args, bool) {
	for _, args := range s.LazyFields.fields() {
		if args.Index() == idx {
			return args, true
		}
	}
	return nil, false
}

func (s *StructArgs) IsSkip() bool {
	return s.optional && s.skip
}

func (s *StructArgs) Write(p thrift.TProtocol) error {
	if s.IsSkip() {
		return nil
	}

	if err := p.WriteStructBegin(s.TypeName); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", s, err)
	}
	fields := s.LazyFields.fields()
	for _, arg := range fields {
		if arg.IsSkip() {
			continue
		}

		if err := p.WriteFieldBegin(arg.ThriftName(), arg.GetType(), arg.Index()); err != nil {
			return fmt.Errorf("%T write field begin error %d:groupName: %s", arg, arg.Index(), err)
		}
	}
	if err := p.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := p.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (s *StructArgs) Read(p thrift.TProtocol) error {
	if _, err := p.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", s), err)
	}
	for {
		_, fieldTypeId, fieldId, err := p.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", s, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		arg, ok := s.findSubArg(fieldId)
		if !ok {
			if err := p.Skip(fieldTypeId); err != nil {
				return err
			}
			continue
		}
		if err = arg.Read(p); err != nil {
			return err
		}
		if err := p.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := p.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", s), err)
	}
	return nil
}

func (s *StructArgs) GetType() thrift.TType {
	return thrift.STRUCT
}

func (s *StructArgs) BindValue(o interface{}) error {
	switch any := o.(type) {
	case jsoniter.Any:
		if err := any.LastError(); err != nil {
			s.skip = true
			inline.WithFields("err", err).Errorln("bind err")
			return nil
		}
		args := s.LazyFields.fields()
		for _, arg := range args {
			path := arg.JSONPath()
			subAny := any.Get(path)
			if err := subAny.LastError(); err != nil {
				inline.WithFields("jsonpath", path).Infoln("skip json path")
				continue
			}
			if err := arg.BindValue(subAny); err != nil {
				return inline.PrependErrorFmt(err, "bind value err %s", path)
			}
		}
	default:
		return inline.NewError(ErrUnsupportType, "unknown type of o %+v", o)
	}

	return nil
}

type LazyCacheArgs struct {
	caches   []Args
	LazyArgs LazyArgs
}

func (l *LazyCacheArgs) fields() []Args {
	if l.caches == nil {
		args := l.LazyArgs()
		l.caches = args
	}
	return l.caches
}

type LazyArgs func() []Args

type StructParser struct {
	ctx generic.ThriftContext

	model generic.ThriftFieldModel
}

func (s *StructParser) Parse() (Args, error) {
	model := s.model
	ctx := s.ctx
	
	structModel, err := ctx.Ptr(model.Base, model.StructTypeName)
	if err != nil {return nil, inline.PrependErrorFmt(err, "struct type name %s", model.StructTypeName)}

	return &StructArgs{
		TypeName:   structModel.StructName,
		jsonPath:   model.FieldName,
		thriftName: model.FieldName,
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				args := make([]Args, 0)
				for _, fieldModel := range structModel.FieldMap {
					parser := NewParser(ctx, *fieldModel)
					subArgs, err := parser.Parse()
					if err != nil {
						inline.WithFields("parser", parser).Errorln("parse fail err %s", err)
						continue
					}
					args = append(args, subArgs)
				}
				return args
			},
		},
		optional:   model.Optional,
		ID:         model.Idx,
		skip:       false,
	}, nil
}

func NewStructParser(ctx generic.ThriftContext, model generic.ThriftFieldModel) *StructParser {
	return &StructParser{
		ctx:         ctx,
		model: 		model,
	}
}


