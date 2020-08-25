package main

import (
	"github.com/941112341/avalon/gateway/service"
	"github.com/941112341/avalon/sdk/inline"
	"net/http"
)

func Transfer(writer http.ResponseWriter, request *http.Request) {

	errResponse := &service.BaseResp{
		Code:    -1,
		Message: "",
		Data:    nil,
	}
	defer func() {
		inline.WithFields("request", request, "response", errResponse, "path", request.URL.Path).Infoln("finish")
		writer.Write(inline.ToJsonBytes(errResponse))
	}()
	data, err := service.TransferServiceInstance.Transfer(request)
	if err != nil {
		errResponse.Message = err.Error()
		return
	}
	errResponse = data.BaseResp
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	for k, v := range data.Header {
		writer.Header().Set(k, v)
	}
	return
}
