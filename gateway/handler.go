package main

import (
	"context"
	"errors"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/gateway/service"
	"github.com/941112341/avalon/sdk/inline"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type Handler interface {
	Upload(request *http.Request) (*Response, error)
	Test() (*Response, error)
	Registry(request *http.Request) (*Response, error)
	Transfer(request *http.Request) (*Response, error)
}

type Response struct {
	Code    int
	Message string
}

func (r *Response) write(writer http.ResponseWriter, err error) error {
	if r == nil {
		if err == nil {
			return errors.New("err && response all nil")
		}
		r = &Response{
			Code:    500,
			Message: err.Error(),
		}
	}

	_, err = writer.Write(inline.ToJsonBytes(r))
	return err
}

var handler = &DefaultHandler{}

func init() {
	_ = registry.Registry("", handler)
}

type DefaultHandler struct {
	Gateway model.Gateway `inject:"Gateway"`
}

func (d *DefaultHandler) Upload(request *http.Request) (*Response, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	r := service.SaveGroupContentRequest{}
	err = jsoniter.Unmarshal(body, &r)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "unmarshal %s", body)
	}
	err = d.Gateway.AddUploader(context.Background(), &r)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "add uploader %+v", r)
	}
	return &Response{Message: "success"}, nil
}

func (d *DefaultHandler) Test() (*Response, error) {
	return &Response{Message: "hello world"}, nil
}

func (d *DefaultHandler) Registry(request *http.Request) (*Response, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	r := service.MapperData{}
	err = jsoniter.Unmarshal(body, &r)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "unmarshal %s", body)
	}
	err = d.Gateway.AddMapper(context.Background(), &r)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "add mapper %+v", r)
	}
	return &Response{Message: "success"}, nil
}

func (d *DefaultHandler) Transfer(request *http.Request) (*Response, error) {
	response, err := d.Gateway.Transfer(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return &Response{
		Code:    response.HTTPCode,
		Message: response.Body,
	}, nil
}
