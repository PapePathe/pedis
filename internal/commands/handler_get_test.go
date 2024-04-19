package commands

import (
	"fmt"
	"pedis/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGetHandlerPersistend(t *testing.T)  {
  h := GetHandler{}
  assert.True(t, h.Persistent(&MockClient{}))
}

func TestGetHandler(t *testing.T) {
	type authtest struct {
		req *MockClient
		rep  []string
    err []string
		name string
	}
	tests := []*authtest{
    { 
      name: "get a key that does not exist",
      req: &MockClient{
        body: []string{ "key:404" },
        store: storage.MockStorage{
          GetFn: func(k string) (string, error){
            return "", fmt.Errorf("key %s not found", k)
          },
        },
      },
      err: []string{"key key:404 not found"},
    },
    { 
      name: "get a key that does not exist",
      req: &MockClient{
        body: []string{ "key:200" },
        store: storage.MockStorage{
          GetFn: func(k string) (string, error){
            return "two hundred thousand xof", nil 
          },
        },
      },
      rep: []string{"two hundred thousand xof"},
    },
  }


  for _, test := range tests {
    h := GetHandler{}
    h.Handle(test.req)

    if len(test.err) > 0 {
      assert.Equal(t, test.err, test.req.errors)
    } else {
      assert.Equal(t, test.rep, test.req.response)
    }
  }
}
