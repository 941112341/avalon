package server

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/avalon"
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

func (s *ThriftServer) wrap(handler interface{}, wrapper avalon.Wrapper) interface{} {
	processorMap := ProcessorMap{processorMap: map[string]*ProcessorFunction{}}
	val := reflect.ValueOf(handler)
	if val.NumMethod() == 0 {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < typ.NumMethod(); i++ {
		method := val.Method(i)

		var call = func(ctx context.Context, invoke *avalon.Invoke) error {
			result := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(invoke.Request)})
			err, _ := result[1].Interface().(error)
			if err == nil {
				invoke.Response = result[0].Interface()
			}
			return err
		}

		call = wrapper.Middleware(call)
		methodType := method.Type()
		methodName := typ.Method(i).Name

		processorMap.processorMap[methodName] = &ProcessorFunction{
			call:         call,
			methodName:   methodName,
			requestType:  methodType.In(1).Elem(),
			responseType: methodType.Out(0).Elem(),
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
