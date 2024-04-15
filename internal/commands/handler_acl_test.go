package commands

import (
	"bytes"
	"pedis/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAclHandlerSetUser(t *testing.T) {
	body := []byte("*x\r\n$3\r\nacl\r\n$7\r\nsetuser\r\n$5\r\npathe\r\n$2\r\non\r\n$10\r\n>password:\r\n")
	h := AclHandler{}
	r := NewMockClient(RawRequest(body), bytes.Split(body, []byte{13, 10})).WithMockedStore(
		func() storage.Storage {
			return &storage.MockStorage{
				SetUserFunc: func(string, []storage.AclRule) error {
					return nil
				},
			}
		})

	h.Handle(r)

	assert.Equal(t, []string{"OK"}, r.response)
}

func TestAclHandlerListUsers(t *testing.T) {
	body := []byte("*2\r\n$3\r\nacl\r\n$5\r\nusers\r\n")
	h := AclHandler{}
	r := NewMockClient(RawRequest(body), bytes.Split(body, []byte{13, 10})).WithMockedStore(
		func() storage.Storage {
			return &storage.MockStorage{
				UsersFunc: func() []string {
					return []string{"pedis"}
				},
			}
		})

	h.Handle(r)

	assert.Equal(t, []string{"pedis"}, r.response)
}

//func TestAclHandlerPermissions(t *testing.T) {
//	h := AclHandler{}
//	perms := h.Permissions(&MockClient{})
//	assert.Equal(t, []string{"fast", "connection"}, perms)
//}

func TestAclHandlerAuthorize(t *testing.T) {
	h := AclHandler{}
	require.NoError(t, h.Authorize(&MockClient{}))
}

func TestAclHandlerPersistent(t *testing.T) {
	h := AclHandler{}
	require.False(t, h.Persistent(&MockClient{}))
}
