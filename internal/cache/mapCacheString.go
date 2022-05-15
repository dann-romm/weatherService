package cache

import "sync"

type MapCacheString struct {
	Storage map[string]string
	mutex   sync.RWMutex
}

func (m *MapCacheString) Set(key, value string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Storage[key] = value
	return nil
}

func (m *MapCacheString) Get(key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Storage[key]

	if ok {
		return value, nil
	} else {
		return "", ErrorNotFound
	}
}

func (m *MapCacheString) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.Storage, key)
	return nil
}

func (m *MapCacheString) Clear() error {
	for key := range m.Storage {
		delete(m.Storage, key)
	}
	return nil
}
