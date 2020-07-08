package inline

import jsoniter "github.com/json-iterator/go"

// dest need ptr
// src need not be Recursive
func Copy(src, destPtr interface{}) error {
	s, _ := jsoniter.MarshalToString(src)
	return jsoniter.UnmarshalFromString(s, destPtr)
}
