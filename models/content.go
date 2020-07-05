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
	Path      string `bson:"path"`
}

func (db *DB) AddContent(content Content) (*mongo.InsertOneResult, error) {
	collection := db.Database("personal-site").Collection("content")

	if content.Name == "" || content.ParentDir == "" || content.Type == "" || content.Path == "" {
		return nil, errors.New("name, parentDir, type, path are required")
	}

	res, err := collection.InsertOne(context.TODO(), content)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *DB) GetContentByName(name string) (*Content, error) {
	collection := db.Database("personal-site").Collection("content")

	content := Content{}

	result := collection.FindOne(context.TODO(), bson.M{"name": name})

	if result.Err() != nil {
		return nil, result.Err()
	}

	decodeError := result.Decode(&content)

	if decodeError != nil {
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

// gets the parent dir by the name of one of its children
func (db *DB) GetContentsParentByChild(name string) (*Content, error) {
	collection := db.Database("personal-site").Collection("content")
	child := collection.FindOne(context.TODO(), bson.M{"name": name})

	if child.Err() != nil {
		return nil, child.Err()
	}

	var childContent Content
	_ = child.Decode(&childContent)

	parent := collection.FindOne(context.TODO(), bson.M{"name": childContent.ParentDir})

	if parent.Err() != nil {
		return nil, parent.Err()
	}

	var parentContent Content
	_ = parent.Decode(&parentContent)

	return &parentContent, nil
}

func (db *DB) GetFileContent(name string) string {
	collection := db.Database("personal-site").Collection("content")
	file := Content{}

	result := collection.FindOne(context.TODO(), bson.M{"name": name, "type": "file"})

	_ = result.Decode(&file)

	return file.Content
}
