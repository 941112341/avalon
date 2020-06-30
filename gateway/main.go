package main

import "net/http"

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})
	err := http.ListenAndServeTLS(":443", "resources/1_www.jiangshihao.cn_bundle.crt", "resources/2_www.jiangshihao.cn.key", nil)
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
