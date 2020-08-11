package server

import "context"

type Call func(ctx context.Context, request interface{}) (interface{}, error)

type Advance func(call Call) Call

func DefaultAdvance() []Advance {
	return []Advance{
		ConvertResponseAdvance,
		SetBase2Consistent,
		Log,
	}
}
