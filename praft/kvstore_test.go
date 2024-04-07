// Copyright 2016 The etcd Authors
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

//func Test_kvstore_snapshot(t *testing.T) {
//	tm := map[string]string{"foo": "bar"}
//	s := &kvstore{kvStore: tm}
//
//	v, _ := s.Lookup("foo")
//	if v != "bar" {
//		t.Fatalf("foo has unexpected value, got %s", v)
//	}
//
//	data, err := s.getSnapshot()
//	if err != nil {
//		t.Fatal(err)
//	}
//	s.kvStore = nil
//
//	if err := s.recoverFromSnapshot(data); err != nil {
//		t.Fatal(err)
//	}
//	v, _ = s.Lookup("foo")
//	if v != "bar" {
//		t.Fatalf("foo has unexpected value, got %s", v)
//	}
//	if !reflect.DeepEqual(s.kvStore, tm) {
//		t.Fatalf("store expected %+v, got %+v", tm, s.kvStore)
//	}
//}

import (
	"context"
	"pedis/internal/storage"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerSetAndGet(t *testing.T) {
	storageProposeChan := make(chan storage.StorageData)

	s := NewPedisServer(
		"localhost:6379",
		storage.NewSimpleStorage(storageProposeChan),
	)

	go s.StartPedis()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	t.Run("Can set a value with expiration date", func(t *testing.T) {
		err := client.Set(context.Background(), "key", "value", 2*time.Minute).Err()
		require.NoError(t, err)
	})

	t.Run("Can set a value and retrieve it", func(t *testing.T) {
		err := client.Set(context.Background(), "key", "value", 0).Err()
		require.NoError(t, err)

		result, err := client.Get(context.Background(), "key").Result()
		require.NoError(t, err)
		assert.Equal(t, "value", result)
	})

	t.Run("Cannot set a key with empty value", func(t *testing.T) {
		err := client.Set(context.Background(), "key", "", 0).Err()

		assert.Equal(t, err.Error(), "ERR value is empty")
	})

	t.Run("Cannot set a empty string as key", func(t *testing.T) {
		err := client.Set(context.Background(), "", "value", 0).Err()
		assert.Equal(t, err.Error(), "ERR key is empty")
	})

	t.Run("Cannot get a key that does not exist", func(t *testing.T) {
		_, err := client.Get(context.Background(), "key:not:found").Result()
		assert.Equal(t, err.Error(), "ERR key not found")
	})
}

func TestServerHSetAndHGet(t *testing.T) {
	storageProposeChan := make(chan storage.StorageData)
	s := NewPedisServer(
		"localhost:6379",
		storage.NewSimpleStorage(storageProposeChan),
	)

	go s.StartPedis()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	t.Run("Can set and get a hash", func(t *testing.T) {
		//	m := map[string]interface{}{"key-one": "one value", "key-two": "two value"}
		//		err = client.HMSet(context.Background(), "myhash", m, 0).Err()
		ctx := context.Background()

		result, err := client.HSet(context.Background(), "user", "name", "Pathe", "country", "Senegal", 221).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(3), result)

		name, err := client.HGet(context.Background(), "user", "name").Result()
		require.NoError(t, err)
		assert.Equal(t, "Pathe", name)

		name, err = client.HGet(context.Background(), "user", "country").Result()
		require.NoError(t, err)
		assert.Equal(t, "Senegal", name)

		_, err = client.HGet(context.Background(), "user", "not-a-field").Result()
		assert.Equal(t, redis.Nil, err)

		_, err = client.HGet(context.Background(), "not-a-key", "country").Result()
		assert.Equal(t, redis.Nil, err)

		name, err = client.HGet(context.Background(), "user", "221").Result()
		require.NoError(t, err)
		assert.Equal(t, "", name)

		l, err := client.HLen(ctx, "user").Result()
		require.NoError(t, err)
		assert.Equal(t, int64(3), l)

		l, err = client.HLen(ctx, "not-a-key").Result()
		require.NoError(t, err)
		assert.Equal(t, int64(0), l)
	})
}
