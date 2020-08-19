package client

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/common/gen/idgenerator"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/generic/invoker"
	"github.com/941112341/avalon/sdk/inline"
	"os"
	"testing"
)

const args = `
	{
		"request": {
			"count": 10,
			"base": {
				"psm":"a.b.c",
				"ip":"localhost",
				"time":"10000000",
				"extra": {
					"hello": "world"
				 }
			}
		}
	}

`

var IDLMap = map[string]string{
	"base": `namespace go base


struct Base {
    1: string psm
    2: string ip
    3: i64 time
    4: map<string, string> extra
    5: optional Base base
}

struct BaseResp {
    1: i32 code
    2: string message
}`,
	"idgenerator": `namespace go idgenerator

include "base.thrift"

struct IDRequest {
    1: i32 count

    255: base.Base base
}

struct IDResponse {
    1: list<i64> IDs

    255: base.BaseResp baseResp
}

service IDGenerator {
    IDResponse GenIDs(1: IDRequest request)

}`,
}

func TestCacheIDClient_GenID(t *testing.T) {

	os.Setenv("base", "../../example/idgenerator/base.yaml")
	cli := NewCacheIDClient(5)
	for i := 0; i < 1; i++ {
		fmt.Println(cli.GenID())
	}
}


func TestInvoker(t *testing.T) {

	grp, err := generic.NewThriftGroup(IDLMap)
	if err != nil {
		panic(err)
	}

	invoker, err := invoker.CreateInvoker(grp, "idgenerator", "IDGenerator", "GenIDs")
	if err != nil {
		panic(err)
	}
	input := inline.JSONAny(args)
	data, err := invoker.Invoke(context.Background(), client.Client_(), input)
	if err != nil {
	    panic(err)
	}
	fmt.Println(inline.ToJsonString(data))
}



func TestReflectInvoker(t *testing.T) {

	invoker, err := invoker.ReflectInvoker("GenIDs", &idgenerator.IDGeneratorGenIDsArgs{}, &idgenerator.IDGeneratorGenIDsResult{}, )
	if err != nil {
		panic(err)
	}
	input := inline.JSONAny(args)
	data, err := invoker.Invoke(context.Background(), client.Client_(), input)
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(data))
}

