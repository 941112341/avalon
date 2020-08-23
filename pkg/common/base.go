package common

import (
	"github.com/941112341/avalon/common/gen/base"
	"github.com/941112341/avalon/sdk/inline"
	jsoniter "github.com/json-iterator/go"
)

func GetValue(base *base.Base, key string, o interface{}) error {
	value := GetExtraValue(base, key)
	return jsoniter.UnmarshalFromString(value, o)
}

func GetExtraValue(base *base.Base, key string) string {
	if base == nil {
		return ""
	}

	value, ok := base.Extra[key]
	if ok {
		return value
	}
	return GetExtraValue(base.Base, key)
}

func SetValue(base *base.Base, key string, o interface{}) {
	base.Extra[key] = inline.ToJsonString(o)
}
