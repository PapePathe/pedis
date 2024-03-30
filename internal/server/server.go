package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"pedis/internal/commands"
	"pedis/internal/storage"
	"time"

	"github.com/hashicorp/raft"
)

type RedisCommand interface {
	Run([]byte, net.Conn, storage.Storage)
}

type RedisServer struct {
	handlers map[string]RedisCommand
	store    storage.Storage
	config   Config
	raft     *raft.Raft
}

type Config struct {
	Bootstrap      bool
	JoinAddr       string
	MembershipAddr string
	RaftAddr       string
	ServerAddr     string
	ServerId       string
}

func NewRedisServer(store storage.Storage, serverCfg Config) (*RedisServer, error) {
	server := RedisServer{
		handlers: make(map[string]RedisCommand),
		store:    store,
		config:   serverCfg,
	}

	_ = server.AddHandler("*", commands.RequestHandler{})

	return &server, nil
}

func (rs *RedisServer) StartRaft() error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(rs.config.ServerId)
	addr, err := net.ResolveTCPAddr("tcp", rs.config.RaftAddr)
	if err != nil {
		return fmt.Errorf("error resolving raft address %v", err)
	}
	transport, err := raft.NewTCPTransport(rs.config.RaftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return fmt.Errorf("error creating raft transport %v", err)
	}

	snapshots, err := raft.NewFileSnapshotStore("./raft-tmp", 2, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()
	ra, err := raft.NewRaft(config, (*fsm)(rs), logStore, stableStore, snapshots, transport)

	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}

	if ra == nil {
		return fmt.Errorf("raft is nil: %s", err)
	}

	rs.raft = ra

	//	if rs.config.Bootstrap {
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}
	err = rs.raft.BootstrapCluster(configuration).Error()

	if err != nil {
		return fmt.Errorf("error bootstrapping cluster %v", err)
	}
	//	}

	return nil
}

func (rs *RedisServer) Start() error {
	listener, err := net.Listen("tcp", rs.config.ServerAddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
		}

		go rs.handleConnection(conn)
	}
}

func (rs *RedisServer) AddHandler(firstByte string, c RedisCommand) error {
	rs.handlers[firstByte] = c

	return nil
}

func (rs *RedisServer) handleConnection(conn net.Conn) {
	for {
		b := make([]byte, 1024)

		size, err := conn.Read(b)

		if err != nil || size == 0 {
			conn.Close()
			break
		}

		commandId := string(b[0])
		handler, commandNotFound := rs.handlers[commandId]

		if !commandNotFound {
			log.Println(err)
			continue
		}

		handler.Run(b[1:size], conn, rs.store)
	}
}

type fsm RedisServer

func (f *fsm) Apply(l *raft.Log) interface{} {
	return nil
}

func (f *fsm) Restore(rc io.ReadCloser) error {
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (rs *RedisServer) Leave(id string) error {
	if rs.raft == nil {
		return fmt.Errorf("leave aborted raft server not running")
	}
	removeFuture := rs.raft.RemoveServer(raft.ServerID(id), 0, 0)

	return removeFuture.Error()
}

func (n *RedisServer) Join(nodeID, addr string) error {
	log.Println("node id", nodeID, "addr", addr)
	if n.raft == nil {
		return fmt.Errorf("join aborted raft server not running")
	}

	log.Println("debug raft", n.raft)

	configFuture := n.raft.GetConfiguration()

	if err := configFuture.Error(); err != nil {
		return fmt.Errorf("error configuring raft %v+", err)
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				log.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
			}

			future := n.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}
	f := n.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	log.Printf("node %s at %s joined successfully", nodeID, addr)

	return nil
}
