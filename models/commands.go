package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Command struct {
	Name            string `json:"name"`
	ExecutionStatus bool   `json:"executionStatus"`
	Date            string `json:"date"`
	ClientIP        string `json:"clientIP"`
}

func (db *DB) SaveCommand(command string, status bool, clientIP string) (*mongo.InsertOneResult, error) {
	collection := db.Database("personal-site").Collection("commands")

	commandToSave := Command{command, status, time.Now().String(), clientIP}

	savedCommand, err := collection.InsertOne(context.TODO(), commandToSave)

	if err != nil {
		log.Println("Error saving command: ", err)

		return nil, err
	}

	return savedCommand, nil
}
