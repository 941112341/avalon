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

	gen, err := repo.R.Get()
	if err != nil {
		panic(err)
	}

	gen.MaxID = 2000
	rows, err := repo.R.UpdateVersion(*gen)
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
	fmt.Println(inline.ToJsonString(gen))
}

func TestSave(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}
	err = repo.R.Save(IdGenerator{
		ID:      0,
		MaxID:   100,
		Length:  100,
		BizID:   "flow",
		Version: 0,
	})
	if err != nil {
		panic(err)
	}

}

func TestIdGeneratorRepository_FindByMaxIDBetween(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	gen, err := repo.R.FindByMaxIDBetween(1000, 2000)
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(gen))

}
