package main

import (
	"fmt"
	"github.com/941112341/avalon/gateway/conf"
	"github.com/941112341/avalon/gateway/initial"
	_ "github.com/941112341/avalon/gateway/model/impl"
	_ "github.com/941112341/avalon/gateway/service/impl"
	"github.com/941112341/avalon/sdk/inline"
	"net/http"
)

func main() {
	var err error
	err = initial.InitAll()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		resp, err := handler.Test()
		if err = resp.write(writer, err); err != nil {
			inline.WithFields("response", resp, "request", request).Errorln("test err %s", err)
		}
	})
	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		resp, err := handler.Upload(request)
		if err = resp.write(writer, err); err != nil {
			inline.WithFields("response", resp, "request", request).Errorln("upload %s", err)
		}
	})
	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		resp, err := handler.Transfer(request)
		if err = resp.write(writer, err); err != nil {
			inline.WithFields("response", resp, "request", request).Errorln("transfer %s", err)
		}
	})
	http.HandleFunc("/registry", func(writer http.ResponseWriter, request *http.Request) {
		resp, err := handler.Registry(request)
		if err = resp.write(writer, err); err != nil {
			inline.WithFields("response", resp, "request", request).Errorln("registry %s", err)
		}
	})
	//err = http.ListenAndServeTLS(fmt.Sprintf(":%d", conf.Config.Https.Port), "resource/1_www.jiangshihao.cn_bundle.crt", "resource/2_www.jiangshihao.cn.key", nil)
	//if err != nil {
	//	panic(err)
	//}
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Config.Http.Port), nil)
	if err != nil {
		panic(err)
	}
}
