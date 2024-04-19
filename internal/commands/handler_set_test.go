package commands

import(
  "testing"
  "fmt"
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
      name: "get a key that does not exist",
      req: &MockClient{
        body: []string{ "key:404", "a value", "ex", "1000" },
        store: storage.MockStorage{
          GetFn: func(k string) (string, error){
            return "", fmt.Errorf("key %s not found", k)
          },
        },
      },
      rep: []string{"OK"},
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
