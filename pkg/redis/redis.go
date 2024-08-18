package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
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

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	return nil
}
