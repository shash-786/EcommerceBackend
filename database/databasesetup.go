package database

import "go.mongodb.org/mongo-driver/mongo"

func DBSet() *mongo.Client {
}

var Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
}

func ProductData(clien *mongo.Client, collectionName string) *mongo.Collection {
}
