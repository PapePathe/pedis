package storage

type SimpleStorage struct {
	data map[string]string
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		data: make(map[string]string),
	}
}

func (ss *SimpleStorage) Set(key string, value string) error {
	ss.data[key] = value
	return nil
}

func (ss *SimpleStorage) Get(key string) (string, error) {
	return ss.data[key], nil
}
