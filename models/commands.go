package models

import (
	"context"
	"log"
	"time"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/types"

	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) SaveCommand(command string, status bool, clientIP string, userAgent string) (*mongo.InsertOneResult, error) {
	collection := db.Database("personal-site").Collection("commands")

	commandToSave := types.Command{
		Name:            command,
		ExecutionStatus: status,
		Date:            time.Now().String(),
		ClientIP:        clientIP,
		UserAgent:       userAgent,
	}

	savedCommand, err := collection.InsertOne(context.TODO(), commandToSave)

	if err != nil {
		log.Println("Error saving command: ", err)

		return nil, err
	}

	return savedCommand, nil
}
