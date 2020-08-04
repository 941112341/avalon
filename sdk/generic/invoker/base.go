package invoker

import (
	"fmt"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/json-iterator/go"
	"reflect"
	"strconv"
	"strings"
)

const (
	ErrUnsupportType inline.AvalonErrorCode = iota + 500
)

type BaseArgs struct {
	ArgsMeta
	data interface{}
}

func (b *BaseArgs) Meta() ArgsMeta {
	return b.ArgsMeta
}

func (b *BaseArgs) Data() interface{} {
	return b.data
}

func (b *BaseArgs) Write(p thrift.TProtocol) error {
	switch b.Type() {
	case thrift.STRING:
		str, _ := b.data.(string)
		if err := p.WriteString(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.ID(), err)
		}
	case thrift.DOUBLE:
		str, _ := b.data.(float64)
		if err := p.WriteDouble(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.ID(), err)
		}
	case thrift.BOOL:
		v, _ := b.data.(bool)

		if err := p.WriteBool(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.ID(), err)
		}
	case thrift.BYTE:
		v, _ := b.data.(int8)
		if err := p.WriteByte(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.ID(), err)
		}
	case thrift.I16:
		v, _ := b.data.(int16)
		if err := p.WriteI16(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.ID(), err)
		}
	case thrift.I32:
		v, _ := b.data.(int32)
		if err := p.WriteI32(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.ID(), err)
		}
	case thrift.I64:
		v, _ := b.data.(int64)
		if err := p.WriteI64(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", b, b.ID(), err)
		}
	default:
		return inline.NewError(ErrUnsupportType, "write unknown type %s", b.Type())
	}
	return nil
}

func (b *BaseArgs) Read(p thrift.TProtocol) error {
	switch b.Type() {
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
		return inline.NewError(ErrUnsupportType, "read unknown type %s", b.Type())
	}
	return nil
}

func (b *BaseArgs) BindValue(o interface{}) error {
	switch any := o.(type) {
	case jsoniter.Any:

		if err := any.LastError(); err != nil {
			if !b.Optional() {
				inline.WithFields("path", b.JsonPath(), "typeName", b.TypeName()).Infoln("skip field")
				//return inline.Error("bind data is nil %+v", b)
				return nil
			}
			return nil
		}
		var data interface{}
		switch b.Type() {
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
		if b.Type() == thrift.I08 {
			b.data = any
		} else if !b.Optional() {
			b.data = int8(0)
		}
	case int16:
		if b.Type() == thrift.I16 {
			b.data = any
		} else if !b.Optional() {
			b.data = int16(0)
		}
	case int32:
		if b.Type() == thrift.I32 {
			b.data = any
		} else if !b.Optional() {
			b.data = int32(0)
		}
	case int64:
		if b.Type() == thrift.I64 {
			b.data = any
		} else if !b.Optional() {
			b.data = int64(0)
		}
	case bool:
		if b.Type() == thrift.BOOL {
			b.data = any
		} else if !b.Optional() {
			b.data = false
		}
	case float32:
		if b.Type() == thrift.DOUBLE {
			b.data = float64(any)
		} else if !b.Optional() {
			b.data = float64(0)
		}
	case float64:
		if b.Type() == thrift.DOUBLE {
			b.data = any
		} else if !b.Optional() {
			b.data = float64(0)
		}
	case string:
		if b.Type() == thrift.STRING {
			b.data = any
		} else if !b.Optional() {
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
	return b.Optional() && b.data == nil
}

type BaseParser struct {
	ArgsMetaParser
}

func (b *BaseParser) Parse() (Args, error) {
	return &BaseArgs{
		ArgsMeta: b.ArgsMeta(),
	}, nil
}

func NewBaseParser(parser ArgsMetaParser) *BaseParser {
	return &BaseParser{parser}
}

type BaseReflectParser struct {
	Field reflect.StructField
}

func (b *BaseReflectParser) Parse() (Args, error) {
	return nil, nil

}

type CommonModel struct {
	goType     reflect.Type
	optional   bool
	typeName   string
	jsonPath   string
	thriftName string
	idx        int16
	ttype      thrift.TType
}

func (m *CommonModel) ID() int16 {
	return m.idx
}

func (m *CommonModel) Type() thrift.TType {
	return m.ttype
}

func (m *CommonModel) Optional() bool {
	return m.optional
}

func (m *CommonModel) TypeName() string {
	return m.typeName
}

func (m *CommonModel) JsonPath() string {
	return m.jsonPath
}

func (m *CommonModel) ThriftName() string {
	return m.thriftName
}

func (m *CommonModel) Elem() *CommonModel {
	if m.ttype != thrift.LIST {
		return nil
	}
	goType := m.goType.Elem()
	return &CommonModel{ttype: generic.Type2thrift(goType), typeName: goType.Name(), goType: goType}
}

func (m *CommonModel) KVElem() (*CommonModel, *CommonModel) {
	if m.ttype != thrift.MAP {
		return nil, nil
	}

	types := inline.Redirect(m.goType)
	keyType, valueType := types.Key(), types.Elem()
	return &CommonModel{ttype: generic.Type2thrift(keyType), typeName: keyType.Name(), goType: keyType},
		&CommonModel{ttype: generic.Type2thrift(valueType), typeName: valueType.Name(), goType: valueType}
}

func NewModel(field reflect.StructField) (*CommonModel, error) {
	ttype := field.Type
	tag := field.Tag
	thriftTag := strings.Split(tag.Get("thrift"), ",")
	idx, err := strconv.ParseInt(thriftTag[1], 10, 16)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "tag %+v", tag)
	}
	jsonTag := strings.Split(tag.Get("json"), ",")
	optional := len(jsonTag) > 1
	jsonPath := jsonTag[0]
	return &CommonModel{
		goType:     ttype,
		optional:   optional,
		typeName:   ttype.Name(),
		jsonPath:   jsonPath,
		thriftName: thriftTag[1],
		idx:        int16(idx),
		ttype:      generic.Type2thrift(ttype),
	}, nil
}
