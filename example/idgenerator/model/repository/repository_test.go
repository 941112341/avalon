package repository

import (
	"fmt"
	"github.com/941112341/avalon/example/idgenerator/initial"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/sdk/inline"
	"testing"
)

type RepoTest struct {
	R IdGeneratorRepository `inject:"IdGeneratorRepository"`
}

var repo RepoTest

func init() {
	_ = registry.Registry("", &repo)
}

func TestGet(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	gen, err := repo.R.Get()
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(gen))
}

func TestUpdate(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

}
