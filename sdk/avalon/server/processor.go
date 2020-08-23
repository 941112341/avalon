package server

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
)

type ProcessorMap struct {
	processorMap map[string]*ProcessorFunction
}

func (p *ProcessorMap) Process(ctx context.Context, in, out thrift.TProtocol) (bool, thrift.TException) {
	name, _, seqId, err := in.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.processorMap[name]; ok {
		return processor.Process(ctx, seqId, in, out)
	}
	_ = in.Skip(thrift.STRUCT)
	_ = in.ReadMessageEnd()
	x4 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	_ = out.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	_ = x4.Write(out)
	_ = out.WriteMessageEnd()
	_ = out.Flush(ctx)
	return false, x4
}

type ProcessorFunction struct {
	call                      Call
	requestType, responseType reflect.Type
	methodName                string
}

type Args struct {
	Request     thrift.TStruct `thrift:"request,1" db:"request" json:"request"`
	requestType reflect.Type
	method      string
}

func (p *Args) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(fmt.Sprintf("%s_args", p.method)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Args) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("request", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:request: ", p), err)
	}
	if err := p.Request.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Request), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:request: ", p), err)
	}
	return err
}

func (p *Args) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRUCT {
				if err := p.ReadField1(iprot); err != nil {
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

func (p *Args) ReadField1(iprot thrift.TProtocol) error {
	request := reflect.New(p.requestType).Interface()
	p.Request = request.(thrift.TStruct)
	if err := p.Request.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Request), err)
	}
	return nil
}

type Result struct {
	Success      thrift.TStruct `thrift:"success,0" db:"success" json:"success,omitempty"`
	responseType reflect.Type
	method       string
}

func (p *Result) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if fieldTypeId == thrift.STRUCT {
				if err := p.ReadField0(iprot); err != nil {
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

func (p *Result) ReadField0(iprot thrift.TProtocol) error {
	response := reflect.New(p.responseType).Interface()
	p.Success = response.(thrift.TStruct)
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *Result) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(fmt.Sprintf("%s_result", p.method)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Result) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *Result) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ProcessorFunction) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := Args{method: p.methodName, requestType: p.requestType}
	if err = args.Read(iprot); err != nil {
		_ = iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		_ = oprot.WriteMessageBegin(p.methodName, thrift.EXCEPTION, seqId)
		_ = x.Write(oprot)
		_ = oprot.WriteMessageEnd()
		_ = oprot.Flush(ctx)
		return false, err
	}

	_ = iprot.ReadMessageEnd()

	result := Result{method: p.methodName, responseType: p.responseType}
	var err2 error
	invoke := &Invoke{MethodName: p.methodName, Request: args.Request, Response: reflect.New(p.responseType).Interface()}
	if err2 = p.call(ctx, invoke); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing GenIDs: "+err2.Error())
		_ = oprot.WriteMessageBegin(p.methodName, thrift.EXCEPTION, seqId)
		_ = x.Write(oprot)
		_ = oprot.WriteMessageEnd()
		_ = oprot.Flush(ctx)
		return true, err2
	} else {
		result.Success, _ = invoke.Response.(thrift.TStruct)
	}
	if err2 = oprot.WriteMessageBegin(p.methodName, thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(ctx); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}
