package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	for i := 0; i < 2; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := client.Set(context.Background(), key, "value", 0).Err()
		if err != nil {
			log.Println(err)
		}

		result, err := client.Get(context.Background(), key).Result()

		if err != nil {
			log.Println(err)
		}

		log.Println("Get result", result)
	}
}
