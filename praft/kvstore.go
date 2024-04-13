// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package praft

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"pedis/internal/commands"
	"pedis/internal/storage"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
)

type RedisCommand interface {
	Run(commands.ClientRequest)
}

type PedisServerOpts struct {
	DefaultUser     string
	DefaultPassword string
}

// a key-value store backed by raft
type PedisServer struct {
	proposeC           chan<- string // channel for proposing updates
	clusterChangesChan chan<- raftpb.ConfChange
	mu                 sync.RWMutex
	kvStore            map[string]string // current committed key-value pairs
	snapshotter        *snap.Snapshotter

	handlers map[string]RedisCommand
	store    storage.Storage
	addr     string

	storageProposeChan chan storage.StorageData

	logger   zerolog.Logger
	listener net.Listener
}

func NewPedisServer(
	pedisAddr string,
	store storage.Storage,
) *PedisServer {
	s := &PedisServer{
		handlers: make(map[string]RedisCommand),
		addr:     pedisAddr,
		store:    store,
		logger: zerolog.New(
			zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		).With().Timestamp().Logger(),
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	_ = s.AddHandler("*", commands.DefaultRequestHandler())
	return s
}

func NewKVStore(
	snapshotter *snap.Snapshotter,
	proposeC chan<- string,
	commitC <-chan *commit,
	errorC <-chan error,
	store storage.Storage,
	pedisAddr string,
	storageProposeChan chan storage.StorageData,
	clusterConfChan chan raftpb.ConfChange,
) *PedisServer {
	s := &PedisServer{
		proposeC:           proposeC,
		kvStore:            make(map[string]string),
		snapshotter:        snapshotter,
		handlers:           make(map[string]RedisCommand),
		addr:               pedisAddr,
		store:              store,
		storageProposeChan: storageProposeChan,
		clusterChangesChan: clusterConfChan,
		logger: zerolog.New(
			zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		).With().Timestamp().Logger(),
	}
	snapshot, err := s.loadSnapshot()
	if err != nil {
		log.Panic(err)
	}
	if snapshot != nil {
		log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
		if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
			log.Panic(err)
		}
	}

	_ = s.AddHandler("*", commands.DefaultRequestHandler())
	// read commits from raft into kvStore map until error
	go s.readCommits(commitC, errorC)
	go s.readProposeChan()
	return s
}

func (rs *PedisServer) AddHandler(firstByte string, c RedisCommand) error {
	rs.handlers[firstByte] = c

	return nil
}

func (s *PedisServer) StartPedis() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	defer listener.Close()
	s.listener = listener

	for {
		s.logger.Debug().Msg("received new connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
		}

		go s.handleConnection(conn)
	}
}

func (rs *PedisServer) handleConnection(conn net.Conn) {
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

		request := commands.ClientRequest{
			Conn:               conn,
			Data:               bytes.Split(b[1:size], []byte{13, 10}),
			Store:              rs.store,
			Logger:             rs.logger,
			DataRaw:            commands.RawRequest(b[0:size]),
			ClusterChangesChan: rs.clusterChangesChan,
		}

		handler.Run(request)
	}
}

func (s *PedisServer) Lookup(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.kvStore[key]
	return v, ok
}

func (s *PedisServer) Propose(data storage.StorageData) {
	var buf strings.Builder
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		log.Fatal(err)
	}

	s.proposeC <- buf.String()
}

func (s *PedisServer) readProposeChan() {
	for propose := range s.storageProposeChan {
		s.Propose(propose)
	}
}
func (s *PedisServer) readCommits(commitC <-chan *commit, errorC <-chan error) {
	log.Println("reading commits")
	for commit := range commitC {
		if commit == nil {
			// signaled to load snapshot
			snapshot, err := s.loadSnapshot()
			if err != nil {
				log.Panic(err)
			}
			if snapshot != nil {
				log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
				if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		for _, data := range commit.data {
			var dataKv storage.StorageData
			dec := gob.NewDecoder(bytes.NewBufferString(data))
			if err := dec.Decode(&dataKv); err != nil {
				log.Fatalf("raftexample: could not decode message (%v)", err)
			}
			err := s.store.WriteFromRaft(dataKv)
			if err != nil {
				log.Println("error writing proposed change from raft", err)
			}
		}
		close(commit.applyDoneC)
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (s *PedisServer) GetSnapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.kvStore)
}

func (s *PedisServer) PedisAddr() (*net.TCPAddr, error) {
	if s.listener == nil {
		return nil, errors.New("listener is not started or assigned")
	}
	return s.listener.Addr().(*net.TCPAddr), nil
}

func (s *PedisServer) loadSnapshot() (*raftpb.Snapshot, error) {
	snapshot, err := s.snapshotter.Load()
	if err == snap.ErrNoSnapshot {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (s *PedisServer) recoverFromSnapshot(snapshot []byte) error {
	var store map[string]string
	if err := json.Unmarshal(snapshot, &store); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.kvStore = store
	return nil
}
