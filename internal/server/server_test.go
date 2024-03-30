package server

import (
	"context"
	"pedis/internal/storage"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerSetAndGet(t *testing.T) {
	s, err := NewRedisServer(storage.NewSimpleStorage(), Config{ServerAddr: "localhost:6379"})

	require.NoError(t, err)

	go s.Start()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	t.Run("Can set a value and retrieve it", func(t *testing.T) {
		err = client.Set(context.Background(), "key", "value", 0).Err()
		require.NoError(t, err)

		result, err := client.Get(context.Background(), "key").Result()
		require.NoError(t, err)
		assert.Equal(t, "value", result)
	})

	t.Run("Cannot set a key with empty value", func(t *testing.T) {
		err = client.Set(context.Background(), "key", "", 0).Err()

		assert.Equal(t, err.Error(), "ERR value is empty")
	})

	t.Run("Cannot set a empty string as key", func(t *testing.T) {
		err = client.Set(context.Background(), "", "value", 0).Err()
		assert.Equal(t, err.Error(), "ERR key is empty")
	})
}
