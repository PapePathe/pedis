package commands

import (
	"pedis/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHValsHandlerPersistent(t *testing.T) {
	h := HValsHandler{}

	assert.False(t, h.Persistent(&MockClient{}))
}

func TestHValsHandler(t *testing.T) {
	type hgettest struct {
		cli  *MockClient
		rep  []string
		err  []string
		name string
	}

	store := &storage.MockStorage{}
	tests := []*hgettest{
		{
			name: "with valid params",
			cli: &MockClient{
				body:  []string{"user:101"},
				store: store,
			},
			rep: []string{},
		},
	}
	h := HValsHandler{}
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
