package main

import (
	"context"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/example/idgenerator/service"
)

var handler Handler

func init() {
	_ = registry.Registry("", &handler)
}

type Handler struct {
	GenIdsService service.GenIdsService `inject:"GenIdsService"`
}

func (h Handler) GenIDs(ctx context.Context, request *idgenerator.IDRequest) (r *idgenerator.IDResponse, err error) {
	return h.GenIdsService.GenIDs(ctx, request)
}
