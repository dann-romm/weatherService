package cache

import "sync"

type MapCache struct {
	Storage map[string]interface{}
	mutex   sync.RWMutex
}

func NewMapCache() *MapCache {
	return &MapCache{
		Storage: make(map[string]interface{}),
	}
}

func (m *MapCache) Set(key string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Storage[key] = value
	return nil
}

func (m *MapCache) Get(key string) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Storage[key]

	if ok {
		return value, nil
	} else {
		return nil, ErrorNotFound
	}
}

func (m *MapCache) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.Storage, key)
	return nil
}

func (m *MapCache) Clear() error {
	for key := range m.Storage {
		delete(m.Storage, key)
	}
	return nil
}
