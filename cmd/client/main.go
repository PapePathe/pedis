package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	//	for i := 0; i < 100; i++ {
	log.Println("creating client")
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println(client.Conn())
	err := client.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		log.Println(err)
	}
	// }
}
