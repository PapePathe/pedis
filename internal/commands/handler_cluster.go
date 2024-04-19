package commands

import (
	"strconv"

	"go.etcd.io/etcd/raft/v3/raftpb"
)

type ClusterHandler struct{}

func (ch ClusterHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch ClusterHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch ClusterHandler) Persistent(IClientRequest) bool {
	return true 
}

func (ch ClusterHandler) Handle(r IClientRequest) {
	data := r.Body()

  if(len(data)) < 2 {
    r.WriteError("you must supply a subcommand and args")
    return
  }

	switch string(data[0]) {
	case "forget":
		id, err := strconv.Atoi(data[1])
		if err != nil {
      r.WriteError("node id is not a number")
		}

		cc := raftpb.ConfChange{
			Type:   raftpb.ConfChangeRemoveNode,
			NodeID: uint64(id),
		}
		r.WriteOK()
		r.SendClusterConfigChange(cc)
	case "meet":
		id, err := strconv.Atoi(data[1])

		if err != nil {
      r.WriteError("node id is not a number")
      return
		}

		cc := raftpb.ConfChange{
			Type:    raftpb.ConfChangeAddNode,
			NodeID:  uint64(id),
			Context: []byte(data[2]),
		}
		r.WriteOK()
		r.SendClusterConfigChange(cc)
	default:
		r.WriteError("subcommand not found")
	}

}

func init() {
	RegisterCommand("cluster", ClusterHandler{})
}
