package models

import (
	"context"
	"errors"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) AddContent(content types.Content) (*mongo.InsertOneResult, error) {
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

func (db *DB) GetContentByName(name string) (*types.Content, error) {
	collection := db.Database("personal-site").Collection("content")

	content := types.Content{}

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
		var c types.Content

		_ = result.Decode(&c)

		if c.Type == "dir" {
			c.Name = c.Name + "/"
		}

		content = append(content, c.Name)
	}

	return content, nil
}

func (db *DB) GetContentByPath(path string) (*types.Content, error) {
	collection := db.Database("personal-site").Collection("content")
	result := collection.FindOne(context.TODO(), bson.M{"path": path})
	content := types.Content{}

	if result.Err() != nil {
		return nil, result.Err()
	}

	_ = result.Decode(&content)

	return &content, nil
}

// gets the parent dir by the name of one of its children
func (db *DB) GetContentsParentByChild(name string) (*types.Content, error) {
	collection := db.Database("personal-site").Collection("content")
	child := collection.FindOne(context.TODO(), bson.M{"name": name})

	if child.Err() != nil {
		return nil, child.Err()
	}

	var childContent types.Content
	_ = child.Decode(&childContent)

	parent := collection.FindOne(context.TODO(), bson.M{"name": childContent.ParentDir})

	if parent.Err() != nil {
		return nil, parent.Err()
	}

	var parentContent types.Content
	_ = parent.Decode(&parentContent)

	return &parentContent, nil
}

func (db *DB) GetFileContent(name string) string {
	collection := db.Database("personal-site").Collection("content")
	file := types.Content{}

	result := collection.FindOne(context.TODO(), bson.M{"name": name, "type": "file"})

	_ = result.Decode(&file)

	return file.Content
}
