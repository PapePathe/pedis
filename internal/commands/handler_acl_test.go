package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"pedis/internal/storage"
)

func TestAclHandler(t *testing.T) {
	type acltest struct {
		cli  *MockClient
		rep  []string
		err  []string
		name string
	}
	store := &storage.MockStorage{
		SetUserFn: func(k string, r []storage.AclRule) error { return nil },
		UsersFn:   func() []string { return []string{"pedis", "pathe"} },
	}
	tests := []*acltest{
		{
			name: "with a user that is offf",
			cli:  &MockClient{body: []string{"setuser", "pathe", "off"}, store: store},
			rep:  []string{"OK"},
			err:  []string{},
		},
		{
			name: "with valid user and no password",
			cli:  &MockClient{body: []string{"setuser", "pathe", "nopass"}, store: store},
			rep:  []string{"OK"},
			err:  []string{},
		},
		{
			name: "when listing users",
			cli:  &MockClient{body: []string{"users"}, store: store},
			rep:  []string{"pedis", "pathe"},
			err:  []string{},
		},
		{
			name: "with unsupported acl rule",
			cli:  &MockClient{body: []string{"setuser", "pathe", "notanaclrule"}, store: store},
			rep:  []string{},
			err:  []string{"acl rule (notanaclrule) not supported"},
		},
	}
	h := AclHandler{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.Handle(test.cli)
			if len(test.err) > 0 {
				assert.Equal(t, test.cli.errors, test.err)
			} else {
				assert.Equal(t, test.cli.response, test.rep)
			}
		})
	}
}
