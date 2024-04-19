package commands

import(
  "testing"
  "github.com/stretchr/testify/assert"
  "pedis/internal/storage"
)

func TestSetHandlerPersistent(t *testing.T)  {
  h := SetHandler{}
  assert.True(t, h.Persistent(&MockClient{}))
}

func TestSetHandler(t *testing.T)  {
  type settest struct {
    req *MockClient
    rep  []string
    err []string
    name string
  }
  tests := []*settest{
    { 
      name: "set a key with expiration date",
      req: &MockClient{
        body: []string{ "key:with:exp", "a value", "ex", "1000" },
        store: storage.MockStorage{},
      },
      rep: []string{"OK"},
    },
   { 
      name: "set with an empty key",
      req: &MockClient{
        body: []string{ "" },
        store: storage.MockStorage{ },
      },
      err: []string{"key is empty"},
    },
   { 
      name: "set with an no value",
      req: &MockClient{
        body: []string{ "key:xxx" },
        store: storage.MockStorage{ },
      },
      err: []string{"value is required"},
    },
    { 
      name: "set with an empty value",
      req: &MockClient{
        body: []string{ "key:xxx", "" },
        store: storage.MockStorage{ },
      },
      err: []string{"value is empty"},
    },
  } 

  for _, test := range tests {
    h := SetHandler{}
    h.Handle(test.req)

    if len(test.err) > 0 {
      assert.Equal(t, test.err, test.req.errors)
    } else {
      assert.Equal(t, test.rep, test.req.response)
    }
  }

}
