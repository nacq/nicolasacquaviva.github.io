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

func (db *DB) GetContentByName(name string) (*Content, error) {
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

func (db *DB) GetContentByParentDir(parentDir string) ([]string, error) {
	collection := db.Database("personal-site").Collection("content")

	content := []string{}

	// result is a cursor
	result, err := collection.Find(context.TODO(), bson.M{"parentDir": parentDir})

	if err != nil {
		log.Println("Error finding content:", err.Error())
		return nil, err
	}

	defer result.Close(context.TODO())

	for result.Next(context.TODO()) {
		var c Content

		_ = result.Decode(&c)

		if c.Type == "dir" {
			c.Name = c.Name + "/"
		}

		content = append(content, c.Name)
	}

	return content, nil
}

func (db *DB) GetFileContent(name string) string {
	collection := db.Database("personal-site").Collection("content")
	file := Content{}

	result := collection.FindOne(context.TODO(), bson.M{"name": name, "type": "file"})

	_ = result.Decode(&file)

	return file.Content
}
