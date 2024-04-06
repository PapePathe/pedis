package storage

type Storage interface {
	// Simple strings
	Set(key string, value string, expires int64) error
	Get(key string) (string, error)

	// Maps
	HGet(key string) ([]byte, error)
	HSet(key string, value []byte, expires int64) (int, error)

	// Raft
	WriteFromRaft(StorageData) error
}

type StorageData struct {
	K string
	T rune
	D []byte
}
