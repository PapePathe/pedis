package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingHandle(t *testing.T) {
	body := []byte("*1\r\n$4\r\nping\r\n")
	h := PingHandler{}
	r := MockClient{
		r: RawRequest(body),
		d: bytes.Split(body, []byte{13, 10}),
	}

	h.Handle(&r)

	assert.Equal(t, []string{"PONG"}, r.response)
}

func TestCustomPingHandle(t *testing.T) {
	body := []byte("*1\r\n$4\r\nping\r\n$4\r\nnangadef\r\n")
	h := PingHandler{}
	r := MockClient{
		r: RawRequest(body),
		d: bytes.Split(body, []byte{13, 10}),
	}

	h.Handle(&r)

	assert.Equal(t, []string{"nangadef"}, r.response)
}

func TestPingPermissions(t *testing.T) {
	h := PingHandler{}
	perms := h.Permissions(&MockClient{})
	assert.Equal(t, []string{"fast", "connection"}, perms)
}

func TestPingAuthorize(t *testing.T) {
	h := PingHandler{}
	require.NoError(t, h.Authorize(&MockClient{}))
}

func TestPingPersistent(t *testing.T) {
	h := PingHandler{}
	require.False(t, h.Persistent(&MockClient{}))
}
