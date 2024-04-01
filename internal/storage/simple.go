package storage

import (
	"sync"
	"time"
)

type SimpleStorage struct {
	data    map[string]string
	exp     map[string]time.Time
	expLock sync.RWMutex
	sync.RWMutex
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		data: make(map[string]string),
		exp:  make(map[string]time.Time),
	}
}

func (ss *SimpleStorage) Set(key string, value string, expires int64) error {
	ss.Lock()
	ss.data[key] = value
	ss.Unlock()

	if expires > 0 {
		ss.expLock.Lock()
		ss.exp[key] = time.Now().Add(time.Duration(expires))
		ss.expLock.Unlock()
	}

	return nil
}

func (ss *SimpleStorage) Get(key string) (string, error) {
	ss.RLock()
	defer ss.RUnlock()
	return ss.data[key], nil
}
