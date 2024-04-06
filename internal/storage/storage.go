package storage

type Storage interface {
	Set(key string, value string, expires int64) error
	Get(key string) (string, error)

	// Raft related methods

	WriteFromRaft(StorageData) error
}
