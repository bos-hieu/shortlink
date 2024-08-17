package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	client *mongo.Client
	once   sync.Once
)

// GetClient returns a mongo client instance
func GetClient() *mongo.Client {
	if client == nil {
		once.Do(func() {
			serverAPI := options.ServerAPI(options.ServerAPIVersion1)
			url := "mongodb://localhost:27017"
			opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
			newClient, err := mongo.Connect(context.TODO(), opts)
			if err != nil {
				panic(err)
			}

			client = newClient
		})
	}

	return client
}