package invoker

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/json-iterator/go"
	"reflect"
)

type ListArgs struct {
	ArgsMeta

	SubArgs Args
	data    []interface{}
}

func (l *ListArgs) Meta() ArgsMeta {
	return l.ArgsMeta
}

func (l *ListArgs) Data() interface{} {
	return l.data
}

func (l *ListArgs) Write(p thrift.TProtocol) error {
	subArgs := l.SubArgs
	if err := p.WriteListBegin(subArgs.Meta().Type(), len(l.data)); err != nil {
		return err
	}

	for _, datum := range l.data {
		if err := subArgs.BindValue(datum); err != nil {
			inline.WithFields("datum", datum).Errorln("bind value err %s", err)
			continue
		}

		if subArgs.IsSkip() {
			continue
		}

		if err := subArgs.Write(p); err != nil {
			return inline.PrependErrorFmt(err, "any %+v", l.Data())
		}
	}

	if err := p.WriteListEnd(); err != nil {
		return fmt.Errorf("error writing list end: %s", err)
	}
	return nil
}

func (l *ListArgs) Read(p thrift.TProtocol) error {
	_, size, err := p.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}

	subArg := l.SubArgs
	ifaces := make([]interface{}, 0)
	for i := 0; i < size; i++ {

		if err = subArg.Read(p); err != nil {
			return err
		}
		ifaces = append(ifaces, subArg.Data())
	}
	l.data = ifaces
	if err := p.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (l *ListArgs) GetType() thrift.TType {
	return thrift.LIST
}

func (l *ListArgs) BindValue(o interface{}) error {
	switch any := o.(type) {
	case jsoniter.Any:
		if err := any.LastError(); err != nil {
			if !l.Optional() {
				l.data = []interface{}{}
			}
			return nil
		}
		l.data = []interface{}{}
		for i := 0; i < any.Size(); i++ {
			o := any.Get(i)

			if err := l.SubArgs.BindValue(o); err != nil {
				inline.WithFields("o", o.ToString()).Errorln("bind value %+v", l.SubArgs)
				continue
			} else {
				l.data = append(l.data, l.SubArgs.Data())
			}
		}
	default:
		value := reflect.ValueOf(o)
		if value.Kind() != reflect.Slice {
			return inline.NewError(ErrUnsupportType, "unsupport type %+v", any)
		}
		size := value.Len()
		for i := 0; i < size; i++ {
			o := value.Index(i)
			if err := l.SubArgs.BindValue(o); err != nil {
				inline.WithFields("o", o.Interface()).Errorln("bind value %+v", l.SubArgs)
				continue
			} else {
				l.data = append(l.data, l.SubArgs.Data())
			}
		}
	}

	return nil
}

func (l *ListArgs) IsSkip() bool {
	return l.Optional() && l.data == nil
}

type ListParser struct {
	ArgsMetaParser
}

func (l *ListParser) Parse() (Args, error) {

	subParser := l.ElemParsers()
	subArg, err := subParser.Parse()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "parse err %+v", subParser)
	}
	return &ListArgs{

		SubArgs:  subArg,
		ArgsMeta: l.ArgsMeta(),
	}, nil
}

func NewListParser(parser ArgsMetaParser) *ListParser {
	return &ListParser{parser}
}
