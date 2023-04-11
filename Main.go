package main

import (
	"context"
	"fmt"
	"github.com/jackjack-iot/supercharged-chainsaw-go/pkg/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func main() {
	var err error
	err = MongoDbSetup()
	if err != nil {
		log.Fatal(err)
	}

	err = RabbitMqSetup()
	if err != nil {
		log.Fatal(err)
	}
}

func MongoDbSetup() error {
	client, err := utilities.ConnectMongo("mongodb://localhost:27017", 3, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return err
	}

	// Use the client to interact with the database
	// printing out all databases by name in your Mongo collection
	databases, err := client.ListDatabases(context.Background(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to get databases from MongoDB: %v", err)
		return err
	}
	for i := 0; i < len(databases.Databases); i++ {
		fmt.Println(databases.Databases[i].Name)
	}
	// ...

	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Failed to disconnect to MongoDB: %v", err)
		}
	}(client, context.Background())
	return err
}

func RabbitMqSetup() error {
	mq, err := utilities.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to create RabbitMQ connection: %v", err)
		return err
	}

	defer mq.Close()

	err = mq.DeclareExchange("test_exchange", "direct", true)
	if err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
		return err
	}


	queue, err := mq.DeclareQueue("test_queue", true)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
		return err
	}


	err = mq.BindQueue(queue.Name, "test_exchange", "test_routing_key")
	if err != nil {
		log.Fatalf("failed to bind queue: %v", err)
		return err
	}


	msgs, err := mq.Consume(queue.Name, "test_consumer", false)
	if err != nil {
		log.Fatalf("failed to start consuming messages: %v", err)
		return err
	}

	for msg := range msgs {
		log.Printf("received message: %s", msg.Body)
		msg.Ack(false)
	}

	time.Sleep(time.Second * 10)
	return nil
}