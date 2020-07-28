package generic

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type lazyField struct {
	lazy  func() []*CommonTStruct
	cache []*CommonTStruct
}

func (l *lazyField) lazyFields() []*CommonTStruct {
	if l.cache == nil {
		l.cache = l.lazy()
	}
	return l.cache
}

func (l *lazyField) fields() []*CommonTStruct {
	return l.cache
}

type CommonTStruct struct {
	ID         int16
	StructName string // 类型名，用于写入struct,只有结构体有
	FieldName  string // thrift 上的name，用于写入field, 一般是小写
	JSONPath   string // 用于json解析

	Type  thrift.TType
	Value interface{} `json:",omitempty"`

	ArrayStruct    *CommonTStruct `json:",omitempty"` // list
	MapKeyStruct   *CommonTStruct `json:",omitempty"` // map
	MapValueStruct *CommonTStruct `json:",omitempty"` // map
	FieldMap       lazyField      `json:",omitempty"` // struct
}

func (c *CommonTStruct) findSubField(id int16) *CommonTStruct {

	fields := c.FieldMap.lazyFields()
	for _, tStruct := range fields {
		if tStruct == nil {
			continue
		}
		if tStruct.ID == id {
			return tStruct
		}
	}
	return nil
}

func (c *CommonTStruct) Write(p thrift.TProtocol) error {
	s := time.Now()
	defer func() {
		inline.WithFields("spend", time.Since(s).String()).Debugln("write spend")
	}()

	switch c.Type {
	case thrift.STRUCT:
		if err := p.WriteStructBegin(c.StructName); err != nil {
			return fmt.Errorf("%T write struct begin error: %s", c, err)
		}
		fieldMap := c.FieldMap.fields()
		for _, tStruct := range fieldMap {
			if tStruct == nil {
				continue
			}

			if tStruct.Type == thrift.STRUCT && len(fieldMap) == 0 {
				continue
			}
			if err := p.WriteFieldBegin(tStruct.FieldName, tStruct.Type, tStruct.ID); err != nil {
				return fmt.Errorf("%T write field begin error %d:groupName: %s", tStruct, tStruct.ID, err)
			}

			if err := tStruct.Write(p); err != nil {
				return err
			}
			if err := p.WriteFieldEnd(); err != nil {
				return fmt.Errorf("%T write field end error %d:groupName: %s", tStruct, tStruct.ID, err)
			}
			inline.WithFields("fieldName", tStruct.FieldName, "type", tStruct.Type, "id", tStruct.ID).Debugln("write field success")
		}
		if err := p.WriteFieldStop(); err != nil {
			return fmt.Errorf("write field stop error: %s", err)
		}
		if err := p.WriteStructEnd(); err != nil {
			return fmt.Errorf("write struct stop error: %s", err)
		}
		inline.WithFields("struct", c).Debugln("write struct success")
	case thrift.STRING:
		str, ok := c.Value.(string)
		if !ok {
			return fmt.Errorf("c.value %+v is not string", c.Value)
		}
		if err := p.WriteString(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.DOUBLE:
		str, ok := c.Value.(float64)
		if !ok {
			return fmt.Errorf("c.value %+v is not double", c.Value)
		}
		if err := p.WriteDouble(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.BOOL:
		v, ok := c.Value.(bool)
		if !ok {
			return fmt.Errorf("c.value %+v is not bool", c.Value)
		}
		if err := p.WriteBool(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.BYTE:
		v, ok := c.Value.(int8)
		if !ok {
			return fmt.Errorf("c.value %+v is not  byte", c.Value)
		}
		if err := p.WriteByte(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.I16:
		v, ok := c.Value.(int16)
		if !ok {
			return fmt.Errorf("c.value %+v is not i16", c.Value)
		}
		if err := p.WriteI16(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.I32:
		v, ok := c.Value.(int32)
		if !ok {
			return fmt.Errorf("c.value %+v is not i32", c.Value)
		}
		if err := p.WriteI32(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.I64:
		v, ok := c.Value.(int64)
		if !ok {
			return fmt.Errorf("c.value %+v is not i64", c.Value)
		}
		if err := p.WriteI64(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.LIST:
		vs, _ := c.Value.([]interface{})
		typ := c.ArrayStruct
		if err := p.WriteListBegin(typ.Type, len(vs)); err != nil {
			return err
		}
		for _, v := range vs {
			typ.Value = v
			if err := typ.Write(p); err != nil {
				return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
			}
		}
		if err := p.WriteListEnd(); err != nil {
			return fmt.Errorf("error writing list end: %s", err)
		}
	case thrift.MAP:
		vmap, _ := c.Value.(map[interface{}]interface{})
		ks, vs := c.MapKeyStruct, c.MapValueStruct

		if err := p.WriteMapBegin(ks.Type, vs.Type, len(vmap)); err != nil {
			return err
		}
		for k, v := range vmap {
			ks.Value = k
			if err := ks.Write(p); err != nil {
				return err
			}

			vs.Value = v
			if err := vs.Write(p); err != nil {
				return err
			}
		}
		if err := p.WriteMapEnd(); err != nil {
			return err
		}

	}
	return nil
}

func (c *CommonTStruct) Read(p thrift.TProtocol) error {
	s := time.Now()
	defer func() {
		inline.WithFields("spend", time.Since(s).String()).Debugln("read spend")
	}()

	switch c.Type {
	case thrift.STRUCT:
		if _, err := p.ReadStructBegin(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T read error: ", c), err)
		}
		for {
			_, fieldTypeId, fieldId, err := p.ReadFieldBegin()
			if err != nil {
				return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", c, fieldId), err)
			}
			if fieldTypeId == thrift.STOP {
				break
			}
			field := c.findSubField(fieldId)
			if field == nil {
				if err := p.Skip(fieldTypeId); err != nil {
					return err
				}
				continue
			}
			if err = field.Read(p); err != nil {
				return err
			}
			if err := p.ReadFieldEnd(); err != nil {
				return err
			}
		}
		if err := p.ReadStructEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", c), err)
		}
	case thrift.STOP: // do nothing
	case thrift.STRING:
		if v, err := p.ReadString(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			c.Value = v
		}
	case thrift.BOOL:
		if v, err := p.ReadBool(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			c.Value = v
		}
	case thrift.BYTE:
		if v, err := p.ReadByte(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			c.Value = v
		}
	case thrift.I16:
		if v, err := p.ReadI16(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			c.Value = v
		}
	case thrift.I32:
		if v, err := p.ReadI32(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			c.Value = v
		}
	case thrift.I64:
		if v, err := p.ReadI64(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			c.Value = v
		}
	case thrift.DOUBLE:
		if v, err := p.ReadDouble(); err != nil {
			return thrift.PrependError("error reading field: ", err)
		} else {
			c.Value = v
		}
	case thrift.LIST:
		_, size, err := p.ReadListBegin()
		if err != nil {
			return thrift.PrependError("error reading list begin: ", err)
		}
		l := c.ArrayStruct
		ifaces := make([]interface{}, 0)
		for i := 0; i < size; i++ {

			if err = l.Read(p); err != nil {
				return err
			}
			ifaces = append(ifaces, l.Value)
		}
		c.Value = ifaces
		if err := p.ReadListEnd(); err != nil {
			return thrift.PrependError("error reading list end: ", err)
		}
	case thrift.MAP:
		_, _, size, err := p.ReadMapBegin()
		if err != nil {
			return thrift.PrependError("error reading map begin: ", err)
		}
		k, v := c.MapKeyStruct, c.MapValueStruct
		imaps := make(map[interface{}]interface{})
		for i := 0; i < size; i++ {
			if err := k.Read(p); err != nil {
				return thrift.PrependError("error reading map key: ", err)
			}

			if err := v.Read(p); err != nil {
				return thrift.PrependError("error reading map key: ", err)
			}
			imaps[k.Value] = v.Value
		}
		c.Value = imaps
		if err = p.ReadMapEnd(); err != nil {
			return thrift.PrependError("error reading map end: ", err)
		}
	}
	return nil
}

func (c *CommonTStruct) WithValue(args jsoniter.Any) {
	get := args.Get(c.JSONPath)
	switch c.Type {
	case thrift.BOOL:
		c.Value = get.ToBool()
	case thrift.I08:
		c.Value = int8(get.ToInt())
	case thrift.I16:
		c.Value = int16(get.ToInt())
	case thrift.I32:
		c.Value = get.ToInt32()
	case thrift.I64:
		c.Value = get.ToInt64()
	case thrift.DOUBLE:
		c.Value = get.ToFloat32()
	case thrift.STRING:
		c.Value = get.ToString()
	case thrift.LIST:
		ifaces := make([]interface{}, 0)
		for i := 0; i < get.Size(); i++ {
			ifaces = append(ifaces, get.Get(i).GetInterface())
		}
		c.Value = ifaces
	case thrift.MAP:
		m := make(map[interface{}]interface{})
		for _, key := range get.Keys() {
			m[key] = get.Get(key).GetInterface()
		}
		c.Value = m
	case thrift.STRUCT:
		any := args
		if c.JSONPath != "" {
			any = get
		}
		if any.LastError() != nil {
			return
		}
		fields := c.FieldMap.lazyFields()
		for _, tStruct := range fields {
			tStruct.WithValue(any)
		}
	default:
		inline.WithFields("fieldName", c.FieldName).Debugln("unsupport type")
	}
}
