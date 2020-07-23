package main

import (
	"context"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/example/idgenerator/service"
)

func init() {
	_ = registry.Registry("", &Handler{})
}

type Handler struct {
	GenIdsService service.GenIdsService `inject:""`
}

func (h Handler) GenIDs(ctx context.Context, request *idgenerator.IDRequest) (r *idgenerator.IDResponse, err error) {
	return h.GenIdsService.GenIDs(ctx, request)
}
