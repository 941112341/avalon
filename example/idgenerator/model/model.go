package model

import (
	"github.com/941112341/avalon/example/idgenerator/model/repository"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/sdk/inline"
	"sync"
)

type Generator interface {
	Assign(cnt int64, bizId string) ([]int64, error)
	GetIds() []int64
}

type GeneratorFactory interface {
	Create() (Generator, error)
}

type factory struct {
	gen  Generator
	once sync.Once
	R    repository.IdGeneratorRepository `inject:"IdGeneratorRepository"`
}

func (f *factory) Create() (gen Generator, err error) {
	if f.gen != nil {
		return f.gen, nil
	}
	f.once.Do(func() {
		gen, err = NewGeneratorModel(f.R)
		if err != nil {
			inline.WithFields("err", err).Error("generator err")
		}
	})
	return
}

func init() {
	_ = registry.Registry("GeneratorFactory", &factory{})
}
