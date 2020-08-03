package invoker

import (
	"fmt"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"testing"
)

func TestInvoker(t *testing.T) {

	grp, err := generic.NewThriftGroup(map[string]string{
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
	})
	if err != nil {
	    panic(err)
	}

	invoker, err := CreateInvoker(grp, "idgenerator", "IDGenerator", "GenIDs")
	if err != nil {
	    panic(err)
	}
	fmt.Println(inline.ToJsonString(invoker))

}
