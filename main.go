package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/models"
	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/server"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type HttpResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Env struct {
	db models.Datastore
}

func (env *Env) content(w http.ResponseWriter, r *http.Request) {
	var content models.Content

	if r.Method == "GET" {
		params, ok := r.URL.Query()["name"]

		if !ok || len(params[0]) < 1 {
			err := errors.New("'name' query param is required")
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		name := params[0]

		content, err := env.db.GetContentByName(name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(content)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	} else if r.Method == "POST" {
		_ = json.NewDecoder(r.Body).Decode(&content)

		newContent, err := env.db.AddContent(content)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		data, err := json.Marshal(newContent)

		if err != nil {
			log.Println("Error creating content:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	response := HttpResponse{}
	response.Success = true
	response.Message = "Api up and running"

	w.Header().Set("access-control-allow-origin", "*")

	data, err := json.Marshal(response)

	if err != nil {
		log.Printf("Error parsing json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (env *Env) ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = server.CheckOrigin

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Error upgrading:", err)

		return
	}

	defer c.Close()

	c.WriteMessage(websocket.TextMessage, []byte("connection:status:1"))

	for {
		messageType, message, err := c.ReadMessage()

		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		messageParts := strings.Split(string(message), ":")

		if server.IsProduction() {
			env.db.SaveCommand(
				messageParts[1],
				messageParts[0] == "command",
				server.GetIPFromRequest(r), r.Header.Get("user-agent"),
			)
		}

		commandResponse := env.executeCommand(messageParts)

		err = c.WriteMessage(messageType, []byte(commandResponse))

		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func (env *Env) executeCommand(input []string) string {
	dir := input[0]
	command := input[1]
	params := input[2]

	ChangeDirectory := server.ChangeDirectory(env.db)
	ListDirectory := server.NewListDirectory(env.db)
	PrintFileContent := server.NewPrintFileContent(env.db)
	Help := server.NewHelp(env.db)

	switch command {
	case "cd":
		return ChangeDirectory(dir, params)
	case "ls":
		return ListDirectory(dir, params)
	case "cat":
		return PrintFileContent(params)
	case "help":
		return Help()
	case "clear":
		return ""
	default:
		return "command not found: " + command + ". Try using the 'help'"
	}
}

func (env *Env) attachHttpHandlers() {
	http.HandleFunc("/content", env.content)
	http.HandleFunc("/health", health)
	http.HandleFunc("/ws", env.ws)
}

func main() {
	db, err := models.NewDB(os.Getenv("MONGODB_URI"))

	if err != nil {
		log.Fatal("DB error", err)
	}

	env := &Env{db}

	env.attachHttpHandlers()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
