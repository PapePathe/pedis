package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type SimpleStorage struct {
	data map[string]StorageDataInternal
	sync.RWMutex
	acl         map[string]User
	aclLock     sync.RWMutex
	exp         map[string]time.Time
	expLock     sync.RWMutex
	proposeChan chan StorageData
}

func NewSimpleStorage(proposeChan chan StorageData) *SimpleStorage {
	return &SimpleStorage{
		data:        make(map[string]StorageDataInternal),
		acl:         make(map[string]User),
		exp:         make(map[string]time.Time),
		proposeChan: proposeChan,
	}
}

func (ss *SimpleStorage) HSet(key string, value []byte, expires int64) (int, error) {
	go ss.proposeRaftChange(StorageData{T: StorageDataTypeMap, D: []byte(value), K: key, Op: StorageOperationCreate})
	data := StorageDataInternal{T: StorageDataTypeMap, D: value}

	ss.Lock()
	ss.data[key] = data
	ss.Unlock()

	return 0, nil
}

func (ss *SimpleStorage) HGet(key string) ([]byte, error) {
	ss.RLock()
	defer ss.RUnlock()
	val, ok := ss.data[key]

	if !ok {
		return nil, errors.New("hget key not found in store")
	}

	return val.D, nil
}

func (ss *SimpleStorage) GetUser(username string) (*User, error) {
	ss.aclLock.Lock()
	u, ok := ss.acl[username]
	ss.aclLock.Unlock()

	if !ok {
		return nil, fmt.Errorf("User %s not found in storage", username)
	}

	return &u, nil
}
func (ss *SimpleStorage) SetUser(username string, _ []AclRule) error {
	ss.aclLock.Lock()
	_, ok := ss.acl[username]
	ss.aclLock.Unlock()

	if !ok {
		ss.aclLock.Lock()
		ss.acl[username] = User{}
		ss.aclLock.Unlock()

		return nil
	}

	return nil
}

func (ss *SimpleStorage) DelUser(key string) error {
	ss.aclLock.Lock()
	_, ok := ss.acl[key]
	ss.aclLock.Unlock()

	if !ok {
		return fmt.Errorf("User (%s) Not found", key)
	}

	ss.aclLock.Lock()
	delete(ss.acl, key)
	ss.aclLock.Unlock()

	return nil
}

func (ss *SimpleStorage) Users() []string {
	users := []string{}

	ss.aclLock.RLock()
	for u, _ := range ss.acl {
		users = append(users, u)
	}
	ss.aclLock.RUnlock()

	return users
}

func (ss *SimpleStorage) Del(key string) error {
	_, err := ss.Get(key)

	if err != nil {
		return err
	}

	ss.Lock()
	delete(ss.data, key)
	ss.Unlock()

	return nil
}

// Endpoint for redis SET command
func (ss *SimpleStorage) Set(key string, value string, expires int64) error {
	go ss.proposeRaftChange(StorageData{
		T:  StorageDataTypeString,
		D:  []byte(value),
		K:  key,
		Op: StorageOperationCreate,
	})
	ss.Lock()
	ss.data[key] = StorageDataInternal{T: StorageDataTypeString, D: []byte(value)}
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
	switch d.Op {
	case StorageOperationCreate:
		ss.data[d.K] = StorageDataInternal{T: d.T, D: []byte(d.D)}
	case StorageOperationDelete:
		delete(ss.data, d.K)
	}
	ss.Unlock()

	return nil
}

func (ss *SimpleStorage) proposeRaftChange(data StorageData) {
	ss.proposeChan <- data
}
