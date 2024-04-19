package commands

import (
	"testing"

	"pedis/internal/storage"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

type clustertest struct {
	cli  *MockClient
	rep  []string
	err  []string
	name string
}

func TestClusterHandlerForget(t *testing.T) {
	store := &storage.MockStorage{}
	tests := []*clustertest{
		{
			name: "with only forget subcommand",
			cli:  &MockClient{body: []string{"forget"}, store: store},
			err:  []string{"you must supply a subcommand and args"},
		},
		{
			name: "with meet subcommand and valid parameters",
			cli:  &MockClient{body: []string{"forget", "19"}, store: store},
			rep:  []string{"OK"},
		},
		{
			name: "with meet subcommand and invalid node id",
			cli:  &MockClient{body: []string{"forget", "x19"}, store: store},
			err:  []string{"node id is not a number"},
		},
	}
	h := ClusterHandler{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.Handle(test.cli)

			if len(test.err) > 0 {
				assert.Equal(t, test.err, test.cli.errors)
			} else {
				assert.Equal(t, test.rep, test.cli.response)
				assert.Equal(t, uint64(19), test.cli.confChange.NodeID)
				assert.Equal(t, raftpb.ConfChangeRemoveNode, test.cli.confChange.Type)
			}
		})
	}
}
func TestClusterHandlerMeet(t *testing.T) {
	store := &storage.MockStorage{}
	tests := []*clustertest{
		{
			name: "with no subcommand",
			cli:  &MockClient{body: []string{}, store: store},
			err:  []string{"you must supply a subcommand and args"},
		},
		{
			name: "with only forget subcommand",
			cli:  &MockClient{body: []string{"forget"}, store: store},
			err:  []string{"you must supply a subcommand and args"},
		},
		{
			name: "with only meet subcommand",
			cli:  &MockClient{body: []string{"meet"}, store: store},
			err:  []string{"you must supply a subcommand and args"},
		},
		{
			name: "with meet subcommand and a node id that is not a number",
			cli:  &MockClient{body: []string{"meet", "x", "http://local.cluster"}, store: store},
			err:  []string{"node id is not a number"},
		},
		{
			name: "with meet subcommand and valid parameters",
			cli:  &MockClient{body: []string{"meet", "19", "http://local.cluster"}, store: store},
			rep:  []string{"OK"},
		},
	}
	h := ClusterHandler{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.Handle(test.cli)

			if len(test.err) > 0 {
				assert.Equal(t, test.err, test.cli.errors)
			} else {
				assert.Equal(t, test.rep, test.cli.response)
				assert.Equal(t, uint64(19), test.cli.confChange.NodeID)
				assert.Equal(t, []byte("http://local.cluster"), test.cli.confChange.Context)
				assert.Equal(t, raftpb.ConfChangeAddNode, test.cli.confChange.Type)
			}
		})
	}
}

func TestClusterHandlerPersistent(t *testing.T) {
	h := ClusterHandler{}
	assert.Equal(t, true, h.Persistent(&MockClient{}))
}
