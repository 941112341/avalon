package client

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/common/gen/test"
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/inline"
	"testing"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	ctx = both.SetConsistentValue(ctx, "test", "hello world")
	resp, err := TestCli.GetLittleCat(ctx, &test.LittleCatRequest{
		Cat:  &test.Cat{Name: inline.StringPtr("joker")},
		Base: nil,
	})
	if err != nil {
	    panic(err)
	}
	fmt.Println(inline.ToJsonString(resp))
}


func TestClient2(t *testing.T) {
	ctx := context.Background()
	resp, err := TestCli.GetCat(ctx, &test.CatRequest{
		ID:   []int64{11},
		Base: nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(resp))
}
