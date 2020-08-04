package invoker

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

type MapArgs struct {
	ArgsMeta

	KeyArgs   Args
	ValueArgs Args
	data      map[interface{}]interface{}
}

func (m *MapArgs) Meta() ArgsMeta {
	return m.ArgsMeta
}

func (m *MapArgs) Data() interface{} {
	return m.data
}

func (m *MapArgs) Write(p thrift.TProtocol) error {
	ks, vs := m.KeyArgs, m.ValueArgs
	vmap := m.data
	if err := p.WriteMapBegin(ks.Meta().Type(), vs.Meta().Type(), len(vmap)); err != nil {
		return err
	}
	for key, value := range vmap {
		if err := ks.BindValue(key); err != nil {
			inline.WithFields("key", key, "ks", ks).Errorln("bind key err %s", err)
			continue
		}
		if err := vs.BindValue(value); err != nil {
			inline.WithFields("value", value, "vs", vs).Errorln("bind value err %s", err)
			continue
		}

		if ks.IsSkip() || vs.IsSkip() {
			continue
		}
		if err := ks.Write(p); err != nil {
			return err
		}
		if err := vs.Write(p); err != nil {
			return err
		}
	}
	if err := p.WriteMapEnd(); err != nil {
		return err
	}
	return nil
}

func (m *MapArgs) Read(p thrift.TProtocol) error {
	_, _, size, err := p.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	k, v := m.KeyArgs, m.ValueArgs
	imaps := make(map[interface{}]interface{})
	for i := 0; i < size; i++ {
		if err := k.Read(p); err != nil {
			return thrift.PrependError("error reading map key: ", err)
		}

		if err := v.Read(p); err != nil {
			return thrift.PrependError("error reading map key: ", err)
		}
		imaps[k.Data()] = v.Data()
	}
	m.data = imaps
	if err = p.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (m *MapArgs) GetType() thrift.TType {
	return thrift.MAP
}

func (m *MapArgs) BindValue(o interface{}) error {
	switch any := o.(type) {
	case jsoniter.Any:
		if any.LastError() != nil {
			if !m.Optional() {
				m.data = map[interface{}]interface{}{}
			}
			return nil
		}
		maps := make(map[interface{}]interface{})
		keys := any.Keys()
		for _, key := range keys {
			value := any.Get(key)
			maps[key] = value.GetInterface()
		}
		m.data = maps
	default:
		value := reflect.ValueOf(o)
		if value.Kind() != reflect.Map {
			return inline.NewError(ErrUnsupportType, "unsupport type %+v", any)
		}
		maps := make(map[interface{}]interface{})
		for _, key := range value.MapKeys() {
			value := value.MapIndex(key)
			maps[key.Interface()] = value.Interface()
		}
		m.data = maps
	}
	return nil
}

func (m *MapArgs) IsSkip() bool {
	return m.Optional() && m.data == nil
}

type MapParser struct {
	ArgsMetaParser
}

func (m *MapParser) Parse() (Args, error) {

	kParser, vParser := m.KVParsers()
	keyArgs, err := kParser.Parse()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "parse key")
	}
	valueArgs, err := vParser.Parse()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "parse value")
	}

	return &MapArgs{
		KeyArgs:   keyArgs,
		ValueArgs: valueArgs,
		ArgsMeta:  m.ArgsMeta(),
	}, nil
}

func NewMapParser(parser ArgsMetaParser) *MapParser {
	return &MapParser{
		parser,
	}
}
