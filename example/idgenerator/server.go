package main

import (
	"fmt"
	"github.com/941112341/avalon/example/idgenerator/initial"
	"github.com/941112341/avalon/sdk/avalon/server"
	"os"
)

func main() {

	fmt.Println(os.Getwd())

	var err error
	err = initial.InitAll()
	if err != nil {
		panic(err)
	}

	err = server.DefaultServer().Run(&handler)

	if err != nil {
		panic(err)
	}
}
