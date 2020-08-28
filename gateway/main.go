package main

import (
	"fmt"
	"github.com/941112341/avalon/gateway/initial"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"net/http"
	"os"
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

	// sentry 激活暂时不写到配置文件里 等待重构
	if os.Getenv("env") == "online" {
		err := sentry.Init(sentry.ClientOptions{
			// Either set your DSN here or set the SENTRY_DSN environment variable.
			Dsn: "http://77c1184d08904731813a00071efa1358@book.jiangshihao.cn:9000/1",
			// Enable printing of SDK debug messages.
			// Useful when getting started or trying to figure something out.
			Debug: false,
		})
		if err != nil {
			panic(err)
		}

	}

	env := inline.GetEnv("env", "dev")
	if env == "dev" {
		err = http.ListenAndServe(hostport, nil)
		if err != nil {
			panic(err)
		}
	} else {
		err = http.ListenAndServeTLS(hostport, "resource/1_book.jiangshihao.cn_bundle.crt", "resource/2_book.jiangshihao.cn.key", nil)
		if err != nil {
			panic(err)
		}
	}

}
