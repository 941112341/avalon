package main

import (
	"github.com/941112341/avalon/gateway/service"
	"github.com/941112341/avalon/sdk/inline"
	"net/http"
)

type BaseResp struct {
	Code    int32
	Message string
	Data    interface{}
}

func Transfer(writer http.ResponseWriter, request *http.Request) {

	errResponse := &BaseResp{
		Code:    -1,
		Message: "",
		Data:    nil,
	}
	defer func() {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err := writer.Write(inline.ToJsonBytes(errResponse)); err != nil {
			panic(err)
		}
	}()
	data, err := service.TransferServiceInstance.Transfer(request)
	if err != nil {
		errResponse.Message = err.Error()
		return
	}
	errResponse.Data = data.Data
	for k, v := range data.Header {
		writer.Header().Set(k, v)
	}

	return
}
