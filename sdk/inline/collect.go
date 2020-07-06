package inline

import "sync"

type SaveMapList struct {
	lock sync.RWMutex

	m map[interface{}]interface{}
}

func (s *SaveMapList) Get(key interface{}) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.m[key]
	return v, ok
}

func (s *SaveMapList) Put(key, value interface{}) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	v := s.m[key]
	s.m[key] = value
	return v
}

func (s *SaveMapList) Append(key, value interface{}) bool {
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
	list = append(list, value)
	s.m[key] = list
	return true
}

type SyncMap struct {
	sync.Map
}

func (s *SyncMap) SetString(key interface{}, value string) {
	s.Store(key, value)
}

func (s *SyncMap) GetString(key interface{}) string {
	v, ok := s.Load(key)
	if !ok {
		return ""
	}
	vs, ok := v.(string)
	return vs
}

func (s *SyncMap) Contains(key interface{}) bool {
	_, ok := s.Load(key)
	return ok
}

type Tree interface {
	Data() interface{}
	Put(key, value interface{})
	Get(key interface{}) interface{}
}

type TreeImpl struct {
	o    interface{}
	maps map[interface{}]*TreeImpl
}

func (t *TreeImpl) Data() interface{} {
	return t.o
}

func (t *TreeImpl) Put(key, value interface{}) {
	t.maps[key] = NewTree(value)
}

func NewTree(value interface{}) *TreeImpl {
	return &TreeImpl{
		o:    value,
		maps: map[interface{}]*TreeImpl{},
	}
}

func (t *TreeImpl) Get(key interface{}) interface{} {
	return t.maps[key]
}
