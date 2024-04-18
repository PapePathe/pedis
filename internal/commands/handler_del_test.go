package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
  "pedis/internal/storage"
)

func TestDelHandler(t *testing.T) {
	type pingtest struct {
		cli  *MockClient
		rep  []string
		name string
	}
  store := &storage.MockStorage{ 
    DelUserFn: func(k string) error { return nil },
  }
	tests := []*pingtest{
    &pingtest{ 
      name: "with no user to delete", 
      cli: &MockClient{body: []string{}, store: store}, 
      rep: []string{"0"},
    },
		&pingtest{
      name: "one user", 
      cli: &MockClient{body: []string{"user-1"}, store: store}, 
      rep: []string{"1"},
    },
		&pingtest{
      name: "multiple users", 
      cli: &MockClient{body: []string{"user-1", "user-2", "user-3"}, store: store}, 
      rep: []string{"3"},
    },
	}
	h := DelHandler{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.Handle(test.cli)
			assert.Equal(t, test.cli.response, test.rep)
		})
	}
}


func TestDelHandlerPersistent(t *testing.T) {
	h := DelHandler{}
  assert.Equal(t, true, h.Persistent(&MockClient{}))
}
