package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongodb *mongo.Client

func ConnectDB(uri string) (*mongo.Client, error) {

	// const uri = "mongodb://127.0.0.1:27017/?retryWrites=true&w=majority"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return client, err
	}
	fmt.Println("you successfully connected to MongoDB.\n")

	return client, nil
}
