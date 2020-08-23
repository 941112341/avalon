package server

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

type ThriftServer struct {
	Timeout string `default:"2s"`
	timeout time.Duration
	Port    int `default:"8080"`

	server thrift.TServer
}

func (s *ThriftServer) Key() string {
	return "server"
}

func (s *ThriftServer) wrap(handler interface{}, wrapper Wrapper) interface{} {
	processorMap := ProcessorMap{processorMap: map[string]*ProcessorFunction{}}
	val := reflect.ValueOf(handler).Elem()
	for i := 0; i < val.NumMethod(); i++ {
		method := val.Method(i)

		var call Call = func(ctx context.Context, invoke *Invoke) error {
			result := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(invoke.Request)})
			invoke.Response = result[0].Interface()
			err := result[1].Interface().(error)
			return err
		}

		call = wrapper.Middleware(call)
		methodType := method.Type()
		methodName := methodType.Name()
		processorMap.processorMap[methodName] = &ProcessorFunction{
			call:         call,
			methodName:   methodName,
			requestType:  methodType.In(1),
			responseType: methodType.Out(1),
		}

	}
	return &processorMap
}

func (s *ThriftServer) Initial() error {
	s.timeout = inline.Parse(s.Timeout)

	return nil
}

func (s *ThriftServer) Destroy() error {
	if s.server == nil {
		return nil
	}
	return s.server.Stop()
}

func (s *ThriftServer) Run(handler interface{}) error {
	processor, ok := handler.(thrift.TProcessor)
	if !ok {
		return errors.New("handler invalid")
	}
	serverTransport, err := thrift.NewTServerSocketTimeout(fmt.Sprintf(":%d", s.Port), s.timeout)
	if err != nil {
		return errors.WithMessage(err, "new socket")
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	tServer := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	s.server = tServer
	return tServer.Serve()
}
