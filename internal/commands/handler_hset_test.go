package commands

import (
	"pedis/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHSetHandlerPersistent(t *testing.T) {
	h := HSetHandler{}

	assert.True(t, h.Persistent(&MockClient{}))
}

func TestHSetHandler(t *testing.T) {
	type hsettest struct {
		cli  *MockClient
		rep  []string
		err  []string
		name string
	}

	store := &storage.MockStorage{}
	tests := []*hsettest{
		{
			name: "with valid params",
			cli: &MockClient{
				body:  []string{"user:101", "name", "Pathe", "country", "Senegal", "phone", "221"},
				store: store,
			},
			rep: []string{"3"},
		},
		{
			name: "with key only",
			cli: &MockClient{
				body:  []string{"user:101"},
				store: store,
			},
			err: []string{"you must provide values"},
		},
	}
	h := HSetHandler{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.Handle(test.cli)

			if len(test.err) > 0 {
				assert.Equal(t, test.err, test.cli.errors)
			} else {
				assert.Equal(t, test.rep, test.cli.response)
			}
		})
	}
}
