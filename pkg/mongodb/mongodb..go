package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var (
	client *mongo.Client
)

// GetClient returns a mongo client instance
func GetClient() *mongo.Client {
	return client
}

// InitClient initializes a mongo client instance
func InitClient() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	url := "mongodb://localhost:27017"
	if os.Getenv("MONGO_URL") != "" {
		url = os.Getenv("MONGO_URL")
	}
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	newClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	client = newClient
	return nil
}
