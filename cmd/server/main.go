package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/models"
)

const (
	PORT int = 8080
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var players = make(map[*websocket.Conn]*models.Player)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

type ClientMessage struct {
	Type string
	X	 int
	Y    int
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading: ", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	// TODO: Get from input
	player := models.NewPlayer("randika", math.Vector2{
		X: math.RandRange(0, 100),
		Y: math.RandRange(0, 100),
	})
	players[conn] = &player
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(players, conn)
			mutex.Unlock()
			break
		}

		if err = json.Unmarshal(message, &msg); err != nil {
			log.Printf("Couldn't process msg: %s", err)
			continue
		}
		log.Printf("Message received: %s\n", message)
	}
}

func gameLoop() {
	for {
		message := <-broadcast

		mutex.Lock()
		for client := range players {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(players, client)
			}
		}
		mutex.Unlock()
	}
}

func main() {
	_ = models.NewPlayer("NeiKer", math.Vector2{
		X: 150,
		Y: 250,
	})

	http.Handle("/", http.FileServer(http.Dir("./cmd/server/static")))
	http.HandleFunc("/ws", wsHandler)

	go gameLoop()

	log.Printf("WebSocket server started on port %d", PORT)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
