package main

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/example/idgenerator/service"
	"github.com/941112341/avalon/sdk/inline"
)

var handler Handler

func init() {
	_ = registry.Registry("", &handler)
}

type Handler struct {
	GenIdsService service.GenIdsService `inject:"GenIdsService"`
}

func (h Handler) GenIDs(ctx context.Context, request *idgenerator.IDRequest) (r *idgenerator.IDResponse, err error) {
	fmt.Println("receive a request " + inline.ToJsonString(request))
	return h.GenIdsService.GenIDs(ctx, request)
}
