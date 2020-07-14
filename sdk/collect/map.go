package collect

import (
	"github.com/941112341/avalon/sdk/inline"
	"sync"
)

type Map interface {
	Put(key, value interface{}) interface{}
	Get(key interface{}) (interface{}, bool)
	Delete(key interface{}) bool
	Contains(key interface{}) bool
	Range(func(key, value interface{}) bool)
}

type LenMap interface {
	Map
	Length() int
}

type MapList interface {
	Map
	GetList(key interface{}) []interface{}
	Append(key, value interface{}) bool
	Remove(value interface{}) bool
	Len(key interface{}) int
}

type SyncMap struct {
	lock sync.RWMutex

	m map[interface{}]interface{}
}

func (s *SyncMap) Length() int {
	return len(s.m)
}

func (s *SyncMap) Get(key interface{}) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.m[key]
	return v, ok
}

func (s *SyncMap) GetList(key interface{}) []interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	i, ok := s.Get(key)
	if !ok {
		return []interface{}{}
	}
	list, ok := i.([]interface{})
	if !ok {
		return nil
	}
	return list
}

func (s *SyncMap) Put(key, value interface{}) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	v := s.m[key]
	s.m[key] = value
	return v
}

func (s *SyncMap) Append(key, value interface{}) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	v, ok := s.m[key]
	if !ok {
		v = make([]interface{}, 0)
	}
	list, ok := v.([]interface{})
	if !ok {
		return false
	}
	s.m[key] = append(list, value)
	return true
}

func (s *SyncMap) Remove(value interface{}) (r bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for k, v := range s.m {
		if v == value {
			delete(s.m, k)
			r = true
		}
	}
	return
}

func (s *SyncMap) Len(key interface{}) int {
	list := s.GetList(key)
	if list == nil {
		return -1
	}
	return len(list)
}

func (s *SyncMap) Delete(key interface{}) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.m[key]
	if !ok {
		return false
	}
	delete(s.m, key)
	return true
}

func (s *SyncMap) Contains(key interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.m[key]
	return ok
}

func (s *SyncMap) Range(f func(key, value interface{}) bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for k, v := range s.m {
		if !f(k, v) {
			break
		}
	}
}

func (s *SyncMap) GetString(key interface{}) string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	any, ok := s.m[key]
	if !ok {
		return ""
	}
	str, ok := any.(string)
	if !ok {
		return ""
	}
	return str
}

func (s *SyncMap) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return inline.ToJsonString(s.m)
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		lock: sync.RWMutex{},
		m:    map[interface{}]interface{}{},
	}
}
