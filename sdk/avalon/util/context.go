package util

import "context"

const ConsistentPrefix = "__Avalon__"

func SetValue(ctx context.Context, key string, value interface{}) context.Context {
	cache, ok := ctx.Value(ConsistentPrefix).(map[string]interface{})
	if !ok {
		ctx = context.WithValue(ctx, ConsistentPrefix, make(map[string]interface{}))
		return SetValue(ctx, key, value)
	}
	cache[key] = value
	return ctx
}

func GetValue(ctx context.Context, key string) interface{} {
	cache, ok := ctx.Value(ConsistentPrefix).(map[string]interface{})
	if !ok {
		return nil
	}
	return cache[key]
}
