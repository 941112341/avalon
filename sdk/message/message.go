// Autogenerated by Thrift Compiler (0.13.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package message

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = reflect.DeepEqual
var _ = bytes.Equal

// Attributes:
//  - Header
//  - Body
//  - MethodName
//  - URL
type MessageRequest struct {
	Header     map[string]string `thrift:"header,1" db:"header" json:"header"`
	Body       []byte            `thrift:"body,2" db:"body" json:"body"`
	MethodName string            `thrift:"methodName,3" db:"methodName" json:"methodName"`
	URL        string            `thrift:"url,4" db:"url" json:"url"`
}

func NewMessageRequest() *MessageRequest {
	return &MessageRequest{}
}

func (p *MessageRequest) GetHeader() map[string]string {
	return p.Header
}

func (p *MessageRequest) GetBody() []byte {
	return p.Body
}

func (p *MessageRequest) GetMethodName() string {
	return p.MethodName
}

func (p *MessageRequest) GetURL() string {
	return p.URL
}
func (p *MessageRequest) Read(iprot thrift.TProtocol) error {
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
			if fieldTypeId == thrift.MAP {
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
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField3(iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(fieldTypeId); err != nil {
					return err
				}
			}
		case 4:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField4(iprot); err != nil {
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

func (p *MessageRequest) ReadField1(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Header = tMap
	for i := 0; i < size; i++ {
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
		p.Header[_key0] = _val1
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *MessageRequest) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Body = v
	}
	return nil
}

func (p *MessageRequest) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.MethodName = v
	}
	return nil
}

func (p *MessageRequest) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.URL = v
	}
	return nil
}

func (p *MessageRequest) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("MessageRequest"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
			return err
		}
		if err := p.writeField4(oprot); err != nil {
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

func (p *MessageRequest) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("header", thrift.MAP, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:header: ", p), err)
	}
	if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Header)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Header {
		if err := oprot.WriteString(string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteString(string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:header: ", p), err)
	}
	return err
}

func (p *MessageRequest) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("body", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:body: ", p), err)
	}
	if err := oprot.WriteBinary(p.Body); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.body (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:body: ", p), err)
	}
	return err
}

func (p *MessageRequest) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("methodName", thrift.STRING, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:methodName: ", p), err)
	}
	if err := oprot.WriteString(string(p.MethodName)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.methodName (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:methodName: ", p), err)
	}
	return err
}

func (p *MessageRequest) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("url", thrift.STRING, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:url: ", p), err)
	}
	if err := oprot.WriteString(string(p.URL)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.url (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:url: ", p), err)
	}
	return err
}

func (p *MessageRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MessageRequest(%+v)", *p)
}

// Attributes:
//  - Header
//  - Body
//  - Status
type MessageResponse struct {
	Header map[string]string `thrift:"header,1" db:"header" json:"header"`
	Body   []byte            `thrift:"body,2" db:"body" json:"body"`
	Status int32             `thrift:"status,3" db:"status" json:"status"`
}

func NewMessageResponse() *MessageResponse {
	return &MessageResponse{}
}

func (p *MessageResponse) GetHeader() map[string]string {
	return p.Header
}

func (p *MessageResponse) GetBody() []byte {
	return p.Body
}

func (p *MessageResponse) GetStatus() int32 {
	return p.Status
}
func (p *MessageResponse) Read(iprot thrift.TProtocol) error {
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
			if fieldTypeId == thrift.MAP {
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
			if fieldTypeId == thrift.I32 {
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

func (p *MessageResponse) ReadField1(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Header = tMap
	for i := 0; i < size; i++ {
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
		p.Header[_key2] = _val3
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *MessageResponse) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Body = v
	}
	return nil
}

func (p *MessageResponse) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Status = v
	}
	return nil
}

func (p *MessageResponse) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("MessageResponse"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(oprot); err != nil {
			return err
		}
		if err := p.writeField2(oprot); err != nil {
			return err
		}
		if err := p.writeField3(oprot); err != nil {
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

func (p *MessageResponse) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("header", thrift.MAP, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:header: ", p), err)
	}
	if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Header)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Header {
		if err := oprot.WriteString(string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteString(string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:header: ", p), err)
	}
	return err
}

func (p *MessageResponse) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("body", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:body: ", p), err)
	}
	if err := oprot.WriteBinary(p.Body); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.body (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:body: ", p), err)
	}
	return err
}

func (p *MessageResponse) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("status", thrift.I32, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:status: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Status)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.status (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:status: ", p), err)
	}
	return err
}

func (p *MessageResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MessageResponse(%+v)", *p)
}

type MessageService interface {
	// Parameters:
	//  - Request
	MessageDispatcher(ctx context.Context, request *MessageRequest) (r *MessageResponse, err error)
}

type MessageServiceClient struct {
	c thrift.TClient
}

func NewMessageServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *MessageServiceClient {
	return &MessageServiceClient{
		c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
	}
}

func NewMessageServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *MessageServiceClient {
	return &MessageServiceClient{
		c: thrift.NewTStandardClient(iprot, oprot),
	}
}

func NewMessageServiceClient(c thrift.TClient) *MessageServiceClient {
	return &MessageServiceClient{
		c: c,
	}
}

func (p *MessageServiceClient) Client_() thrift.TClient {
	return p.c
}

// Parameters:
//  - Request
func (p *MessageServiceClient) MessageDispatcher(ctx context.Context, request *MessageRequest) (r *MessageResponse, err error) {
	var _args4 MessageServiceMessageDispatcherArgs
	_args4.Request = request
	var _result5 MessageServiceMessageDispatcherResult
	if err = p.Client_().Call(ctx, "MessageDispatcher", &_args4, &_result5); err != nil {
		return
	}
	return _result5.GetSuccess(), nil
}

type MessageServiceProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      MessageService
}

func (p *MessageServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *MessageServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *MessageServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewMessageServiceProcessor(handler MessageService) *MessageServiceProcessor {

	self6 := &MessageServiceProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self6.processorMap["MessageDispatcher"] = &messageServiceProcessorMessageDispatcher{handler: handler}
	return self6
}

func (p *MessageServiceProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x7 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x7.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush(ctx)
	return false, x7

}

type messageServiceProcessorMessageDispatcher struct {
	handler MessageService
}

func (p *messageServiceProcessorMessageDispatcher) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := MessageServiceMessageDispatcherArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("MessageDispatcher", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush(ctx)
		return false, err
	}

	iprot.ReadMessageEnd()
	result := MessageServiceMessageDispatcherResult{}
	var retval *MessageResponse
	var err2 error
	if retval, err2 = p.handler.MessageDispatcher(ctx, args.Request); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing MessageDispatcher: "+err2.Error())
		oprot.WriteMessageBegin("MessageDispatcher", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush(ctx)
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("MessageDispatcher", thrift.REPLY, seqId); err2 != nil {
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

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Request
type MessageServiceMessageDispatcherArgs struct {
	Request *MessageRequest `thrift:"request,1" db:"request" json:"request"`
}

func NewMessageServiceMessageDispatcherArgs() *MessageServiceMessageDispatcherArgs {
	return &MessageServiceMessageDispatcherArgs{}
}

var MessageServiceMessageDispatcherArgs_Request_DEFAULT *MessageRequest

func (p *MessageServiceMessageDispatcherArgs) GetRequest() *MessageRequest {
	if !p.IsSetRequest() {
		return MessageServiceMessageDispatcherArgs_Request_DEFAULT
	}
	return p.Request
}
func (p *MessageServiceMessageDispatcherArgs) IsSetRequest() bool {
	return p.Request != nil
}

func (p *MessageServiceMessageDispatcherArgs) Read(iprot thrift.TProtocol) error {
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

func (p *MessageServiceMessageDispatcherArgs) ReadField1(iprot thrift.TProtocol) error {
	p.Request = &MessageRequest{}
	if err := p.Request.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Request), err)
	}
	return nil
}

func (p *MessageServiceMessageDispatcherArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("MessageDispatcher_args"); err != nil {
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

func (p *MessageServiceMessageDispatcherArgs) writeField1(oprot thrift.TProtocol) (err error) {
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

func (p *MessageServiceMessageDispatcherArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MessageServiceMessageDispatcherArgs(%+v)", *p)
}

// Attributes:
//  - Success
type MessageServiceMessageDispatcherResult struct {
	Success *MessageResponse `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func NewMessageServiceMessageDispatcherResult() *MessageServiceMessageDispatcherResult {
	return &MessageServiceMessageDispatcherResult{}
}

var MessageServiceMessageDispatcherResult_Success_DEFAULT *MessageResponse

func (p *MessageServiceMessageDispatcherResult) GetSuccess() *MessageResponse {
	if !p.IsSetSuccess() {
		return MessageServiceMessageDispatcherResult_Success_DEFAULT
	}
	return p.Success
}
func (p *MessageServiceMessageDispatcherResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *MessageServiceMessageDispatcherResult) Read(iprot thrift.TProtocol) error {
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

func (p *MessageServiceMessageDispatcherResult) ReadField0(iprot thrift.TProtocol) error {
	p.Success = &MessageResponse{}
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *MessageServiceMessageDispatcherResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("MessageDispatcher_result"); err != nil {
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

func (p *MessageServiceMessageDispatcherResult) writeField0(oprot thrift.TProtocol) (err error) {
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

func (p *MessageServiceMessageDispatcherResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("MessageServiceMessageDispatcherResult(%+v)", *p)
}
