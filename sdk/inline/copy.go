package inline

import "github.com/json-iterator/go"

func Copy(src, destPtr interface{}) error {
	s, _ := jsoniter.MarshalToString(src)
	return jsoniter.UnmarshalFromString(s, destPtr)
}

func MustCopy(src, destPtr interface{}) {
	_ = Copy(src, destPtr)
}
