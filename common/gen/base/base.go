// Autogenerated by Thrift Compiler (0.13.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package base

import(
	"bytes"
	"context"
	"reflect"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = reflect.DeepEqual
var _ = bytes.Equal

// Attributes:
//  - Psm
//  - IP
//  - Time
//  - Extra
//  - Base
type Base struct {
  Psm string `thrift:"psm,1" db:"psm" json:"psm"`
  IP string `thrift:"ip,2" db:"ip" json:"ip"`
  Time int64 `thrift:"time,3" db:"time" json:"time"`
  Extra map[string]string `thrift:"extra,4" db:"extra" json:"extra"`
  Base *Base `thrift:"base,5" db:"base" json:"base,omitempty"`
}

func NewBase() *Base {
  return &Base{}
}


func (p *Base) GetPsm() string {
  return p.Psm
}

func (p *Base) GetIP() string {
  return p.IP
}

func (p *Base) GetTime() int64 {
  return p.Time
}

func (p *Base) GetExtra() map[string]string {
  return p.Extra
}
var Base_Base_DEFAULT *Base
func (p *Base) GetBase() *Base {
  if !p.IsSetBase() {
    return Base_Base_DEFAULT
  }
return p.Base
}
func (p *Base) IsSetBase() bool {
  return p.Base != nil
}

func (p *Base) Read(iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField3(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField4(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.STRUCT {
        if err := p.ReadField5(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *Base)  ReadField1(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Psm = v
}
  return nil
}

func (p *Base)  ReadField2(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.IP = v
}
  return nil
}

func (p *Base)  ReadField3(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Time = v
}
  return nil
}

func (p *Base)  ReadField4(iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin()
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]string, size)
  p.Extra =  tMap
  for i := 0; i < size; i ++ {
var _key0 string
    if v, err := iprot.ReadString(); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key0 = v
}
var _val1 string
    if v, err := iprot.ReadString(); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val1 = v
}
    p.Extra[_key0] = _val1
  }
  if err := iprot.ReadMapEnd(); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *Base)  ReadField5(iprot thrift.TProtocol) error {
  p.Base = &Base{}
  if err := p.Base.Read(iprot); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Base), err)
  }
  return nil
}

func (p *Base) Write(oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin("Base"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(oprot); err != nil { return err }
    if err := p.writeField2(oprot); err != nil { return err }
    if err := p.writeField3(oprot); err != nil { return err }
    if err := p.writeField4(oprot); err != nil { return err }
    if err := p.writeField5(oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *Base) writeField1(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("psm", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:psm: ", p), err) }
  if err := oprot.WriteString(string(p.Psm)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.psm (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:psm: ", p), err) }
  return err
}

func (p *Base) writeField2(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("ip", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:ip: ", p), err) }
  if err := oprot.WriteString(string(p.IP)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ip (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:ip: ", p), err) }
  return err
}

func (p *Base) writeField3(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("time", thrift.I64, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:time: ", p), err) }
  if err := oprot.WriteI64(int64(p.Time)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.time (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:time: ", p), err) }
  return err
}

func (p *Base) writeField4(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("extra", thrift.MAP, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:extra: ", p), err) }
  if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Extra)); err != nil {
    return thrift.PrependError("error writing map begin: ", err)
  }
  for k, v := range p.Extra {
    if err := oprot.WriteString(string(k)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    if err := oprot.WriteString(string(v)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
  }
  if err := oprot.WriteMapEnd(); err != nil {
    return thrift.PrependError("error writing map end: ", err)
  }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:extra: ", p), err) }
  return err
}

func (p *Base) writeField5(oprot thrift.TProtocol) (err error) {
  if p.IsSetBase() {
    if err := oprot.WriteFieldBegin("base", thrift.STRUCT, 5); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:base: ", p), err) }
    if err := p.Base.Write(oprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Base), err)
    }
    if err := oprot.WriteFieldEnd(); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 5:base: ", p), err) }
  }
  return err
}

func (p *Base) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("Base(%+v)", *p)
}

// Attributes:
//  - Code
//  - Message
//  - Extra
type BaseResp struct {
  Code int32 `thrift:"code,1" db:"code" json:"code"`
  Message string `thrift:"message,2" db:"message" json:"message"`
  Extra map[string]string `thrift:"extra,3" db:"extra" json:"extra"`
}

func NewBaseResp() *BaseResp {
  return &BaseResp{}
}


func (p *BaseResp) GetCode() int32 {
  return p.Code
}

func (p *BaseResp) GetMessage() string {
  return p.Message
}

func (p *BaseResp) GetExtra() map[string]string {
  return p.Extra
}
func (p *BaseResp) Read(iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField1(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.MAP {
        if err := p.ReadField3(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *BaseResp)  ReadField1(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Code = v
}
  return nil
}

func (p *BaseResp)  ReadField2(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Message = v
}
  return nil
}

func (p *BaseResp)  ReadField3(iprot thrift.TProtocol) error {
  _, _, size, err := iprot.ReadMapBegin()
  if err != nil {
    return thrift.PrependError("error reading map begin: ", err)
  }
  tMap := make(map[string]string, size)
  p.Extra =  tMap
  for i := 0; i < size; i ++ {
var _key2 string
    if v, err := iprot.ReadString(); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _key2 = v
}
var _val3 string
    if v, err := iprot.ReadString(); err != nil {
    return thrift.PrependError("error reading field 0: ", err)
} else {
    _val3 = v
}
    p.Extra[_key2] = _val3
  }
  if err := iprot.ReadMapEnd(); err != nil {
    return thrift.PrependError("error reading map end: ", err)
  }
  return nil
}

func (p *BaseResp) Write(oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin("BaseResp"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(oprot); err != nil { return err }
    if err := p.writeField2(oprot); err != nil { return err }
    if err := p.writeField3(oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *BaseResp) writeField1(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("code", thrift.I32, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:code: ", p), err) }
  if err := oprot.WriteI32(int32(p.Code)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.code (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:code: ", p), err) }
  return err
}

func (p *BaseResp) writeField2(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("message", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:message: ", p), err) }
  if err := oprot.WriteString(string(p.Message)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.message (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:message: ", p), err) }
  return err
}

func (p *BaseResp) writeField3(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("extra", thrift.MAP, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:extra: ", p), err) }
  if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Extra)); err != nil {
    return thrift.PrependError("error writing map begin: ", err)
  }
  for k, v := range p.Extra {
    if err := oprot.WriteString(string(k)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
    if err := oprot.WriteString(string(v)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err) }
  }
  if err := oprot.WriteMapEnd(); err != nil {
    return thrift.PrependError("error writing map end: ", err)
  }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:extra: ", p), err) }
  return err
}

func (p *BaseResp) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("BaseResp(%+v)", *p)
}

