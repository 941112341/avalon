package server

import (
	"errors"
	"fmt"
	"testing"
)

type TestResponseS struct {
	BaseResp *BaseResp
}

func TestResponse(t *testing.T) {
	var s *TestResponseS
	err := errors.New("not ex")
	fmt.Println(ConvertErr2Resp(s, err))

	s = &TestResponseS{BaseResp: &BaseResp{
		Code:    10,
		Message: "",
	}}
	fmt.Println(ConvertErr2Resp(s, nil))
}
