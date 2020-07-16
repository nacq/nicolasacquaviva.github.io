package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/types"

	"github.com/gorilla/websocket"
)

type Env struct {
	DB types.Datastore
}

var upgrader = websocket.Upgrader{}

func (env *Env) content(w http.ResponseWriter, r *http.Request) {
	var content types.Content

	if r.Method == "GET" {
		params, ok := r.URL.Query()["name"]

		if !ok || len(params[0]) < 1 {
			err := errors.New("'name' query param is required")
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		name := params[0]

		content, err := env.DB.GetContentByName(name)

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

		newContent, err := env.DB.AddContent(content)

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

func (env *Env) executeCommand(input []string) string {
	dir := input[0]
	command := input[1]
	params := input[2]

	ChangeDirectory := NewChangeDirectory(env.DB)
	Display := NewDisplayImage(env.DB)
	Help := NewHelp(env.DB)
	ListDirectory := NewListDirectory(env.DB)
	PrintFileContent := NewPrintFileContent(env.DB)

	switch command {
	case "cd":
		return ChangeDirectory(dir, params)
	case "ls":
		return ListDirectory(dir, params)
	case "cat":
		return PrintFileContent(params)
	case "help":
		return Help()
	case "display":
		return Display(dir + "/" + params)
	case "clear":
		return ""
	default:
		return "command not found: " + command + ". Try using the 'help'"
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	response := types.HttpResponse{}
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
	upgrader.CheckOrigin = CheckOrigin

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

		if IsProduction() {
			env.DB.SaveCommand(
				messageParts[1],
				messageParts[0] == "command",
				GetIPFromRequest(r), r.Header.Get("user-agent"),
			)
		}

		commandResponse := env.executeCommand(messageParts)

		err = c.WriteMessage(messageType, []byte(commandResponse))

		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func StartHttpServer(env *Env) {
	http.HandleFunc("/content", env.content)
	http.HandleFunc("/health", health)
	http.HandleFunc("/ws", env.ws)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
