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

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"pedis/internal/storage"
	"pedis/praft"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	join := flag.Bool("join", false, "join an existing cluster")
	pedis := flag.String("pedis", "0.0.0.0:6379", "port where pedis server is running")
	healthcheck := flag.Bool("health", false, "allow docker and kubernetes to check the health")
	flag.Parse()

	if *healthcheck {
		client := redis.NewClient(&redis.Options{
			Addr:     *pedis,
			Password: "",
			DB:       0,
		})

		log.Println("Running ping command")
		pong, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if strings.ToLower(pong) != "pong" {
			log.Fatalf("Received unexpected response (%v)", pong)
			os.Exit(1)
		}

		log.Println("Server healthy")
		os.Exit(0)
	}

	proposeC := make(chan string)
	defer close(proposeC)

	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	storageProposeChan := make(chan storage.StorageData)
	defer close(storageProposeChan)

	var kvs *praft.PedisServer
	getSnapshot := func() ([]byte, error) { return kvs.GetSnapshot() }
	commitC, errorC, snapshotterReady := praft.NewRaftNode(
		*id,
		strings.Split(*cluster, ","),
		*join,
		getSnapshot,
		proposeC,
		confChangeC,
	)

	kvs = praft.NewKVStore(
		<-snapshotterReady,
		proposeC,
		commitC,
		errorC,
		storage.NewSimpleStorage(storageProposeChan),
		*pedis,
		storageProposeChan,
		confChangeC,
	)

	go func() {
		if err := kvs.StartPedis(); err != nil {
			log.Fatal(err)
		}
	}()
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}
