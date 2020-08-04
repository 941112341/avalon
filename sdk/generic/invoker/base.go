package invoker

import (
	"fmt"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/json-iterator/go"
	"reflect"
)

const (
	ErrUnsupportType inline.AvalonErrorCode = iota + 500
)

type BaseArgs struct {
	data interface{}

	optional   bool
	ttype      thrift.TType
	jsonPath   string
	id         int16
	thriftName string
}

func (b *BaseArgs) Data() interface{} {
	return b.data
}

func (b *BaseArgs) JSONPath() string {
	return b.jsonPath
}

func (b *BaseArgs) Write(p thrift.TProtocol) error {
	switch b.GetType() {
	case thrift.STRING:
		str, ok := b.data.(string)
		if !ok {
			return fmt.Errorf("b.data %+v is not string", b.data)
		}
		if err := p.WriteString(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.Index(), err)
		}
	case thrift.DOUBLE:
		str, ok := b.data.(float64)
		if !ok {
			return fmt.Errorf("b.data %+v is not double", b.data)
		}
		if err := p.WriteDouble(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.Index(), err)
		}
	case thrift.BOOL:
		v, ok := b.data.(bool)
		if !ok {
			return fmt.Errorf("b.data %+v is not bool", b.data)
		}
		if err := p.WriteBool(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.Index(), err)
		}
	case thrift.BYTE:
		v, ok := b.data.(int8)
		if !ok {
			return fmt.Errorf("b.data %+v is not  byte", b.data)
		}
		if err := p.WriteByte(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.Index(), err)
		}
	case thrift.I16:
		v, ok := b.data.(int16)
		if !ok {
			return fmt.Errorf("b.data %+v is not i16", b.data)
		}
		if err := p.WriteI16(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.Index(), err)
		}
	case thrift.I32:
		v, ok := b.data.(int32)
		if !ok {
			return fmt.Errorf("b.data %+v is not i32", b.data)
		}
		if err := p.WriteI32(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.Index(), err)
		}
	case thrift.I64:
		v, ok := b.data.(int64)
		if !ok {
			return fmt.Errorf("b.data %+v is not i64", b.data)
		}
		if err := p.WriteI64(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.Index(), err)
		}
	default:
		return inline.NewError(ErrUnsupportType, "write unknown type %s", b.GetType())
	}
	return nil
}

func (b *BaseArgs) Read(p thrift.TProtocol) error {
	switch b.GetType() {
	case thrift.STOP: // do nothing
	case thrift.STRING:
		if v, err := p.ReadString(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			b.data = v
		}
	case thrift.BOOL:
		if v, err := p.ReadBool(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			b.data = v
		}
	case thrift.BYTE:
		if v, err := p.ReadByte(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			b.data = v
		}
	case thrift.I16:
		if v, err := p.ReadI16(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			b.data = v
		}
	case thrift.I32:
		if v, err := p.ReadI32(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			b.data = v
		}
	case thrift.I64:
		if v, err := p.ReadI64(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			b.data = v
		}
	case thrift.DOUBLE:
		if v, err := p.ReadDouble(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			b.data = v
		}
	default:
		return inline.NewError(ErrUnsupportType, "read unknown type %s", b.GetType())
	}
	return nil
}

func (b *BaseArgs) GetType() thrift.TType {
	return b.ttype
}

func (b *BaseArgs) BindValue(o interface{}) error {
	switch any := o.(type) {
	case jsoniter.Any:

		if err := any.LastError(); err != nil {
			if !b.optional {
				return inline.Error("bind data is nil %+v", b)
			}
			return nil
		}
		var data interface{}
		switch b.GetType() {
		case thrift.BOOL:
			data = any.ToBool()
		case thrift.I08:
			data = int8(any.ToInt())
		case thrift.I16:
			data = int16(any.ToInt())
		case thrift.I32:
			data = any.ToInt32()
		case thrift.I64:
			data = any.ToInt64()
		case thrift.DOUBLE:
			data = any.ToFloat64()
		case thrift.STRING:
			data = any.ToString()
		default:
			inline.WithFields("any", any).Warnln("unknown field type")
			data = any.GetInterface()
		}
		b.data = data
	case int8:
		if b.GetType() == thrift.I08 {
			b.data = any
		} else if !b.optional {
			b.data = int8(0)
		}
	case int16:
		if b.GetType() == thrift.I16 {
			b.data = any
		} else if !b.optional {
			b.data = int16(0)
		}
	case int32:
		if b.GetType() == thrift.I32 {
			b.data = any
		} else if !b.optional {
			b.data = int32(0)
		}
	case int64:
		if b.GetType() == thrift.I64 {
			b.data = any
		} else if !b.optional {
			b.data = int64(0)
		}
	case bool:
		if b.GetType() == thrift.BOOL {
			b.data = any
		} else if !b.optional {
			b.data = false
		}
	case float32:
		if b.GetType() == thrift.DOUBLE {
			b.data = float64(any)
		} else if !b.optional {
			b.data = float64(0)
		}
	case float64:
		if b.GetType() == thrift.DOUBLE {
			b.data = any
		} else if !b.optional {
			b.data = float64(0)
		}
	case string:
		if b.GetType() == thrift.STRING {
			b.data = any
		} else if !b.optional {
			b.data = ""
		}
	default: // 一般为nil/ptr
		value := reflect.ValueOf(o)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
			return b.BindValue(value.Elem().Interface())
		}
		inline.WithFields("o", o).Warnln("unknown type of o")
		b.data = any
	}

	return nil
}

func (b *BaseArgs) IsSkip() bool {
	return b.optional && b.data == nil
}

func (b *BaseArgs) Index() int16 {
	return b.id
}

func (b *BaseArgs) ThriftName() string {
	return b.thriftName
}

type BaseParser struct {
	model generic.ThriftFieldModel
}

func (b *BaseParser) Parse() (Args, error) {

	return &BaseArgs{
		optional:   b.model.Optional,
		ttype:      b.model.Type,
		jsonPath:   b.model.FieldName,
		id:         b.model.Idx,
		thriftName: b.model.FieldName,
	}, nil
}

func NewBaseParser(model generic.ThriftFieldModel) *BaseParser {
	return &BaseParser{model: model}
}
