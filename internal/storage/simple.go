package storage

import "sync"

type SimpleStorage struct {
	data map[string]string
	sync.RWMutex
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		data: make(map[string]string),
	}
}

func (ss *SimpleStorage) Set(key string, value string) error {
	ss.Lock()
	ss.data[key] = value
	ss.Unlock()
	return nil
}

func (ss *SimpleStorage) Get(key string) (string, error) {
	ss.RLock()
	defer ss.RUnlock()
	return ss.data[key], nil
}
