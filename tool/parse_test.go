package tool

import (
	"fmt"
	"testing"
)

const data = `
namespace go idgenerator

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

}

`

func TestParse(t *testing.T) {
	scanner := NewScanner("idgenerator.thrift")
	info, err := scanner.parse(data)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(inline.ToJsonString(info))

	str, err := build(*info, "1.0.0")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
}
