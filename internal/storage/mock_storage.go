package storage

type MockStorage struct {
	GetUserFunc func(string) (*User, error)
	SetUserFunc func(string, []AclRule) error
	UsersFunc   func() []string
	SetFunc     func(string, string, int64) error
}

func (s MockStorage) GetUser(key string) (*User, error) {
	return nil, nil
}

func (ms MockStorage) SetUser(key string, rule []AclRule) error {
	if ms.SetUserFunc == nil {
		panic("mock caller must provide implementation of SetUserFunc")
	}
	return ms.SetUserFunc(key, rule)
}

func (ms MockStorage) DelUser(key string) error {
	return nil
}

func (ms MockStorage) Users() []string {
	if ms.UsersFunc == nil {
		panic("mock caller must provide implementation of UsersFunc")
	}
	return ms.UsersFunc()
}

func (ms MockStorage) Set(key string, value string, expires int64) error {
	return nil
}

func (ms MockStorage) Get(key string) (string, error) {
	return "", nil
}

func (ms MockStorage) Del(key string) error {
	return nil
}

func (ms MockStorage) HGet(key string) ([]byte, error) {
	return nil, nil
}

func (ms MockStorage) HSet(key string, value []byte, expires int64) (int, error) {
	return 0, nil
}

func (ms MockStorage) WriteFromRaft(StorageData) error {
	return nil
}
