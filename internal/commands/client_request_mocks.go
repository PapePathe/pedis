package commands

import (
	"errors"
	"log"
	"pedis/internal/storage"

	"go.etcd.io/etcd/raft/v3/raftpb"
)

var (
	NotImplementedError error = errors.New("mock caller must provide implementation")
)

type MockClient struct {
	r        RawRequest
	d        [][]byte
	body        []string
	response []string
	errors   []string
  store  storage.Storage
}

func (mock MockClient) WriteError(e string) error {
	mock.errors = append(mock.errors, e)
	return nil
}

func (mock *MockClient) WriteString(s string) error {
	mock.response = []string{s} 
	return nil
}

func (mock *MockClient) WriteNumber(s string) error {
  log.Println("string", s)
	mock.response = []string{s} 
	return nil
}

func (mock *MockClient) WriteArray(s []string) error {
	mock.response = s
	return nil
}

func (mock *MockClient) WriteOK() error {
	mock.response = []string{"OK"}
	return nil
}

func (mock *MockClient) WriteNil() error {
	mock.response = []string{"NIL"}
	return nil
}

func (mock *MockClient) Write([]byte) (int, error) {
	return 0, NotImplementedError
}

func (mock MockClient) Data() [][]byte {
	return mock.d
}

func (mock MockClient) Body() []string {
	return mock.body
}

func (mock MockClient) DataRaw() RawRequest {
	return mock.r
}

func (mock MockClient) Store() storage.Storage {
  return mock.store
}

func (mock MockClient) SendClusterConfigChange(raftpb.ConfChange) {
	panic("mock caller must provide implementation of SendClusterConfigChange()")
}
