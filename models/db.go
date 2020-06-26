package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Datastore interface {
	SaveCommand(string, bool, string, string) (*mongo.InsertOneResult, error)
}

type DB struct {
	*mongo.Client
}

// Uppercased function mame means exportable
func NewDB(uri string) (*DB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error connecting to DB: ", err)

		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Cannot reach DB: ", err)

		return nil, err
	}

	return &DB{client}, nil
}
