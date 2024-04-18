package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingHandle(t *testing.T) {
  type pingtest struct { 
    cli MockClient 
    rep []string 
    name string 
  }
  tests := []pingtest{
    {name: "with empty body", cli:  MockClient{ body: []string{}}, rep: []string{"PONG" }},
    {name: "with custom pong", cli:  MockClient{ body: []string{"nangadef"}}, rep: []string{"nangadef" }},
  }

	h := PingHandler{}
  for _, test := range  tests {
     t.Run(test.name, func(t *testing.T){
        h.Handle(&test.cli)
        assert.Equal(t, test.rep, test.cli.response)
      })
  }
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
