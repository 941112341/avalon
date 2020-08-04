package invoker

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/json-iterator/go"
)

type StructArgs struct {
	ArgsMeta

	skip       bool
	LazyFields LazyCacheArgs `json:",omitempty"` // struct
}

func (s *StructArgs) Meta() ArgsMeta {
	return s.ArgsMeta
}

func (s *StructArgs) Data() interface{} {
	m := make(map[string]interface{})
	for _, args := range s.LazyFields.fields() {
		m[args.Meta().JsonPath()] = args.Data()
	}
	return m
}

func (s *StructArgs) findSubArg(idx int16) (Args, bool) {
	for _, args := range s.LazyFields.fields() {
		if args.Meta().ID() == idx {
			return args, true
		}
	}
	return nil, false
}

func (s *StructArgs) IsSkip() bool {
	return s.Optional() && s.skip
}

func (s *StructArgs) Write(p thrift.TProtocol) error {

	if err := p.WriteStructBegin(s.TypeName()); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", s, err)
	}

	fields := s.LazyFields.fields()
	if s.skip { // 空循环
		fields = nil
	}
	for _, arg := range fields {
		if arg.IsSkip() {
			continue
		}

		if err := p.WriteFieldBegin(arg.Meta().ThriftName(), arg.Meta().Type(), arg.Meta().ID()); err != nil {
			return fmt.Errorf("%T write field begin error %d:groupName: %s", arg, arg.Meta().ID(), err)
		}

		if err := arg.Write(p); err != nil {
			return inline.PrependErrorFmt(err, "write arg %+v", arg)
		}
		if err := p.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error %d:groupName: %s", arg, arg.Meta().ID(), err)
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
			if !s.Optional() {
				inline.WithFields("err", err).Errorln("bind err")
			}
			s.skip = true
			return nil
		}
		args := s.LazyFields.fields()
		for _, arg := range args {
			path := arg.Meta().JsonPath()
			subAny := any.Get(path)
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
	ArgsMetaParser
}

func (s *StructParser) Parse() (Args, error) {

	return &StructArgs{
		LazyFields: LazyCacheArgs{
			LazyArgs: func() []Args {
				return s.FieldsParsers()
			},
		},
		ArgsMeta: s.ArgsMeta(),
		skip:     false,
	}, nil
}

func NewStructParser(parser ArgsMetaParser) *StructParser {
	return &StructParser{
		parser,
	}
}
