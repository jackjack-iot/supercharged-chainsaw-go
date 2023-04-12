package main

import (
	"context"
	"encoding/base32"
	"fmt"
	"github.com/jackjack-iot/supercharged-chainsaw-go/pkg/utilities"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Unable to find env file")
		return
	}

	err = MongoDbSetup()
	if err != nil {
		log.Fatal(err)
	}

	err = RabbitMqSetup()
	if err != nil {
		log.Fatal(err)
	}

	TOTP := GenerateTOTP()
	fmt.Printf("OTP: %06d\n", TOTP)
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

func GenerateTOTP() int {
	secret := os.Getenv("SECRET_KEY") // Your secret key
	if secret == "" {
		fmt.Println("Error: no secret key provided")
		return -1
	}
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println("Error decoding secret key:", err)
		return -2
	}
	counter := time.Now().Unix() / utilities.OtpInterval
	OTPU := utilities.NewOtpUtilities(key)

	otp, err := OTPU.TOTPToken(counter, utilities.SixDigits)
	if err != nil {
		fmt.Println("Error generating OTP:", err)
		return -3
	}

	return otp
}