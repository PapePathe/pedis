package main

import (
	"flag"
	"log"
	"pedis/node"
	"strings"
)

var addr string
var raftAddr string
var joinAddr string
var startJoinAddrs string
var id string
var membershipAddr string
var bootstrap bool

func init() {
	flag.StringVar(&membershipAddr, "serf-addr", "localhost:6359", "the address where the cluster membership server is listening")
	flag.StringVar(&raftAddr, "raft-addr", "localhost:6389", "the address where the raft server is listening")
	flag.StringVar(&addr, "addr", "localhost:6379", "the address where the server will listen")
	flag.StringVar(&id, "id", "primary", "the unique id of the server in the cluster")
	flag.StringVar(&joinAddr, "join-addr", "", "set address of the leader of the cluster")
	flag.StringVar(&startJoinAddrs, "sjoin-addr", "", "set addresses of cluster members to join when starting")
	flag.BoolVar(&bootstrap, "bootstrap", false, "start as bootstrap node")
}

func main() {
	flag.Parse()

	joinAddrs := []string{}

	if startJoinAddrs != "" {
		joinAddrs = strings.Split(startJoinAddrs, ",")
	}

	node, err := node.NewNode(node.Config{
		Bootstrap:      bootstrap,
		JoinAddr:       joinAddr,
		RaftAddr:       raftAddr,
		MembershipAddr: membershipAddr,
		ServerAddr:     addr,
		ServerId:       id,
		StartJoinAddrs: joinAddrs,
	})

	if err != nil {
		log.Fatalf("error creating node %v", err)
	}

	err = node.Start()

	if err != nil {
		log.Fatalf("error starting pedis node %v", err)
	}
}
