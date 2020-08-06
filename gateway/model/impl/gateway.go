package impl

import (
	"context"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/gateway/service"
	"github.com/941112341/avalon/sdk/inline"
	"net/http"
)

func init() {
	_ = registry.Registry("Gateway", &IGateway{})
}

// singleton
type IGateway struct {
	Mappers model.MapperRules
}

func (I *IGateway) AddMapper(ctx context.Context, request *service.MapperData) error {
	err := service.ServiceContainer.MapperService.AddMapperRule(request)
	if err == nil {
		I.ClearRules()
	}
	return err
}

func (I *IGateway) AddUploader(ctx context.Context, request *service.SaveGroupContentRequest) error {
	return service.ServiceContainer.UploadService.SaveGroupContent(request)
}

func (I *IGateway) GetMapperRules() (model.MapperRules, error) {
	rules, err := service.ServiceContainer.MapperService.FetchMapperList()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "fetch rules fail")
	}

	mapperRules := Ruler{}
	for _, rule := range rules {
		mapperRules.Rules = append(mapperRules.Rules, NewRule(rule))
	}
	return &mapperRules, nil
}

func (I *IGateway) ClearRules() {
	I.Mappers = nil
}

func (I *IGateway) Transfer(ctx context.Context, request *http.Request) (*model.HttpResponse, error) {
	if I.Mappers == nil {
		rules, err := I.GetMapperRules()
		if err != nil {
			return nil, inline.PrependErrorFmt(err, "get mapper err")
		}
		I.Mappers = rules
	}

	rules := I.Mappers
	application, err := rules.GetApplication(request)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "get application %+v", request)
	}
	response, err := application.Invoker(ctx, request)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "invoker request %+v", request)
	}
	return response, nil
}
