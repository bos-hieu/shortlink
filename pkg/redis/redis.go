package redis

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	client *redis.Client
	once   sync.Once
)

// GetClient returns a redis client instance
func GetClient() *redis.Client {
	if client == nil {
		once.Do(func() {
			client = redis.NewClient(&redis.Options{
				Addr: "localhost:6379",
			})
		})
	}

	return client
}
