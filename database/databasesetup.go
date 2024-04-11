package database

import (
	"context"
	"log"
  "fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://development:testpassword@localhost:27017")
  if err != nil {
    log.Fatal(err)
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

  err = client.Connect(ctx)
  if err != nil {
    log.Fatal(err)
  }

  err = client.Ping(context.TODO(), nil)
  if err != nil {
    log.Println("Failed to connect to the DB")
    return
  }

  fmt.Prinprintln("Successfullly Connected to the client")
  return client
}

var Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
}

func ProductData(clien *mongo.Client, collectionName string) *mongo.Collection {
}
