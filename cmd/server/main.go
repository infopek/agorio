package main

import (
	"time"
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
	port int = 8080
	fps int = 60
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
	// TODO: Get username from input
	player := models.NewPlayer(fmt.Sprintf("randika%d", math.RandRange(1000, 10000)), math.Vector2{
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

		var msg ClientMessage
		if err = json.Unmarshal(message, &msg); err != nil {
			log.Printf("Couldn't process msg: %s", err)
			continue
		}

		mutex.Lock()
		if player, ok := players[conn]; ok {
			switch msg.Type {
			case "move":
				player.Target.X = msg.X
				player.Target.Y = msg.Y
			}
			log.Printf("Message received: %v\n", msg)
		}
		mutex.Unlock()
	}
}

func gameLoop() {
	ticker := time.NewTicker(time.Duration(1000 / fps) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		mutex.Lock()

		// Update
		for _, p := range players {
			fmt.Printf("Name: %s\n", p.Name)
			//err := client.WriteMessage(websocket.TextMessage, message)
			//if err != nil {
			//	cljent.Close()
			//	delete(players, client)
			//}
		}

		// Handle collisions

		// Respawn pellets

		// Broadcast state
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

	log.Printf("WebSocket server started on port %d", port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
