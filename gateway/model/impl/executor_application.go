package impl

import (
	"context"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/gateway/repository"
	"github.com/941112341/avalon/gateway/service"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/generic/invoker"
	"github.com/941112341/avalon/sdk/inline"
)

type Executor struct {
	converter model.Converter
	psm       string
	version   string
	method    string
	base      string
}

func (e *Executor) Invoker(ctx context.Context, request *model.HttpRequest) (*model.HttpResponse, error) {
	value, err := e.converter.ConvertRequest(request)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "convert request %+v", request)
	}

	data := inline.ToJsonString(value)

	idlMap, err := service.UploadBuilder.UploadService.GroupContent(repository.UploadGroupKey{
		PSM:     e.psm,
		Version: e.version,
	})
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "group content err")
	}

	grp, err := generic.NewThriftGroup(idlMap)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "parse idl map %+v", idlMap)
	}
	i, err := invoker.CreateInvoker(grp, e.base, "", e.method)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "create invoker %+v", e)
	}
	result, err := i.Invoke(ctx, avalon.NewClient(e.psm), data)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "invoke %s", data)
	}

	return e.converter.ConvertResponse(result)
}
