package main

import (
	"log"
	"os"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/models"
	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/server"
)

func main() {
	db, err := models.NewDB(os.Getenv("MONGODB_URI"))

	if err != nil {
		log.Fatal("DB error", err)
	}

	env := &server.Env{DB: db}

	server.StartHttpServer(env)
}
