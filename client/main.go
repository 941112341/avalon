package main

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/message"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"net"
)

func main() {
	tSocket, err := thrift.NewTSocket(net.JoinHostPort("localhost", "8888"))
	if err != nil {
		log.Fatalln("tSocket error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport, err := transportFactory.GetTransport(tSocket)
	if err != nil {
		panic(err)
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := message.NewMessageServiceClientFactory(transport, protocolFactory)

	if err := transport.Open(); err != nil {
		log.Fatalln("Error opening:", "localhostg:8888")
	}
	defer transport.Close()

	data := message.MessageRequest{}
	d, err := client.MessageDispatcher(context.Background(), &data)
	fmt.Println(err)
	fmt.Println(string(d.Body))
}
