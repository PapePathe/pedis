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

	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:             "localhost:6379",
		Password:         "",
		DB:               0,
		DisableIndentity: true,
	})

	t.Run("Can set a value with expiration date", func(t *testing.T) {
		err := client.Set(context.Background(), "key", "value", 2*time.Minute).Err()
		require.NoError(t, err)
	})

	t.Run("Can set a value ,retrieve and delete it ", func(t *testing.T) {
		err := client.Set(context.Background(), "key", "value", 0).Err()
		require.NoError(t, err)

		result, err := client.Get(context.Background(), "key").Result()
		require.NoError(t, err)
		assert.Equal(t, "value", result)

		resultDel, err := client.Del(ctx, "key", "key-2").Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), resultDel)

	})

	t.Run("DEL", func(t *testing.T) {
		t.Parallel()

		t.Run("All specified keys exists", func(t *testing.T) {
			_ = client.Set(context.Background(), "del:key", "value", 0).Err()
			_ = client.Set(context.Background(), "del:key-1", "value", 0).Err()
			resultDel, err := client.Del(ctx, "del:key", "del:key-1").Result()

			require.NoError(t, err)
			assert.Equal(t, int64(2), resultDel)
		})

		t.Run("Only one key exists", func(t *testing.T) {
			_ = client.Set(context.Background(), "del:key", "value", 0).Err()
			_ = client.Set(context.Background(), "del:key-1", "value", 0).Err()
			resultDel, err := client.Del(ctx, "del:key", "del:key-2").Result()

			require.NoError(t, err)
			assert.Equal(t, int64(1), resultDel)
		})
	})

	t.Run("ACL", func(t *testing.T) {
		t.Parallel()

		t.Run("CAT", func(t *testing.T) {
			t.SkipNow()
			_, err := client.Do(ctx, "acl", "cat").Result()

			require.NoError(t, err)
		})

		t.Run("CAT dangerous", func(t *testing.T) {
			t.SkipNow()
			_, err := client.Do(ctx, "acl", "cat", "dangerous").Result()

			require.NoError(t, err)
		})

		t.Run("SETUSER", func(t *testing.T) {
			_, err := client.Do(ctx, "acl", "setuser", "pathe-s").Result()

			require.NoError(t, err)
		})

		t.Run("GETUSER", func(t *testing.T) {
			t.SkipNow()
			_, err := client.Do(ctx, "acl", "getuser", "pathe").Result()

			require.NoError(t, err)
		})

		t.Run("USERS-1", func(t *testing.T) {
			list, err := client.Do(ctx, "acl", "users").Result()

			require.NoError(t, err)
			assert.Equal(t, []interface{}{}, list)
		})

		t.Run("USERS-2", func(t *testing.T) {
			_, _ = client.Do(ctx, "acl", "setuser", "acl-user-1").Result()
			_, _ = client.Do(ctx, "acl", "setuser", "acl-user-2").Result()

			list, err := client.Do(ctx, "acl", "users").Result()

			require.NoError(t, err)
			assert.Equal(t, []interface{}([]interface{}{"acl-user-1", "acl-user-2"}), list)

		})

		t.Run("DELUSER-1", func(t *testing.T) {
			_, _ = client.Do(ctx, "acl", "setuser", "pathe-1").Result()
			_, _ = client.Do(ctx, "acl", "setuser", "mado-1").Result()
			count, err := client.Do(ctx, "acl", "deluser", "pathe-1", "mado-1").Result()

			require.NoError(t, err)
			assert.Equal(t, int64(2), count)
		})

		t.Run("DELUSER-2", func(t *testing.T) {
			_, _ = client.Do(ctx, "acl", "setuser", "pathe-2").Result()
			count, err := client.Do(ctx, "acl", "deluser", "pathe-2", "mado").Result()

			require.NoError(t, err)
			assert.Equal(t, int64(1), count)
		})

		t.Run("DELUSER-3", func(t *testing.T) {
			count, err := client.Do(ctx, "acl", "deluser", "pathe-3", "mado").Result()

			require.NoError(t, err)
			assert.Equal(t, int64(0), count)
		})

		t.Run("DRYRUN", func(t *testing.T) {
			t.SkipNow()
			_, err := client.Do(ctx, "acl", "dryrun", "pathe", "get", "foo").Result()

			require.NoError(t, err)
		})
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
		Addr:             "localhost:6379",
		Password:         "",
		DB:               0,
		DisableIndentity: true,
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

		keys, err := client.HKeys(ctx, "user").Result()
		require.NoError(t, err)
		assert.Equal(t, []string{"name", "country", "221"}, keys)

		keys, err = client.HKeys(ctx, "not-a-key").Result()
		require.NoError(t, err)
		assert.Equal(t, []string{}, keys)

		keys, err = client.HVals(ctx, "user").Result()
		require.NoError(t, err)
		assert.Equal(t, []string{"Pathe", "Senegal", ""}, keys)

		keys, err = client.HVals(ctx, "not-a-key").Result()
		require.NoError(t, err)
		assert.Equal(t, []string{}, keys)

		exists, err := client.HExists(ctx, "user", "name").Result()
		require.NoError(t, err)
		assert.Equal(t, true, exists)

		exists, err = client.HExists(ctx, "user", "not-a-field").Result()
		require.NoError(t, err)
		assert.Equal(t, false, exists)

		exists, err = client.HExists(ctx, "key", "not-a-field").Result()
		require.NoError(t, err)
		assert.Equal(t, false, exists)
	})
}
