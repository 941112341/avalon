package service

import (
	"context"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/model"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/pkg/errors"
)

func init() {
	_ = registry.Registry("GenIdsService", &genIdsService{})
}

type GenIdsService interface {
	GenIDs(ctx context.Context, request *idgenerator.IDRequest) (r *idgenerator.IDResponse, err error)
}

type genIdsService struct {
	F model.GeneratorFactory `inject:"GeneratorFactory"`
}

func (g genIdsService) GenIDs(ctx context.Context, request *idgenerator.IDRequest) (r *idgenerator.IDResponse, err error) {
	generator, err := g.F.Create()
	if err != nil {
		return nil, errors.Wrap(err, "get generator error")
	}
	ids, err := generator.Assign(int64(request.Count), request.Base.Psm)
	if err != nil {
		return nil, err
	}
	return &idgenerator.IDResponse{IDs: ids}, nil
}
