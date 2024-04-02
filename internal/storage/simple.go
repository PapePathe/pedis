package storage

import (
	"errors"
	"sync"
	"time"
)

type datatype struct {
	t    rune
	data []byte
}

type SimpleStorage struct {
	data    map[string]datatype
	exp     map[string]time.Time
	expLock sync.RWMutex
	sync.RWMutex
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		data: make(map[string]datatype),
		exp:  make(map[string]time.Time),
	}
}

func (ss *SimpleStorage) Set(key string, value string, expires int64) error {
	ss.Lock()
	ss.data[key] = datatype{t: 's', data: []byte(value)}
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

	v, ok := ss.data[key]

	if !ok {
		return "", errors.New("key not found")
	}

	return string(v.data), nil
}

func (ss *SimpleStorage) HSet(key string, value []byte, expires int64) (int, error) {
	data := datatype{t: 'm', data: value}

	ss.Lock()
	ss.data[key] = data
	ss.Unlock()

	return 0, nil
}

func (ss *SimpleStorage) HGet(key string) ([]byte, error) {
	return []byte{}, nil
}
