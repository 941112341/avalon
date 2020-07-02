package main

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/message"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
)

type Handler struct {
}

func (Handler Handler) MessageDispatcher(ctx context.Context, r *message.MessageRequest) (*message.MessageResponse, error) {
	return &message.MessageResponse{
		Body: []byte("hello world"),
	}, nil
}

func main() {
	processor := message.NewMessageServiceProcessor(&Handler{})
	serverTransport, err := thrift.NewTServerSocket("localhost:8888")
	if err != nil {
		log.Fatalln("Error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Running at:", "localhost:8888")
	server.Serve()
}
