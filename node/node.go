package node

import (
	"fmt"
	"log"
	"pedis/discovery"
	"pedis/internal/server"
	"pedis/internal/storage"
)

type Node struct {
	config     Config
	pedis      *server.RedisServer
	membership *discovery.Membership
}

func NewNode(serverConfig Config) (*Node, error) {
	pds, err := server.NewRedisServer(
		storage.NewSimpleStorage(),
		server.Config{
			ServerAddr: serverConfig.ServerAddr,
			Bootstrap:  serverConfig.Bootstrap,
			JoinAddr:   serverConfig.JoinAddr,
			RaftAddr:   serverConfig.RaftAddr,
			ServerId:   serverConfig.ServerId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating redis server  %v", err)
	}
	node := &Node{config: serverConfig, pedis: pds}

	log.Println("start join addrs", serverConfig.StartJoinAddrs)
	tags := make(map[string]string)
	tags["rpc_addr"] = serverConfig.RaftAddr
	membership, err := discovery.NewMembership(node.pedis, discovery.Config{
		BindAddr:       serverConfig.MembershipAddr,
		NodeName:       serverConfig.ServerId,
		StartJoinAddrs: serverConfig.StartJoinAddrs,
		Tags:           tags,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating membership server  %v", err)
	}

	node.membership = membership

	return node, nil
}

func (n Node) Start() error {
	if err := n.pedis.StartRaft(); err != nil {
		return fmt.Errorf("pedis raft server errored %v+", err)
	}

	if err := n.pedis.Start(); err != nil {
		return fmt.Errorf("pedis key value server errored %v+", err)
	}

	return nil
}
