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
	"net/http"
)

type ExecutorFactory struct {
	ExecutorData
	MapperArgs map[string]interface{}
}

func (e *ExecutorFactory) GetApplication() model.Application {
	return &Executor{
		converter: &BaseConverter{
			mapperArgs: e.MapperArgs,
		},
		ExecutorData: e.ExecutorData,
	}
}

type ExecutorData struct {
	psm     string
	version string
	method  string
	base    string
}

type Executor struct {
	converter model.Converter
	ExecutorData
}

func (e *Executor) Invoker(ctx context.Context, request *http.Request) (resp *model.HttpResponse, err error) {
	defer func() {
		r, ok := recover().(error)
		if ok {
			err = r
			inline.WithFields("recover", r).Errorln("panic !!")
		}
	}()

	value, err := e.converter.ConvertRequest(request)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "convert request %+v", request)
	}

	data := inline.ToJsonString(value)

	idlMap, err := service.ServiceContainer.UploadService.GroupContent(repository.UploadGroupKey{
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
