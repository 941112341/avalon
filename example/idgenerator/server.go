package main

import (
	"fmt"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/initial"
	"os"
)

func main() {

	fmt.Println(os.Getwd())

	var err error
	err = initial.InitAll()
	if err != nil {
		panic(err)
	}

	err = idgenerator.Run("avalon.test.idgenerator", &handler)
	if err != nil {
		panic(err)
	}
}
