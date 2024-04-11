package commands

import (
	"log"
	"strconv"

	"go.etcd.io/etcd/raft/v3/raftpb"
)

func ClusterHandler(r ClientRequest) {
	data := r.DataRaw.ReadArray()
	log.Println(data)

	switch string(data[0]) {
	case "forget":
		id, err := strconv.Atoi(data[1])

		if err != nil {
			r.WriteError(err.Error())
		}
		cc := raftpb.ConfChange{
			Type:   raftpb.ConfChangeRemoveNode,
			NodeID: uint64(id),
		}
		r.WriteOK()
		r.ClusterChangesChan <- cc

	case "meet":
		id, err := strconv.Atoi(data[1])

		if err != nil {
			r.WriteError(err.Error())
		}

		cc := raftpb.ConfChange{
			Type:    raftpb.ConfChangeAddNode,
			NodeID:  uint64(id),
			Context: []byte(data[2]),
		}
		r.WriteOK()
		r.ClusterChangesChan <- cc
	default:
		r.WriteError("subcommand not found")
	}
}
