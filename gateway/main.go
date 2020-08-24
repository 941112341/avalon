package main

import (
	"fmt"
	"github.com/941112341/avalon/gateway/initial"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	var err error
	err = initial.InitAll()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/", Transfer)

	port := viper.GetInt("conf.port")
	hostport := fmt.Sprintf(":%d", port)

	env := inline.GetEnv("env", "dev")
	if env == "dev" {
		err = http.ListenAndServe(hostport, nil)
		if err != nil {
			panic(err)
		}
	} else {
		err = http.ListenAndServeTLS(hostport, "resource/1_www.jiangshihao.cn_bundle.crt", "resource/2_www.jiangshihao.cn.key", nil)
		if err != nil {
			panic(err)
		}
	}

}
