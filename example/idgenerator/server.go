package main

import (
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/initial"
)

func main() {

	var err error
	err = initial.InitAll()
	if err != nil {
		panic(err)
	}

	err = idgenerator.Run(Handler{})
	if err != nil {
		panic(err)
	}
}
