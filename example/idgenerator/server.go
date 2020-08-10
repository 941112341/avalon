package main

import (
	"fmt"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/initial"
	"github.com/941112341/avalon/sdk/avalon/server"
	"os"
	"time"
)

func main() {

	fmt.Println(os.Getwd())

	var err error
	err = initial.InitAll()
	if err != nil {
		panic(err)
	}

	err = server.Builder().
		Timeout(2 * time.Second).
		Hostport(server.NewZkDiscoverBuilder().
			Port(8889).
			PSM("avalon.test.idgenerator").
			Build()).
		Build().Run(idgenerator.NewIDGeneratorProcessor(&handler))
	if err != nil {
		panic(err)
	}
}
