package main

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/gateway/conf"
	"net/http"
)

func main() {
	var err error
	ctx := context.Background()
	err = conf.InitConfig(ctx)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})
	err = http.ListenAndServeTLS(fmt.Sprintf(":%d", conf.Config.Https.Port), "resource/1_www.jiangshihao.cn_bundle.crt", "resource/2_www.jiangshihao.cn.key", nil)
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Config.Http.Port), nil)
	if err != nil {
		panic(err)
	}
}
