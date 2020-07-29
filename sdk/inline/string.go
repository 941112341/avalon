package inline

import (
	"fmt"
	"github.com/json-iterator/go"
	"strings"
	"unicode"
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

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func IsEmpty(str string) bool {
	s := strings.Trim(str, " ")
	s = strings.Trim(str, "\t")
	return len(s) == 0
}
