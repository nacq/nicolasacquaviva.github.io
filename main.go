package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/nicolasacquaviva/nicolasacquaviva.github.io/models"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type HttpResponse struct {
    Success bool `json:"success"`
    Message string `json:"message,omitempty"`
}

type Env struct {
    db models.Datastore
}

func health(w http.ResponseWriter, r *http.Request) {
    response := HttpResponse{}
    response.Success = true
    response.Message = "Api up and running"

    data, err := json.Marshal(response)

    if err != nil {
        log.Printf("Error parsing json")
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func checkOrigin(r *http.Request) bool {
    if os.Getenv("MODE") == "production" {
        origin := r.Header.Get("Origin")

        return origin == "https://nicolasacquaviva.com" || origin == "https://www.nicolasacquaviva.com"
    }

    return true
}

func getIP(r *http.Request) string {
    forwardedFor := r.Header.Get("x-forwarded-for")

    if forwardedFor != "" {
        return forwardedFor
    }

    return r.RemoteAddr
}

func (env *Env) ws(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = checkOrigin

    c, err := upgrader.Upgrade(w, r, nil)

    if err != nil {
        log.Print("Error upgrading: ", err)

        return
    }

    defer c.Close()

    for {
        _, message, err := c.ReadMessage()

        if err != nil {
            log.Println("Error reading message: ", err)
            break
        }

        log.Printf("Received: %s", message)

        messageParts := strings.Split(string(message), ":")

        env.db.SaveCommand(messageParts[1], messageParts[0] == "command", getIP(r))

        // err = c.WriteMessage(mt, message)

        // if err != nil {
            // log.Println("Error writing message: ", err)
            // break
        // }
    }
}

func (env *Env) attachHttpHandlers() {
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

    log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
