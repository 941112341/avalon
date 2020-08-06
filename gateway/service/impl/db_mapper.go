package impl

import (
	"errors"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/gateway/repository"
	"github.com/941112341/avalon/gateway/service"
	"github.com/941112341/avalon/pkg/mygorm"
)

func init() {
	_ = registry.Registry("MapperService", &DBMapperService{})
}

type DBMapperService struct {
	Repo repository.MapperRepository `inject:"MapperRepository"`
}

func (D *DBMapperService) AddMapperRule(mapper *service.MapperData) error {
	if mapper == nil {
		return errors.New("mapper is nil")
	}
	return D.Repo.AddMapper(repository.MapperList{repository.MapperVo{
		Model:   mygorm.Model{},
		URL:     mapper.URL,
		Type:    mapper.Type,
		Domain:  mapper.Domain,
		PSM:     mapper.PSM,
		Base:    mapper.Base,
		Method:  mapper.Method,
		Version: mapper.Version,
	}})
}

func (D *DBMapperService) FetchMapperList() (repository.MapperList, error) {
	return D.Repo.AllMapper()
}
