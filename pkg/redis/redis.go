package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	client *redis.Client
)

// GetClient returns a redis client instance
func GetClient() *redis.Client {
	return client
}

// InitClient initializes a redis client instance
func InitClient() error {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	result, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	log.Println("Redis client is connected: ", result)
	return nil
}
