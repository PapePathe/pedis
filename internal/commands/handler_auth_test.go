package commands

import (
	"fmt"
	"pedis/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandlerPersistent(t *testing.T)  {
  h := AuthHandler{}
  require.False(t, h.Persistent(&MockClient{}))
  
}

func TestAuthHandler(t *testing.T) {
	type authtest struct {
		req *MockClient
		rep  []string
    err []string
		name string
	}
	tests := []*authtest{
    { 
      name: "with a user that can login without password",
      req: &MockClient{
        body: []string{"pathe"}, 
        store: &storage.MockStorage{
          GetUserFn: func(k string) (*storage.User, error) {
            return &storage.User{AnyPassword: true}, nil
          },
        },
      },
      rep: []string{"OK"},
      err:  []string{},
    },
    {
      name: "with a user that does not exist",
      req: &MockClient{
        body: []string{"pathe"},
        store: &storage.MockStorage{
          GetUserFn: func(k string) (*storage.User, error) {
            return nil, fmt.Errorf("user %s not found", k)
          },
        },
      },
      rep: []string{},
      err:  []string{"user pathe not found"},
    },
    { 
      name: "with a user that cannot login without password",
      req: &MockClient{
        body: []string{"pathe"}, 
        store: &storage.MockStorage{
          GetUserFn: func(k string) (*storage.User, error) {
            return &storage.User{AnyPassword: false}, nil
          },
        },
      },
      rep: []string{},
      err:  []string{"Password must be supplied"},
    },
    { 
      name: "with a user that has an invalid password",
      req: &MockClient{
        body: []string{"pathe", "badpassword"}, 
        store: &storage.MockStorage{
          GetUserFn: func(k string) (*storage.User, error) {
            return &storage.User{AnyPassword: false}, nil
          },
        },
      },
      rep: []string{},
      err:  []string{"Password auth failed"},
    },
    { 
      name: "with a user that has a valid password",
      req: &MockClient{
        body: []string{"pathe", "validpassword"}, 
        store: &storage.MockStorage{
          GetUserFn: func(k string) (*storage.User, error) {
            return &storage.User{AnyPassword: false, Passwords: []string{"validpassword"}}, nil
          },
        },
      },
      rep: []string{"OK"},
      err:  []string{},
    },
  }

  for _, test := range tests {
    h := AuthHandler{}
    h.Handle(test.req)

    if len(test.err) > 0 {
      assert.Equal(t, test.err, test.req.errors)
    } else {
      assert.Equal(t, test.rep, test.req.response)
    }
  }
}
