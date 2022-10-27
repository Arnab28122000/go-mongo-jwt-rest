package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx         context.Context
	mongoclient *mongo.Client
)

func DBinstance() (*mongo.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUrl := os.Getenv("MONGODB_URL")
	if mongoUrl == "" {
		log.Fatal("Unable to get MONGODB_URL")
	}
	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI(mongoUrl)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
		return nil, err
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
		return nil, err
	}

	fmt.Println("mongo connection established")

	return mongoclient, nil
}

var Client, err = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("userdb").Collection(collectionName)
	return collection
}
