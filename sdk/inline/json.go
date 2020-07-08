package inline

import (
	"fmt"
	"github.com/json-iterator/go"
)

func JsonString(o interface{}) string {
	body, _ := jsoniter.MarshalToString(o)
	return body
}

func VString(o interface{}) string {
	return fmt.Sprintf("%+v", o)
}
