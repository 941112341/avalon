package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

// order 约小越先执行
var Initializers InitializersSorted

type Initializer struct {
	Order   int
	Name    string
	Initial func(ctx context.Context) error
}

type InitializersSorted []Initializer

func (sorted InitializersSorted) Len() int {
	return len(sorted)
}

func (sorted InitializersSorted) Swap(i, j int) {
	sorted[i], sorted[j] = sorted[j], sorted[i]
}

func (sorted InitializersSorted) Less(i, j int) bool {
	return sorted[i].Order < sorted[j].Order
}

func init() {

}

func main() {

	ctx := context.Background()
	sort.Sort(Initializers)
	for _, initializer := range Initializers {
		if err := initializer.Initial(ctx); err != nil {
			logrus.Errorf("initial err, task name %s, err %s", initializer.Name, err)
			return
		}
	}

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
