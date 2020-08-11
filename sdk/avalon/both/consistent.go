package both

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon/util"
	"github.com/941112341/avalon/sdk/inline"
	"time"
)

type Base struct {
	Psm   string            `thrift:"psm,1" db:"psm" json:"psm"`
	IP    string            `thrift:"ip,2" db:"ip" json:"ip"`
	Time  int64             `thrift:"time,3" db:"time" json:"time"`
	Extra map[string]string `thrift:"extra,4" db:"extra" json:"extra"`
	Base  *Base             `thrift:"base,5" db:"base" json:"base,omitempty"`
}

func (b *Base) Wrap(other *Base) *Base {
	return &Base{
		Psm:   b.Psm,
		IP:    b.IP,
		Time:  time.Now().Unix(),
		Extra: b.Extra,
		Base: &Base{
			Psm:   other.Psm,
			IP:    other.IP,
			Time:  time.Now().Unix(),
			Extra: other.Extra,
			Base:  other.Base,
		},
	}
}

const BaseKey = "__Base__"

func GetBase(ctx context.Context) *Base {
	base, _ := util.GetValue(ctx, BaseKey).(*Base)
	return base
}

func SetBase(ctx context.Context, base *Base) context.Context {
	b := GetBase(ctx)
	if b == nil {
		return util.SetValue(ctx, BaseKey, base)
	} else {
		return util.SetValue(ctx, BaseKey, b.Wrap(base))
	}
}

func SetConsistentValue(ctx context.Context, key string, value interface{}) context.Context {
	base := GetBase(ctx)
	if base == nil {
		base = &Base{
			Psm:   "",
			IP:    "",
			Time:  0,
			Extra: map[string]string{},
			Base:  nil,
		}
	}
	base.Extra[key] = inline.ToJsonString(value)
	return SetBase(ctx, base)
}

func GetConsistentValue(ctx context.Context, key string, value interface{}) {
	base := GetBase(ctx)
	if base == nil {
		return
	}
	vstring, ok := base.Extra[key]
	for !ok && base != nil {
		vstring, ok = base.Base.Extra[key]
		base = base.Base
	}
	inline.MustUnmarshal(vstring, value)
}
