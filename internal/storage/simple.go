package storage

import (
	"errors"
	"sync"
	"time"
)

type StorageData struct {
	K string
	T rune
	D []byte
}

type SimpleStorage struct {
	data        map[string]StorageData
	exp         map[string]time.Time
	proposeChan chan StorageData
	expLock     sync.RWMutex
	sync.RWMutex
}

func NewSimpleStorage(proposeChan chan StorageData) *SimpleStorage {
	return &SimpleStorage{
		data:        make(map[string]StorageData),
		exp:         make(map[string]time.Time),
		proposeChan: proposeChan,
	}
}

// Endpoint for redis SET command
func (ss *SimpleStorage) Set(key string, value string, expires int64) error {
	go ss.proposeRaftChange(StorageData{T: 's', D: []byte(value), K: key})
	ss.Lock()
	ss.data[key] = StorageData{T: 's', D: []byte(value)}
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

	return string(v.D), nil
}

func (ss *SimpleStorage) WriteFromRaft(d StorageData) error {
	ss.Lock()
	ss.data[d.K] = StorageData{T: d.T, D: []byte(d.D)}
	ss.Unlock()

	return nil
}

func (ss *SimpleStorage) proposeRaftChange(data StorageData) {
	ss.proposeChan <- data
}
