package service

import (
	"context"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/model"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/pkg/errors"
)

func init() {
	_ = registry.Registry("", &genIdsService{})
}

type GenIdsService interface {
	GenIDs(ctx context.Context, request *idgenerator.IDRequest) (r *idgenerator.IDResponse, err error)
}

type genIdsService struct {
	F model.GeneratorFactory `inject:""`
}

func (g genIdsService) GenIDs(ctx context.Context, request *idgenerator.IDRequest) (r *idgenerator.IDResponse, err error) {
	generator, err := g.F.Create()
	if err != nil {
		return nil, errors.Wrap(err, "get generator error")
	}
	return &idgenerator.IDResponse{IDs: generator.GetIds()}, nil
}