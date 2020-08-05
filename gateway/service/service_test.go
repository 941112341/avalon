package service

import (
	"fmt"
	"github.com/941112341/avalon/gateway/initial"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/sdk/inline"
	"testing"
)

type ServiceTest struct {
	Uploader UploadService `inject:"CacheUploadService"`
}

var service ServiceTest

func init() {
	_ = registry.Registry("", &service)
}

func TestUpload(t *testing.T) {
	var err error
	err = initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	err = service.Uploader.Upload(&UploadVoVo{
		PSM:      "a.b.c",
		Filename: "test",
		Body:     "hello world",
	})
	if err != nil {
		panic(err)
	}

	err = service.Uploader.Upload(&UploadVoVo{
		PSM:      "a.b.c",
		Filename: "test",
		Body:     "hello world 3333",
	})
	if err != nil {
		panic(err)
	}

}

func TestGet(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}
	vo, err := service.Uploader.Get(&UploadVoVo{
		PSM:      "a.b.c",
		Filename: "test",
		Body:     "",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(vo))
}
