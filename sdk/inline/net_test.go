package inline

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {

	s := "https://www.jiangshihao.cn/api/registry?color=blue"
	url, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(url.Host)
	fmt.Println(url.Path)
	fmt.Println(url.RawPath)
	fmt.Println(url.RawQuery)
}
