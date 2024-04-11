package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var reqCount int

func main() {
	flag.IntVar(&reqCount, "req", 100, "")
	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	for i := 0; i < reqCount; i++ {
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

		key = fmt.Sprintf("hset:key-%d", i)
		err = client.HSet(context.Background(), key, "name", "pathe").Err()
		if err != nil {
			log.Println(err)
		}

		resultHset, err := client.HGet(context.Background(), key, "name").Result()

		if err != nil {
			log.Println(err)
		}

		log.Println("Hset result", resultHset)
	}
}
