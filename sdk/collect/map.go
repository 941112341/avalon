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
	Range(func(key, value interface{}))
}

type SyncMap struct {
	lock sync.RWMutex

	m map[interface{}]interface{}
}

func (s *SyncMap) Get(key interface{}) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.m[key]
	return v, ok
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

func (s *SyncMap) Range(f func(key, value interface{})) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for k, v := range s.m {
		f(k, v)
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
