package commands

import (
	"errors"
	"pedis/internal/storage"

	"go.etcd.io/etcd/raft/v3/raftpb"
)

var (
	NotImplementedError error = errors.New("mock caller must provide implementation")
)

type MockClient struct {
	r        RawRequest
	d        [][]byte
	response []string
	errors   []string
	storeFn  func() storage.Storage
}

func NewMockClient(r RawRequest, d [][]byte) *MockClient {
	return &MockClient{r: r, d: d}
}

func (mock *MockClient) WithMockedStore(fn func() storage.Storage) *MockClient {
	mock.storeFn = fn

	return mock
}

func (mock *MockClient) WriteError(e string) error {
	mock.errors = append(mock.errors, e)
	return nil
}

func (mock *MockClient) WriteString(s string) error {
	mock.response = append(mock.response, s)
	return nil
}

func (mock *MockClient) WriteNumber(s string) error {
	mock.response = append(mock.response, s)
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

func (mock MockClient) DataRaw() RawRequest {
	return mock.r
}

func (mock MockClient) Store() storage.Storage {
	if mock.storeFn == nil {
		panic("mock caller must provide implementation of Store()")
	}
	return mock.storeFn()
}

func (mock MockClient) SendClusterConfigChange(raftpb.ConfChange) {
	panic("mock caller must provide implementation of SendClusterConfigChange()")
}

func (c MockClient) ArrayResponse([]string) []byte {
	return nil
}

func (c MockClient) StringResponse(string) []byte {
	return nil
}

func (c MockClient) OkResponse(string) []byte {
	return nil
}

func (c MockClient) ErrorResponse(string) []byte {
	return nil
}

func (c MockClient) NilResponse(string) []byte {
	return nil
}
