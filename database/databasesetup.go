package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	connection_str, flag := os.LookupEnv("MONGO_CONN_STRING")
	if !flag {
		log.Fatal("ENV CONN STR NULL")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	clientoption := options.Client().ApplyURI(connection_str)
	client, err := mongo.Connect(ctx, clientoption)
	if err != nil {
		log.Fatalf("Error Connecting to the MongoDB: %v", err)
	}

	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	log.Fatalf("Error Pinging MongoDB: %v", err)
	// }

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// Send a ping to confirm a successful connection
	if err := client.Database("Admin").RunCommand(context.TODO(), bson.D{primitive.E{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}

var Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Ecommerce").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Ecommerce").Collection(collectionName)
	return collection
}
