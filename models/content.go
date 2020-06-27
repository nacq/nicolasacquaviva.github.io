package models

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Content struct {
	ParentDir string `bson:"parentDir"`
	Name      string `bson:"name"`
	Type      string `bson:"type"`
	Content   string `bson:"content"`
}

func (db *DB) AddContent(content Content) (*mongo.InsertOneResult, error) {
	collection := db.Database("personal-site").Collection("content")

	if content.Name == "" || content.ParentDir == "" || content.Type == "" {
		return nil, errors.New("name, parentDir and type are required")
	}

	res, err := collection.InsertOne(context.TODO(), content)

	if err != nil {
		log.Println("Error inserting content:", err)
		return nil, err
	}

	return res, nil
}

func (db *DB) GetContent(name string) (*Content, error) {
	collection := db.Database("personal-site").Collection("content")

	content := Content{}

	result := collection.FindOne(context.TODO(), bson.M{"name": name})

	if result.Err() != nil {
		log.Println("Error finding content:", result.Err())
		return nil, result.Err()
	}

	decodeError := result.Decode(&content)

	if decodeError != nil {
		log.Println("Error decoding content:", decodeError)
		return nil, decodeError
	}

	return &content, nil
}
