package tool

import (
	"fmt"
	"testing"
)

const data = `
	namespace go base

include "base.thrift"

struct MessageRequest {
    1: map<string, string> header
    2: binary body
    3: string methodName
    4: string url

    255: base.Base base
}

struct MessageResponse {
    1: map<string, string> header
    2: binary body
    3: i32 status

    255: base.BaseResp baseResp
}

service messageService {
    MessageResponse MessageDispatcher(1: MessageRequest request)
}

`

func TestParse(t *testing.T) {
	scanner := NewScanner("message.thrift")
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
