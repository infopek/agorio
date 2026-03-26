package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/models"
)

const (
	port int = 8080
	fps  int = 60

	// Game-related
	gameAreaWidth  int = 800
	gameAreaHeight int = 600
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]*models.Player)
var mutex = &sync.RWMutex{}
var worldMutex = &sync.RWMutex{}

type ClientMessage struct {
	Type string
	X    int
	Y    int
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading: ", err)
		return
	}
	defer conn.Close()

	// TODO: Get username from input
	mutex.Lock()
	player := models.NewPlayer(
		fmt.Sprintf("randika%d", math.RandRange(1000, 10000)),
		math.Vector2{
			X: math.RandRange(0, 100),
			Y: math.RandRange(0, 100),
		})
	clients[conn] = &player
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		var msg ClientMessage
		if err = json.Unmarshal(message, &msg); err != nil {
			log.Printf("Couldn't process msg: %s", err)
			continue
		}

		mutex.Lock()
		if player, ok := clients[conn]; ok {
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

func gameBroadcastState() error {
	state := models.WorldState{}

	// Player related
	for _, player := range clients {
		// Positions
		state.PlayerPositions = append(state.PlayerPositions, player.Pos)
	}

	data, err := json.Marshal(state)
	if err != nil {
		log.Println("Error encoding world state: ", err)
		return err
	}

	for conn, _ := range clients {
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
	
	 return nil
}

func gameLoop() {
	ticker := time.NewTicker(time.Duration(1000/fps) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// Update
		worldMutex.RLock()
		for _, _ = range clients {
			//err := client.WriteMessage(websocket.TextMessage, message)
			//if err != nil {
			//	cljent.Close()
			//	delete(players, client)
			//}
		}

		// Handle collisions

		// Respawn pellets
		worldMutex.RUnlock()

		// Broadcast state
		mutex.Lock()
		err := gameBroadcastState()
		if err != nil {
			log.Printf("Error while broadcasting: %s", err)
		}
		mutex.Unlock()
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./cmd/server/static")))
	http.HandleFunc("/ws", wsHandler)

	go gameLoop()

	log.Printf("WebSocket server started on port %d", port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
