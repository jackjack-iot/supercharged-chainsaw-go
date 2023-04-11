package utilities

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(uri string, maxRetries int, retryDelay time.Duration) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *mongo.Client
	var err error
	for i := 0; i < maxRetries; i++ {
		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			return client, nil
		}
		log.Printf("Attempt %d to connect to MongoDB failed: %v", i+1, err)
		time.Sleep(retryDelay)
	}

	return nil, err
}