package generic

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type CommonTStruct struct {
	ID    int16
	Name  string
	Type  thrift.TType
	Value interface{}

	ArrayStruct    *CommonTStruct           // list
	MapKeyStruct   *CommonTStruct           // map
	MapValueStruct *CommonTStruct           // map
	FieldMap       map[int16]*CommonTStruct // struct
}

func (c *CommonTStruct) Write(p thrift.TProtocol) error {

	switch c.Type {
	case thrift.STRUCT:
		if err := p.WriteStructBegin(c.Name); err != nil {
			return fmt.Errorf("%T write struct begin error: %s", c, err)
		}
		for _, tStruct := range c.FieldMap {
			if err := p.WriteFieldBegin(tStruct.Name, thrift.STRING, tStruct.ID); err != nil {
				return fmt.Errorf("%T write field begin error %d:groupName: %s", tStruct, tStruct.ID, err)
			}
			if err := tStruct.Write(p); err != nil {
				return err
			}
			if err := p.WriteFieldEnd(); err != nil {
				return fmt.Errorf("%T write field end error %d:groupName: %s", tStruct, tStruct.ID, err)
			}

		}
		if err := p.WriteFieldStop(); err != nil {
			return fmt.Errorf("write field stop error: %s", err)
		}
		if err := p.WriteStructEnd(); err != nil {
			return fmt.Errorf("write struct stop error: %s", err)
		}
	case thrift.STRING:
		str, ok := c.Value.(string)
		if !ok {
			return nil
		}
		if err := p.WriteString(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.DOUBLE:
		str, ok := c.Value.(float64)
		if !ok {
			return nil
		}
		if err := p.WriteDouble(str); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.BOOL:
		v, ok := c.Value.(bool)
		if !ok {
			return nil
		}
		if err := p.WriteBool(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.BYTE:
		v, ok := c.Value.(int8)
		if !ok {
			return nil
		}
		if err := p.WriteByte(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.I16:
		v, ok := c.Value.(int16)
		if !ok {
			return nil
		}
		if err := p.WriteI16(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.I32:
		v, ok := c.Value.(int32)
		if !ok {
			return nil
		}
		if err := p.WriteI32(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.I64:
		v, ok := c.Value.(int64)
		if !ok {
			return nil
		}
		if err := p.WriteI64(v); err != nil {
			return fmt.Errorf("%T (%d) field write error: %s", c, c.ID, err)
		}
	case thrift.LIST:
		vs, _ := c.Value.([]interface{})
		if len(vs) == 0 {
			return nil
		}
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
		if len(vmap) == 0 {
			return nil
		}
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
			field, ok := c.FieldMap[fieldId]
			if !ok {
				if err := p.Skip(fieldTypeId); err != nil {
					return err
				}
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
