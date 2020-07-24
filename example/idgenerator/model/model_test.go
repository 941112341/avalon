package model

import (
	"fmt"
	"github.com/941112341/avalon/example/idgenerator/initial"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var model ModelTest

func init() {
	_ = registry.Registry("", &model)
}

type ModelTest struct {
	F GeneratorFactory `inject:"GeneratorFactory"`
}

func TestGenerator(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	g, err := model.F.Create()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 25; i++ {
		go func() {
			arr, err := g.Assign(100, "flow")
			assert.Nil(t, err)
			fmt.Println(inline.ToJsonString(arr))
		}()
	}

	time.Sleep(20 * time.Second)
}
