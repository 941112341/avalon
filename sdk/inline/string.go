package inline

import (
	"fmt"
	"github.com/json-iterator/go"
	"strings"
)

func ToJsonString(o interface{}) string {
	body, _ := jsoniter.MarshalToString(o)
	return body
}

func VString(o interface{}) string {
	return fmt.Sprintf("%+v", o)
}

func JoinPath(paths ...string) string {
	return strings.Join(paths, "/")
}

func String(o interface{}) string {
	return fmt.Sprintf("%s", o)
}
