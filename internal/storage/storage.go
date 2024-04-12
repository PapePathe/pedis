package storage

type AclRule struct {
	Command    string
	KeyPattern string
}
type User struct {
	Passwords []string
	Rules     []AclRule
	Active    bool
}

type Storage interface {
	// Users
	GetUser(key string) (*User, error)
	SetUser(key string, rule []AclRule) error
	DelUser(key string) error
	Users() []string

	// Simple strings
	Set(key string, value string, expires int64) error
	Get(key string) (string, error)
	Del(key string) error

	// Maps
	HGet(key string) ([]byte, error)
	HSet(key string, value []byte, expires int64) (int, error)

	// Raft
	WriteFromRaft(StorageData) error
}

type StorageOperation rune

const (
	StorageOperationCreate StorageOperation = 'c'
	StorageOperationDelete StorageOperation = 'd'
)

type StorageDataType rune

const (
	StorageDataTypeMap    StorageDataType = 'm'
	StorageDataTypeString StorageDataType = 's'
	StorageDataTypeList   StorageDataType = 'l'
	StorageDataTypeJson   StorageDataType = 'j'
)

type StorageData struct {
	D  []byte
	K  string
	Op StorageOperation
	T  StorageDataType
}

type StorageDataInternal struct {
	D []byte
	T StorageDataType
}
