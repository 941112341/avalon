package client

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/common/gen/test"
	"github.com/941112341/avalon/sdk/generic/invoker"
	"github.com/941112341/avalon/sdk/inline"
	"testing"
	"time"
)

func TestCall(t *testing.T) {

	i, err := invoker.FileInvoker([]string{"../idl/base.thrift", "../idl/test.thrift"}, "GetCat")
	if err != nil {
	    panic(err)
	}
	resp, err := i.Invoke(context.Background(), TestCli.Client_(), inline.JSONAny(&test.CatServiceGetCatArgs{Request: &test.CatRequest{
		ID:   []int64{1, 2, 3},
	}}))
	if err != nil {
	    panic(err)
	}
	fmt.Println(inline.ToJsonString(resp))

}

func TestCall2(t *testing.T) {
	i, err := invoker.FileInvoker([]string{"../idl/base.thrift", "../idl/test.thrift"}, "GetLittleCat")
	if err != nil {
		panic(err)
	}
	resp, err := i.Invoke(context.Background(), TestCli.Client_(), inline.JSONAny(&test.CatServiceGetLittleCatArgs{Request: &test.LittleCatRequest{
		Cat:  &test.Cat{
			Age:    19,
			Name:   inline.StringPtr("any"),
			Babies: []*test.LittleCat{
				{
					Cat:   nil,
					Age:   18,
					Color: 0,
					Ids: map[int64][]*test.Foo{
						1: {
							{Love: true},
						},
					},
				},
			},
		},
		Base: nil,
	}}))
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(resp))
}

func TestCall3(t *testing.T) {
	i, err := invoker.ReflectInvoker("GetCat", &test.CatServiceGetCatArgs{}, &test.CatServiceGetCatResult{})
	if err != nil {
	    panic(err)
	}

	resp, err := i.Invoke(context.Background(), TestCli.Client_(), inline.JSONAny(&test.CatServiceGetCatArgs{Request: &test.CatRequest{
		ID:   []int64{1, 2, 3},
	}}))
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(resp))
}

func TestCall4(t *testing.T) {
	i, err := invoker.ReflectInvoker("GetLittleCat", &test.CatServiceGetLittleCatArgs{}, &test.CatServiceGetLittleCatResult{})
	if err != nil {
		panic(err)
	}
	resp, err := i.Invoke(context.Background(), TestCli.Client_(), inline.JSONAny(&test.CatServiceGetLittleCatArgs{Request: &test.LittleCatRequest{
		Cat:  &test.Cat{
			Age:    19,
			Name:   inline.StringPtr("any"),
			Babies: []*test.LittleCat{
				{
					Cat:   nil,
					Age:   18,
					Color: 0,
					Ids: map[int64][]*test.Foo{
						1: {
							{Love: true},
						},
					},
				},
			},
		},
		Base: nil,
	}}))
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(resp))
}

func TestCall5(t *testing.T) {
	funcErr()
}

func funcName(idx int) {
	resp, err := TestCli.GetCat(context.Background(), &test.CatRequest{
		ID:   []int64{2,},
		Base: nil,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(inline.ToJsonString(resp))
}

func funcErr() {
	resp, err := TestCli.GetLittleCat(context.Background(), &test.LittleCatRequest{
		Cat:  nil,
		Base: nil,
	})
	fmt.Println(resp, err)
}

func TestCall6(t *testing.T) {
	for i := 0; i < 40; i++ {
		go funcName(i)
	}

	time.Sleep(10 * time.Second)
	for i := 0; i < 40; i++ {
		go funcName(i)
	}
	time.Sleep(10 * time.Second)
}