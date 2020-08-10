package _map

import "github.com/941112341/avalon/sdk/collect"

type LimitMap interface {
	collect.MapList
	Limit() int
	SetLimit(limit int)
	GetOrSet(key, value interface{}) interface{}
}

// not thread safe
type BaseLimitMap struct {
	mapList collect.MapList

	limit int
}

func (b *BaseLimitMap) Put(key, value interface{}) interface{} {
	return b.mapList.Put(key, value)
}

func (b *BaseLimitMap) Get(key interface{}) (interface{}, bool) {
	return b.mapList.Get(key)
}

func (b *BaseLimitMap) Delete(key interface{}) bool {
	return b.mapList.Delete(key)
}

func (b *BaseLimitMap) Contains(key interface{}) bool {
	return b.mapList.Contains(key)
}

func (b *BaseLimitMap) Range(f func(key interface{}, value interface{}) bool) {
	b.mapList.Range(f)
}

func (b *BaseLimitMap) GetList(key interface{}) []interface{} {
	list := b.mapList.GetList(key)
	if len(list) > b.limit {
		list = list[:b.limit]
	}
	return list
}

func (b *BaseLimitMap) Remove(value interface{}) bool {
	return b.mapList.Remove(value)
}

func (b *BaseLimitMap) Len(key interface{}) int {
	return b.mapList.Len(key)
}

func (b *BaseLimitMap) SetLimit(limit int) {
	b.limit = limit
}

func (b *BaseLimitMap) Limit() int {
	return b.limit
}

func (b *BaseLimitMap) Append(key, value interface{}) bool {
	l := b.Len(key)
	if l > b.limit {
		list := b.mapList.GetList(key)
		list = append(list[1:], value)
		b.Put(key, value)
		return true
	}
	return b.mapList.Append(key, value)
}

func (b *BaseLimitMap) GetOrSet(key, value interface{}) interface{} {
	old, ok := b.Get(key)
	if !ok {

		b.Put(key, value)
		return value
	}
	return old
}

func NewLimitMap(limit int) LimitMap {
	return &BaseLimitMap{limit: limit, mapList: collect.NewSyncMap()}
}
